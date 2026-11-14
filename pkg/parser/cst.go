package parser

import (
	"fmt"
	"pt_project/pkg/ast"
	"strings"
)

// ParseNode represents a node in the Concrete Syntax Tree (Parse Tree).
type ParseNode struct {
	Symbol   string         `json:"symbol"`
	Token    string         `json:"token,omitempty"`
	Span     ast.SourceSpan `json:"span"`
	Children []*ParseNode   `json:"children,omitempty"`
}

// CSTParser builds a Concrete Syntax Tree preserving grammar non-terminals and tokens.
type CSTParser struct {
	l         *Lexer
	curToken  Token
	peekToken Token
}

// NewCSTParser creates a new CSTParser from a Lexer.
func NewCSTParser(l *Lexer) *CSTParser {
	p := &CSTParser{l: l}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *CSTParser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *CSTParser) makeTerminal(tok Token) *ParseNode {
	pos := ast.Position{Line: tok.Line, Column: tok.Column}
	return &ParseNode{
		Symbol: string(tok.Type),
		Token:  tok.Literal,
		Span:   ast.SourceSpan{Start: pos, End: pos},
	}
}

// ParseCST parses the token stream into a full Concrete Syntax Tree.
func (p *CSTParser) ParseCST() *ParseNode {
	root := &ParseNode{Symbol: "Program"}
	for p.curToken.Type != EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			root.Children = append(root.Children, stmt)
		}
		p.nextToken()
	}
	if len(root.Children) > 0 {
		root.Span = ast.SourceSpan{
			Start: root.Children[0].Span.Start,
			End:   root.Children[len(root.Children)-1].Span.End,
		}
	}
	return root
}

func (p *CSTParser) parseStatement() *ParseNode {
	left := p.parseBasicStatement()
	if left == nil {
		return nil
	}

	for p.peekToken.Type == SEMICOLON {
		p.nextToken()
		semiNode := p.makeTerminal(p.curToken)
		p.nextToken()
		right := p.parseBasicStatement()
		if right == nil {
			break
		}
		left = &ParseNode{
			Symbol:   "SequenceStmt",
			Span:     ast.SourceSpan{Start: left.Span.Start, End: right.Span.End},
			Children: []*ParseNode{left, semiNode, right},
		}
	}
	return left
}

func (p *CSTParser) parseBasicStatement() *ParseNode {
	switch p.curToken.Type {
	case BEGIN:
		beginTok := p.makeTerminal(p.curToken)
		p.nextToken()
		body := p.parseStatement()
		var children []*ParseNode
		children = append(children, beginTok)
		if body != nil {
			children = append(children, body)
		}
		endSpan := beginTok.Span.End
		if p.peekToken.Type == END {
			p.nextToken()
			endTok := p.makeTerminal(p.curToken)
			endSpan = endTok.Span.End
			children = append(children, endTok)
		}
		return &ParseNode{
			Symbol:   "BlockStmt",
			Span:     ast.SourceSpan{Start: beginTok.Span.Start, End: endSpan},
			Children: children,
		}
	case IDENT:
		idTok := p.makeTerminal(p.curToken)
		if p.peekToken.Type != ASSIGN {
			return idTok
		}
		p.nextToken()
		assignTok := p.makeTerminal(p.curToken)
		p.nextToken()
		expr := p.parseExpression()
		endSpan := assignTok.Span.End
		children := []*ParseNode{idTok, assignTok}
		if expr != nil {
			endSpan = expr.Span.End
			children = append(children, expr)
		}
		return &ParseNode{
			Symbol:   "AssignStmt",
			Span:     ast.SourceSpan{Start: idTok.Span.Start, End: endSpan},
			Children: children,
		}
	case WHILE:
		whileTok := p.makeTerminal(p.curToken)
		p.nextToken()
		cond := p.parseCondition()
		children := []*ParseNode{whileTok}
		if cond != nil {
			children = append(children, cond)
		}
		if p.peekToken.Type == DO {
			p.nextToken()
			children = append(children, p.makeTerminal(p.curToken))
		}
		p.nextToken()
		body := p.parseBasicStatement()
		endSpan := whileTok.Span.End
		if body != nil {
			endSpan = body.Span.End
			children = append(children, body)
		}
		return &ParseNode{
			Symbol:   "WhileStmt",
			Span:     ast.SourceSpan{Start: whileTok.Span.Start, End: endSpan},
			Children: children,
		}
	case IF:
		ifTok := p.makeTerminal(p.curToken)
		p.nextToken()
		cond := p.parseCondition()
		children := []*ParseNode{ifTok}
		if cond != nil {
			children = append(children, cond)
		}
		if p.peekToken.Type == THEN {
			p.nextToken()
			children = append(children, p.makeTerminal(p.curToken))
		}
		p.nextToken()
		thenStmt := p.parseBasicStatement()
		endSpan := ifTok.Span.End
		if thenStmt != nil {
			endSpan = thenStmt.Span.End
			children = append(children, thenStmt)
		}
		if p.peekToken.Type == ELSE {
			p.nextToken()
			children = append(children, p.makeTerminal(p.curToken))
			p.nextToken()
			elseStmt := p.parseBasicStatement()
			if elseStmt != nil {
				endSpan = elseStmt.Span.End
				children = append(children, elseStmt)
			}
		}
		return &ParseNode{
			Symbol:   "IfStmt",
			Span:     ast.SourceSpan{Start: ifTok.Span.Start, End: endSpan},
			Children: children,
		}
	case SKIP:
		skipTok := p.makeTerminal(p.curToken)
		return &ParseNode{
			Symbol:   "SkipStmt",
			Span:     skipTok.Span,
			Children: []*ParseNode{skipTok},
		}
	default:
		return p.makeTerminal(p.curToken)
	}
}

