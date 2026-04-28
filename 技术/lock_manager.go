package lock

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

// ErrLockNotHeld is returned when Release is called by a non-holder.
var ErrLockNotHeld = errors.New("lock not held by the given holder")

// ErrLockAcquireTimeout is returned when lock acquisition times out.
var ErrLockAcquireTimeout = errors.New("lock acquisition timeout")

// LockKey uniquely identifies a lock target (namespace + release).
type LockKey struct {
	Namespace string
	Release   string
}

// String returns a human-readable key string.
func (k LockKey) String() string {
	return fmt.Sprintf("%s:%s", k.Namespace, k.Release)
}

// WaiterResult carries the result of a wait-for-lock operation.
type WaiterResult struct {
	Locked bool
	Err    error
}

// LockEntry holds the state of a lock and its waiting goroutines.
type LockEntry struct {
	mu      sync.Mutex
	holder  string
	expAt   int64  // Unixnano; 0 means not locked
	waiters []chan WaiterResult
}

// LockManager manages exclusive locks per (namespace, release).
// It is safe for concurrent use by multiple goroutines.
type LockManager struct {
	mu    sync.Mutex
	locks map[LockKey]*LockEntry

	cleanupInterval time.Duration
	stopCleanup     chan struct{}
}

// NewLockManager creates a LockManager.
// The cleanup goroutine runs until Shutdown is called.
func NewLockManager(cleanupInterval time.Duration) *LockManager {
	if cleanupInterval <= 0 {
		cleanupInterval = 30 * time.Second
	}
	m := &LockManager{
		locks:           make(map[LockKey]*LockEntry),
		cleanupInterval: cleanupInterval,
		stopCleanup:     make(chan struct{}),
	}
	go m.cleanupLoop()
	return m
}

// Shutdown stops the background cleanup goroutine.
// After Shutdown, the LockManager should not be used.
func (m *LockManager) Shutdown() {
	close(m.stopCleanup)
}

// cleanupLoop periodically removes expired locks.
func (m *LockManager) cleanupLoop() {
	ticker := time.NewTicker(m.cleanupInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			m.mu.Lock()
			now := time.Now().UnixNano()
			for key, entry := range m.locks {
				entry.mu.Lock()
				if entry.holder == "" || entry.expAt == 0 || now < entry.expAt {
					entry.mu.Unlock()
					continue
				}
				entry.holder = ""
				entry.expAt = 0
				for _, ch := range entry.waiters {
					ch <- WaiterResult{Locked: false, Err: ErrLockAcquireTimeout}
					close(ch)
				}
				entry.waiters = nil
				entry.mu.Unlock()
				delete(m.locks, key)
			}
			m.mu.Unlock()
		case <-m.stopCleanup:
			return
		}
	}
}

// Acquire attempts to acquire the lock for (namespace, release).
// If the lock is already held, it waits up to 'timeout'.
// Returns (true, nil) on success, (false, err) on timeout or cancellation.
func (m *LockManager) Acquire(
	ctx context.Context,
	namespace, release string,
	timeout time.Duration,
	holderID string,
) (bool, error) {
	if timeout <= 0 {
		return false, fmt.Errorf("timeout must be positive: %v", timeout)
	}
	key := LockKey{Namespace: namespace, Release: release}
	deadline := time.Now().Add(timeout)

	// Fast path: try to grab an existing lock or create a new one.
	if ok, _ := m.tryAcquire(key, holderID, deadline); ok {
		return true, nil
	}

	// Slow path: register as a waiter and wait.
	resultCh := make(chan WaiterResult, 1)
	func() {
		m.mu.Lock()
		defer m.mu.Unlock()
		entry, exists := m.locks[key]
		if !exists {
			entry = &LockEntry{}
			m.locks[key] = entry
		}
		entry.mu.Lock()
		defer entry.mu.Unlock()
		// Double-check after acquiring entry lock.
		if entry.holder == "" || time.Now().UnixNano() >= entry.expAt {
			entry.holder = holderID
			entry.expAt = deadline.UnixNano()
			entry.mu.Unlock()
			resultCh <- WaiterResult{Locked: true}
			close(resultCh)
			return
		}
		entry.waiters = append(entry.waiters, resultCh)
		entry.mu.Unlock()
	}()

	// Use a timer to avoid time.After goroutine leak when resultCh fires first.
	timer := time.NewTimer(time.Until(deadline))
	defer timer.Stop()

	select {
	case <-ctx.Done():
		m.removeWaiter(key, resultCh)
		return false, ctx.Err()
	case result, ok := <-resultCh:
		if !ok {
			return false, ErrLockAcquireTimeout
		}
		return result.Locked, result.Err
	case <-timer.C:
		m.removeWaiter(key, resultCh)
		return false, ErrLockAcquireTimeout
	}
}

