package lock

import (
	"context"
	"sync"
	"time"
)

type state struct {
	deadline int64
	locked   bool
	cond     *sync.Cond
}

var m = make(map[string]*state) // key = "namespace/release"
var mu sync.Mutex

// Acquire acquires the exclusive lock for (namespace, release).
// Waits up to 'timeout' for the lock. If the timeout is reached, returns false.
func Acquire(namespace, release string, timeout time.Duration) bool {
	key := namespace + "/" + release
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	for {
		mu.Lock()
		e, exists := m[key]
		if !exists {
			// Lock is free; acquire it.
			m[key] = &state{deadline: time.Now().Add(timeout).UnixNano(), locked: true}
			mu.Unlock()
			return true
		}
		if !e.locked || time.Now().UnixNano() >= e.deadline {
			// Lock is free or expired; take it over.
			e.deadline = time.Now().Add(timeout).UnixNano()
			e.locked = true
			mu.Unlock()
			return true
		}
		// Lock is held and not expired; wait.
		if e.cond == nil {
			e.cond = sync.NewCond(&mu)
		}
		mu.Unlock()

		select {
		case <-ctx.Done():
			return false
		case <-time.After(50 * time.Millisecond):
			// Wake up, check deadline again via for-loop.
		}
	}
}

// Release releases the lock for (namespace, release).
func Release(namespace, release string) {
	key := namespace + "/" + release
	mu.Lock()
	if e, ok := m[key]; ok {
		e.locked = false
		e.deadline = 0
		if e.cond != nil {
			e.cond.Broadcast()
		}
		delete(m, key)
	}
	mu.Unlock()
}
