package setup

import (
	"fmt"
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
	updateErr     error
	installErr    error
}

func (m *mockPackageManager) Update() error {
	m.updateCalled = true
	return m.updateErr
}

func (m *mockPackageManager) Install(packages ...string) error {
	m.installCalled = true
	m.installedArgs = packages

	if m.installErr != nil {
		return m.installErr
	}

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

	execution, err := service.Execute(plan)

	if err != nil {
		t.Fatalf("expected execute to succeed, got %v", err)
	}

	if len(execution.Verified.Missing) != 0 {
		t.Fatalf(
			"expected no missing packages after verification, got %v",
			execution.Verified.Missing,
		)
	}

	if len(execution.Verified.Installed) != 2 {
		t.Fatalf(
			"expected 2 installed packages after verification, got %d",
			len(execution.Verified.Installed),
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

	execution, err := service.Execute(plan)

	if err != nil {
		t.Fatalf("expected execute to succeed, got %v", err)
	}

	if len(execution.Verified.Missing) != 0 {
		t.Fatalf(
			"expected no missing packages, got %v",
			execution.Verified.Missing,
		)
	}

	if manager.updateCalled {
		t.Fatal("expected Update not to be called")
	}

	if manager.installCalled {
		t.Fatal("expected Install not to be called")
	}
}

func TestExecuteVerification(t *testing.T) {
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

	execution, err := service.Execute(plan)

	if err != nil {
		t.Fatalf("expected execute to succeed, got %v", err)
	}

	if len(execution.Verified.Missing) != 0 {
		t.Fatalf(
			"expected curl to be installed by mock, got missing %v",
			execution.Verified.Missing,
		)
	}

	if len(execution.Verified.Installed) != 2 {
		t.Fatalf(
			"expected 2 verified installed packages, got %d",
			len(execution.Verified.Installed),
		)
	}
}

func TestExecuteUpdateFailure(t *testing.T) {
	manager := &mockPackageManager{
		installed: map[string]bool{
			"git": true,
		},
		updateErr: fmt.Errorf("update failed"),
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

	_, err := service.Execute(plan)

	if err == nil {
		t.Fatal("expected execute to fail")
	}

	if !manager.updateCalled {
		t.Fatal("expected Update to be called")
	}

	if manager.installCalled {
		t.Fatal("expected Install not to be called after Update failure")
	}
}

func TestExecuteInstallFailure(t *testing.T) {
	manager := &mockPackageManager{
		installed: map[string]bool{
			"git": true,
		},
		installErr: fmt.Errorf("installation failed"),
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

	_, err := service.Execute(plan)

	if err == nil {
		t.Fatal("expected execute to fail")
	}

	if !manager.updateCalled {
		t.Fatal("expected Update to be called")
	}

	if !manager.installCalled {
		t.Fatal("expected Install to be called")
	}
}
