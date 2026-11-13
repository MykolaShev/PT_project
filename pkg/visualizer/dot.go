package visualizer

import (
	"fmt"
	"pt_project/pkg/ast"
	"strings"
)

// DotVisualizer generates a Graphviz DOT representation of the AST.
type DotVisualizer struct {
	nodeCount int
	builder   strings.Builder
}

func NewDotVisualizer() *DotVisualizer {
	return &DotVisualizer{}
}

// Generate returns the DOT script string for the given program.
func (v *DotVisualizer) Generate(program *ast.Program) string {
	v.builder.WriteString("digraph AST {\n")
	v.builder.WriteString("  node [shape=box, style=rounded, fontname=\"Courier\"];\n")
	v.builder.WriteString("  edge [fontname=\"Courier\", fontsize=10];\n")

	v.walk(program)

	v.builder.WriteString("}\n")
	return v.builder.String()
}

func (v *DotVisualizer) addNode(label string) int {
	v.nodeCount++
	id := v.nodeCount
	v.builder.WriteString(fmt.Sprintf("  node%d [label=\"%s\"];\n", id, label))
	return id
}

func (v *DotVisualizer) addEdge(from, to int, label string) {
	if label != "" {
		v.builder.WriteString(fmt.Sprintf("  node%d -> node%d [label=\"%s\"];\n", from, to, label))
	} else {
		v.builder.WriteString(fmt.Sprintf("  node%d -> node%d;\n", from, to))
	}
}

// walk recursively processes AST nodes and returns the generated node ID.
func (v *DotVisualizer) walk(n ast.Node) int {
	if n == nil {
		return v.addNode("nil")
	}

	switch node := n.(type) {
	case *ast.Program:
		id := v.addNode("Program")
		for i, stmt := range node.Statements {
			childId := v.walk(stmt)
			v.addEdge(id, childId, fmt.Sprintf("stmt[%d]", i))
		}
		return id
	case *ast.AssignStmt:
		id := v.addNode(":=")
		idNode := v.addNode(fmt.Sprintf("id(%s)", node.Identifier.Value))
		v.addEdge(id, idNode, "left")
		valId := v.walk(node.Value)
		v.addEdge(id, valId, "right")
		return id
	case *ast.SequenceStmt:
		id := v.addNode(";")
		leftId := v.walk(node.Left)
		rightId := v.walk(node.Right)
		v.addEdge(id, leftId, "left")
		v.addEdge(id, rightId, "right")
		return id
	case *ast.IfStmt:
		id := v.addNode("if")
		condId := v.walk(node.Condition)
		thenId := v.walk(node.ThenBlock)
		v.addEdge(id, condId, "cond")
		v.addEdge(id, thenId, "then")
		if node.ElseBlock != nil {
			elseId := v.walk(node.ElseBlock)
			v.addEdge(id, elseId, "else")
		}
		return id
	case *ast.WhileStmt:
		id := v.addNode("while")
		condId := v.walk(node.Condition)
		bodyId := v.walk(node.Body)
		v.addEdge(id, condId, "cond")
		v.addEdge(id, bodyId, "body")
		return id
	case *ast.BlockStmt:
		id := v.addNode("begin...end")
		bodyId := v.walk(node.Body)
		v.addEdge(id, bodyId, "body")
		return id
	case *ast.SkipStmt:
		return v.addNode("skip")
	case *ast.NumberExpr:
		return v.addNode(fmt.Sprintf("int(%d)", node.Value))
	case *ast.IdentifierExpr:
		return v.addNode(fmt.Sprintf("id(%s)", node.Name))
	case *ast.BinaryExpr:
		id := v.addNode(node.Operator)
		leftId := v.walk(node.Left)
		rightId := v.walk(node.Right)
		v.addEdge(id, leftId, "left")
		v.addEdge(id, rightId, "right")
		return id
	default:
		return v.addNode(fmt.Sprintf("Unknown: %T", node))
	}
}