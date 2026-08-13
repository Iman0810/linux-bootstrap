package packages

import "github.com/Iman0810/linux-bootstrap/internal/runner"

func GetPackageManager(manager Manager, r runner.Runner) PackageManager {
	switch manager {
	case APT:
		return AptManager{
			Runner: r,
		}

	default:
		return nil
	}
}