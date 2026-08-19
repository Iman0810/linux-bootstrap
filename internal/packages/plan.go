package packages

type PackagePlan struct {
	Installed []string
	Missing   []string
}

func BuildPlan(manager PackageManager, desired []string) PackagePlan {
	var plan PackagePlan

	for _, packageName := range desired {
		if manager.IsInstalled(packageName) {
			plan.Installed = append(plan.Installed, packageName)
		} else {
			plan.Missing = append(plan.Missing, packageName)
		}
	}

	return plan
}