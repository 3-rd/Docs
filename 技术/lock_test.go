package lock

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// --- Basic Acquire/Release ---

func TestAcquireRelease_Success(t *testing.T) {
	ok := Acquire("ns1", "r1", 5*time.Second)
	if !ok {
		t.Fatal("expected to acquire lock")
	}
	Release("ns1", "r1")

	// Acquire again should succeed.
	ok = Acquire("ns1", "r1", 5*time.Second)
	if !ok {
		t.Fatal("expected to re-acquire after release")
	}
	Release("ns1", "r1")
}

func TestAcquire_DifferentKeysNoConflict(t *testing.T) {
	Acquire("ns1", "r1", 5*time.Second)
	Acquire("ns1", "r2", 5*time.Second)
	Acquire("ns2", "r1", 5*time.Second)
	Release("ns1", "r1")
	Release("ns1", "r2")
	Release("ns2", "r1")
}

// --- Timeout ---

func TestAcquire_Timeout(t *testing.T) {
	Acquire("ns1", "r1", 5*time.Second)

	ok := Acquire("ns1", "r1", 200*time.Millisecond)
	if ok {
		t.Fatal("expected failure to acquire held lock")
	}
	Release("ns1", "r1")
}

func TestAcquire_TimeoutExceeded(t *testing.T) {
	// Acquire with 5s timeout while held.
	Acquire("ns1", "r1", 5*time.Second)

	ok := Acquire("ns1", "r1", 300*time.Millisecond)
	if ok {
		t.Fatal("expected timeout failure")
	}
	Release("ns1", "r1")
}

// --- Concurrent Acquisition ---

func TestAcquire_Concurrent(t *testing.T) {
	const goroutines = 10
	var acquired int32
	var wg sync.WaitGroup

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ok := Acquire("ns1", "r1", 2*time.Second)
			if ok {
				atomic.AddInt32(&acquired, 1)
				time.Sleep(50 * time.Millisecond)
				Release("ns1", "r1")
			}
		}()
	}

	wg.Wait()
	if acquired != 1 {
		t.Fatalf("expected exactly 1 goroutine to acquire the lock, got %d", acquired)
	}
}

// --- Expired lock is reclaimed ---

func TestAcquire_ExpiredLockReclaimed(t *testing.T) {
	// Simulate an expired lock by directly setting it in the map.
	key := "ns1/r1"
	mu.Lock()
	m[key] = &state{deadline: time.Now().Add(-100 * time.Millisecond).UnixNano()}
	mu.Unlock()

	ok := Acquire("ns1", "r1", 100*time.Millisecond)
	if !ok {
		t.Fatal("expected to acquire expired lock")
	}
	Release("ns1", "r1")
}

// --- Panic Safety ---

func TestAcquire_PanicSafe(t *testing.T) {
	var panicked bool
	func() {
		defer func() {
			if recover() != nil {
				panicked = true
			}
		}()
		ok := Acquire("ns1", "r1", 5*time.Second)
		if !ok {
			t.Fatal("expected to acquire lock")
		}
		panic("upgrade failed")
	}()

	if !panicked {
		t.Fatal("expected panic")
	}

	Release("ns1", "r1")
}

// --- Double Release ---

func TestRelease_DoubleRelease(t *testing.T) {
	Acquire("ns1", "r1", 5*time.Second)
	Release("ns1", "r1")
	// Second release should be a no-op, not panic.
	Release("ns1", "r1")

	// Should be able to re-acquire.
	ok := Acquire("ns1", "r1", 100*time.Millisecond)
	if !ok {
		t.Fatal("expected to acquire after double release")
	}
}

// --- Stress ---

func TestAcquire_Stress(t *testing.T) {
	const n = 50
	var wg sync.WaitGroup

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 20; j++ {
				ok := Acquire("ns1", "r1", 500*time.Millisecond)
				if ok {
					Release("ns1", "r1")
				}
			}
		}()
	}

	wg.Wait()
}

// --- Benchmark ---

func BenchmarkAcquireRelease(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Acquire("ns1", "r1", 5*time.Second)
		Release("ns1", "r1")
	}
}

func BenchmarkAcquireContention(b *testing.B) {
	const goroutines = 10
	var wg sync.WaitGroup
	b.ResetTimer()

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < b.N/goroutines; j++ {
				Acquire("ns1", "r1", 1*time.Second)
				Release("ns1", "r1")
			}
		}()
	}
	wg.Wait()
}

// Internal helpers for testing.

func getState(key string) (*state, bool) {
	mu.Lock()
	defer mu.Unlock()
	e, ok := m[key]
	return e, ok
}
