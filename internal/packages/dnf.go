package packages

import "github.com/Iman0810/linux-bootstrap/internal/runner"

type DnfManager struct {
	Runner runner.Runner
}

func (d DnfManager) Update() error {
	return d.Runner.Run("sudo", "dnf", "makecache")
}

func (d DnfManager) Install(packages ...string) error {
	args := append([]string{"dnf", "install", "-y"}, packages...)

	return d.Runner.Run("sudo", args...)
}

func (d DnfManager) IsInstalled(packageName string) bool {
	_, err := d.Runner.Output("rpm", "-q", packageName)

	return err == nil
}
