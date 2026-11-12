package interpreter

// Environment stores the state of variables during execution.
type Environment struct {
	store map[string]int64
}

// NewEnvironment creates a new, empty variable state.
func NewEnvironment() *Environment {
	return &Environment{
		store: make(map[string]int64),
	}
}

// Get retrieves the value of a variable. Returns 0 if not found.
func (e *Environment) Get(name string) int64 {
	return e.store[name]
}

// Set updates or creates a variable with a new value.
func (e *Environment) Set(name string, val int64) {
	e.store[name] = val
}

// State returns a copy of the current variable state for testing or JSON export.
func (e *Environment) State() map[string]int64 {
	state := make(map[string]int64)
	for k, v := range e.store {
		state[k] = v
	}
	return state
}