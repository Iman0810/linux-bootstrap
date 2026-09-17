package setup

import (
	"github.com/Iman0810/linux-bootstrap/internal/packages"
	"github.com/Iman0810/linux-bootstrap/internal/profile"
	"github.com/Iman0810/linux-bootstrap/internal/runner"
)

type Service struct {
	Manager     packages.PackageManager
	ManagerType packages.Manager
	Runner      runner.Runner
}

type Result struct {
	Profile  profile.Profile
	Plan     packages.PackagePlan
	Verified packages.PackagePlan
}

func (s Service) Prepare(p profile.Profile) (Result, bool) {
	profilePackages, ok := profile.PackagesFor(
		p,
		s.ManagerType,
	)

	if !ok {
		return Result{}, false
	}

	plan := packages.BuildPlan(
		s.Manager,
		profilePackages,
	)

	return Result{
		Profile: p,
		Plan:    plan,
	}, true
}

func (s Service) Execute(plan packages.PackagePlan) (packages.PackagePlan, error) {
	if len(plan.Missing) == 0 {
		return plan, nil
	}

	if err := s.Manager.Update(); err != nil {
		return plan, err
	}

	if err := s.Manager.Install(plan.Missing...); err != nil {
		return plan, err
	}

	if s.Runner.DryRun {
		return plan, nil
	}

	allPackages := append(
		append([]string{}, plan.Installed...),
		plan.Missing...,
	)

	verified := packages.BuildPlan(
		s.Manager,
		allPackages,
	)

	return verified, nil
}
