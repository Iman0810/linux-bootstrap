package profile

import "github.com/Iman0810/linux-bootstrap/internal/packages"

type Status struct {
	Profile Profile
	Plan    packages.PackagePlan
}

func CheckStatus(manager packages.PackageManager, p Profile) Status {
	plan := packages.BuildPlan(manager, p.Packages)

	return Status{
		Profile: p,
		Plan:    plan,
	}
}