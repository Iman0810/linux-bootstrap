package doctor

import (
	"fmt"
	"testing"

	"github.com/Iman0810/linux-bootstrap/internal/hardware"
	"github.com/Iman0810/linux-bootstrap/internal/packages"
	"github.com/Iman0810/linux-bootstrap/internal/system"
)

type fakePackageManager struct {
	installed map[string]bool
}

func (f fakePackageManager) Update() error {
	return nil
}

func (f fakePackageManager) Install(packages ...string) error {
	return nil
}

func (f fakePackageManager) IsInstalled(packageName string) bool {
	return f.installed[packageName]
}

func TestBuildReportBasicInfo(t *testing.T) {
	osInfo := system.OSInfo{
		Name:    "Pop!_OS",
		Version: "24.04 LTS",
		ID:      "pop",
		IDLike:  "ubuntu debian",
	}

	hardwareStatus := hardware.HardwareStatus{}

	report := BuildReport(
		osInfo,
		packages.APT,
		hardwareStatus,
		nil,
	)

	if report.OS != osInfo {
		t.Fatalf("unexpected OS info: got %+v, want %+v", report.OS, osInfo)
	}

	if report.PackageManager != packages.APT {
		t.Fatalf(
			"expected package manager %q, got %q",
			packages.APT,
			report.PackageManager,
		)
	}
}

func TestBuildReportNvidiaInstalled(t *testing.T) {
	gpus := []hardware.GPU{
		{
			Vendor: hardware.Intel,
			Name:   "Intel UHD Graphics",
		},
		{
			Vendor: hardware.NVIDIA,
			Name:   "NVIDIA GeForce RTX 3050 6GB Laptop GPU",
		},
	}

	hardwareStatus := hardware.HardwareStatus{
		GPUs:        gpus,
		NvidiaFound: true,
		Nvidia: &hardware.NvidiaStatus{
			Installed: true,
			Version:   "595.84",
		},
	}

	report := BuildReport(
		system.OSInfo{Name: "Pop!_OS"},
		packages.APT,
		hardwareStatus,
		nil,
	)

	if len(report.GPUs) != 2 {
		t.Fatalf("expected 2 GPUs, got %d", len(report.GPUs))
	}

	if !report.NvidiaFound {
		t.Fatal("expected NVIDIA to be detected")
	}

	if !report.NvidiaInstalled {
		t.Fatal("expected NVIDIA driver to be installed")
	}

	if report.NvidiaVersion != "595.84" {
		t.Fatalf(
			"expected NVIDIA version %q, got %q",
			"595.84",
			report.NvidiaVersion,
		)
	}
}

func TestBuildReportNvidiaMissing(t *testing.T) {
	hardwareStatus := hardware.HardwareStatus{
		GPUs: []hardware.GPU{
			{
				Vendor: hardware.NVIDIA,
				Name:   "NVIDIA GeForce RTX 3050",
			},
		},
		NvidiaFound: true,
		Nvidia: &hardware.NvidiaStatus{
			Installed: false,
		},
	}

	report := BuildReport(
		system.OSInfo{Name: "Pop!_OS"},
		packages.APT,
		hardwareStatus,
		nil,
	)

	if !report.NvidiaFound {
		t.Fatal("expected NVIDIA to be detected")
	}

	if report.NvidiaInstalled {
		t.Fatal("expected NVIDIA driver to be missing")
	}

	if report.NvidiaVersion != "" {
		t.Fatalf(
			"expected empty NVIDIA version, got %q",
			report.NvidiaVersion,
		)
	}
}

func TestBuildReportIntelOnly(t *testing.T) {
	hardwareStatus := hardware.HardwareStatus{
		GPUs: []hardware.GPU{
			{
				Vendor: hardware.Intel,
				Name:   "Intel UHD Graphics",
			},
		},
		NvidiaFound: false,
		Nvidia:      nil,
	}

	report := BuildReport(
		system.OSInfo{Name: "Ubuntu"},
		packages.APT,
		hardwareStatus,
		nil,
	)

	if len(report.GPUs) != 1 {
		t.Fatalf("expected 1 GPU, got %d", len(report.GPUs))
	}

	if report.NvidiaFound {
		t.Fatal("did not expect NVIDIA to be detected")
	}

	if report.NvidiaInstalled {
		t.Fatal("did not expect NVIDIA driver to be installed")
	}

	if report.NvidiaVersion != "" {
		t.Fatalf(
			"expected empty NVIDIA version, got %q",
			report.NvidiaVersion,
		)
	}
}

func TestBuildReportProfiles(t *testing.T) {
	profiles := []ProfileStatus{
		{
			Name:    "essentials",
			Ready:   true,
			Missing: nil,
		},
		{
			Name:    "development",
			Ready:   true,
			Missing: nil,
		},
		{
			Name:    "multimedia",
			Ready:   false,
			Missing: []string{"ffmpeg", "vlc"},
		},
	}

	report := BuildReport(
		system.OSInfo{Name: "Pop!_OS"},
		packages.APT,
		hardware.HardwareStatus{},
		profiles,
	)

	if len(report.Profiles) != 3 {
		t.Fatalf("expected 3 profiles, got %d", len(report.Profiles))
	}

	if !report.Profiles[0].Ready {
		t.Fatal("expected essentials profile to be ready")
	}

	if !report.Profiles[1].Ready {
		t.Fatal("expected development profile to be ready")
	}

	if report.Profiles[2].Ready {
		t.Fatal("expected multimedia profile to be missing packages")
	}

	if len(report.Profiles[2].Missing) != 2 {
		t.Fatalf(
			"expected 2 missing multimedia packages, got %d",
			len(report.Profiles[2].Missing),
		)
	}
}

