package doctor

import (
	"github.com/Iman0810/linux-bootstrap/internal/hardware"
	"github.com/Iman0810/linux-bootstrap/internal/packages"
	"github.com/Iman0810/linux-bootstrap/internal/profile"
	"github.com/Iman0810/linux-bootstrap/internal/runner"
	"github.com/Iman0810/linux-bootstrap/internal/system"
)

type Report struct {
	OS              system.OSInfo
	PackageManager  packages.Manager
	GPUs            []hardware.GPU
	NvidiaFound     bool
	NvidiaInstalled bool
	NvidiaVersion   string
	Profiles        []ProfileStatus
}

type ProfileStatus struct {
	Name    string
	Ready   bool
	Missing []string
}

func Run() (Report, error) {
	osInfo, err := system.GetOSInfo()
	if err != nil {
		return Report{}, err
	}

	packageManager := packages.DetectManager(osInfo)

	r := runner.Runner{}

	manager := packages.GetPackageManager(packageManager, r)

	report := Report{
		OS:             osInfo,
		PackageManager: packageManager,
	}

	// Hardware
	hardwareStatus := hardware.DetectHardware()

	report.GPUs = hardwareStatus.GPUs
	report.NvidiaFound = hardwareStatus.NvidiaFound

	if hardwareStatus.NvidiaFound && hardwareStatus.Nvidia != nil {
		report.NvidiaInstalled = hardwareStatus.Nvidia.Installed
		report.NvidiaVersion = hardwareStatus.Nvidia.Version
	}

	// Profile status
	if manager != nil {
		for _, p := range profile.List() {
			status := profile.CheckStatus(manager, p)

			report.Profiles = append(report.Profiles, ProfileStatus{
				Name:    p.Name,
				Ready:   len(status.Plan.Missing) == 0,
				Missing: status.Plan.Missing,
			})
		}
	}

	return report, nil
}