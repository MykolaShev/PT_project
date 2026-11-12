package interpreter

import (
	"fmt"
	"pt_project/pkg/ast"
)

// Interpreter executes the SIPL Abstract Syntax Tree.
type Interpreter struct {
	Env *Environment
}

// New initializes a new interpreter with an empty environment.
func New() *Interpreter {
	return &Interpreter{
		Env: NewEnvironment(),
	}
}

// Interpret is the main entry point for executing a Program.
func (i *Interpreter) Interpret(program *ast.Program) error {
	for _, stmt := range program.Statements {
		if err := i.ExecuteStatement(stmt); err != nil {
			return err
		}
	}
	return nil
}

// ExecuteStatement recursively evaluates SIPL statements and modifies the environment.
func (i *Interpreter) ExecuteStatement(node ast.Statement) error {
	switch stmt := node.(type) {
	case *ast.AssignStmt:
		val, err := i.EvaluateExpression(stmt.Value)
		if err != nil {
			return err
		}
		i.Env.Set(stmt.Identifier.Value, val)

	case *ast.SequenceStmt:
		if err := i.ExecuteStatement(stmt.Left); err != nil {
			return err
		}
		if err := i.ExecuteStatement(stmt.Right); err != nil {
			return err
		}

	case *ast.IfStmt:
		cond, err := i.EvaluateCondition(stmt.Condition)
		if err != nil {
			return err
		}
		if cond {
			return i.ExecuteStatement(stmt.ThenBlock)
		} else if stmt.ElseBlock != nil {
			return i.ExecuteStatement(stmt.ElseBlock)
		}

	case *ast.WhileStmt:
		for {
			cond, err := i.EvaluateCondition(stmt.Condition)
			if err != nil {
				return err
			}
			if !cond {
				break
			}
			if err := i.ExecuteStatement(stmt.Body); err != nil {
				return err
			}
		}

	case *ast.BlockStmt:
		return i.ExecuteStatement(stmt.Body)

	case *ast.SkipStmt:
		// Do nothing
		return nil

	default:
		return fmt.Errorf("unknown statement type: %T", stmt)
	}
	return nil
}

// EvaluateExpression computes the integer result of an arithmetic expression.
func (i *Interpreter) EvaluateExpression(node ast.Expression) (int64, error) {
	switch expr := node.(type) {
	case *ast.NumberExpr:
		return expr.Value, nil
	case *ast.IdentifierExpr:
		return i.Env.Get(expr.Name), nil
	case *ast.BinaryExpr:
		left, err := i.EvaluateExpression(expr.Left)
		if err != nil {
			return 0, err
		}
		right, err := i.EvaluateExpression(expr.Right)
		if err != nil {
			return 0, err
		}

		switch expr.Operator {
		case "+": return left + right, nil
		case "-": return left - right, nil
		case "*": return left * right, nil
		case "/":
			if right == 0 {
				return 0, fmt.Errorf("division by zero at line %d, col %d", expr.Span().Start.Line, expr.Span().Start.Column)
			}
			return left / right, nil
		default:
			return 0, fmt.Errorf("unknown arithmetic operator: %s", expr.Operator)
		}
	}
	return 0, fmt.Errorf("unknown expression type: %T", node)
}

// EvaluateCondition computes the boolean result of a logical expression.
func (i *Interpreter) EvaluateCondition(node ast.Expression) (bool, error) {
	binExpr, ok := node.(*ast.BinaryExpr)
	if !ok {
		return false, fmt.Errorf("condition must be a binary expression (e.g., X > 0)")
	}

	left, err := i.EvaluateExpression(binExpr.Left)
	if err != nil {
		return false, err
	}
	right, err := i.EvaluateExpression(binExpr.Right)
	if err != nil {
		return false, err
	}

	switch binExpr.Operator {
	case ">": return left > right, nil
	case "<": return left < right, nil
	case "=": return left == right, nil
	default:
		return false, fmt.Errorf("unknown relational operator: %s", binExpr.Operator)
	}
}