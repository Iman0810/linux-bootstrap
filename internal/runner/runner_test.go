package runner

import (
	"strings"
	"testing"
)

func TestRunSuccessfulCommand(t *testing.T) {
	r := Runner{}

	err := r.Run("true")

	if err != nil {
		t.Fatalf("expected command to succeed, got error: %v", err)
	}
}

func TestRunFailedCommand(t *testing.T) {
	r := Runner{}

	err := r.Run("false")

	if err == nil {
		t.Fatal("expected command to fail, got nil")
	}

	expected := `command "false" failed with exit code 1`

	if err.Error() != expected {
		t.Fatalf("expected error %q, got %q", expected, err.Error())
	}
}

func TestRunDryRun(t *testing.T) {
	r := Runner{
		DryRun: true,
	}

	// "false" would fail if actually executed.
	// Dry-run should return nil without executing it.
	err := r.Run("false")

	if err != nil {
		t.Fatalf("expected dry-run to succeed, got error: %v", err)
	}
}

func TestOutputSuccessfulCommand(t *testing.T) {
	r := Runner{}

	output, err := r.Output("printf", "hello")

	if err != nil {
		t.Fatalf("expected command to succeed, got error: %v", err)
	}

	if output != "hello" {
		t.Fatalf("expected output %q, got %q", "hello", output)
	}
}

func TestOutputFailedCommand(t *testing.T) {
	r := Runner{}

	output, err := r.Output("sh", "-c", "printf 'failure output'; exit 2")

	if err == nil {
		t.Fatal("expected command to fail, got nil")
	}

	if output != "failure output" {
		t.Fatalf("expected output %q, got %q", "failure output", output)
	}

	expected := `command "sh" failed with exit code 2`

	if err.Error() != expected {
		t.Fatalf("expected error %q, got %q", expected, err.Error())
	}
}

func TestOutputMissingCommand(t *testing.T) {
	r := Runner{}

	_, err := r.Output("this-command-definitely-does-not-exist")

	if err == nil {
		t.Fatal("expected missing command to return an error")
	}

	if !strings.Contains(err.Error(), "failed to execute") {
		t.Fatalf("expected execution error, got %q", err.Error())
	}
}
