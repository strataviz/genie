package events

import (
	"context"
	"sync"
)

// Manager is responsible for managing the lifecycle of all events.
type Manager struct {
	// add logging
	events   []*Event
	stopOnce sync.Once

	sync.Mutex
}

// NewManager returns a new event manager.
func NewManager() *Manager {
	return &Manager{
		events: make([]*Event, 0),
	}
}

// Start starts an event generator.
// TODO: There's a possibility that I'll allow the configuration of
// multiple sinks in the future.  It's just a single for now.
func (m *Manager) Start(ctx context.Context, event *Event, sendChan chan<- []byte) {
	m.Lock()
	defer m.Unlock()

	m.events = append(m.events, event)
	event.Start(ctx, sendChan)
}

// Stop stops all event generators.
func (m *Manager) Stop() {
	m.Lock()
	defer m.Unlock()

	m.stopOnce.Do(func() {
		for _, g := range m.events {
			g.Stop()
		}
	})
}
