package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/Iman0810/linux-bootstrap/internal/doctor"
	"github.com/Iman0810/linux-bootstrap/internal/hardware"
	"github.com/Iman0810/linux-bootstrap/internal/packages"
	"github.com/Iman0810/linux-bootstrap/internal/profile"
	"github.com/Iman0810/linux-bootstrap/internal/prompt"
	"github.com/Iman0810/linux-bootstrap/internal/recommendation"
	"github.com/Iman0810/linux-bootstrap/internal/runner"
	"github.com/Iman0810/linux-bootstrap/internal/setup"
	"github.com/Iman0810/linux-bootstrap/internal/system"
)

func Run() error {
	if len(os.Args) < 2 {
		printUsage()
		return nil
	}

	command := os.Args[1]

	switch command {
	case "info":
		return runInfo()

	case "profiles":
		runProfiles()

	case "status":
		return runStatus()

	case "doctor":
		return runDoctor()

	case "setup":
		return runSetup(os.Args[2:])

	default:
		fmt.Printf("Unknown command: %s\n\n", command)
		printUsage()
		return fmt.Errorf("unknown command: %s", command)
	}

	return nil
}

func runInfo() error {
	osInfo, err := system.GetOSInfo()
	if err != nil {
		return err
	}

	packageManager := packages.DetectManager(osInfo)

	fmt.Println("Linux Bootstrap")
	fmt.Println("----------------")
	fmt.Println("OS:", osInfo.Name)
	fmt.Println("Version:", osInfo.Version)
	fmt.Println("ID:", osInfo.ID)
	fmt.Println("Based on:", osInfo.IDLike)
	fmt.Println("Package Manager:", packageManager)

	fmt.Println()
	fmt.Println("Hardware")
	fmt.Println("--------")

	hardwareStatus := hardware.DetectHardware()

	if len(hardwareStatus.GPUs) == 0 {
		fmt.Println("GPU: Not detected")
	} else {
		for _, gpu := range hardwareStatus.GPUs {
			fmt.Printf("GPU: %s (%s)\n", gpu.Name, gpu.Vendor)
		}
	}

	if hardwareStatus.NvidiaFound && hardwareStatus.Nvidia != nil {
		fmt.Println()
		fmt.Println("NVIDIA Driver")
		fmt.Println("-------------")

		if hardwareStatus.Nvidia.Installed {
			fmt.Println("Status:", "Installed")
			fmt.Println("Version:", hardwareStatus.Nvidia.Version)
		} else {
			fmt.Println("Status:", "Not detected")
		}
	}
	return nil
}

func runProfiles() {
	fmt.Println("Linux Bootstrap Profiles")
	fmt.Println("-------------------------")
	fmt.Println()

	for _, p := range profile.List() {
		fmt.Println(p.Name)
		fmt.Println("  " + p.Description)
		fmt.Println()
	}
}

func runSetup(args []string) error {
	setupFlags := flag.NewFlagSet("setup", flag.ContinueOnError)

	dryRun := setupFlags.Bool(
		"dry-run",
		false,
		"Show commands without executing them",
	)

	profileName := setupFlags.String(
		"profile",
		"essentials",
		"Profile to install",
	)

	if err := setupFlags.Parse(args); err != nil {
		return err
	}

	osInfo, err := system.GetOSInfo()
	if err != nil {
		return err
	}

	packageManager := packages.DetectManager(osInfo)

	r := runner.Runner{
		DryRun: *dryRun,
	}

	manager := packages.GetPackageManager(packageManager, r)

	if manager == nil {
		return fmt.Errorf(
			"unsupported package manager: %s",
			packageManager,
		)
	}

	selectedProfile, ok := profile.Get(*profileName)
	if !ok {
		return fmt.Errorf(
			"unknown profile: %s",
			*profileName,
		)
	}

	service := setup.Service{
		Manager:     manager,
		ManagerType: packageManager,
		Runner:      r,
	}

	result, ok := service.Prepare(selectedProfile)
	if !ok {
		return fmt.Errorf(
			"profile %q is not supported for package manager: %s",
			selectedProfile.Name,
			packageManager,
		)
	}

	plan := result.Plan

	fmt.Println("Linux Bootstrap Setup")
	fmt.Println("----------------------")
	fmt.Println("Profile:", selectedProfile.Name)
	fmt.Println("Description:", selectedProfile.Description)
	fmt.Println("OS:", osInfo.Name)
	fmt.Println("Package Manager:", packageManager)
	fmt.Println("Dry Run:", *dryRun)
	fmt.Println()

	fmt.Println("Package Check")
	fmt.Println("-------------")

	for _, packageName := range plan.Installed {
		fmt.Println("✓", packageName)
	}

	for _, packageName := range plan.Missing {
		fmt.Println("✗", packageName)
	}

	if len(plan.Missing) == 0 {
		fmt.Println("\nEverything is already installed.")
		return nil
	}

	fmt.Printf("\nPackages to install: %d\n", len(plan.Missing))

	if *dryRun {
		fmt.Println("Dry-run mode enabled. No changes will be made.")
	} else {
		confirmed := prompt.Confirm(
			"These operations will modify your system. Continue?",
		)

		if !confirmed {
			fmt.Println("Setup cancelled.")
			return nil
		}
	}

	execution, err := service.Execute(plan)
	if err != nil {
		return fmt.Errorf("setup failed: %w", err)
	}

	if *dryRun {
		fmt.Println("\nDry-run mode enabled. No changes were made.")
		return nil
	}

	if len(execution.Verified.Missing) > 0 {
		fmt.Println("\nSetup completed, but some packages are still missing:")

		for _, packageName := range execution.Verified.Missing {
			fmt.Println("✗", packageName)
		}

		return fmt.Errorf(
			"setup verification failed: %d package(s) still missing",
			len(execution.Verified.Missing),
		)
	}
	fmt.Println("\nSetup completed successfully.")
	return nil
}

