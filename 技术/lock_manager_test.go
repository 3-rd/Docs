package lock

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// --- Basic Acquire/Release ---

func TestAcquireRelease_Success(t *testing.T) {
	m := NewLockManager(100 * time.Millisecond)
	defer m.Shutdown()

	ok, err := m.Acquire(context.Background(), "ns1", "r1", 5*time.Second)
	if err != nil {
		t.Fatalf("Acquire failed: %v", err)
	}
	if !ok {
		t.Fatal("expected to acquire lock")
	}
	if !m.IsLocked("ns1", "r1") {
		t.Fatal("expected IsLocked to be true")
	}

	err = m.Release("ns1", "r1")
	if err != nil {
		t.Fatalf("Release failed: %v", err)
	}
	if m.IsLocked("ns1", "r1") {
		t.Fatal("expected IsLocked to be false after Release")
	}
}

func TestAcquire_DifferentKeysNoConflict(t *testing.T) {
	m := NewLockManager(100 * time.Millisecond)
	defer m.Shutdown()

	ok1, _ := m.Acquire(context.Background(), "ns1", "r1", 5*time.Second)
	ok2, _ := m.Acquire(context.Background(), "ns1", "r2", 5*time.Second)
	ok3, _ := m.Acquire(context.Background(), "ns2", "r1", 5*time.Second)

	if !ok1 || !ok2 || !ok3 {
		t.Fatal("different keys should not conflict")
	}
}

func TestRelease_NotHeld(t *testing.T) {
	m := NewLockManager(100 * time.Millisecond)
	defer m.Shutdown()

	err := m.Release("ns1", "r1")
	if !errors.Is(err, ErrLockNotHeld) {
		t.Fatalf("expected ErrLockNotHeld, got: %v", err)
	}
}

// --- Timeout ---

func TestAcquire_Timeout(t *testing.T) {
	m := NewLockManager(100 * time.Millisecond)
	defer m.Shutdown()

	m.Acquire(context.Background(), "ns1", "r1", 5*time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	ok, err := m.Acquire(ctx, "ns1", "r1", 500*time.Millisecond)
	if ok {
		t.Fatal("expected failure to acquire held lock")
	}
	if !errors.Is(err, ErrLockAcquireTimeout) && !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected timeout error, got: %v", err)
	}
}

func TestAcquire_ContextCancelled(t *testing.T) {
	m := NewLockManager(100 * time.Millisecond)
	defer m.Shutdown()

	m.Acquire(context.Background(), "ns1", "r1", 5*time.Second)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()

	ok, err := m.Acquire(ctx, "ns1", "r1", 5*time.Second)
	if ok {
		t.Fatal("expected failure after context cancel")
	}
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got: %v", err)
	}
}

// --- Concurrent Acquisition ---

func TestAcquire_Concurrent(t *testing.T) {
	m := NewLockManager(100 * time.Millisecond)
	defer m.Shutdown()

	const goroutines = 10
	var acquired int32
	var wg sync.WaitGroup
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ok, _ := m.Acquire(ctx, "ns1", "r1", 2*time.Second)
			if ok {
				atomic.AddInt32(&acquired, 1)
				time.Sleep(50 * time.Millisecond)
				m.Release("ns1", "r1")
			}
		}()
	}

	wg.Wait()
	if acquired != 1 {
		t.Fatalf("expected exactly 1 goroutine to acquire the lock, got %d", acquired)
	}
}

// --- All goroutines eventually acquire ---

func TestAcquire_AllEventuallyAcquire(t *testing.T) {
	m := NewLockManager(100 * time.Millisecond)
	defer m.Shutdown()

	m.Acquire(context.Background(), "ns1", "r1", 5*time.Second)

	const n = 5
	results := make([]bool, n)
	var wg sync.WaitGroup

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			ok, _ := m.Acquire(ctx, "ns1", "r1", 3*time.Second)
			results[idx] = ok
			if ok {
				time.Sleep(10 * time.Millisecond)
				m.Release("ns1", "r1")
			}
		}(i)
	}

	time.Sleep(50 * time.Millisecond)
	m.Release("ns1", "r1")

	wg.Wait()

	acquired := 0
	for _, v := range results {
		if v {
			acquired++
		}
	}
	if acquired != n {
		t.Errorf("expected all %d to acquire, got %d", n, acquired)
	}
}

// --- Panic Safety ---

func TestAcquire_PanicSafe(t *testing.T) {
	m := NewLockManager(100 * time.Millisecond)
	defer m.Shutdown()

	var panicked bool
	func() {
		defer func() {
			if recover() != nil {
				panicked = true
			}
		}()
		ok, err := m.Acquire(context.Background(), "ns1", "r1", 5*time.Second)
		if !ok || err != nil {
			t.Fatalf("Acquire failed: %v", err)
		}
		panic("upgrade failed")
	}()

	if !panicked {
		t.Fatal("expected panic")
	}

	// Caller's defer would call Release here.
	err := m.Release("ns1", "r1")
	if err != nil {
		t.Fatalf("Release failed: %v", err)
	}
	if m.IsLocked("ns1", "r1") {
		t.Fatal("lock should not be held")
	}
}

func TestRelease_DoubleRelease(t *testing.T) {
	m := NewLockManager(100 * time.Millisecond)
	defer m.Shutdown()

	m.Acquire(context.Background(), "ns1", "r1", 5*time.Second)

	m.Release("ns1", "r1")
	err := m.Release("ns1", "r1")
	if !errors.Is(err, ErrLockNotHeld) {
		t.Fatalf("expected ErrLockNotHeld, got: %v", err)
	}
}

// --- Expired Lock Cleanup ---

func TestExpiredLock_CleanedUp(t *testing.T) {
	m := NewLockManager(50 * time.Millisecond)
	defer m.Shutdown()

	// Simulate an expired lock by directly manipulating the entry.
	key := LockKey{Namespace: "ns1", Release: "r1"}
	func() {
		m.mu.Lock()
		defer m.mu.Unlock()
		entry, exists := m.locks[key]
		if !exists {
			entry = &LockEntry{}
			m.locks[key] = entry
		}
		entry.mu.Lock()
		entry.holder = "locked"
		entry.expAt = time.Now().Add(-100 * time.Millisecond).UnixNano()
		entry.mu.Unlock()
	}()

	time.Sleep(200 * time.Millisecond)

	if m.IsLocked("ns1", "r1") {
		t.Fatal("expected expired lock to be cleaned up")
	}
}

// --- Invalid Input ---

func TestAcquire_ZeroTimeout(t *testing.T) {
	m := NewLockManager(100 * time.Millisecond)
	defer m.Shutdown()

	_, err := m.Acquire(context.Background(), "ns1", "r1", 0)
	if err == nil {
		t.Fatal("expected error for zero timeout")
	}
}

// --- IsLocked on non-existent key ---

func TestIsLocked_NonExistent(t *testing.T) {
	m := NewLockManager(100 * time.Millisecond)
	defer m.Shutdown()

	if m.IsLocked("ns1", "r1") {
		t.Fatal("expected false for non-existent lock")
	}
}