// tryAcquire attempts a single optimistic lock acquisition.
// Returns (true) if acquired, (false) if contention exists.
func (m *LockManager) tryAcquire(key LockKey, holderID string, deadline time.Time) (bool, error) {
	m.mu.Lock()
	entry, exists := m.locks[key]
	if !exists {
		entry = &LockEntry{}
		m.locks[key] = entry
		m.mu.Unlock()
		entry.mu.Lock()
		entry.holder = holderID
		entry.expAt = deadline.UnixNano()
		entry.mu.Unlock()
		return true, nil
	}
	m.mu.Unlock()

	entry.mu.Lock()
	defer entry.mu.Unlock()
	if entry.holder == "" || time.Now().UnixNano() >= entry.expAt {
		entry.holder = holderID
		entry.expAt = deadline.UnixNano()
		return true, nil
	}
	return false, nil
}

// removeWaiter removes a waiting channel from the entry's waiters list.
// Must be called while holding m.mu.
func (m *LockManager) removeWaiter(key LockKey, ch chan WaiterResult) {
	m.mu.Lock()
	defer m.mu.Unlock()
	entry, exists := m.locks[key]
	if !exists {
		return
	}
	entry.mu.Lock()
	defer entry.mu.Unlock()
	for i, w := range entry.waiters {
		if w == ch {
			l := len(entry.waiters)
			entry.waiters[i] = entry.waiters[l-1]
			entry.waiters = entry.waiters[:l-1]
			break
		}
	}
}

// Release releases the lock for (namespace, release) held by holderID.
// Returns nil on success, ErrLockNotHeld if the holder does not own the lock.
func (m *LockManager) Release(namespace, release, holderID string) error {
	key := LockKey{Namespace: namespace, Release: release}
	m.mu.Lock()
	entry, exists := m.locks[key]
	if !exists {
		m.mu.Unlock()
		return ErrLockNotHeld
	}
	entry.mu.Lock()
	defer entry.mu.Unlock()

	if entry.holder != holderID {
		m.mu.Unlock()
		return ErrLockNotHeld
	}

	entry.holder = ""
	entry.expAt = 0
	for _, ch := range entry.waiters {
		ch <- WaiterResult{Locked: false, Err: nil}
		close(ch)
	}
	entry.waiters = nil
	m.mu.Unlock()
	return nil
}

// ForceRelease forcefully releases a lock regardless of holder, waking all waiters.
func (m *LockManager) ForceRelease(namespace, release string) {
	key := LockKey{Namespace: namespace, Release: release}
	m.mu.Lock()
	entry, exists := m.locks[key]
	if !exists {
		m.mu.Unlock()
		return
	}
	entry.mu.Lock()
	entry.holder = ""
	entry.expAt = 0
	for _, ch := range entry.waiters {
		ch <- WaiterResult{Locked: false, Err: errors.New("force released")}
		close(ch)
	}
	entry.waiters = nil
	entry.mu.Unlock()
	delete(m.locks, key)
	m.mu.Unlock()
}

// IsLocked reports whether the lock for (namespace, release) is currently held.
func (m *LockManager) IsLocked(namespace, release string) bool {
	key := LockKey{Namespace: namespace, Release: release}
	m.mu.Lock()
	entry, exists := m.locks[key]
	if !exists {
		m.mu.Unlock()
		return false
	}
	entry.mu.Lock()
	defer entry.mu.Unlock()
	m.mu.Unlock()
	return entry.holder != "" && (entry.expAt == 0 || time.Now().UnixNano() < entry.expAt)
}
