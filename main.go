package main

import (


	"github.com/Iman0810/linux-bootstrap/internal/cli"

)

func main() {
	cli.Run()
}

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