func (p *CSTParser) parseCondition() *ParseNode {
	left := p.parseExpression()
	if left == nil {
		return nil
	}
	switch p.peekToken.Type {
	case EQ, NOTEQ, LT, LTE, GT, GTE:
		p.nextToken()
		opTok := p.makeTerminal(p.curToken)
		p.nextToken()
		right := p.parseExpression()
		endSpan := opTok.Span.End
		children := []*ParseNode{left, opTok}
		if right != nil {
			endSpan = right.Span.End
			children = append(children, right)
		}
		return &ParseNode{
			Symbol:   "Condition",
			Span:     ast.SourceSpan{Start: left.Span.Start, End: endSpan},
			Children: children,
		}
	}
	return &ParseNode{
		Symbol:   "Condition",
		Span:     left.Span,
		Children: []*ParseNode{left},
	}
}

func (p *CSTParser) parseExpression() *ParseNode {
	left := p.parseTerm()
	if left == nil {
		return nil
	}
	node := &ParseNode{
		Symbol:   "Expression",
		Span:     left.Span,
		Children: []*ParseNode{left},
	}
	for p.peekToken.Type == PLUS || p.peekToken.Type == MINUS {
		p.nextToken()
		opTok := p.makeTerminal(p.curToken)
		p.nextToken()
		right := p.parseTerm()
		node.Children = append(node.Children, opTok)
		if right != nil {
			node.Children = append(node.Children, right)
			node.Span.End = right.Span.End
		}
	}
	return node
}

func (p *CSTParser) parseTerm() *ParseNode {
	left := p.parseFactor()
	if left == nil {
		return nil
	}
	node := &ParseNode{
		Symbol:   "Term",
		Span:     left.Span,
		Children: []*ParseNode{left},
	}
	for p.peekToken.Type == ASTERISK || p.peekToken.Type == SLASH {
		p.nextToken()
		opTok := p.makeTerminal(p.curToken)
		p.nextToken()
		right := p.parseFactor()
		node.Children = append(node.Children, opTok)
		if right != nil {
			node.Children = append(node.Children, right)
			node.Span.End = right.Span.End
		}
	}
	return node
}

func (p *CSTParser) parseFactor() *ParseNode {
	switch p.curToken.Type {
	case INT, IDENT:
		tok := p.makeTerminal(p.curToken)
		return &ParseNode{
			Symbol:   "Factor",
			Span:     tok.Span,
			Children: []*ParseNode{tok},
		}
	case LPAREN:
		lparen := p.makeTerminal(p.curToken)
		p.nextToken()
		expr := p.parseExpression()
		children := []*ParseNode{lparen}
		if expr != nil {
			children = append(children, expr)
		}
		endSpan := lparen.Span.End
		if p.peekToken.Type == RPAREN {
			p.nextToken()
			rparen := p.makeTerminal(p.curToken)
			endSpan = rparen.Span.End
			children = append(children, rparen)
		}
		return &ParseNode{
			Symbol:   "Factor",
			Span:     ast.SourceSpan{Start: lparen.Span.Start, End: endSpan},
			Children: children,
		}
	default:
		return p.makeTerminal(p.curToken)
	}
}

// FormatText returns an indented text tree of the Concrete Syntax Tree.
func (n *ParseNode) FormatText() string {
	var b strings.Builder
	n.writeText(&b, "", true, true)
	return b.String()
}

func (n *ParseNode) writeText(b *strings.Builder, indent string, isLast bool, isRoot bool) {
	label := n.Symbol
	if n.Token != "" {
		label = fmt.Sprintf("%s (%q)", n.Symbol, n.Token)
	}
	pos := ""
	if n.Span.Start.Line > 0 {
		pos = fmt.Sprintf(" [line %d:%d]", n.Span.Start.Line, n.Span.Start.Column)
	}

	if isRoot {
		b.WriteString(label + pos + "\n")
	} else {
		marker := "├── "
		if isLast {
			marker = "└── "
		}
		b.WriteString(indent + marker + label + pos + "\n")
	}

	childIndent := indent
	if !isRoot {
		if isLast {
			childIndent += "    "
		} else {
			childIndent += "│   "
		}
	}

	for i, child := range n.Children {
		child.writeText(b, childIndent, i == len(n.Children)-1, false)
	}
}