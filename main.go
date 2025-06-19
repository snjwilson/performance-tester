package main

import (
	"log"
	"os"
	"os/exec"
)

func main() {
	// Create a log file
	logFile, err := os.OpenFile("performance-tester.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()

	// Set log output to the file
	log.SetOutput(logFile)

	log.Println("Performance tester initializing...")
	cmd := exec.Command("k6", "run", "load-test.ts")

	// Capture the output of the command
	output, err := cmd.Output()
	if err != nil {
		log.Printf("Error while running k6 binary: %v", err)
		return
	}

	// Log the output
	log.Println("Command Output:")
	log.Println(string(output))

	log.Println("Performance Tester successfully completed!")
}