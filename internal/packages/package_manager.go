package packages

type PackageManager interface {
	Update() error
	Install(packages ...string) error
	IsInstalled(packageName string) bool
}