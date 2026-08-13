package main

import (
	"fmt"
	"log"

	"github.com/Iman0810/linux-bootstrap/internal/packages"
	"github.com/Iman0810/linux-bootstrap/internal/runner"
	"github.com/Iman0810/linux-bootstrap/internal/system"
)

func main() {
	osInfo, err := system.GetOSInfo()
	if err != nil {
		log.Fatal(err)
	}

	packageManager := packages.DetectManager(osInfo)

	r := runner.Runner{
		DryRun: true,
	}

	manager := packages.GetPackageManager(packageManager, r)

	fmt.Println("Linux Bootstrap")
	fmt.Println("----------------")
	fmt.Println("OS:", osInfo.Name)
	fmt.Println("Version:", osInfo.Version)
	fmt.Println("ID:", osInfo.ID)
	fmt.Println("Based on:", osInfo.IDLike)
	fmt.Println("Package Manager:", packageManager)

	if manager == nil {
		log.Fatal("Unsupported package manager")
	}

	fmt.Println("Package manager initialized successfully")

	fmt.Println("\nTesting package manager...")

	err = manager.Update()
	if err != nil {
		log.Fatal(err)
	}
	err = manager.Install("git", "curl", "wget")
	if err != nil {
	log.Fatal(err)
}
}