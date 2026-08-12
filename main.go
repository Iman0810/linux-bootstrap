package main

import (
	"fmt"
	"log"

	
	"github.com/Iman0810/linux-bootstrap/internal/system"
	"github.com/Iman0810/linux-bootstrap/internal/packages"
)

func main() {
	
	osInfo, err := system.GetOSInfo()
	if err != nil {
		log.Fatal(err)
	}

	packageManager := packages.DetectManager(osInfo)

	manager := packages.GetPackageManager(packageManager)
	
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
}

