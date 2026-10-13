package packages

import (
	"fmt"

	"github.com/Iman0810/linux-bootstrap/internal/runner"
)

type AptManager struct {
	Runner runner.Executor
}

func (a AptManager) Update() error {
	if err := a.Runner.Run("sudo", "apt", "update"); err != nil {
		return fmt.Errorf("apt: failed to update package lists: %w", err)
	}

	return nil
}

func (a AptManager) Install(packages ...string) error {
	args := append([]string{"apt", "install", "-y"}, packages...)

	if err := a.Runner.Run("sudo", args...); err != nil {
		return fmt.Errorf(
			"apt: failed to install packages %v: %w",
			packages,
			err,
		)
	}

	return nil
}

func (a AptManager) IsInstalled(packageName string) bool {
	_, err := a.Runner.Output("dpkg", "-s", packageName)

	return err == nil
}
