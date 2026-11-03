package doctor

import (
	"fmt"

	"github.com/Iman0810/linux-bootstrap/internal/hardware"
	"github.com/Iman0810/linux-bootstrap/internal/packages"
	"github.com/Iman0810/linux-bootstrap/internal/profile"
	"github.com/Iman0810/linux-bootstrap/internal/runner"
	"github.com/Iman0810/linux-bootstrap/internal/system"
)

type Report struct {
	OS                      system.OSInfo
	PackageManager          packages.Manager
	GPUs                    []hardware.GPU
	NvidiaFound             bool
	NvidiaInstalled         bool
	NvidiaVersion           string
	Profiles                []ProfileStatus
	PackageManagerAvailable bool
	PackageManagerError     string
}

type ProfileStatus struct {
	Name    string
	Ready   bool
	Missing []string
}

func Run() (Report, error) {
	r := runner.Runner{}

	return RunWithDependenciesAndValidator(
		system.GetOSInfo,
		hardware.DetectHardware,
		func(manager packages.Manager) packages.PackageManager {
			return packages.GetPackageManager(manager, r)
		},
		packages.ValidateManagerAvailable,
	)
}

func RunWithDependencies(
	getOSInfo func() (system.OSInfo, error),
	detectHardware func() hardware.HardwareStatus,
	getPackageManager func(packages.Manager) packages.PackageManager,
) (Report, error) {
	return RunWithDependenciesAndValidator(
		getOSInfo,
		detectHardware,
		getPackageManager,
		func(packages.Manager) error {
			return nil
		},
	)
}

func RunWithDependenciesAndValidator(
	getOSInfo func() (system.OSInfo, error),
	detectHardware func() hardware.HardwareStatus,
	getPackageManager func(packages.Manager) packages.PackageManager,
	validateManager func(packages.Manager) error,
) (Report, error) {
	osInfo, err := getOSInfo()
	if err != nil {
		return Report{}, err
	}

	packageManager := packages.DetectManager(osInfo)
	hardwareStatus := detectHardware()
	manager := getPackageManager(packageManager)

	var profiles []ProfileStatus
	managerAvailable := false
	managerError := ""

	if manager == nil {
		if packageManager == packages.Unknown {
			managerError = "unsupported distribution: no supported package manager mapping found"
		} else {
			managerError = fmt.Sprintf(
				"no package manager implementation available for %s",
				packageManager,
			)
		}
	} else if err := validateManager(packageManager); err != nil {
		managerError = err.Error()
	} else {
		managerAvailable = true

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

	report := BuildReport(
		osInfo,
		packageManager,
		hardwareStatus,
		profiles,
	)

	report.PackageManagerAvailable = managerAvailable
	report.PackageManagerError = managerError

	return report, nil
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
