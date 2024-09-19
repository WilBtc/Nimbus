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
	// Load configuration settings
	config, err := core.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	logger := utils.NewLogger(config.LogLevel, config.LogFilePath, config.RotateLogs, config.RotationInterval)
	defer logger.Close()

	// Initialize Security Gateway
	securityGateway := core.NewSecurityGateway(config)
	if err := securityGateway.Initialize(); err != nil {
		logger.Error("Error initializing Security Gateway:", err)
		return
	}

	// Initialize DESS (Distributed Edge Secondary Server)
	dessServer := core.NewAtSecondaryServer(config, securityGateway)
	if err := dessServer.Start(); err != nil {
		logger.Error("Failed to start DESS server:", err)
		return
	}
	defer dessServer.Stop()

	// Initialize Access Control
	accessControl := modules.NewAccessControl(log.Default(), dessServer.Server())
	// Additional initialization if needed

	// Initialize Analytics Engine
	analyticsEngine := modules.NewAnalyticsEngine(config.AnomalyThreshold, log.Default(), func(msg string) {
		// Handle alerts (e.g., send email or log)
		logger.Warn(msg)
	}, dessServer.Server())
	analyticsEngine.StartProcessing(config.AnalyticsInterval)

	// Setup signal handling for graceful shutdown
	setupSignalHandler(dessServer, securityGateway, analyticsEngine, logger)

	logger.Info("Nimbus Edge Server is running...")

	// Keep the server running indefinitely
	select {}
}

func setupSignalHandler(dessServer *core.AtSecondaryServer, securityGateway *core.SecurityGateway, analyticsEngine *modules.AnalyticsEngine, logger *utils.Logger) {
	// Capture OS signals for termination
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		logger.Info("Shutdown signal received. Shutting down Nimbus Edge Server...")

		// Stop Analytics Engine
		analyticsEngine.Stop()

		// Stop Security Gateway
		securityGateway.Stop()

		// Stop DESS server
		if err := dessServer.Stop(); err != nil {
			logger.Error("Error stopping DESS server:", err)
		}

		logger.Info("Nimbus Edge Server shutdown complete.")
		os.Exit(0)
	}()
}
