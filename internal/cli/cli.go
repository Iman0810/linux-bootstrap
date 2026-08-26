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
		runInfo()

	case "profiles":
		runProfiles()

	case "status":
		runStatus()

	case "doctor":
		runDoctor()

	case "setup":
		runSetup(os.Args[2:])

	default:
		fmt.Printf("Unknown command: %s\n\n", command)
		printUsage()
		return fmt.Errorf("unknown command: %s", command)
	}

	return nil
}

func runInfo() {
	osInfo, err := system.GetOSInfo()
	if err != nil {
		fmt.Println("Error:", err)
		return
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

func runSetup(args []string) {
	setupFlags := flag.NewFlagSet("setup", flag.ExitOnError)

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

	setupFlags.Parse(args)

	osInfo, err := system.GetOSInfo()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	packageManager := packages.DetectManager(osInfo)

	r := runner.Runner{
		DryRun: *dryRun,
	}

	manager := packages.GetPackageManager(packageManager, r)

	if manager == nil {
		fmt.Println("Unsupported package manager:", packageManager)
		return
	}

	selectedProfile, ok := profile.Get(*profileName)
	if !ok {
		fmt.Println("Unknown profile:", *profileName)
		return
	}

	plan := packages.BuildPlan(manager, selectedProfile.Packages)

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
		return
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
			return
		}
	}

	err = manager.Update()
	if err != nil {
		fmt.Println("Update failed:", err)
		return
	}

	err = manager.Install(plan.Missing...)
	if err != nil {
		fmt.Println("Installation failed:", err)
		return
	}

	if *dryRun {
		fmt.Println("\nDry-run mode enabled. No changes were made.")
	} else {
		fmt.Println("\nSetup completed successfully.")
	}
}

func runStatus() {
	osInfo, err := system.GetOSInfo()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	packageManager := packages.DetectManager(osInfo)

	r := runner.Runner{}

	manager := packages.GetPackageManager(packageManager, r)

	if manager == nil {
		fmt.Println("Unsupported package manager:", packageManager)
		return
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
		status := profile.CheckStatus(manager, p)

		if len(status.Plan.Missing) == 0 {
			fmt.Printf("✓ %-15s Ready\n", p.Name)
		} else {
			fmt.Printf(
				"✗ %-15s Missing %d package(s)\n",
				p.Name,
				len(status.Plan.Missing),
			)

			for _, packageName := range status.Plan.Missing {
				fmt.Println("    -", packageName)
			}
		}
	}
}

func runDoctor() {
	report, err := doctor.Run()
	if err != nil {
		fmt.Println("Error:", err)
		return
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
		return
	}

	for _, rec := range recommendations {
		fmt.Println()
		fmt.Println("→", rec.Title)
		fmt.Println(" ", rec.Description)
		fmt.Println(" ", rec.Command)
	}
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
