package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/Iman0810/linux-bootstrap/internal/packages"
	"github.com/Iman0810/linux-bootstrap/internal/profile"
	"github.com/Iman0810/linux-bootstrap/internal/prompt"
	"github.com/Iman0810/linux-bootstrap/internal/runner"
	"github.com/Iman0810/linux-bootstrap/internal/system"
)

func Run() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	command := os.Args[1]

	switch command {
	case "info":
		runInfo()

	case "setup":
		runSetup(os.Args[2:])

	default:
		fmt.Printf("Unknown command: %s\n\n", command)
		printUsage()
	}
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

func printUsage() {
	fmt.Println("Linux Bootstrap")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  linux-bootstrap info")
	fmt.Println("  linux-bootstrap setup [--profile <name>] [--dry-run]")
	fmt.Println()
	fmt.Println("Profiles:")
	fmt.Println("  essentials")
	fmt.Println("  development")
}
