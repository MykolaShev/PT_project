package parser

import (
	"strings"
	"testing"
)

// TestParseCST verifies Concrete Syntax Tree construction across all SIPL constructs.
func TestParseCST(t *testing.T) {
	src := `begin
		X := ( 1 + 2 ) * 3 - 4 / 2 ;
		if X > 0 then skip else X := 0 ;
		if X = 7 then skip ;
		while X > 0 do X := X - 1
	end`

	cst := NewCSTParser(NewLexer(src)).ParseCST()
	if cst == nil {
		t.Fatal("expected non-nil CST root")
	}
	if cst.Symbol != "Program" {
		t.Fatalf("expected root symbol Program, got %s", cst.Symbol)
	}
	if cst.Span.Start.Line != 1 {
		t.Errorf("expected start line 1, got %d", cst.Span.Start.Line)
	}
	if len(cst.Children) != 1 || cst.Children[0].Symbol != "BlockStmt" {
		t.Fatalf("expected single BlockStmt child, got %+v", cst.Children)
	}

	block := cst.Children[0]
	if len(block.Children) != 3 {
		t.Fatalf("expected 3 children in BlockStmt (begin, body, end), got %d", len(block.Children))
	}
	if block.Children[0].Token != "begin" || block.Children[2].Token != "end" {
		t.Errorf("expected begin/end terminals, got %q and %q", block.Children[0].Token, block.Children[2].Token)
	}

	// Verify nested SequenceStmt chain for the 4 statements inside the block.
	seq := block.Children[1]
	if seq.Symbol != "SequenceStmt" {
		t.Fatalf("expected SequenceStmt inside BlockStmt, got %s", seq.Symbol)
	}
	if seq.Children[1].Token != ";" {
		t.Errorf("expected ';' terminal token in SequenceStmt, got %q", seq.Children[1].Token)
	}

	// Verify FormatText output includes non-terminals, terminals, and source spans.
	rendered := cst.FormatText()
	expectedFragments := []string{
		"Program [line 1:1]",
		"BlockStmt [line 1:1]",
		"BEGIN (\"begin\")",
		"AssignStmt",
		"Expression",
		"Term",
		"Factor",
		"( (\"(\")",
		") (\")\")",
		"IfStmt",
		"WhileStmt",
		"SkipStmt",
		"END (\"end\")",
	}
	for _, frag := range expectedFragments {
		if !strings.Contains(rendered, frag) {
			t.Errorf("expected FormatText output to contain %q, got:\n%s", frag, rendered)
		}
	}
}

// TestCSTConditionsAndEdgeCases covers relational operators and malformed token fallbacks.
func TestCSTConditionsAndEdgeCases(t *testing.T) {
	cases := []string{
		"begin if 1 != 2 then skip else skip end",
		"begin if 1 < 2 then skip else skip end",
		"begin if 1 <= 2 then skip else skip end",
		"begin if 2 >= 1 then skip else skip end",
		"begin if 1 then skip else skip end",
		"begin X ; @ end",
		"begin X := @ end",
	}

	for _, src := range cases {
		cst := NewCSTParser(NewLexer(src)).ParseCST()
		if cst == nil || cst.Symbol != "Program" {
			t.Errorf("failed to build CST root for input: %s", src)
		}
	}
}