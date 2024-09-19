// src/main.go

package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	logger := utils.NewLogger(config.LogLevel)

	// Initialize Security Gateway
	securityGateway := core.NewSecurityGateway(config)
	err = securityGateway.Initialize()
	if err != nil {
		logger.Error("Error initializing Security Gateway:", err)
		return
	}

	// Initialize DESS (Distributed Edge Secondary Server)
	dessServer := core.NewDessServer(config)
	err = dessServer.Start()
	if err != nil {
		logger.Error("Failed to start DESS server:", err)
		return
	}

	// Initialize Analytics Engine for real-time data analysis
	analyticsEngine := modules.NewAnalyticsEngine(config)
	err = analyticsEngine.Start()
	if err != nil {
		logger.Error("Failed to start Analytics Engine:", err)
		return
	}

	// Initialize Routing Manager for managing data flows and decision logic
	routingManager := modules.NewRoutingManager(config)
	err = routingManager.Start()
	if err != nil {
		logger.Error("Failed to start Routing Manager:", err)
		return
	}

	// Initialize Access Control for managing permissions and secure access
	accessControl := modules.NewAccessControl(config)
	err = accessControl.Initialize()
	if err != nil {
		logger.Error("Failed to initialize Access Control:", err)
		return
	}

	// Setup signal handling for graceful shutdown
	setupSignalHandler(dessServer, analyticsEngine, routingManager, logger)

	logger.Info("NimBus Edge Server is running...")

	// Keep the server running indefinitely
	select {}
}

// setupSignalHandler handles graceful shutdown of server components
func setupSignalHandler(dessServer *core.DessServer, analyticsEngine *modules.AnalyticsEngine, routingManager *modules.RoutingManager, logger *utils.Logger) {
	// Capture OS signals for termination
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		logger.Info("Shutdown signal received. Shutting down NimBus Edge Server...")

		// Stop DESS server
		if err := dessServer.Stop(); err != nil {
			logger.Error("Error stopping DESS server:", err)
		}

		// Stop Analytics Engine
		if err := analyticsEngine.Stop(); err != nil {
			logger.Error("Error stopping Analytics Engine:", err)
		}

		// Stop Routing Manager
		if err := routingManager.Stop(); err != nil {
			logger.Error("Error stopping Routing Manager:", err)
		}

		logger.Info("NimBus Edge Server shutdown complete.")
		os.Exit(0)
	}()
}
