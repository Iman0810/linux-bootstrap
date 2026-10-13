package runner

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

type Executor interface {
	Run(command string, args ...string) error
	Output(command string, args ...string) (string, error)
}

type Runner struct {
	DryRun bool
}

type CommandError struct {
	Command  string
	ExitCode int
	Err      error
}

func (e *CommandError) Error() string {
	return fmt.Sprintf(
		"command %q failed with exit code %d",
		e.Command,
		e.ExitCode,
	)
}

func (e *CommandError) Unwrap() error {
	return e.Err
}

func commandError(command string, err error) error {
	var exitError *exec.ExitError
	if errors.As(err, &exitError) {
		return &CommandError{
			Command:  command,
			ExitCode: exitError.ExitCode(),
			Err:      err,
		}
	}

	return fmt.Errorf("failed to execute %q: %w", command, err)
}

func (r Runner) Run(command string, args ...string) error {
	if r.DryRun {
		fmt.Printf("[DRY RUN] %s %v\n", command, args)
		return nil
	}

	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return commandError(command, err)
	}

	return nil
}

func (r Runner) Output(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		return string(output), commandError(command, err)
	}

	return string(output), nil
}
