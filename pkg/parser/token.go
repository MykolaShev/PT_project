package parser

import "fmt"

// TokenType defines the type of a lexical token.
type TokenType string

const (
	// ILLEGAL represents an unrecognized character.
	ILLEGAL TokenType = "ILLEGAL"
	// EOF marks the end of the source file.
	EOF TokenType = "EOF"

	// Identifiers and literals
	IDENT TokenType = "IDENT"
	INT   TokenType = "INT"

	// Operators
	ASSIGN   TokenType = ":="
	PLUS     TokenType = "+"
	MINUS    TokenType = "-"
	ASTERISK TokenType = "*"
	SLASH    TokenType = "/"

	// Relational Operators
	EQ    TokenType = "="
	NOTEQ TokenType = "!="
	LT    TokenType = "<"
	LTE   TokenType = "<="
	GT    TokenType = ">"
	GTE   TokenType = ">="

	// Delimiters
	SEMICOLON TokenType = ";"
	LPAREN    TokenType = "("
	RPAREN    TokenType = ")"

	// Keywords
	SKIP  TokenType = "SKIP"
	IF    TokenType = "IF"
	THEN  TokenType = "THEN"
	ELSE  TokenType = "ELSE"
	END   TokenType = "END"
	WHILE TokenType = "WHILE"
	DO    TokenType = "DO"
	BEGIN TokenType = "BEGIN"
)

// Represents a lexical unit with its source location.
type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

// String provides a readable string representation of the token for debugging.
func (t Token) String() string {
	return fmt.Sprintf("Token{%s, %q, %d:%d}", t.Type, t.Literal, t.Line, t.Column)
}

// keywords maps string identifiers to their respective token types.
var keywords = map[string]TokenType{
	"skip":  SKIP,
	"if":    IF,
	"then":  THEN,
	"else":  ELSE,
	"end":   END,
	"while": WHILE,
	"do":    DO,
	"begin": BEGIN,
}

// LookupIdent checks if an identifier is a reserved keyword.
func LookupIdent(ident string) TokenType {
	if tok, isKeyword := keywords[ident]; isKeyword {
		return tok
	}
	return IDENT
}