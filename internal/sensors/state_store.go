package sensors

import (
	"sort"
	"sync"
	"time"
)

type StateStore struct {
	mu     sync.RWMutex
	states map[string]SensorState
}

func NewStateStore(definitions []SensorDefinition) *StateStore {
	store := &StateStore{states: make(map[string]SensorState, len(definitions))}
	for _, definition := range definitions {
		if definition.ID == "" {
			continue
		}
		store.states[definition.ID] = SensorState{
			ID:     definition.ID,
			Unit:   definition.Unit,
			Min:    definition.Min,
			Max:    definition.Max,
			Status: StatusUnknown,
		}
	}
	return store
}

func (s *StateStore) SetValue(id string, value float64, unit string, updatedAt time.Time) SensorState {
	s.mu.Lock()
	defer s.mu.Unlock()
	state := s.states[id]
	state.ID = id
	state.Value = value
	if unit != "" {
		state.Unit = unit
	}
	state.Status = StatusOK
	state.Error = ""
	state.UpdatedAt = updatedAt
	s.states[id] = state
	return state
}

func (s *StateStore) SetError(id string, readErr error, updatedAt time.Time) SensorState {
	s.mu.Lock()
	defer s.mu.Unlock()
	state := s.states[id]
	state.ID = id
	state.Status = StatusError
	state.Error = ""
	if readErr != nil {
		state.Error = readErr.Error()
	}
	state.UpdatedAt = updatedAt
	s.states[id] = state
	return state
}

func (s *StateStore) Get(id string) (SensorState, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	state, ok := s.states[id]
	return state, ok
}

func (s *StateStore) GetWithStale(id string, staleAfter time.Duration, now time.Time) (SensorState, bool) {
	state, ok := s.Get(id)
	if !ok {
		return SensorState{}, false
	}
	return withStaleStatus(state, staleAfter, now), true
}

func (s *StateStore) Snapshot() []SensorState {
	s.mu.RLock()
	defer s.mu.RUnlock()
	states := make([]SensorState, 0, len(s.states))
	for _, state := range s.states {
		states = append(states, state)
	}
	sort.Slice(states, func(i, j int) bool {
		return states[i].ID < states[j].ID
	})
	return states
}

func (s *StateStore) SnapshotWithStale(staleAfter time.Duration, now time.Time) []SensorState {
	states := s.Snapshot()
	for i := range states {
		states[i] = withStaleStatus(states[i], staleAfter, now)
	}
	return states
}

func withStaleStatus(state SensorState, staleAfter time.Duration, now time.Time) SensorState {
	if staleAfter <= 0 || state.UpdatedAt.IsZero() || state.Status != StatusOK {
		return state
	}
	if now.Sub(state.UpdatedAt) > staleAfter {
		state.Status = StatusStale
	}
	return state
}
