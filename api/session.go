package api

import (
	"sync"
	"time"
)

// SessionStore is an interface to an arbitrary session backend.
type SessionStore interface {
	// Create returns a new sessionID. If the backend malfunctions,
	// sessionID will be an empty string and err will be non-nil.
	Create() (sessionID string, err error)

	// Check returns whether or not sessionID is a valid session. If the backend
	// malfunctions, status will be false and err will be non-nil.
	Check(sessionID string) (status bool, err error)
}

//MemorySessionStore represents a SessionStore that uses an in-memory map
type MemorySessionStore struct {
	store    map[string]time.Time
	duration time.Duration
	mu       *sync.Mutex
}

//scavenge removes stale records every 10 minutes
func scavenge(m *MemorySessionStore) {
	for {
		time.Sleep(10 * time.Minute)
		now := time.Now()
		m.mu.Lock()
		for id, t := range m.store {
			if t.Before(now) {
				delete(m.store, id)
			}
		}
		m.mu.Unlock()
	}
}

//NewMemorySessionStore returns a new MemorySessionStore with the given expiration duration.
func NewMemorySessionStore(duration time.Duration) *MemorySessionStore {
	m := &MemorySessionStore{
		store:    make(map[string]time.Time),
		duration: duration,
		mu:       new(sync.Mutex),
	}
	go scavenge(m)
	return m
}

// Create returns a new sessionID. err will always be nil.
func (m *MemorySessionStore) Create() (sessionID string, err error) {
	id := randString(128)
	m.mu.Lock()
	m.store[id] = time.Now().Add(m.duration)
	m.mu.Unlock()
	return id, nil
}

// Check returns whether or not sessionID is a valid session. err will always be nil.
func (m *MemorySessionStore) Check(sessionID string) (status bool, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if t, ok := m.store[sessionID]; ok {
		if t.After(time.Now()) {
			return true, nil
		}
		delete(m.store, sessionID)
	}
	return false, nil
}
