package main

import (
	"fmt"
	"os/exec"
)

func main() {
	fmt.Println("Performance tester initializing...")
	cmd := exec.Command("k6", "run", "load-test.ts")

	// Capture the output of the command
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error while running k6 binary:", err)
		return
	}

	// Print the output
	fmt.Println("Command Output:")
	fmt.Println(string(output))

	fmt.Println("Performance Tester successfully completed!")
}