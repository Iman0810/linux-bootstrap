package packages

import "github.com/Iman0810/linux-bootstrap/internal/runner"

type PacmanManager struct {
	Runner runner.Runner
}

func (p PacmanManager) Update() error {
	return p.Runner.Run("sudo", "pacman", "-Sy")
}

func (p PacmanManager) Install(packages ...string) error {
	args := append([]string{"pacman", "-S", "--noconfirm"}, packages...)

	return p.Runner.Run("sudo", args...)
}

func (p PacmanManager) IsInstalled(packageName string) bool {
	_, err := p.Runner.Output("pacman", "-Q", packageName)

	return err == nil
}
