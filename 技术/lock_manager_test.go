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

	ok, err := m.Acquire(context.Background(), "ns1", "r1", 5*time.Second, "holder1")
	if err != nil {
		t.Fatalf("Acquire failed: %v", err)
	}
	if !ok {
		t.Fatal("expected to acquire lock")
	}
	if !m.IsLocked("ns1", "r1") {
		t.Fatal("expected IsLocked to be true")
	}

	err = m.Release("ns1", "r1", "holder1")
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

	ok1, err := m.Acquire(context.Background(), "ns1", "r1", 5*time.Second, "holder1")
	if err != nil || !ok1 {
		t.Fatal("failed to acquire ns1/r1")
	}
	ok2, err := m.Acquire(context.Background(), "ns1", "r2", 5*time.Second, "holder2")
	if err != nil || !ok2 {
		t.Fatal("failed to acquire ns1/r2")
	}
	ok3, err := m.Acquire(context.Background(), "ns2", "r1", 5*time.Second, "holder3")
	if err != nil || !ok3 {
		t.Fatal("failed to acquire ns2/r1")
	}

	m.Release("ns1", "r1", "holder1")
	m.Release("ns1", "r2", "holder2")
	m.Release("ns2", "r1", "holder3")
}

func TestRelease_NotHeld(t *testing.T) {
	m := NewLockManager(100 * time.Millisecond)
	defer m.Shutdown()

	err := m.Release("ns1", "r1", "wrong_holder")
	if !errors.Is(err, ErrLockNotHeld) {
		t.Fatalf("expected ErrLockNotHeld, got: %v", err)
	}
}

// --- Timeout ---

func TestAcquire_Timeout(t *testing.T) {
	m := NewLockManager(100 * time.Millisecond)
	defer m.Shutdown()

	_, err := m.Acquire(context.Background(), "ns1", "r1", 5*time.Second, "holder1")
	if err != nil {
		t.Fatalf("first Acquire failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	ok, err := m.Acquire(ctx, "ns1", "r1", 500*time.Millisecond, "holder2")
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

	_, err := m.Acquire(context.Background(), "ns1", "r1", 5*time.Second, "holder1")
	if err != nil {
		t.Fatalf("first Acquire failed: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()

	ok, err := m.Acquire(ctx, "ns1", "r1", 5*time.Second, "holder2")
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
		go func(id int) {
			defer wg.Done()
			ok, _ := m.Acquire(ctx, "ns1", "r1", 2*time.Second, fmt.Sprintf("holder-%d", id))
			if ok {
				atomic.AddInt32(&acquired, 1)
				time.Sleep(50 * time.Millisecond)
				m.Release("ns1", "r1", fmt.Sprintf("holder-%d", id))
			}
		}(i)
	}

	wg.Wait()
	if acquired != 1 {
		t.Fatalf("expected exactly 1 goroutine to acquire the lock, got %d", acquired)
	}
}

// --- FIFO (all goroutines eventually acquire) ---

func TestAcquire_AllGoroutinesEventuallyAcquire(t *testing.T) {
	m := NewLockManager(100 * time.Millisecond)
	defer m.Shutdown()

	_, err := m.Acquire(context.Background(), "ns1", "r1", 5*time.Second, "holder1")
	if err != nil {
		t.Fatalf("Acquire failed: %v", err)
	}

	const n = 5
	results := make([]bool, n)
	var wg sync.WaitGroup

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			ok, _ := m.Acquire(ctx, "ns1", "r1", 3*time.Second, fmt.Sprintf("holder-%d", idx))
			results[idx] = ok
			if ok {
				time.Sleep(10 * time.Millisecond)
				m.Release("ns1", "r1", fmt.Sprintf("holder-%d", idx))
			}
		}(i)
	}

	time.Sleep(50 * time.Millisecond)
	m.Release("ns1", "r1", "holder1")

	wg.Wait()

	acquiredCount := 0
	for _, v := range results {
		if v {
			acquiredCount++
		}
	}
	if acquiredCount != n {
		t.Errorf("expected all %d goroutines to acquire eventually, got %d", n, acquiredCount)
	}
}

// --- Panic Safety ---

func TestAcquireRelease_PanicSafe(t *testing.T) {
	m := NewLockManager(100 * time.Millisecond)
	defer m.Shutdown()

	var panicked bool
	func() {
		defer func() {
			if r := recover(); r != nil {
				panicked = true
			}
		}()
		ok, err := m.Acquire(context.Background(), "ns1", "r1", 5*time.Second, "holder1")
		if !ok || err != nil {
			t.Fatalf("Acquire failed: %v", err)
		}
		panic("simulated upgrade panic")
	}()

	if !panicked {
		t.Fatal("expected panic to have occurred")
	}

	// With real defer+recover in calling code, Release would be called here.
	err := m.Release("ns1", "r1", "holder1")
	if err != nil {
		t.Fatalf("Release failed after panic defer: %v", err)
	}
	if m.IsLocked("ns1", "r1") {
		t.Fatal("lock should not be held after Release")
	}
}

func TestRelease_DoubleRelease(t *testing.T) {
	m := NewLockManager(100 * time.Millisecond)
	defer m.Shutdown()

	_, err := m.Acquire(context.Background(), "ns1", "r1", 5*time.Second, "holder1")
	if err != nil {
		t.Fatalf("Acquire failed: %v", err)
	}

	err = m.Release("ns1", "r1", "holder1")
	if err != nil {
		t.Fatalf("first Release failed: %v", err)
	}

	err = m.Release("ns1", "r1", "holder1")
	if !errors.Is(err, ErrLockNotHeld) {
		t.Fatalf("expected ErrLockNotHeld on double release, got: %v", err)
	}
}

// --- Expired Lock Cleanup ---

func TestExpiredLock_CleanedUp(t *testing.T) {
	m := NewLockManager(50 * time.Millisecond)
	defer m.Shutdown()

	// Simulate a lock held by a crashed holder.
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
		entry.holder = "dead-holder"
		entry.expAt = time.Now().Add(-100 * time.Millisecond).UnixNano()
		entry.mu.Unlock()
	}()

	time.Sleep(200 * time.Millisecond)

	if m.IsLocked("ns1", "r1") {
		t.Fatal("expected expired lock to be cleaned up")
	}
}

// --- ForceRelease ---

func TestForceRelease(t *testing.T) {
	m := NewLockManager(100 * time.Millisecond)
	defer m.Shutdown()

	_, err := m.Acquire(context.Background(), "ns1", "r1", 5*time.Second, "holder1")
	if err != nil {
		t.Fatalf("Acquire failed: %v", err)
	}

	m.ForceRelease("ns1", "r1")
	if m.IsLocked("ns1", "r1") {
		t.Fatal("expected lock to be force-released")
	}

	ok, err := m.Acquire(context.Background(), "ns1", "r1", 100*time.Millisecond, "new-holder")
	if err != nil {
		t.Fatalf("Acquire after ForceRelease failed: %v", err)
	}
	if !ok {
		t.Fatal("expected to acquire after ForceRelease")
	}
}

// --- Invalid Input ---

func TestAcquire_ZeroTimeout(t *testing.T) {
	m := NewLockManager(100 * time.Millisecond)
	defer m.Shutdown()

	_, err := m.Acquire(context.Background(), "ns1", "r1", 0, "holder1")
	if err == nil {
		t.Fatal("expected error for zero timeout")
	}
}
