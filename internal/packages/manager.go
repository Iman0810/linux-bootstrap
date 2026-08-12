package packages

import (
	"fmt"
	
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
	switch osInfo.ID {
	case "ubuntu", "debian", "pop":
		return APT

	case "fedora", "rhel", "centos", "rocky", "almalinux":
		return DNF

	case "arch", "manjaro", "endeavouros":
		return Pacman
	}

	return Unknown
}

func (m Manager) String() string {
	if m == Unknown {
		return fmt.Sprintf("%q", string(m))
	}

	return string(m)
}