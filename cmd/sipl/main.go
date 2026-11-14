package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"pt_project/pkg/interpreter"
	"pt_project/pkg/parser"
	"pt_project/pkg/visualizer"
)

// ConsoleTracer implements the Tracer interface to print events to stdout.
type ConsoleTracer struct{}

// OnEvent is called by the interpreter for every traced step.
func (t *ConsoleTracer) OnEvent(event interpreter.TraceEvent) {
	switch event.Type {
	case interpreter.EventStatementEnter:
		fmt.Printf("[>] Line %d: %s\n", event.Span.Start.Line, event.Message)
	case interpreter.EventVariableChange:
		fmt.Printf("    [*] %s\n", event.Message)
	case interpreter.EventConditionEvaluated:
		fmt.Printf("    [?] %s\n", event.Message)
	}
}

func main() {
	jsonFlag := flag.Bool("json", false, "Output final state as JSON")
	flag.Parse()

	args := flag.Args()
	if len(args) < 2 {
		fmt.Println("Usage: sipl [--json] <command> <file.sipl>")
		fmt.Println("Commands: run, trace, parse, check, ast, cst, dot, text, web")
		os.Exit(1)
	}

	command := args[0]
	filename := args[1]

	// 1. Read source code
	source, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	// 2. Lexical and Syntax Analysis (Frontend)
	l := parser.NewLexer(string(source))
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		fmt.Println("Syntax errors:")
		for _, msg := range p.Errors() {
			fmt.Println(" -", msg)
		}
		os.Exit(1) // Exit with error code
	}

	if command == "parse" || command == "check" {
		fmt.Println("Parsing successful. AST built.")
		os.Exit(0)
	}

	if command == "ast" {
		astJSON, _ := json.MarshalIndent(program, "", "  ")
		fmt.Println(string(astJSON))
		os.Exit(0)
	}

	if command == "dot" {
		v := visualizer.NewDotVisualizer()
		fmt.Println(v.Generate(program))
		os.Exit(0)
	}

	if command == "text" {
		v := visualizer.NewTextVisualizer()
		fmt.Print(v.Generate(program, "", true))
		os.Exit(0)
	}

	if command == "web" {
		htmlContent := visualizer.GenerateHTML(program)
		err := os.WriteFile("tree.html", []byte(htmlContent), 0644)
		if err != nil {
			fmt.Printf("Error writing tree.html: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Generated tree.html successfully.")
		os.Exit(0)
	}

	// 3. Execution (Backend)
	interp := interpreter.New()

	var jsonTracer *interpreter.JSONTracer
	if command == "trace" {
		if *jsonFlag {
			jsonTracer = interpreter.NewJSONTracer()
			interp.Tracer = jsonTracer
		} else {
			interp.Tracer = &ConsoleTracer{}
		}
	}

	if command == "run" || command == "trace" {
		err := interp.Interpret(program)
		if err != nil {
			fmt.Printf("Runtime error: %v\n", err)
			os.Exit(1)
		}

		// 4. Output Final State / Trace
		if *jsonFlag {
			if command == "trace" {
				report := interpreter.TraceReport{
					Events: jsonTracer.Events,
					State:  interp.Env.State(),
				}
				out, _ := json.MarshalIndent(report, "", "  ")
				fmt.Println(string(out))
			} else {
				stateJSON, _ := json.MarshalIndent(interp.Env.State(), "", "  ")
				fmt.Println(string(stateJSON))
			}
		} else {
			fmt.Println("\n--- Final State ---")
			for k, v := range interp.Env.State() {
				fmt.Printf("%s = %d\n", k, v)
			}
		}
	} else {
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}
}