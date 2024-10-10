// server/src/main.go
package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"server/core"
	"server/modules"
	"server/utils"
)

func main() {
	// Load configuration settings from environment or config file
	config, err := core.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger with rotation and log level settings
	logger := utils.NewLogger(config.LogLevel, config.LogFilePath, config.RotateLogs, config.RotationInterval)
	defer logger.Close()

	// Initialize the Security Gateway for handling encryption, IDS/IPS
	securityGateway := core.NewSecurityGateway(config)
	if err := securityGateway.Initialize(); err != nil {
		logger.Error("Error initializing Security Gateway:", err)
		return
	}

	// Start the DESS (Distributed Edge Secondary Server) for secure communication
	dessServer := core.NewAtSecondaryServer(config, securityGateway)
	if err := dessServer.Start(); err != nil {
		logger.Error("Failed to start DESS server:", err)
		return
	}
	defer dessServer.Stop() // Ensure the DESS server stops gracefully

	// Initialize Access Control for managing roles and permissions
	accessControl := modules.NewAccessControl(log.Default(), dessServer.Server())
	// Add any required access rules or additional logic if needed

	// Initialize the Analytics Engine for processing real-time data and detecting anomalies
	analyticsEngine := modules.NewAnalyticsEngine(config.AnomalyThreshold, log.Default(), func(msg string) {
		// Handle anomaly alerts (e.g., log or send notifications)
		logger.Warn(msg)
	}, dessServer.Server())
	analyticsEngine.StartProcessing(config.AnalyticsInterval) // Start real-time analytics

	// Set up signal handling for graceful shutdown
	setupSignalHandler(dessServer, securityGateway, analyticsEngine, logger)

	// Log that the server is running
	logger.Info("Nimbus Edge Server is running...")

	// Block forever to keep the server running
	select {}
}

// setupSignalHandler captures OS signals (e.g., SIGINT, SIGTERM) to gracefully shut down the server
func setupSignalHandler(dessServer *core.AtSecondaryServer, securityGateway *core.SecurityGateway, analyticsEngine *modules.AnalyticsEngine, logger *utils.Logger) {
	// Create a channel to listen for termination signals
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		// Wait for a signal to trigger the shutdown process
		<-sigs
		logger.Info("Shutdown signal received. Shutting down Nimbus Edge Server...")

		// Stop the Analytics Engine to halt real-time data processing
		analyticsEngine.Stop()

		// Stop the Security Gateway (e.g., stop monitoring and encryption services)
		securityGateway.Stop()

		// Stop the DESS server gracefully
		if err := dessServer.Stop(); err != nil {
			logger.Error("Error stopping DESS server:", err)
		}

		logger.Info("Nimbus Edge Server shutdown complete.")
		os.Exit(0) // Exit the program cleanly
	}()
}
