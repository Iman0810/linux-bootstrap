package cli

import "testing"

func TestRunArgsUnknownCommand(t *testing.T) {
	err := RunArgs([]string{"banana"})

	if err == nil {
		t.Fatal("expected error for unknown command")
	}
}

func TestRunArgsNoCommand(t *testing.T) {
	err := RunArgs([]string{})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestRunArgsInvalidSetupFlag(t *testing.T) {
	err := RunArgs([]string{
		"setup",
		"--banana",
	})

	if err == nil {
		t.Fatal("expected error for invalid setup flag")
	}
}
