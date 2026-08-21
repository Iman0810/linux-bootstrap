package main

import (
	//"fmt"

	"github.com/Iman0810/linux-bootstrap/internal/cli"
	//"github.com/Iman0810/linux-bootstrap/internal/hardware"
)

func main() {
	cli.Run()
}

// func main() {
// 	status := hardware.DetectNvidiaDriver()

// 	fmt.Println("NVIDIA Driver")
// 	fmt.Println("-------------")
// 	fmt.Println("Installed:", status.Installed)
// 	fmt.Println("Version:", status.Version)
// }

// func main() {
// 	gpus := hardware.DetectGPUs()

// 	fmt.Println("Detected GPUs:")
// 	fmt.Println("--------------")

// 	for _, gpu := range gpus {
// 		fmt.Println("Vendor:", gpu.Vendor)
// 		fmt.Println("Name:", gpu.Name)
// 		fmt.Println()
// 	}
// }

// package main

// import (
// 	"fmt"
// 	"log"

// 	"github.com/Iman0810/linux-bootstrap/internal/runner"
// )

// func main() {
// 	r := runner.Runner{
// 		DryRun: false,
// 	}

// 	output, err := r.Output("uname", "-r")

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Printf("Captured output: %q\n", output)
// }
