package packages

import "testing"

type fakeManager struct {
	installed map[string]bool
}

func (f fakeManager) Update() error {
	return nil
}

func (f fakeManager) Install(packages ...string) error {
	return nil
}

func (f fakeManager) IsInstalled(packageName string) bool {
	return f.installed[packageName]
}

func TestBuildPlan(t *testing.T) {
	manager := fakeManager{
		installed: map[string]bool{
			"git":  true,
			"curl": true,
		},
	}

	desired := []string{
		"git",
		"curl",
		"wget",
		"unzip",
	}

	plan := BuildPlan(manager, desired)

	if len(plan.Installed) != 2 {
		t.Fatalf(
			"expected 2 installed packages, got %d",
			len(plan.Installed),
		)
	}

	if len(plan.Missing) != 2 {
		t.Fatalf(
			"expected 2 missing packages, got %d",
			len(plan.Missing),
		)
	}
}
func TestBuildPlanAllInstalled(t *testing.T) {
	manager := fakeManager{
		installed: map[string]bool{
			"git":  true,
			"curl": true,
			"wget": true,
		},
	}

	desired := []string{
		"git",
		"curl",
		"wget",
	}

	plan := BuildPlan(manager, desired)

	if len(plan.Installed) != 3 {
		t.Fatalf(
			"expected 3 installed packages, got %d",
			len(plan.Installed),
		)
	}

	if len(plan.Missing) != 0 {
		t.Fatalf(
			"expected 0 missing packages, got %d",
			len(plan.Missing),
		)
	}
}

func TestBuildPlanAllMissing(t *testing.T) {
	manager := fakeManager{
		installed: map[string]bool{},
	}

	desired := []string{
		"git",
		"curl",
		"wget",
	}

	plan := BuildPlan(manager, desired)

	if len(plan.Installed) != 0 {
		t.Fatalf(
			"expected 0 installed packages, got %d",
			len(plan.Installed),
		)
	}

	if len(plan.Missing) != 3 {
		t.Fatalf(
			"expected 3 missing packages, got %d",
			len(plan.Missing),
		)
	}
}

func TestBuildPlanEmptyDesired(t *testing.T) {
	manager := fakeManager{
		installed: map[string]bool{
			"git": true,
		},
	}

	plan := BuildPlan(manager, []string{})

	if len(plan.Installed) != 0 {
		t.Fatalf(
			"expected 0 installed packages, got %d",
			len(plan.Installed),
		)
	}

	if len(plan.Missing) != 0 {
		t.Fatalf(
			"expected 0 missing packages, got %d",
			len(plan.Missing),
		)
	}
}
