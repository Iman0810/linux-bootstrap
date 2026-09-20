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
	r := runner.Runner{}

	return RunWithDependencies(
		system.GetOSInfo,
		hardware.DetectHardware,
		func(manager packages.Manager) packages.PackageManager {
			return packages.GetPackageManager(manager, r)
		},
	)
}

func RunWithDependencies(
	getOSInfo func() (system.OSInfo, error),
	detectHardware func() hardware.HardwareStatus,
	getPackageManager func(packages.Manager) packages.PackageManager,
) (Report, error) {
	osInfo, err := getOSInfo()
	if err != nil {
		return Report{}, err
	}

	packageManager := packages.DetectManager(osInfo)

	hardwareStatus := detectHardware()

	manager := getPackageManager(packageManager)

	var profiles []ProfileStatus

	if manager != nil {
		for _, p := range profile.List() {
			status, supported := profile.CheckStatus(
				manager,
				packageManager,
				p,
			)

			if !supported {
				profiles = append(profiles, ProfileStatus{
					Name: p.Name,
				})
				continue
			}

			profiles = append(profiles, ProfileStatus{
				Name:    p.Name,
				Ready:   len(status.Plan.Missing) == 0,
				Missing: status.Plan.Missing,
			})
		}
	}

	return BuildReport(
		osInfo,
		packageManager,
		hardwareStatus,
		profiles,
	), nil
}

func BuildReport(
	osInfo system.OSInfo,
	packageManager packages.Manager,
	hardwareStatus hardware.HardwareStatus,
	profiles []ProfileStatus,
) Report {
	report := Report{
		OS:             osInfo,
		PackageManager: packageManager,
		GPUs:           hardwareStatus.GPUs,
		NvidiaFound:    hardwareStatus.NvidiaFound,
		Profiles:       profiles,
	}

	if hardwareStatus.NvidiaFound && hardwareStatus.Nvidia != nil {
		report.NvidiaInstalled = hardwareStatus.Nvidia.Installed
		report.NvidiaVersion = hardwareStatus.Nvidia.Version
	}

	return report
}
