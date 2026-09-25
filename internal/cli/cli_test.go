package cli

import (
	"errors"
	"testing"

	"github.com/Iman0810/linux-bootstrap/internal/doctor"
	"github.com/Iman0810/linux-bootstrap/internal/hardware"
	"github.com/Iman0810/linux-bootstrap/internal/packages"
	"github.com/Iman0810/linux-bootstrap/internal/runner"
	"github.com/Iman0810/linux-bootstrap/internal/system"
)

type statusFakePackageManager struct {
	installed         map[string]bool
	updateErr         error
	installErr        error
	updated           bool
	installedPackages []string
}

func (f *statusFakePackageManager) Update() error {
	f.updated = true
	return f.updateErr
}

func (f *statusFakePackageManager) Install(packages ...string) error {
	if f.installErr != nil {
		return f.installErr
	}

	f.installedPackages = append(
		f.installedPackages,
		packages...,
	)

	for _, packageName := range packages {
		f.installed[packageName] = true
	}

	return nil
}

func (f *statusFakePackageManager) IsInstalled(packageName string) bool {
	return f.installed[packageName]
}

func TestRunArgsUnknownCommand(t *testing.T) {
	err := RunArgs([]string{"banana"})

	if err == nil {
		t.Fatal("expected error for unknown command")
	}
}

