package packages

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/Iman0810/linux-bootstrap/internal/system"
)

type Manager string

const (
	APT     Manager = "apt"
	DNF     Manager = "dnf"
	Pacman  Manager = "pacman"
	Unknown Manager = "unknown"
)

func DetectManager(osInfo system.OSInfo) Manager {
	id := strings.ToLower(strings.TrimSpace(osInfo.ID))

	if manager := detectManagerByID(id); manager != Unknown {
		return manager
	}

	for _, idLike := range strings.Fields(strings.ToLower(osInfo.IDLike)) {
		if manager := detectManagerByID(idLike); manager != Unknown {
			return manager
		}
	}

	return Unknown
}

func detectManagerByID(id string) Manager {
	switch id {
	case "ubuntu", "debian", "pop", "linuxmint", "zorin", "kali", "elementary",
		"devuan", "mx", "neon":
		return APT

	case "fedora", "rhel", "centos", "rocky", "almalinux", "ol", "amzn":
		return DNF

	case "arch", "manjaro", "endeavouros", "garuda", "artix":
		return Pacman

	default:
		return Unknown
	}
}
func validateManagerAvailable(
	manager Manager,
	lookPath func(string) (string, error),
) error {
	var command string

	switch manager {
	case APT:
		command = "apt"
	case DNF:
		command = "dnf"
	case Pacman:
		command = "pacman"
	default:
		return fmt.Errorf("unsupported package manager: %s", manager)
	}

	if _, err := lookPath(command); err != nil {
		return fmt.Errorf(
			"package manager %q is not available in PATH: %w",
			command,
			err,
		)
	}

	return nil
}

func ValidateManagerAvailable(manager Manager) error {
	return validateManagerAvailable(manager, exec.LookPath)
}

func (m Manager) String() string {
	if m == Unknown {
		return fmt.Sprintf("%q", string(m))
	}

	return string(m)
}
