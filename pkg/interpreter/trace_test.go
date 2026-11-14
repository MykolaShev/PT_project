package interpreter

import (
	"encoding/json"
	"pt_project/pkg/parser"
	"testing"
)

// TestJSONTracer verifies structured event recording and JSON serialization.
func TestJSONTracer(t *testing.T) {
	src := "begin X := 2 ; if X > 5 then X := 10 else skip end"
	program := parser.New(parser.NewLexer(src)).ParseProgram()

	tracer := NewJSONTracer()
	interp := New()
	interp.Tracer = tracer

	if err := interp.Interpret(program); err != nil {
		t.Fatalf("unexpected runtime error: %v", err)
	}

	if len(tracer.Events) == 0 {
		t.Fatal("expected JSONTracer to record execution events")
	}

	var hasVarChange, hasCondEval bool
	for _, ev := range tracer.Events {
		if ev.Type == EventVariableChange && ev.Variable == "X" && ev.Value == 2 {
			hasVarChange = true
		}
		if ev.Type == EventConditionEvaluated && !ev.CondRes {
			hasCondEval = true
		}
	}

	if !hasVarChange {
		t.Error("expected VariableChange event for X = 2")
	}
	if !hasCondEval {
		t.Error("expected ConditionEvaluated event with result false")
	}

	report := TraceReport{
		Events: tracer.Events,
		State:  interp.Env.State(),
	}
	data, err := json.Marshal(report)
	if err != nil || len(data) == 0 {
		t.Fatalf("failed to marshal TraceReport: %v", err)
	}
}