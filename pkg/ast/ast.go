package ast

// Position represents a single location in source code
type Position struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

// SourceSpan defines the start and end positions of an AST node
type SourceSpan struct {
	Start Position `json:"start"`
	End   Position `json:"end"`
}

// Base interface for all elements in the AST
type Node interface {
	Span() SourceSpan
}

// Represents a SIPL executable statement
type Statement interface {
	Node
	stmtNode() // Dummy to restrict interface impl
}

// Represents a boolean condition or an arithmetic expression
type Expression interface {
	Node
	exprNode() // Dummy to restrict interface impl
}

// Represent the "skip" statement
type SkipStmt struct {
	SourceSpan SourceSpan `json:"span"`
}

func (s *SkipStmt) Span() SourceSpan { return s.SourceSpan }
func (s *SkipStmt) stmtNode()        {}

// Stands for assignment operation ("Ident." := "Expression")
type AssignStmt struct {
	SourceSpan  SourceSpan `json:"span"`
	Identrifier string     `json:"indentifier"`
	Value       Expression `json:"value"`
}

func (s *AssignStmt) Span() SourceSpan { return s.SourceSpan }
func (s *AssignStmt) stmtNode()        {}

// Stands for sequential execution ("Stmt1 ; Stmt2")
type SequenceStmt struct {
	SourceSpan SourceSpan `json:"span"`
	Left       Statement  `json:"left"`
	Right      Statement  `json:"right"`
}

func (s *SequenceStmt) Span() SourceSpan { return s.SourceSpan }
func (s *SequenceStmt) stmtNode()        {}

// Stands for if/else construction
type IfStmt struct {
	SourceSpan SourceSpan `json:"span"`
	Condition  Expression `json:"condition"`
	ThenBlock  Statement  `json:"then_block"`
	ElseBlock  Statement  `json:"else_block"`
}

func (s *IfStmt) Span() SourceSpan { return s.SourceSpan }
func (s *IfStmt) stmtNode()        {}

// Stands for while construction
type WhileStmt struct {
	SourceSpan SourceSpan `json:"span"`
	Condition  Expression `json:"condition"`
	Body       Statement  `json:"body"`
}

func (s *WhileStmt) Span() SourceSpan { return s.SourceSpan }
func (s *WhileStmt) stmtNode()        {}

// Represents a "begin ... end" block
type BlockStmt struct {
	SourceSpan SourceSpan `json:"span"`
	Body       Statement  `json:"body"`
}

func (s *BlockStmt) Span() SourceSpan { return s.SourceSpan }
func (s *BlockStmt) stmtNode()        {}

// Represents an integer literal
type NumberExpr struct {
	SourceSpan SourceSpan `json:"span"`
	Value      int64      `json:"value"`
}

func (s *NumberExpr) Span() SourceSpan { return s.SourceSpan }
func (s *NumberExpr) stmtNode()        {}

// Represents a variable lookup
type IdentifierExpr struct {
	SourceSpan SourceSpan `json:"span"`
	Name       string     `json:"name"`
}

func (s *IdentifierExpr) Span() SourceSpan { return s.SourceSpan }
func (s *IdentifierExpr) stmtNode()        {}

// Represents binary operations for both arithmetic and conditions
type BinaryExpr struct {
	SourceSpan SourceSpan `json:"span"`
	Left       Expression `json:"left"`
	Operator   string     `json:"operator"`
	Right      Expression `json:"right"`
}

func (s *BinaryExpr) Span() SourceSpan { return s.SourceSpan }
func (s *BinaryExpr) stmtNode()        {}