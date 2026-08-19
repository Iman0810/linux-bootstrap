package runner

import (
	"fmt"
	"os"
	"os/exec"
)

type Runner struct {
	DryRun bool
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

	err := cmd.Run()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return fmt.Errorf(
				"command %q failed with exit code %d",
				command,
				exitError.ExitCode(),
			)
		}

		return fmt.Errorf("failed to execute %q: %w", command, err)
	}

	return nil
}

func (r Runner) Output(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)

	output, err := cmd.CombinedOutput()

	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return string(output), fmt.Errorf(
				"command %q failed with exit code %d",
				command,
				exitError.ExitCode(),
			)
		}

		return string(output), fmt.Errorf(
			"failed to execute %q: %w",
			command,
			err,
		)
	}

	return string(output), nil
}
