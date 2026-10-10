package prompt

import (
	"os"
	"strings"
	"testing"
)

func TestConfirm(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{name: "lowercase yes", input: "y\n", want: true},
		{name: "uppercase yes", input: "Y\n", want: true},
		{name: "full yes", input: "yes\n", want: true},
		{name: "uppercase full yes", input: "YES\n", want: true},
		{name: "yes with whitespace", input: "  yes  \n", want: true},
		{name: "no", input: "n\n", want: false},
		{name: "full no", input: "no\n", want: false},
		{name: "empty input", input: "\n", want: false},
		{name: "unexpected input", input: "maybe\n", want: false},
		{name: "EOF without input", input: "", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldStdin := os.Stdin
			oldStdout := os.Stdout
			defer func() {
				os.Stdin = oldStdin
				os.Stdout = oldStdout
			}()

			stdinR, stdinW, err := os.Pipe()
			if err != nil {
				t.Fatalf("create stdin pipe: %v", err)
			}

			stdoutR, stdoutW, err := os.Pipe()
			if err != nil {
				stdinR.Close()
				stdinW.Close()
				t.Fatalf("create stdout pipe: %v", err)
			}

			os.Stdin = stdinR
			os.Stdout = stdoutW

			if tt.input != "" {
				if _, err := stdinW.WriteString(tt.input); err != nil {
					t.Fatalf("write stdin: %v", err)
				}
			}
			stdinW.Close()

			got := Confirm("Continue?")

			stdoutW.Close()
			output := make([]byte, 256)
			n, readErr := stdoutR.Read(output)
			if readErr != nil {
				t.Fatalf("read prompt output: %v", readErr)
			}

			stdinR.Close()
			stdoutR.Close()

			if got != tt.want {
				t.Errorf("Confirm() = %v, want %v", got, tt.want)
			}

			if !strings.Contains(string(output[:n]), "Continue? [y/N]: ") {
				t.Errorf("unexpected prompt output: %q", string(output[:n]))
			}
		})
	}
}
