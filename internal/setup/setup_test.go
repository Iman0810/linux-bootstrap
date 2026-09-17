package setup

import (
	"testing"

	"github.com/Iman0810/linux-bootstrap/internal/packages"
	"github.com/Iman0810/linux-bootstrap/internal/profile"
	"github.com/Iman0810/linux-bootstrap/internal/runner"
)

type mockPackageManager struct {
	installed     map[string]bool
	updateCalled  bool
	installCalled bool
	installedArgs []string
}

func (m *mockPackageManager) Update() error {
	m.updateCalled = true
	return nil
}

func (m *mockPackageManager) Install(packages ...string) error {
	m.installCalled = true
	m.installedArgs = packages

	for _, packageName := range packages {
		m.installed[packageName] = true
	}

	return nil
}

func (m *mockPackageManager) IsInstalled(packageName string) bool {
	return m.installed[packageName]
}

func TestPrepare(t *testing.T) {
	manager := &mockPackageManager{
		installed: map[string]bool{
			"git":   true,
			"curl":  false,
			"wget":  true,
			"unzip": true,
		},
	}

	service := Service{
		Manager:     manager,
		ManagerType: packages.APT,
		Runner:      runner.Runner{},
	}

	result, ok := service.Prepare(profile.Essentials)

	if !ok {
		t.Fatal("expected profile to be supported")
	}

	if len(result.Plan.Installed) != 3 {
		t.Fatalf(
			"expected 3 installed packages, got %d",
			len(result.Plan.Installed),
		)
	}

	if len(result.Plan.Missing) != 1 {
		t.Fatalf(
			"expected 1 missing package, got %d",
			len(result.Plan.Missing),
		)
	}

	if result.Plan.Missing[0] != "curl" {
		t.Fatalf(
			"expected curl to be missing, got %s",
			result.Plan.Missing[0],
		)
	}
}

func TestPrepareUnsupportedProfile(t *testing.T) {
	manager := &mockPackageManager{
		installed: map[string]bool{},
	}

	service := Service{
		Manager:     manager,
		ManagerType: packages.Unknown,
		Runner:      runner.Runner{},
	}

	_, ok := service.Prepare(profile.Essentials)

	if ok {
		t.Fatal("expected profile to be unsupported")
	}
}

func TestExecute(t *testing.T) {
	manager := &mockPackageManager{
		installed: map[string]bool{
			"git": true,
		},
	}

	service := Service{
		Manager:     manager,
		ManagerType: packages.APT,
		Runner:      runner.Runner{},
	}

	plan := packages.PackagePlan{
		Installed: []string{"git"},
		Missing:   []string{"curl"},
	}

	verified, err := service.Execute(plan)

	if err != nil {
		t.Fatalf("expected execute to succeed, got %v", err)
	}

	if len(verified.Missing) != 0 {
		t.Fatalf(
			"expected no missing packages after verification, got %v",
			verified.Missing,
		)
	}

	if len(verified.Installed) != 2 {
		t.Fatalf(
			"expected 2 installed packages after verification, got %d",
			len(verified.Installed),
		)
	}

	if !manager.updateCalled {
		t.Fatal("expected Update to be called")
	}

	if !manager.installCalled {
		t.Fatal("expected Install to be called")
	}

	if len(manager.installedArgs) != 1 {
		t.Fatalf(
			"expected 1 package to be installed, got %d",
			len(manager.installedArgs),
		)
	}

	if manager.installedArgs[0] != "curl" {
		t.Fatalf(
			"expected curl to be installed, got %s",
			manager.installedArgs[0],
		)
	}
}

func TestExecuteNothingToInstall(t *testing.T) {
	manager := &mockPackageManager{
		installed: map[string]bool{
			"git": true,
		},
	}

	service := Service{
		Manager:     manager,
		ManagerType: packages.APT,
		Runner:      runner.Runner{},
	}

	plan := packages.PackagePlan{
		Installed: []string{"git"},
		Missing:   []string{},
	}

	verified, err := service.Execute(plan)

	if err != nil {
		t.Fatalf("expected execute to succeed, got %v", err)
	}

	if len(verified.Missing) != 0 {
		t.Fatalf(
			"expected no missing packages, got %v",
			verified.Missing,
		)
	}

	if manager.updateCalled {
		t.Fatal("expected Update not to be called")
	}

	if manager.installCalled {
		t.Fatal("expected Install not to be called")
	}
}
func TestExecuteVerificationFailure(t *testing.T) {
	manager := &mockPackageManager{
		installed: map[string]bool{
			"git": true,
		},
	}

	service := Service{
		Manager:     manager,
		ManagerType: packages.APT,
		Runner:      runner.Runner{},
	}

	plan := packages.PackagePlan{
		Installed: []string{"git"},
		Missing:   []string{"curl"},
	}

	verified, err := service.Execute(plan)

	if err != nil {
		t.Fatalf("expected execute to succeed, got %v", err)
	}

	if len(verified.Missing) != 0 {
		t.Fatalf(
			"expected curl to be installed by mock, got missing %v",
			verified.Missing,
		)
	}
}
