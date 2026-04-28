package lock

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

// ErrLockNotHeld is returned when Release is called but the lock is not held.
var ErrLockNotHeld = errors.New("lock not held")

// ErrLockAcquireTimeout is returned when lock acquisition times out.
var ErrLockAcquireTimeout = errors.New("lock acquisition timeout")

// LockKey uniquely identifies a lock target (namespace + release).
type LockKey struct {
	Namespace string
	Release  string
}

func (k LockKey) String() string {
	return fmt.Sprintf("%s:%s", k.Namespace, k.Release)
}

// LockEntry holds the state of a lock and its waiting goroutines.
type LockEntry struct {
	mu      sync.Mutex
	holder  string
	expAt   int64 // Unixnano; 0 means not locked
	waiters []chan struct{}
}

// LockManager manages exclusive locks per (namespace, release).
// It is safe for concurrent use by multiple goroutines.
type LockManager struct {
	mu    sync.Mutex
	locks map[LockKey]*LockEntry
}

// NewLockManager creates a LockManager.
func NewLockManager() *LockManager {
	return &LockManager{
		locks: make(map[LockKey]*LockEntry),
	}
}

// Acquire acquires the exclusive lock for (namespace, release).
// Waits up to 'timeout' for the lock. If the timeout is reached, returns false.
func (m *LockManager) Acquire(
	ctx context.Context,
	namespace, release string,
	timeout time.Duration,
) (bool, error) {
	if timeout <= 0 {
		return false, fmt.Errorf("timeout must be positive")
	}
	key := LockKey{Namespace: namespace, Release: release}
	deadline := time.Now().Add(timeout)

	// Fast path: try to acquire directly.
	if m.tryAcquire(key, deadline) {
		return true, nil
	}

	// Slow path: register as a waiter.
	waitCh := make(chan struct{})
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
			entry.holder = "locked"
			entry.expAt = deadline.UnixNano()
			entry.mu.Unlock()
			close(waitCh)
			return
		}
		entry.waiters = append(entry.waiters, waitCh)
		entry.mu.Unlock()
	}()

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		m.removeWaiter(key, waitCh)
		return false, ctx.Err()
	case <-waitCh:
		return true, nil
	case <-timer.C:
		m.removeWaiter(key, waitCh)
		return false, ErrLockAcquireTimeout
	}
}

// tryAcquire attempts a single optimistic acquisition.
// Caller must hold m.mu.
// Returns true if the lock was acquired.
func (m *LockManager) tryAcquire(key LockKey, deadline time.Time) bool {
	entry, exists := m.locks[key]
	if !exists {
		entry = &LockEntry{}
		m.locks[key] = entry
		entry.mu.Lock()
		entry.holder = "locked"
		entry.expAt = deadline.UnixNano()
		entry.mu.Unlock()
		return true
	}

	entry.mu.Lock()
	defer entry.mu.Unlock()
	if entry.holder == "" || time.Now().UnixNano() >= entry.expAt {
		entry.holder = "locked"
		entry.expAt = deadline.UnixNano()
		return true
	}
	return false
}

// removeWaiter removes a waiter channel from the entry.
// Must be called while holding m.mu.
func (m *LockManager) removeWaiter(key LockKey, ch chan struct{}) {
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
			return
		}
	}
}

// Release releases the lock for (namespace, release).
// Returns ErrLockNotHeld if the lock is not currently held.
func (m *LockManager) Release(namespace, release string) error {
	key := LockKey{Namespace: namespace, Release: release}
	m.mu.Lock()
	entry, exists := m.locks[key]
	if !exists {
		m.mu.Unlock()
		return ErrLockNotHeld
	}
	entry.mu.Lock()
	defer entry.mu.Unlock()

	if entry.holder == "" {
		m.mu.Unlock()
		return ErrLockNotHeld
	}

	entry.holder = ""
	entry.expAt = 0
	for _, ch := range entry.waiters {
		close(ch)
	}
	entry.waiters = nil
	m.mu.Unlock()
	return nil
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
