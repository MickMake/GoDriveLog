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
			ID:         definition.ID,
			TypedValue: NewMissingValue("no value has been read yet"),
			Unit:       definition.Unit,
			Min:        definition.Min,
			Max:        definition.Max,
			Status:     StatusUnknown,
			StaleAfter: definition.StaleAfter,
		}
	}
	return store
}

func (s *StateStore) SetValue(id string, value float64, unit string, updatedAt time.Time) SensorState {
	return s.SetTypedValue(id, NewNumericValue(value, unit), updatedAt)
}

func (s *StateStore) SetTypedValue(id string, value Value, updatedAt time.Time) SensorState {
	s.mu.Lock()
	defer s.mu.Unlock()
	state := s.states[id]
	state.ID = id
	state.TypedValue = value
	if numeric, ok := value.Numeric(); ok {
		state.Value = numeric
	}
	if value.Unit != "" {
		state.Unit = value.Unit
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
		state.TypedValue = NewErrorValue(readErr.Error())
	} else {
		state.TypedValue = NewErrorValue("")
	}
	state.UpdatedAt = updatedAt
	s.states[id] = state
	return state
}

func (s *StateStore) MarkStale(id string, now time.Time) (SensorState, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	state, ok := s.states[id]
	if !ok {
		return SensorState{}, false
	}
	staleState := withStaleStatus(state, now)
	if staleState.Status == state.Status {
		return state, false
	}
	s.states[id] = staleState
	return staleState, true
}

func (s *StateStore) Get(id string) (SensorState, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	state, ok := s.states[id]
	return state, ok
}

func (s *StateStore) GetWithStale(id string, now time.Time) (SensorState, bool) {
	state, ok := s.Get(id)
	if !ok {
		return SensorState{}, false
	}
	return withStaleStatus(state, now), true
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

func (s *StateStore) SnapshotWithStale(now time.Time) []SensorState {
	states := s.Snapshot()
	for i := range states {
		states[i] = withStaleStatus(states[i], now)
	}
	return states
}

func withStaleStatus(state SensorState, now time.Time) SensorState {
	if state.StaleAfter <= 0 || state.UpdatedAt.IsZero() || state.Status != StatusOK {
		return state
	}
	if now.Sub(state.UpdatedAt) > state.StaleAfter {
		state.Status = StatusStale
	}
	return state
}
