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