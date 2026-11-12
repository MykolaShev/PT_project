package main_test

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"

	"pt_project/pkg/interpreter"
	"pt_project/pkg/parser"
)

// TestCase matches the structure defined in testdata/suite.json
type TestCase struct {
	Name          string           `json:"name"`
	Source        string           `json:"source"`
	ExpectedState map[string]int64 `json:"expected_state"`
	ExpectError   bool             `json:"expect_error"`
}

func TestSIPLConformanceSuite(t *testing.T) {
	data, err := os.ReadFile("testdata/suite.json")
	if err != nil {
		t.Fatalf("Failed to read conformance suite: %v", err)
	}

	var tests []TestCase
	if err := json.Unmarshal(data, &tests); err != nil {
		t.Fatalf("Failed to parse JSON suite: %v", err)
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			l := parser.NewLexer(tc.Source)
			p := parser.New(l)
			program := p.ParseProgram()

			hasSyntaxErrors := len(p.Errors()) > 0

			if tc.ExpectError {
				// If we expect an error, it could be a syntax error.
				// If it has syntax errors, the test passes.
				if hasSyntaxErrors {
					return
				}
				// If no syntax errors, run interpreter to check for runtime errors
				interp := interpreter.New()
				err := interp.Interpret(program)
				if err == nil {
					t.Errorf("Expected an error but execution succeeded")
				}
				return
			}

			// If we do NOT expect an error
			if hasSyntaxErrors {
				t.Fatalf("Unexpected syntax errors: %v", p.Errors())
			}

			interp := interpreter.New()
			err = interp.Interpret(program)
			if err != nil {
				t.Fatalf("Unexpected runtime error: %v", err)
			}

			actualState := interp.Env.State()

			if !reflect.DeepEqual(tc.ExpectedState, actualState) {
				t.Errorf("State mismatch.\nExpected: %v\nGot:      %v", tc.ExpectedState, actualState)
			}
		})
	}
}