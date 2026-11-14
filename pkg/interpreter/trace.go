package interpreter

import "pt_project/pkg/ast"

// EventType defines what kind of trace event occurred.
type EventType string

const (
	EventStatementEnter     EventType = "StatementEnter"
	EventStatementExit      EventType = "StatementExit"
	EventVariableChange     EventType = "VariableChange"
	EventConditionEvaluated EventType = "ConditionEvaluated"
)

// TraceEvent represents a single step in the execution trace.
type TraceEvent struct {
	Type     EventType      `json:"type"`
	Span     ast.SourceSpan `json:"span,omitempty"`
	Message  string         `json:"message"`
	Variable string         `json:"variable,omitempty"`
	Value    int64          `json:"value,omitempty"`
	CondRes  bool           `json:"cond_res,omitempty"`
}

// Tracer interface allows different implementations to collect execution events.
type Tracer interface {
	OnEvent(event TraceEvent)
}

// JSONTracer collects execution events in memory for structured JSON export.
type JSONTracer struct {
	Events []TraceEvent `json:"events"`
}

// NewJSONTracer initializes an empty JSONTracer.
func NewJSONTracer() *JSONTracer {
	return &JSONTracer{
		Events: []TraceEvent{},
	}
}

// OnEvent appends a trace event to the recorded slice.
func (t *JSONTracer) OnEvent(event TraceEvent) {
	t.Events = append(t.Events, event)
}

// TraceReport represents the structured JSON output for traced execution.
type TraceReport struct {
	Events []TraceEvent     `json:"events"`
	State  map[string]int64 `json:"state"`
}