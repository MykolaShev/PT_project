package parser

import (
	"fmt"
	"pt_project/pkg/ast"
)

const (
	_ int = iota
	LOWEST
	EQUALS      // =
	LESSGREATER // > or <
	SUM         // + or -
	PRODUCT     // * or /
)

var precedences = map[TokenType]int{
	EQ:       EQUALS,
	LT:       LESSGREATER,
	GT:       LESSGREATER,
	PLUS:     SUM,
	MINUS:    SUM,
	ASTERISK: PRODUCT,
	SLASH:    PRODUCT,
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

// Parser implements a top-down recursive descent parser for SIPL.
type Parser struct {
	l         *Lexer
	curToken  Token
	peekToken Token
	errors    []string
}

func New(l *Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	// Read two tokens so curToken and peekToken are populated
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) Errors() []string {
	return p.errors
}

// ParseProgram is the entry point for parsing a SIPL source file.
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

// parseStatement routes to specific statement parsing based on the current token.
func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case IDENT:
		return p.parseAssignStatement()
	case BEGIN:
		return p.parseBlockStatement()
	case WHILE:
		return p.parseWhileStatement()
	case IF:
		return p.parseIfStatement()
	case SKIP:
		return p.parseSkipStatement()
	default:
		p.addError(fmt.Sprintf("unrecognized statement starting with %q at line %d, col %d",
			p.curToken.Literal, p.curToken.Line, p.curToken.Column))
		return nil
	}
}

// parseCompoundStatement handles sequences separated by semicolons (S1 ; S2)
func (p *Parser) parseCompoundStatement() ast.Statement {
	left := p.parseStatement()

	for p.peekToken.Type == SEMICOLON {
		p.nextToken() // move to ';'
		p.nextToken() // move to start of the next statement

		right := p.parseStatement()

		if left != nil && right != nil {
			left = &ast.SequenceStmt{
				SourceSpan: ast.SourceSpan{Start: left.Span().Start, End: right.Span().End},
				Left:       left,
				Right:      right,
			}
		}
	}

	return left
}

// parseAssignStatement handles "variable := expression" constructs.
func (p *Parser) parseAssignStatement() ast.Statement {
	stmt := &ast.AssignStmt{}

	startPos := ast.Position{Line: p.curToken.Line, Column: p.curToken.Column}

	stmt.Identifier = ast.Identifier{
		Value: p.curToken.Literal,
		// Span for the identifier itself
		SourceSpan: ast.SourceSpan{Start: startPos, End: startPos},
	}

	// Expect the next token to be ':='
	if !p.expectPeek(ASSIGN) {
		return nil
	}

	p.nextToken() // Skip ':='

	// Parse the actual expression on the right-hand side
	stmt.Value = p.parseExpression(LOWEST)

	// Calculate end position based on the parsed expression
	if stmt.Value != nil {
		stmt.SourceSpan = ast.SourceSpan{Start: startPos, End: stmt.Value.Span().End}
	} else {
		stmt.SourceSpan = ast.SourceSpan{
			Start: startPos,
			End:   ast.Position{Line: p.curToken.Line, Column: p.curToken.Column},
		} // fallback
	}

	return stmt
}

// parseBlockStatement handles "begin ... end" blocks.
func (p *Parser) parseBlockStatement() ast.Statement {
	stmt := &ast.BlockStmt{}
	startPos := ast.Position{Line: p.curToken.Line, Column: p.curToken.Column}

	p.nextToken() // consume 'begin'

	stmt.Body = p.parseCompoundStatement()

	if !p.expectPeek(END) {
		return nil
	}

	stmt.SourceSpan = ast.SourceSpan{
		Start: startPos,
		End:   ast.Position{Line: p.curToken.Line, Column: p.curToken.Column},
	}
	return stmt
}

// parseWhileStatement handles "while E do S" loops.
func (p *Parser) parseWhileStatement() ast.Statement {
	stmt := &ast.WhileStmt{}
	startPos := ast.Position{Line: p.curToken.Line, Column: p.curToken.Column}

	p.nextToken() // consume 'while'
	stmt.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(DO) {
		return nil
	}

	p.nextToken() // consume 'do'
	stmt.Body = p.parseStatement()

	if stmt.Body != nil {
		stmt.SourceSpan = ast.SourceSpan{Start: startPos, End: stmt.Body.Span().End}
	}
	return stmt
}

