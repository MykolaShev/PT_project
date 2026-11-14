# SIPL Language Toolkit (Go Implementation)

[![SIPL CI](https://github.com/MykolaShev/PT_project/actions/workflows/test.yml/badge.svg)](https://github.com/<твій-github-юзернейм>/<назва-репозиторію>/actions)

An open, extensible toolkit for **SIPL (Simple Programming Language)** implemented in Go. This project provides a robust foundation for formal language analysis, abstract syntax tree (AST) construction, execution tracing, and visual AST inspection.

## 🎯 Selected Scope & Score Breakdown (10.0 / 10.0 Points)

| Task ID | Component Description | Points | Deliverable / Readiness Criteria |
|---|---|---:|---|
| **CORE-01** | Formal SIPL Grammar Specification | **0.5** | Machine-readable EBNF specification + SIPL test programs. |
| **CORE-02** | Typed Parser & AST with `SourceSpan` | **1.5** | Lexer and Recursive-Descent/Pratt Parser producing custom AST nodes with exact source positions and syntax errors. |
| **CORE-03** | Tree-Walk SIPL Interpreter | **2.0** | Core evaluation engine supporting assignments, sequence, `if`, `while`, `begin...end`, `skip`, and state tracking. |
| **CORE-04** | Execution Tracing & Event Model | **1.0** | Programmatic execution traces (`StatementEnter`, `VariableChange`, `ConditionEvaluated`). |
| **CORE-05** | Unified Command Line Interface (CLI) | **1.0** | CLI binary supporting `parse`, `run`, `trace`, `ast`, and `dot` commands with JSON outputs. |
| **CORE-06** | Independent Conformance Test Suite | **1.0** | Data-driven test suite (`suite.json`) verifying parser output, interpreter states, and runtime errors. |
| **VIS-01** | JSON AST Exporter | **1.5** | Generation of deeply nested AST structures in formatted JSON via reflection. |
| **VIS-02** | Graphviz DOT Visualizer | **1.0** | Generator for Graphviz DOT format to render hierarchical visual graphs of the AST. |
| **REL-01** | Continuous Integration (CI) Pipeline | **0.5** | Automated GitHub Actions workflow (`test.yml`) running the conformance suite on push/PR. |
| **TOTAL** | | **10.0** | **Fully compliant with course rubric requirements.** |

## 🏗️ Project Architecture

The architecture enforces a strict separation between core language semantics (frontend/backend) and peripheral utilities (visualization, CLI).

```text
sipl/
├── .github/
│   └── workflows/
│       └── test.yml           # GitHub Actions automated CI workflow (REL-01)
├── cmd/
│   └── sipl/
│       └── main.go            # CLI entrypoint (run, trace, parse, ast, dot) (CORE-05)
├── docs/
│   └── grammar.ebnf           # Formal EBNF language specification (CORE-01)
├── pkg/
│   ├── ast/
│   │   └── ast.go             # Typed AST structures & SourceSpan with JSON tags (CORE-02, VIS-01)
│   ├── parser/
│   │   ├── token.go           # Tokens and Lexer types (CORE-02)
│   │   ├── lexer.go           # Lexical analyzer
│   │   └── parser.go          # Recursive Descent & Pratt Parser (CORE-02)
│   ├── interpreter/
│   │   ├── environment.go     # Variable state management (CORE-03)
│   │   ├── interpreter.go     # Tree-walk execution engine (CORE-03)
│   │   └── trace.go           # Execution observer and event model (CORE-04)
│   └── visualizer/
│       └── dot.go             # Graphviz DOT generator for AST visualization (VIS-02)
├── testdata/
│   ├── README.md              # Test suite documentation
│   └── suite.json             # Data-driven conformance tests (CORE-06)
├── conformance_test.go        # Go Test Runner for suite.json (CORE-06)
├── go.mod                     # Go module definition
├── Makefile                   # Build and execution shortcuts
└── test.sipl                  # Scratchpad for local testing
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

## 📊 JSON Output Schema (CORE-05)
When executing with `--json run` or `ast`, the CLI outputs a stable JSON schema for integrations:
- **Trace/Run State**: Output is a flat key-value mapping `map[string]int64` representing the final environment variables.
- **AST Nodes**: Objects containing standard fields:
  - `span`: `{ "start": {"line": int, "column": int}, "end": {"line": int, "column": int} }`
  - Control flows map specific sub-blocks (e.g., `IfStmt` contains `Condition`, `ThenBlock`, `ElseBlock`).

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
  - Export AST as formatted JSON structure.
  - Build Graphviz DOT generator for visual tree representation.

- [x] **Phase 4: CLI & Conformance Suite (`CORE-05`, `CORE-06`)**
  - Implement CLI router (`make run`, `make trace`, `make ast`, `make dot`).
  - Create data-driven JSON test suite and Go test runner.

- [x] **Phase 5: CI Integration (`REL-01`)**
  - Configure GitHub Actions pipeline for automated `go test` runs.

## ⚙️ Building & Testing

### Prerequisites

- Go `1.22` or higher.
- [Graphviz](https://graphviz.org/) or an online viewer (like WebGraphviz) to render `.dot` exports.

### Quick Commands (Makefile)

```bash
# Execute test.sipl program
make run

# Execute with debug trace events
make trace

# Generate AST JSON representation
make ast

# Generate Graphviz DOT visualization
make dot

# Run the complete conformance test suite
make test
```

## 📄 License

This project is developed for educational purposes as part of the "Theory of Programming" university course.