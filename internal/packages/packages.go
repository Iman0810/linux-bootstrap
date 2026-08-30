package packages

func ResolvePackages(manager Manager, packages []string) []string {
	var resolved []string

	for _, packageName := range packages {
		switch manager {
		case APT:
			resolved = append(resolved, resolveAPT(packageName))

		case DNF:
			resolved = append(resolved, resolveDNF(packageName))

		case Pacman:
			resolved = append(resolved, resolvePacman(packageName))

		default:
			resolved = append(resolved, packageName)
		}
	}

	return resolved
}

func resolveAPT(packageName string) string {
	switch packageName {
	case "build-tools":
		return "build-essential"

	default:
		return packageName
	}
}

func resolveDNF(packageName string) string {
	switch packageName {
	case "build-tools":
		return "gcc gcc-c++ make"

	default:
		return packageName
	}
}

func resolvePacman(packageName string) string {
	switch packageName {
	case "build-tools":
		return "base-devel"

	default:
		return packageName
	}
}