# SIPL Language Toolkit (Go Implementation)

[![SIPL CI](https://github.com/MykolaShev/PT_project/actions/workflows/test.yml/badge.svg)](https://github.com/MykolaShev/PT_project/actions/workflows/test.yml)

An open, extensible toolkit for **SIPL (Simple Programming Language)** implemented in Go. This project provides formal language analysis, Concrete Syntax Tree (Parse Tree) and Abstract Syntax Tree (AST) construction, event-driven execution tracing, and interactive structural visualization.

## 🎯 Selected Scope & Score Breakdown (10.0 / 10.0 Points)

| Task ID | Component Description | Points | Deliverable / Readiness Criteria |
|---|---|---:|---|
| **CORE-01** | Formal SIPL Grammar Specification | **0.5** | Machine-readable EBNF (`docs/grammar.ebnf`) + 10 positive and 10 negative SIPL test cases (reused via `testdata/suite.json`). |
| **CORE-02** | Typed Parser & AST with `SourceSpan` | **1.5** | Lexer and Recursive-Descent/Pratt Parser producing typed AST nodes with exact `[line:col]` source spans and syntax diagnostics. |
| **CORE-03** | Tree-Walk SIPL Interpreter | **2.0** | Evaluation engine supporting `:=`, `;`, `if/then/else`, `while/do`, `begin...end`, `skip`, infinite loop step guard (`MaxSteps = 10000`), and >90% core test coverage. |
| **CORE-04** | Execution Tracing & Event Model | **1.0** | Observer-based `Tracer` emitting `StatementEnter`, `StatementExit`, `VariableChange`, and `ConditionEvaluated` events in CLI and JSON modes. |
| **CORE-05** | Unified Command Line Interface (CLI) | **1.0** | CLI binary supporting `run`, `trace`, `parse`, `check`, `ast`, `cst`, `dot`, `text`, and `web` with `--json` output and integration tests. |
| **CORE-06** | Independent Conformance Test Suite | **1.0** | Portable JSON test suite (`testdata/suite.json`) with documented schema (`testdata/README.md`) verifying states, syntax errors, and runtime limits. |
| **VIS-01** | Parse Tree & AST Visualizers | **1.5** | Multi-format exporters (`text`, `json`, `dot`, `html`) where every node carries its `[line:col]` source position. |
| **VIS-02** | Parse Tree ↔ AST Structural Comparison | **1.0** | Separate Concrete Syntax Tree (`pkg/parser/cst.go`) and AST representations rendered side-by-side in interactive D3.js HTML (`make web`) with synchronized span highlighting. |
| **REL-01** | Continuous Integration (CI) Pipeline | **0.5** | Automated GitHub Actions workflow (`.github/workflows/test.yml`) running the test suite on Go `1.25.3`. |
| **TOTAL** | | **10.0** | **Fully compliant with course rubric requirements.** |

## 🏗️ Project Architecture

The architecture enforces a strict separation between core language semantics (frontend/backend) and peripheral utilities (visualization, CLI).

```text
PT_project/
├── .github/
│   └── workflows/
│       └── test.yml           # GitHub Actions automated CI workflow (REL-01)
├── cmd/
│   └── sipl/
│       ├── main.go            # CLI entrypoint (run, trace, parse, check, ast, cst, dot, text, web) (CORE-05)
│       └── main_test.go       # CLI integration tests for exit codes (CORE-05)
├── docs/
│   └── grammar.ebnf           # Formal EBNF language specification (CORE-01)
├── pkg/
│   ├── ast/
│   │   └── ast.go             # Typed AST structures & SourceSpan with JSON tags (CORE-02, VIS-01)
│   ├── parser/
│   │   ├── token.go           # Tokens and Lexer types (CORE-02)
│   │   ├── lexer.go           # Lexical analyzer (CORE-02)
│   │   ├── parser.go          # Recursive Descent & Pratt AST Parser (CORE-02)
│   │   └── cst.go             # Concrete Syntax Tree (Parse Tree) builder (VIS-01, VIS-02)
│   ├── interpreter/
│   │   ├── environment.go     # Variable state management (CORE-03)
│   │   ├── interpreter.go     # Tree-walk execution engine with MaxSteps guard (CORE-03)
│   │   └── trace.go           # Observer-based Console & JSON event tracers (CORE-04)
│   └── visualizer/
│       ├── dot.go             # Graphviz DOT generator with [line:col] spans (VIS-01)
│       ├── text.go            # ASCII tree generator with [line:col] spans (VIS-01)
│       └── html.go            # Interactive D3.js Parse Tree vs AST visualizer (VIS-02)
├── testdata/
│   ├── README.md              # Conformance suite format documentation (CORE-06)
│   └── suite.json             # 10 positive & 10 negative conformance tests (CORE-01, CORE-06)
├── conformance_test.go        # Data-driven Go test runner for suite.json (CORE-06)
├── go.mod                     # Go module definition (Go 1.25.3)
├── Makefile                   # Build, test, coverage, and visualization shortcuts
└── test.sipl                  # Sample SIPL program for quick CLI runs
```

## 📜 SIPL Syntax Specification (EBNF)

SIPL syntax adheres to the academic specification, utilizing unified expressions for conditions and arithmetic operations to align with standard parser implementations.

```ebnf
Program            ::= "begin" Statement "end"

Statement          ::= AssignStmt
                     | SequenceStmt
                     | IfStmt
                     | WhileStmt
                     | BlockStmt
                     | SkipStmt

SequenceStmt       ::= Statement ";" Statement
AssignStmt         ::= Identifier ":=" Expression
IfStmt             ::= "if" Expression "then" Statement [ "else" Statement ]
WhileStmt          ::= "while" Expression "do" Statement
BlockStmt          ::= "begin" Statement "end"
SkipStmt           ::= "skip"

Expression         ::= MathExpr [ RelationalOp MathExpr ]
MathExpr           ::= Term ( AddOp Term )*
Term               ::= Factor ( MulOp Factor )*
Factor             ::= Number | Identifier | "(" Expression ")"

RelationalOp       ::= "=" | "<" | ">"
AddOp              ::= "+" | "-"
MulOp              ::= "*" | "/"

Identifier         ::= Letter ( Letter | Digit )*
Number             ::= [ "-" ] Digit ( Digit )*
```

## 📊 CLI Commands & JSON Schemas (`CORE-05`)

- `sipl parse <file>` / `sipl check <file>` — validates syntax and builds AST (exit `0` on success, `1` on syntax error).
- `sipl [--json] run <file>` — executes program; `--json` outputs `map[string]int64` final variable state.
- `sipl [--json] trace <file>` — executes with event tracing; `--json` outputs `{"events": [...], "state": {...}}`.
- `sipl ast <file>` — outputs typed AST in JSON format with `span` coordinates.
- `sipl [--json] cst <file>` — outputs Concrete Syntax Tree (Parse Tree) in Text or JSON format.
- `sipl text <file>` / `sipl dot <file>` — outputs AST with `[line:col]` spans in ASCII tree or Graphviz DOT format.
- `sipl web <file>` — generates `tree.html` with side-by-side interactive D3.js comparison of Parse Tree and AST.

## 🛣️ Implementation Roadmap

- [x] **Phase 0: Project Setup & Specification (`CORE-01`)**
  - Initialize Go module structure.
  - Define formal EBNF grammar specification.

- [x] **Phase 1: Lexer, Parser & AST (`CORE-02`)**
  - Implement Lexer and typed AST nodes with `SourceSpan`.
  - Construct Recursive Descent Parser for control flow.
  - Construct Pratt Parser for operator precedence.

- [x] **Phase 2: Interpreter & Execution Tracing (`CORE-03`, `CORE-04`)**
  - Implement Tree-Walk Interpreter and memory Environment.
  - Create execution observer (Tracer) for state change events.

- [x] **Phase 3: Visualizers (`VIS-01`, `VIS-02`)**
  - Build Concrete Syntax Tree (`CSTParser`) alongside typed AST with `[line:col]` spans.
  - Export trees in Text, JSON, and Graphviz DOT formats.
  - Generate interactive D3.js side-by-side Parse Tree vs. AST comparison (`make web`).

- [x] **Phase 4: CLI & Conformance Suite (`CORE-05`, `CORE-06`)**
  - Implement CLI router (`make run`, `make trace`, `make ast`, `make dot`).
  - Create data-driven JSON test suite and Go test runner.

- [x] **Phase 5: CI Integration (`REL-01`)**
  - Configure GitHub Actions pipeline for automated `go test` runs.

## ⚙️ Building & Testing

### Prerequisites

- Go `1.25` or higher.
- [Graphviz](https://graphviz.org/) or an online viewer (like WebGraphviz) to render `.dot` exports.

## ⚙️ Quick Commands (`Makefile`)

```bash
make run      # Execute test.sipl
make trace    # Execute test.sipl with event trace
make parse    # Validate syntax of test.sipl
make ast      # Print AST JSON
make cst      # Print Concrete Parse Tree (Text)
make dot      # Print Graphviz DOT with [line:col] spans
make text     # Print ASCII AST with [line:col] spans
make web      # Generate tree.html (Parse Tree vs AST) and open in browser
make test     # Run conformance suite and CLI integration tests
make cover    # Measure core package statement coverage
```

## 📄 License

This project is developed for educational purposes as part of the "Theory of Programming" university course.