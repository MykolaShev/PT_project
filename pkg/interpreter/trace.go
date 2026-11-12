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