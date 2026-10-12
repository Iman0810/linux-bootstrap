package runner

import (
	"errors"
	"io"
	"os"
	"path/filepath"
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

	markerPath := filepath.Join(t.TempDir(), "command-executed")

	originalStdout := os.Stdout
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create stdout pipe: %v", err)
	}

	os.Stdout = writer

	runErr := r.Run("touch", markerPath)

	_ = writer.Close()
	os.Stdout = originalStdout

	output, readErr := io.ReadAll(reader)
	_ = reader.Close()

	if readErr != nil {
		t.Fatalf("failed to read captured output: %v", readErr)
	}

	if runErr != nil {
		t.Fatalf("expected dry-run to succeed, got error: %v", runErr)
	}

	if !strings.Contains(string(output), "[DRY RUN] touch") {
		t.Fatalf("expected dry-run message, got %q", string(output))
	}

	if _, err := os.Stat(markerPath); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected command not to execute; marker file status error: %v", err)
	}
}
func TestRunMissingCommand(t *testing.T) {
	r := Runner{}

	err := r.Run("this-command-definitely-does-not-exist")

	if err == nil {
		t.Fatal("expected missing command to return an error")
	}

	if !strings.Contains(err.Error(), "failed to execute") {
		t.Fatalf("expected execution error, got %q", err.Error())
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
