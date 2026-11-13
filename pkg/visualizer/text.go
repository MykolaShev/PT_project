package visualizer

import (
	"fmt"
	"pt_project/pkg/ast"
	"strings"
)

// TextVisualizer generates a CLI-friendly tree representation of the AST.
type TextVisualizer struct {
	builder strings.Builder
}

func NewTextVisualizer() *TextVisualizer {
	return &TextVisualizer{}
}

// Generate creates a string with the text tree.
func (v *TextVisualizer) Generate(n ast.Node, indent string, isLast bool) string {
	marker := "├── "
	if isLast {
		marker = "└── "
	}

	span := n.Span()
	posInfo := fmt.Sprintf(" [line %d:%d]", span.Start.Line, span.Start.Column)

	switch node := n.(type) {
	case *ast.Program:
		v.builder.WriteString(indent + "Program\n")
		for i, stmt := range node.Statements {
			v.Generate(stmt, indent+"    ", i == len(node.Statements)-1)
		}
	case *ast.AssignStmt:
		v.builder.WriteString(indent + marker + "AssignStmt" + posInfo + "\n")
		v.Generate(&node.Identifier, indent+getIndentAdd(isLast), false)
		v.Generate(node.Value, indent+getIndentAdd(isLast), true)
	case *ast.SequenceStmt:
		v.builder.WriteString(indent + marker + "SequenceStmt" + posInfo + "\n")
		v.Generate(node.Left, indent+getIndentAdd(isLast), false)
		v.Generate(node.Right, indent+getIndentAdd(isLast), true)
	case *ast.IfStmt:
		v.builder.WriteString(indent + marker + "IfStmt" + posInfo + "\n")
		v.Generate(node.Condition, indent+getIndentAdd(isLast), false)
		v.Generate(node.ThenBlock, indent+getIndentAdd(isLast), node.ElseBlock == nil)
		if node.ElseBlock != nil {
			v.Generate(node.ElseBlock, indent+getIndentAdd(isLast), true)
		}
	case *ast.WhileStmt:
		v.builder.WriteString(indent + marker + "WhileStmt" + posInfo + "\n")
		v.Generate(node.Condition, indent+getIndentAdd(isLast), false)
		v.Generate(node.Body, indent+getIndentAdd(isLast), true)
	case *ast.BlockStmt:
		v.builder.WriteString(indent + marker + "BlockStmt" + posInfo + "\n")
		v.Generate(node.Body, indent+getIndentAdd(isLast), true)
	case *ast.BinaryExpr:
		v.builder.WriteString(indent + marker + "BinaryExpr (" + node.Operator + ")" + posInfo + "\n")
		v.Generate(node.Left, indent+getIndentAdd(isLast), false)
		v.Generate(node.Right, indent+getIndentAdd(isLast), true)
	case *ast.IdentifierExpr:
		v.builder.WriteString(indent + marker + "Identifier (" + node.Name + ")" + posInfo + "\n")
	case *ast.NumberExpr:
		v.builder.WriteString(indent + marker + fmt.Sprintf("Number (%d)", node.Value) + posInfo + "\n")
	case *ast.Identifier:
		v.builder.WriteString(indent + marker + "Identifier (" + node.Value + ")" + posInfo + "\n")
	case *ast.SkipStmt:
		v.builder.WriteString(indent + marker + "SkipStmt" + posInfo + "\n")
	}

	return v.builder.String()
}

func getIndentAdd(isLast bool) string {
	if isLast {
		return "    "
	}
	return "│   "
}