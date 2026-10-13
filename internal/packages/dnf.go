package packages

import (
	"fmt"

	"github.com/Iman0810/linux-bootstrap/internal/runner"
)

type DnfManager struct {
	Runner runner.Executor
}

func (d DnfManager) Update() error {
	if err := d.Runner.Run("sudo", "dnf", "makecache"); err != nil {
		return fmt.Errorf("dnf: failed to refresh package metadata: %w", err)
	}

	return nil
}

func (d DnfManager) Install(packages ...string) error {
	args := append([]string{"dnf", "install", "-y"}, packages...)

	if err := d.Runner.Run("sudo", args...); err != nil {
		return fmt.Errorf(
			"dnf: failed to install packages %v: %w",
			packages,
			err,
		)
	}

	return nil
}

func (d DnfManager) IsInstalled(packageName string) bool {
	_, err := d.Runner.Output("rpm", "-q", packageName)

	return err == nil
}