// parseIfStatement handles "if E then S1 else S2" conditions.
func (p *Parser) parseIfStatement() ast.Statement {
	stmt := &ast.IfStmt{}
	startPos := ast.Position{Line: p.curToken.Line, Column: p.curToken.Column}

	p.nextToken() // consume 'if'
	stmt.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(THEN) {
		return nil
	}

	p.nextToken() // consume 'then'
	stmt.ThenBlock = p.parseStatement()

	if p.peekToken.Type == ELSE {
		p.nextToken() // move to 'else'
		p.nextToken() // move to the start of the else block
		stmt.ElseBlock = p.parseStatement()
	}

	endPos := stmt.ThenBlock.Span().End
	if stmt.ElseBlock != nil {
		endPos = stmt.ElseBlock.Span().End
	}

	stmt.SourceSpan = ast.SourceSpan{Start: startPos, End: endPos}
	return stmt
}

// parseSkipStatement handles "skip".
func (p *Parser) parseSkipStatement() ast.Statement {
	stmt := &ast.SkipStmt{}
	pos := ast.Position{Line: p.curToken.Line, Column: p.curToken.Column}
	stmt.SourceSpan = ast.SourceSpan{Start: pos, End: pos}
	return stmt
}

// parseExpression is the core of the Pratt parser.
func (p *Parser) parseExpression(precedence int) ast.Expression {
	var leftExp ast.Expression

	// Prefix parsing (numbers, identifiers, or grouped expressions)
	switch p.curToken.Type {
	case IDENT:
		leftExp = p.parseIdentifier()
	case INT:
		leftExp = p.parseIntegerLiteral()
	case LPAREN:
		leftExp = p.parseGroupedExpression()
	default:
		p.addError(fmt.Sprintf("no prefix parse function for %s found at line %d, col %d", p.curToken.Type, p.curToken.Line, p.curToken.Column))
		return nil
	}

	// Infix parsing (binary operators)
	for p.peekToken.Type != SEMICOLON && precedence < p.peekPrecedence() {
		switch p.peekToken.Type {
		case PLUS, MINUS, ASTERISK, SLASH, LT, GT, EQ:
			p.nextToken()
			leftExp = p.parseInfixExpression(leftExp)
		default:
			return leftExp
		}
	}

	return leftExp
}

func (p *Parser) parseIdentifier() ast.Expression {
	startPos := ast.Position{Line: p.curToken.Line, Column: p.curToken.Column}
	return &ast.IdentifierExpr{
		SourceSpan: ast.SourceSpan{Start: startPos, End: startPos}, // Simplified span for now
		Name:       p.curToken.Literal,
	}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	startPos := ast.Position{Line: p.curToken.Line, Column: p.curToken.Column}

	// Convert p.curToken.Literal to int64
	var value int64
	_, err := fmt.Sscanf(p.curToken.Literal, "%d", &value)
	if err != nil {
		p.addError(fmt.Sprintf("could not parse %q as integer", p.curToken.Literal))
		return nil
	}

	return &ast.NumberExpr{
		SourceSpan: ast.SourceSpan{Start: startPos, End: startPos}, // Simplified span
		Value:      value,
	}
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()
	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(RPAREN) {
		return nil
	}
	return exp
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	startPos := left.Span().Start

	expression := &ast.BinaryExpr{
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()

	expression.Right = p.parseExpression(precedence)

	// Calculate end position based on the right expression
	if expression.Right != nil {
		expression.SourceSpan = ast.SourceSpan{Start: startPos, End: expression.Right.Span().End}
	}

	return expression
}

// expectPeek asserts the next token type and advances if matched.
func (p *Parser) expectPeek(t TokenType) bool {
	if p.peekToken.Type == t {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

// peekError records an error when the expected token doesn't match reality.
func (p *Parser) peekError(t TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead at line %d, col %d",
		t, p.peekToken.Type, p.peekToken.Line, p.peekToken.Column)
	p.addError(msg)
}

func (p *Parser) addError(msg string) {
	p.errors = append(p.errors, msg)
}