func TestBuildReportNilNvidiaStatus(t *testing.T) {
	hardwareStatus := hardware.HardwareStatus{
		GPUs: []hardware.GPU{
			{
				Vendor: hardware.NVIDIA,
				Name:   "NVIDIA GeForce RTX 3050",
			},
		},
		NvidiaFound: true,
		Nvidia:      nil,
	}

	report := BuildReport(
		system.OSInfo{Name: "Pop!_OS"},
		packages.APT,
		hardwareStatus,
		nil,
	)

	if !report.NvidiaFound {
		t.Fatal("expected NVIDIA to be detected")
	}

	if report.NvidiaInstalled {
		t.Fatal("expected NVIDIA driver to remain false")
	}

	if report.NvidiaVersion != "" {
		t.Fatalf(
			"expected empty NVIDIA version, got %q",
			report.NvidiaVersion,
		)
	}
}
func TestRunWithDependenciesOSInfoError(t *testing.T) {
	expectedErr := fmt.Errorf("os detection failed")

	_, err := RunWithDependencies(
		func() (system.OSInfo, error) {
			return system.OSInfo{}, expectedErr
		},
		func() hardware.HardwareStatus {
			t.Fatal("hardware detection should not be called")
			return hardware.HardwareStatus{}
		},
		func(manager packages.Manager) packages.PackageManager {
			t.Fatal("package manager should not be created")
			return nil
		},
	)

	if err != expectedErr {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}
}

func TestRunWithDependenciesSupportedManager(t *testing.T) {
	hardwareStatus := hardware.HardwareStatus{
		GPUs: []hardware.GPU{
			{
				Vendor: hardware.Intel,
				Name:   "Intel UHD Graphics",
			},
		},
	}

	manager := fakePackageManager{
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

	report, err := RunWithDependencies(
		func() (system.OSInfo, error) {
			return system.OSInfo{
				Name: "Pop!_OS",
				ID:   "pop",
			}, nil
		},
		func() hardware.HardwareStatus {
			return hardwareStatus
		},
		func(managerType packages.Manager) packages.PackageManager {
			if managerType != packages.APT {
				t.Fatalf("expected APT, got %q", managerType)
			}

			return manager
		},
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if report.PackageManager != packages.APT {
		t.Fatalf(
			"expected package manager %q, got %q",
			packages.APT,
			report.PackageManager,
		)
	}

	if len(report.Profiles) != 3 {
		t.Fatalf("expected 3 profiles, got %d", len(report.Profiles))
	}

	for _, p := range report.Profiles {
		if !p.Ready {
			t.Fatalf("expected profile %q to be ready", p.Name)
		}
	}

	if len(report.GPUs) != 1 {
		t.Fatalf("expected 1 GPU, got %d", len(report.GPUs))
	}
}

func TestRunWithDependenciesMissingPackages(t *testing.T) {
	manager := fakePackageManager{
		installed: map[string]bool{
			"git":  true,
			"curl": true,
		},
	}

	report, err := RunWithDependencies(
		func() (system.OSInfo, error) {
			return system.OSInfo{
				Name: "Pop!_OS",
				ID:   "pop",
			}, nil
		},
		func() hardware.HardwareStatus {
			return hardware.HardwareStatus{}
		},
		func(managerType packages.Manager) packages.PackageManager {
			return manager
		},
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(report.Profiles) != 3 {
		t.Fatalf("expected 3 profiles, got %d", len(report.Profiles))
	}

	for _, p := range report.Profiles {
		if p.Ready {
			t.Fatalf("expected profile %q to have missing packages", p.Name)
		}

		if len(p.Missing) == 0 {
			t.Fatalf("expected profile %q to have missing packages", p.Name)
		}
	}
}

func TestRunWithDependenciesUnsupportedManager(t *testing.T) {
	report, err := RunWithDependencies(
		func() (system.OSInfo, error) {
			return system.OSInfo{
				Name: "Unknown Linux",
				ID:   "banana",
			}, nil
		},
		func() hardware.HardwareStatus {
			return hardware.HardwareStatus{}
		},
		func(managerType packages.Manager) packages.PackageManager {
			if managerType != packages.Unknown {
				t.Fatalf(
					"expected unknown manager, got %q",
					managerType,
				)
			}

			return nil
		},
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if report.PackageManager != packages.Unknown {
		t.Fatalf(
			"expected package manager %q, got %q",
			packages.Unknown,
			report.PackageManager,
		)
	}

	if len(report.Profiles) != 0 {
		t.Fatalf(
			"expected no profile statuses, got %d",
			len(report.Profiles),
		)
	}
}
func TestRunWithDependenciesUnsupportedProfileMapping(t *testing.T) {
	manager := fakePackageManager{
		installed: map[string]bool{},
	}

	report, err := RunWithDependencies(
		func() (system.OSInfo, error) {
			return system.OSInfo{
				Name: "Unknown Linux",
				ID:   "banana",
			}, nil
		},
		func() hardware.HardwareStatus {
			return hardware.HardwareStatus{}
		},
		func(managerType packages.Manager) packages.PackageManager {
			if managerType != packages.Unknown {
				t.Fatalf(
					"expected unknown manager, got %q",
					managerType,
				)
			}

			return manager
		},
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(report.Profiles) != 3 {
		t.Fatalf(
			"expected 3 profile statuses, got %d",
			len(report.Profiles),
		)
	}

	for _, p := range report.Profiles {
		if p.Ready {
			t.Fatalf("expected profile %q to not be ready", p.Name)
		}

		if len(p.Missing) != 0 {
			t.Fatalf(
				"expected no missing packages for unsupported profile %q, got %v",
				p.Name,
				p.Missing,
			)
		}
	}
}
