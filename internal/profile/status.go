package profile

import "github.com/Iman0810/linux-bootstrap/internal/packages"

type Status struct {
	Profile Profile
	Plan    packages.PackagePlan
}

func CheckStatus(
	manager packages.PackageManager,
	managerType packages.Manager,
	p Profile,
) Status {
	profilePackages, ok := PackagesFor(p, managerType)
	if !ok {
		return Status{
			Profile: p,
			Plan:    packages.PackagePlan{},
		}
	}

	plan := packages.BuildPlan(manager, profilePackages)

	return Status{
		Profile: p,
		Plan:    plan,
	}
}
