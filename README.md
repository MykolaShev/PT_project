SIPL Language Toolkit (Go Implementation)

An open, extensible toolkit for SIPL (SImple Programming Language) implemented in Go. This project provides a robust foundation for formal language analysis, abstract syntax tree (AST) construction, execution tracing, and visual AST inspection.

🎯 Selected Scope & Score Breakdown (10.0 / 10.0 Points)

Task ID

Component Description

Points

Deliverable / Readiness Criteria

CORE-01

Formal SIPL Grammar Specification

0.5

Machine-readable EBNF specification + $\ge10$ positive and $\ge10$ negative SIPL test programs.

CORE-02

Typed Parser & AST with SourceSpan

1.5

Lexer and Recursive-Descent Parser producing custom AST nodes with exact source positions and rich syntax errors.

CORE-03

Tree-Walk SIPL Interpreter

2.0

Core evaluation engine supporting assignments, sequence, if, while, begin/end, skip, state tracking, and loop step limits.

CORE-04

Execution Tracing & Event Model

1.0

Programmatic and JSON-serializable execution traces (BeforeStatement, AfterStatement, variable changes, condition checks).

CORE-05

Unified Command Line Interface (CLI)

1.0

CLI binary supporting parse, check, run, and trace subcommands with CLI/JSON outputs.

CORE-06

Independent Conformance Test Suite

1.0

Portable test suite with fixtures verifying parser output, interpreter states, and error handling.

VIS-01

Parse Tree Visualizer

1.5

Generation of Parse Tree exports in JSON and Graphviz DOT/SVG formats.

VIS-02

AST Visualizer & Structural Comparison

1.0

Interactive or visual exporter for AST nodes alongside Parse Tree comparison tools.

REL-01

Continuous Integration (CI) Pipeline

0.5

Automated GitHub Actions workflow performing go vet, type checking, and unit testing on every PR/push.

TOTAL



10.0

Fully compliant with course rubric requirements.

🏗️ Project Architecture

The architecture enforces a strict separation between core language semantics and peripheral utilities (visualization, CLI).

sipl/
├── .github/
│   └── workflows/
│       └── ci.yml             # GitHub Actions automated test workflow (REL-01)
├── cmd/
│   └── sipl/
│       └── main.go            # CLI entrypoint (CORE-05)
├── grammar/
│   ├── ebnf.md                # EBNF Formal Syntax definition (CORE-01)
│   └── sipl.ebnf              # Machine-readable EBNF specification (CORE-01)
├── pkg/
│   ├── ast/                   # Typed AST structures & SourceSpan (CORE-02)
│   ├── eval/                  # Interpreter state, evaluation & trace events (CORE-03, CORE-04)
│   ├── parser/                # Lexer, Tokenizer & Parser engine (CORE-02)
│   ├── testsuite/             # Conformance test runner (CORE-06)
│   └── vis/                   # DOT, JSON, and SVG tree visualizers (VIS-01, VIS-02)
├── testdata/
│   ├── valid/                 # Positive test cases (CORE-01, CORE-06)
│   └── invalid/               # Negative test cases with syntax/runtime errors
├── go.mod                     # Go module definition
└── README.md                  # Project documentation & roadmap


📜 SIPL Syntax Specification (EBNF)

SIPL statements and expressions follow standard language structures as defined in the course curriculum:

Program     = Statement ;

Statement   = "skip"
            | Identifier ":=" ArithmeticExpr
            | Statement ";" Statement
            | "if" Condition "then" Statement "else" Statement "end"
            | "while" Condition "do" Statement "end"
            | "begin" Statement "end" ;

Condition   = ArithmeticExpr RelOp ArithmeticExpr ;
RelOp       = "=" | "!=" | "<" | "<=" | ">" | ">=" ;

ArithmeticExpr = Term { ("+" | "-") Term } ;
Term           = Factor { ("*" | "/") Factor } ;
Factor         = Number | Identifier | "(" ArithmeticExpr ")" ;


🛣️ Implementation Roadmap

[x] Phase 0: Project Setup & Specification

Initialize Go module structure.

Define EBNF grammar specification and basic fixtures (CORE-01).

[ ] Phase 1: Lexer, Parser & AST (CORE-02)

Implement position-aware Tokenizer (SourceSpan).

Construct Recursive Descent Parser building custom AST nodes.

[ ] Phase 2: Interpreter & Execution Tracing (CORE-03, CORE-04)

Implement variable state environment.

Create statement execution step model and event emitter (trace).

[ ] Phase 3: Visualizers (VIS-01, VIS-02)

Construct Graphviz DOT/SVG exporter for Concrete Parse Trees and ASTs.

Add AST vs. Parse Tree structural comparison output.

[ ] Phase 4: CLI & Conformance Suite (CORE-05, CORE-06)

Implement CLI subcommands (parse, check, run, trace).

Add standardized JSON output options.

Build fixture runner to validate expected states and syntax errors.

[ ] Phase 5: CI Integration (REL-01)

Configure GitHub Actions pipeline for automated testing and linter checks.

⚙️ Building & Testing

Prerequisites

Go 1.21 or higher.

Graphviz (dot) for SVG/PNG rendering (optional, for visual output inspection).

Running Tests

# Run unit tests across all packages
go test -v ./...

# Run conformance suite
go test -v ./pkg/testsuite/...


📄 License

This project is developed for educational purposes as part of the "Theory of Programming" university course.