func TestRunArgsNoCommand(t *testing.T) {
	err := RunArgs([]string{})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestRunArgsInvalidSetupFlag(t *testing.T) {
	err := RunArgs([]string{
		"setup",
		"--banana",
	})

	if err == nil {
		t.Fatal("expected error for invalid setup flag")
	}
}
func TestRunInfoWithDependenciesOSError(t *testing.T) {
	expectedErr := errors.New("os detection failed")

	getOSInfo := func() (system.OSInfo, error) {
		return system.OSInfo{}, expectedErr
	}

	detectHardware := func() hardware.HardwareStatus {
		t.Fatal("detectHardware should not be called")
		return hardware.HardwareStatus{}
	}

	err := runInfoWithDependencies(getOSInfo, detectHardware)

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected %v, got %v", expectedErr, err)
	}
}
func TestRunInfoWithDependenciesNoGPU(t *testing.T) {
	getOSInfo := func() (system.OSInfo, error) {
		return system.OSInfo{
			Name:    "Test Linux",
			Version: "1.0",
			ID:      "test",
			IDLike:  "linux",
		}, nil
	}

	detectHardware := func() hardware.HardwareStatus {
		return hardware.HardwareStatus{
			GPUs: []hardware.GPU{},
		}
	}

	err := runInfoWithDependencies(getOSInfo, detectHardware)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
func TestRunInfoWithDependenciesNvidiaInstalled(t *testing.T) {
	getOSInfo := func() (system.OSInfo, error) {
		return system.OSInfo{
			Name:    "Test Linux",
			Version: "1.0",
			ID:      "test",
			IDLike:  "linux",
		}, nil
	}

	detectHardware := func() hardware.HardwareStatus {
		return hardware.HardwareStatus{
			GPUs: []hardware.GPU{
				{
					Vendor: hardware.NVIDIA,
					Name:   "Test NVIDIA GPU",
				},
			},
			NvidiaFound: true,
			Nvidia: &hardware.NvidiaStatus{
				Installed: true,
				Version:   "595.84",
			},
		}
	}

	err := runInfoWithDependencies(getOSInfo, detectHardware)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
func TestRunInfoWithDependenciesNvidiaMissing(t *testing.T) {
	getOSInfo := func() (system.OSInfo, error) {
		return system.OSInfo{
			Name:    "Test Linux",
			Version: "1.0",
			ID:      "test",
			IDLike:  "linux",
		}, nil
	}

	detectHardware := func() hardware.HardwareStatus {
		return hardware.HardwareStatus{
			GPUs: []hardware.GPU{
				{
					Vendor: hardware.NVIDIA,
					Name:   "Test NVIDIA GPU",
				},
			},
			NvidiaFound: true,
			Nvidia: &hardware.NvidiaStatus{
				Installed: false,
			},
		}
	}

	err := runInfoWithDependencies(getOSInfo, detectHardware)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
func TestRunStatusWithDependenciesOSError(t *testing.T) {
	expectedErr := errors.New("os detection failed")

	getOSInfo := func() (system.OSInfo, error) {
		return system.OSInfo{}, expectedErr
	}

	getPackageManager := func(managerType packages.Manager) packages.PackageManager {
		t.Fatal("getPackageManager should not be called")
		return nil
	}

	err := runStatusWithDependencies(
		getOSInfo,
		getPackageManager,
	)

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected %v, got %v", expectedErr, err)
	}
}
func TestRunStatusWithDependenciesUnsupportedManager(t *testing.T) {
	getOSInfo := func() (system.OSInfo, error) {
		return system.OSInfo{
			Name:    "Unknown Linux",
			Version: "1.0",
			ID:      "unknown",
		}, nil
	}

	getPackageManager := func(managerType packages.Manager) packages.PackageManager {
		return nil
	}

	err := runStatusWithDependencies(
		getOSInfo,
		getPackageManager,
	)

	if err == nil {
		t.Fatal("expected error for unsupported package manager")
	}
}
func TestRunStatusWithDependenciesReadyProfile(t *testing.T) {
	getOSInfo := func() (system.OSInfo, error) {
		return system.OSInfo{
			Name:    "Test Linux",
			Version: "1.0",
			ID:      "ubuntu",
		}, nil
	}

	getPackageManager := func(managerType packages.Manager) packages.PackageManager {
		return &statusFakePackageManager{
			installed: map[string]bool{
				"git":             true,
				"curl":            true,
				"wget":            true,
				"unzip":           true,
				"build-essential": true,
				"ffmpeg":          true,
				"vlc":             true,
			},
		}
	}

	err := runStatusWithDependencies(
		getOSInfo,
		getPackageManager,
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
func TestRunStatusWithDependenciesMissingPackages(t *testing.T) {
	getOSInfo := func() (system.OSInfo, error) {
		return system.OSInfo{
			Name:    "Test Linux",
			Version: "1.0",
			ID:      "ubuntu",
		}, nil
	}

	getPackageManager := func(managerType packages.Manager) packages.PackageManager {
		return &statusFakePackageManager{
			installed: map[string]bool{
				"git":   true,
				"curl":  true,
				"wget":  false,
				"unzip": false,
			},
		}
	}

	err := runStatusWithDependencies(
		getOSInfo,
		getPackageManager,
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
func TestRunStatusWithDependenciesUnsupportedProfileMapping(t *testing.T) {
	getOSInfo := func() (system.OSInfo, error) {
		return system.OSInfo{
			Name:    "Test Linux",
			Version: "1.0",
			ID:      "unknown",
		}, nil
	}

	getPackageManager := func(managerType packages.Manager) packages.PackageManager {
		return &statusFakePackageManager{
			installed: map[string]bool{},
		}
	}

	err := runStatusWithDependencies(
		getOSInfo,
		getPackageManager,
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
func TestRunDoctorWithDependenciesError(t *testing.T) {
	expectedErr := errors.New("doctor failed")

	runDoctor := func() (doctor.Report, error) {
		return doctor.Report{}, expectedErr
	}

	err := runDoctorWithDependencies(runDoctor)

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected %v, got %v", expectedErr, err)
	}
}
func TestRunDoctorWithDependenciesReport(t *testing.T) {
	runDoctor := func() (doctor.Report, error) {
		return doctor.Report{
			OS: system.OSInfo{
				Name:    "Test Linux",
				Version: "1.0",
			},
			PackageManager: packages.APT,
			GPUs: []hardware.GPU{
				{
					Vendor: hardware.NVIDIA,
					Name:   "Test NVIDIA GPU",
				},
			},
			NvidiaFound:     true,
			NvidiaInstalled: true,
			NvidiaVersion:   "595.84",
			Profiles: []doctor.ProfileStatus{
				{
					Name:  "essentials",
					Ready: true,
				},
				{
					Name:    "development",
					Ready:   false,
					Missing: []string{"build-essential"},
				},
			},
		}, nil
	}

	err := runDoctorWithDependencies(runDoctor)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
func TestRunDoctorWithDependenciesNoGPUAndNvidiaMissing(t *testing.T) {
	runDoctor := func() (doctor.Report, error) {
		return doctor.Report{
			OS: system.OSInfo{
				Name:    "Test Linux",
				Version: "1.0",
			},
			PackageManager:  packages.APT,
			GPUs:            []hardware.GPU{},
			NvidiaFound:     true,
			NvidiaInstalled: false,
			Profiles: []doctor.ProfileStatus{
				{
					Name:    "development",
					Ready:   false,
					Missing: []string{"build-essential"},
				},
			},
		}, nil
	}

	err := runDoctorWithDependencies(runDoctor)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
func TestRunProfiles(t *testing.T) {
	runProfiles()
}
func TestRunDoctorWithDependenciesNoRecommendations(t *testing.T) {
	runDoctor := func() (doctor.Report, error) {
		return doctor.Report{
			OS: system.OSInfo{
				Name:    "Test Linux",
				Version: "1.0",
			},
			PackageManager: packages.APT,
			GPUs: []hardware.GPU{
				{
					Vendor: hardware.Intel,
					Name:   "Test Intel GPU",
				},
			},
			Profiles: []doctor.ProfileStatus{
				{
					Name:  "essentials",
					Ready: true,
				},
			},
		}, nil
	}

	err := runDoctorWithDependencies(runDoctor)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
func TestRunSetupWithDependenciesOSError(t *testing.T) {
	expectedErr := errors.New("os detection failed")

	getOSInfo := func() (system.OSInfo, error) {
		return system.OSInfo{}, expectedErr
	}

	getPackageManager := func(
		managerType packages.Manager,
		r runner.Runner,
	) packages.PackageManager {
		t.Fatal("getPackageManager should not be called")
		return nil
	}

	err := runSetupWithDependencies(
		[]string{},
		getOSInfo,
		getPackageManager,
		func(string) bool {
			return false
		},
	)

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected %v, got %v", expectedErr, err)
	}
}

func TestRunSetupWithDependenciesUnsupportedManager(t *testing.T) {
	getOSInfo := func() (system.OSInfo, error) {
		return system.OSInfo{
			Name:    "Unknown Linux",
			Version: "1.0",
			ID:      "unknown",
		}, nil
	}

	getPackageManager := func(
		managerType packages.Manager,
		r runner.Runner,
	) packages.PackageManager {
		return nil
	}

	err := runSetupWithDependencies(
		[]string{},
		getOSInfo,
		getPackageManager,
		func(string) bool {
			return false
		},
	)

	if err == nil {
		t.Fatal("expected error for unsupported package manager")
	}
}

func TestRunSetupWithDependenciesUnknownProfile(t *testing.T) {
	getOSInfo := func() (system.OSInfo, error) {
		return system.OSInfo{
			Name:    "Test Linux",
			Version: "1.0",
			ID:      "ubuntu",
		}, nil
	}

	getPackageManager := func(
		managerType packages.Manager,
		r runner.Runner,
	) packages.PackageManager {
		return &statusFakePackageManager{
			installed: map[string]bool{},
		}
	}

	err := runSetupWithDependencies(
		[]string{"--profile", "banana"},
		getOSInfo,
		getPackageManager,
		func(string) bool {
			return false
		},
	)

	if err == nil {
		t.Fatal("expected error for unknown profile")
	}
}

func TestRunSetupWithDependenciesUnsupportedProfileMapping(t *testing.T) {
	getOSInfo := func() (system.OSInfo, error) {
		return system.OSInfo{
			Name:    "Test Linux",
			Version: "1.0",
			ID:      "unknown",
		}, nil
	}

	getPackageManager := func(
		managerType packages.Manager,
		r runner.Runner,
	) packages.PackageManager {
		return &statusFakePackageManager{
			installed: map[string]bool{},
		}
	}

	err := runSetupWithDependencies(
		[]string{"--profile", "essentials"},
		getOSInfo,
		getPackageManager,
		func(string) bool {
			return false
		},
	)

	if err == nil {
		t.Fatal("expected error for unsupported profile mapping")
	}
}

func TestRunSetupWithDependenciesDryRun(t *testing.T) {
	getOSInfo := func() (system.OSInfo, error) {
		return system.OSInfo{
			Name:    "Test Linux",
			Version: "1.0",
			ID:      "ubuntu",
		}, nil
	}

	getPackageManager := func(
		managerType packages.Manager,
		r runner.Runner,
	) packages.PackageManager {
		return &statusFakePackageManager{
			installed: map[string]bool{},
		}
	}

	err := runSetupWithDependencies(
		[]string{
			"--profile",
			"essentials",
			"--dry-run",
		},
		getOSInfo,
		getPackageManager,
		func(string) bool {
			return false
		},
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestRunSetupWithDependenciesEverythingInstalled(t *testing.T) {
	getOSInfo := func() (system.OSInfo, error) {
		return system.OSInfo{
			Name:    "Test Linux",
			Version: "1.0",
			ID:      "ubuntu",
		}, nil
	}

	getPackageManager := func(
		managerType packages.Manager,
		r runner.Runner,
	) packages.PackageManager {
		return &statusFakePackageManager{
			installed: map[string]bool{
				"git":   true,
				"curl":  true,
				"wget":  true,
				"unzip": true,
			},
		}
	}

	err := runSetupWithDependencies(
		[]string{"--profile", "essentials"},
		getOSInfo,
		getPackageManager,
		func(string) bool {
			return false
		},
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
func TestRunSetupWithDependenciesCancelled(t *testing.T) {
	getOSInfo := func() (system.OSInfo, error) {
		return system.OSInfo{
			Name:    "Test Linux",
			Version: "1.0",
			ID:      "ubuntu",
		}, nil
	}

	getPackageManager := func(
		managerType packages.Manager,
		r runner.Runner,
	) packages.PackageManager {
		return &statusFakePackageManager{
			installed: map[string]bool{},
		}
	}

	err := runSetupWithDependencies(
		[]string{"--profile", "essentials"},
		getOSInfo,
		getPackageManager,
		func(string) bool {
			return false
		},
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
func TestRunSetupWithDependenciesSuccess(t *testing.T) {
	getOSInfo := func() (system.OSInfo, error) {
		return system.OSInfo{
			Name:    "Test Linux",
			Version: "1.0",
			ID:      "ubuntu",
		}, nil
	}

	var fake *statusFakePackageManager

	getPackageManager := func(
		managerType packages.Manager,
		r runner.Runner,
	) packages.PackageManager {
		fake = &statusFakePackageManager{
			installed: map[string]bool{},
		}

		return fake
	}

	err := runSetupWithDependencies(
		[]string{"--profile", "essentials"},
		getOSInfo,
		getPackageManager,
		func(string) bool {
			return true
		},
	)

	if err != nil {
		t.Fatalf("expected setup to succeed, got %v", err)
	}

	if !fake.updated {
		t.Fatal("expected package manager Update to be called")
	}

	if len(fake.installedPackages) != 4 {
		t.Fatalf(
			"expected 4 packages to be installed, got %d",
			len(fake.installedPackages),
		)
	}

	for _, packageName := range []string{
		"git",
		"curl",
		"wget",
		"unzip",
	} {
		if !fake.installed[packageName] {
			t.Errorf("expected %s to be installed", packageName)
		}
	}
}
func TestRunSetupWithDependenciesUpdateFailure(t *testing.T) {
	expectedErr := errors.New("update failed")

	getOSInfo := func() (system.OSInfo, error) {
		return system.OSInfo{
			Name:    "Test Linux",
			Version: "1.0",
			ID:      "ubuntu",
		}, nil
	}

	var manager *statusFakePackageManager

	getPackageManager := func(
		managerType packages.Manager,
		r runner.Runner,
	) packages.PackageManager {
		manager = &statusFakePackageManager{
			installed: map[string]bool{},
			updateErr: expectedErr,
		}

		return manager
	}

	err := runSetupWithDependencies(
		[]string{"--profile", "essentials"},
		getOSInfo,
		getPackageManager,
		func(string) bool {
			return true
		},
	)

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected %v, got %v", expectedErr, err)
	}

	if !manager.updated {
		t.Fatal("expected Update to be called")
	}

	if len(manager.installedPackages) != 0 {
		t.Fatal("Install should not be called after Update fails")
	}
}
func TestRunSetupWithDependenciesInstallFailure(t *testing.T) {
	expectedErr := errors.New("install failed")

	getOSInfo := func() (system.OSInfo, error) {
		return system.OSInfo{
			Name:    "Test Linux",
			Version: "1.0",
			ID:      "ubuntu",
		}, nil
	}

	var manager *statusFakePackageManager

	getPackageManager := func(
		managerType packages.Manager,
		r runner.Runner,
	) packages.PackageManager {
		manager = &statusFakePackageManager{
			installed:  map[string]bool{},
			installErr: expectedErr,
		}

		return manager
	}

	err := runSetupWithDependencies(
		[]string{"--profile", "essentials"},
		getOSInfo,
		getPackageManager,
		func(string) bool {
			return true
		},
	)

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected %v, got %v", expectedErr, err)
	}

	if !manager.updated {
		t.Fatal("expected Update to be called")
	}

	if len(manager.installedPackages) != 0 {
		t.Fatal("expected installedPackages to remain empty")
	}
}
