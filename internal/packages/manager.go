package packages

import (
	"fmt"
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
	if manager := detectManagerByID(osInfo.ID); manager != Unknown {
		return manager
	}

	for _, id := range strings.Fields(osInfo.IDLike) {
		if manager := detectManagerByID(id); manager != Unknown {
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

func (m Manager) String() string {
	if m == Unknown {
		return fmt.Sprintf("%q", string(m))
	}

	return string(m)
}
