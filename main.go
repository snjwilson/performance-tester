package main

import (
	"fmt"
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

	startPrometheus := func ()  {
		log.Println("Starting Prometheus...")
		cmd := exec.Command("prometheus", "--config.file=./config/prometheus.yml", "--web.enable-remote-write-receiver")
		_, err = cmd.Output()
		if err != nil {
			log.Printf("Error while running prometheus binary: %v", err)
			return
		}
	}
	stopPrometheus := func ()  {
		log.Println("Stopping Prometheus...")
		cmd := exec.Command("pkill", "prometheus")
		_, err = cmd.Output()
		if err != nil {
			log.Printf("Error while stopping prometheus binary: %v", err)
			return
		}
	}

	go startPrometheus()
	defer stopPrometheus()
	// todo: set remote write url k6 env K6_PROMETHEUS_RW_SERVER_URL=http://localhost:9090/api/v1/write
	cmd := exec.Command("k6", "run", "load-test.ts", "--vus", "2", "--duration", "5s", "--out", "experimental-prometheus-rw")
	output, err := cmd.Output()
	if err != nil {
		log.Printf("Error while running k6 binary: %v", err)
		return
	}

	// Log the output
	log.Println("Command Output:")
	log.Println(string(output))
	log.Println("Performance Tester successfully completed!")
	fmt.Println("Press Enter to exit...")
	fmt.Scanln()
}