package main_test

import (
	"os/exec"
	"testing"
)

// Integration test checking finishing errors
func TestCLIExitCodes(t *testing.T) {
	// Scenario 1: exit code 0
	cmdParse := exec.Command("go", "run", "main.go", "parse", "../../test.sipl")
	errParse := cmdParse.Run()
	if errParse != nil {
		t.Fatalf("Expected parse command to succeed (exit code 0), but got error: %v", errParse)
	}

	// Scenario 2: exit code 1
	cmdUnknown := exec.Command("go", "run", "main.go", "invalid_command", "../../test.sipl")
	errUnknown := cmdUnknown.Run()
	if errUnknown == nil {
		t.Fatalf("Expected invalid command to fail (exit code 1), but it succeeded")
	}

	// Scenario 3: exit code 1
	cmdNoFile := exec.Command("go", "run", "main.go", "run", "does_not_exist.sipl")
	errNoFile := cmdNoFile.Run()
	if errNoFile == nil {
		t.Fatalf("Expected missing file to cause failure (exit code 1), but it succeeded")
	}
}