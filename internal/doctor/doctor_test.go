package doctor

import (
	"testing"

	"github.com/Iman0810/linux-bootstrap/internal/hardware"
	"github.com/Iman0810/linux-bootstrap/internal/packages"
	"github.com/Iman0810/linux-bootstrap/internal/system"
)

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
