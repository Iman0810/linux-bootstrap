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

	return cmd.Run()
}