package packages

import (
	"github.com/Iman0810/linux-bootstrap/internal/runner"
)

type AptManager struct {
	Runner runner.Runner
}

func (a AptManager) Update() error {
	return a.Runner.Run("sudo", "apt", "update")
}

func (a AptManager) Install(packages ...string) error {
	args := append([]string{"apt", "install", "-y"}, packages...)

	return a.Runner.Run("sudo", args...)
}

func (a AptManager) IsInstalled(packageName string) bool {
	return false
}