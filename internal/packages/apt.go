package packages

import (
	"os/exec"
)

type AptManager struct{}

func (a AptManager) Update() error {
	cmd := exec.Command("sudo", "apt", "update")

	return cmd.Run()
}

func (a AptManager) Install(packages ...string) error {
	args := append([]string{"apt", "install", "-y"}, packages...)

	cmd := exec.Command("sudo", args...)

	return cmd.Run()
}

func (a AptManager) IsInstalled(packageName string) bool {
	cmd := exec.Command("dpkg", "-s", packageName)

	err := cmd.Run()

	return err == nil
}

func (a AptManager) String() string {
	return "APT"
}