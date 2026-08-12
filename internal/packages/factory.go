package packages

func GetPackageManager(manager Manager) PackageManager {
	switch manager {
	case APT:
		return AptManager{}

	default:
		return nil
	}
}