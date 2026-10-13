package packages

import (
	"fmt"

	"github.com/Iman0810/linux-bootstrap/internal/runner"
)

type PacmanManager struct {
	Runner runner.Executor
}

func (p PacmanManager) Update() error {
	if err := p.Runner.Run("sudo", "pacman", "-Syu"); err != nil {
		return fmt.Errorf("pacman: failed to synchronize and update the system: %w", err)
	}

	return nil
}

func (p PacmanManager) Install(packages ...string) error {
	args := append([]string{"pacman", "-S", "--noconfirm"}, packages...)

	if err := p.Runner.Run("sudo", args...); err != nil {
		return fmt.Errorf(
			"pacman: failed to install packages %v: %w",
			packages,
			err,
		)
	}

	return nil
}

func (p PacmanManager) IsInstalled(packageName string) bool {
	_, err := p.Runner.Output("pacman", "-Q", packageName)

	return err == nil
}
