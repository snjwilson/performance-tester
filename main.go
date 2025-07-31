package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"

	"github.com/snjwilson/performance-tester/config"
)

func main() {
	// Setup structured JSON logging to stdout
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true, // Add file and line number to logs
		Level:     slog.LevelInfo, // Set default log level
	}))
	slog.SetDefault(logger)
	slog.Info("Performance Tester application starting", "version", "1.0.0", "pid", os.Getpid())
	slog.Info("Setting up environment")
	slog.Info("Default application directory", "dir", "~/performance-tester")
	// Check if the application folder exists
	appDir := "performance-tester"
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		slog.Error("Failed to get user home directory", "error", err)
		os.Exit(1)
	}
	appDir = userHomeDir + "/" + appDir

	if appDirInfo, err := os.Stat(appDir); err != nil {
		slog.Error("Application directory does not exist", "dir", appDir)
		slog.Info("Creating application directory", "dir", appDir)
		os.Mkdir(appDir, 0755)
	} else {
		slog.Info("Application directory exists", "dir", appDir, "size", appDirInfo.Size())
	}

	startPrometheus := func ()  {
		slog.Info("Starting Prometheus...")
		cmd := exec.Command("prometheus", "--config.file=./config/prometheus.yml", "--web.enable-remote-write-receiver")
		_, err = cmd.Output()
		if err != nil {
			slog.Error("Error while running prometheus binary: %v", err)
			return
		}
	}
	stopPrometheus := func ()  {
		slog.Info("Stopping Prometheus...")
		cmd := exec.Command("pkill", "prometheus")
		_, err = cmd.Output()
		if err != nil {
			slog.Error("Error while stopping prometheus binary: %v", err)
			return
		}
	}

	go startPrometheus()
	defer stopPrometheus()

	slog.Info("Creating config files...")
	slog.Info("Loading test configuration", "config_file", "config/default.yml")
	config, err := config.ReadConfig()
	if err != nil {
		slog.Error("Error reading config: %v", err)
		slog.Info("Exiting application due to configuration error...")
		os.Exit(1)
	}

	testTarget := "http://localhost:8080"
	// testDuration := "30s"
	// testVUs := 50

	slog.Info("Configuration loaded",
		"target", testTarget,
		"duration", config.Duration,
		"vus", config.Vus,
		)

	// --- Simulate k6 script generation ---
	k6ScriptPath := "./load-test.ts" // In real app, generate content
	slog.Info("Generated k6 script", "path", k6ScriptPath)

	// --- Simulate k6 execution ---
	slog.Info("Starting k6 test", "target", testTarget, "duration", config.Duration)

	// Example k6 command (simplified, real would have more flags for metrics, etc.)
	k6Cmd := exec.Command("k6", "run", k6ScriptPath,
		"--vus", config.Vus,
		"--duration", config.Duration,
		"--out", "experimental-prometheus-rw=http://localhost:9090/api/v1/write") // Replace with your Prometheus remote write URL

	// Capture k6's stdout/stderr
	k6Cmd.Stdout = os.Stdout // Direct k6 output to console
	k6Cmd.Stderr = os.Stderr // Direct k6 errors to console

	err = k6Cmd.Run() // Run the k6 process
	if err != nil {
		slog.Error("k6 test failed", "error", err, "command", k6Cmd.String())
		// Handle k6 specific exit codes if necessary
		if exitError, ok := err.(*exec.ExitError); ok {
			slog.Error("k6 exit error details", "code", exitError.ExitCode(), "stderr", string(exitError.Stderr))
		}
		os.Exit(1)
	}

	slog.Info("k6 test completed successfully")
	slog.Info("Performance Tester successfully completed!")
	slog.Info("Press Enter to exit...")
	fmt.Scanln()
}