func runStatus() error {
	osInfo, err := system.GetOSInfo()
	if err != nil {
		return err
	}

	packageManager := packages.DetectManager(osInfo)

	r := runner.Runner{}

	manager := packages.GetPackageManager(packageManager, r)

	if manager == nil {
		return err
	}

	fmt.Println("Linux Bootstrap Status")
	fmt.Println("----------------------")
	fmt.Println()
	fmt.Println("OS:", osInfo.Name, osInfo.Version)
	fmt.Println("Package Manager:", packageManager)
	fmt.Println()

	fmt.Println("Profiles")
	fmt.Println("--------")

	for _, p := range profile.List() {
		status, supported := profile.CheckStatus(
			manager,
			packageManager,
			p,
		)

		if !supported {
			fmt.Printf(
				"⚠ %-15s Unsupported\n",
				p.Name,
			)
			continue
		}

		if len(status.Plan.Missing) == 0 {
			fmt.Printf(
				"✓ %-15s Ready\n",
				p.Name,
			)
			continue
		}

		fmt.Printf(
			"✗ %-15s Missing %d package(s)\n",
			p.Name,
			len(status.Plan.Missing),
		)

		for _, packageName := range status.Plan.Missing {
			fmt.Println("    -", packageName)
		}
	}
	return nil
}

func runDoctor() error {
	report, err := doctor.Run()
	if err != nil {
		return err
	}

	fmt.Println("Linux Bootstrap Doctor")
	fmt.Println("----------------------")
	fmt.Println()

	fmt.Println("OS")
	fmt.Println("--")
	fmt.Println("✓", report.OS.Name, report.OS.Version)

	fmt.Println()
	fmt.Println("Package Manager")
	fmt.Println("---------------")
	fmt.Println("✓", report.PackageManager)

	fmt.Println()
	fmt.Println("Hardware")
	fmt.Println("--------")

	if len(report.GPUs) == 0 {
		fmt.Println("✗ No GPU detected")
	} else {
		for _, gpu := range report.GPUs {
			fmt.Printf("✓ %s (%s)\n", gpu.Name, gpu.Vendor)
		}
	}

	if report.NvidiaFound {
		fmt.Println()
		fmt.Println("NVIDIA Driver")
		fmt.Println("-------------")

		if report.NvidiaInstalled {
			fmt.Println("✓ Driver installed")
			fmt.Println("  Version:", report.NvidiaVersion)
		} else {
			fmt.Println("✗ NVIDIA driver not detected")
		}
	}

	fmt.Println()
	fmt.Println("Profiles")
	fmt.Println("--------")

	for _, p := range report.Profiles {
		if p.Ready {
			fmt.Printf("✓ %-15s Ready\n", p.Name)
		} else {
			fmt.Printf(
				"✗ %-15s Missing %d package(s)\n",
				p.Name,
				len(p.Missing),
			)

			for _, packageName := range p.Missing {
				fmt.Println("    -", packageName)
			}
		}
	}
	recommendations := recommendation.Generate(report)

	fmt.Println()
	fmt.Println("Recommendations")
	fmt.Println("---------------")

	if len(recommendations) == 0 {
		fmt.Println("✓ No issues found.")
		return nil
	}

	for _, rec := range recommendations {
		fmt.Println()
		fmt.Println("→", rec.Title)
		fmt.Println(" ", rec.Description)
		fmt.Println(" ", rec.Command)
	}
	return nil
}

func printUsage() {
	fmt.Println("Linux Bootstrap")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  linux-bootstrap info")
	fmt.Println("  linux-bootstrap profiles")
	fmt.Println("  linux-bootstrap status")
	fmt.Println("  linux-bootstrap doctor")
	fmt.Println("  linux-bootstrap setup [--profile <name>] [--dry-run]")
	fmt.Println()
	fmt.Println("Profiles:")
	fmt.Println("  essentials")
	fmt.Println("  development")
	fmt.Println("  multimedia")
}
