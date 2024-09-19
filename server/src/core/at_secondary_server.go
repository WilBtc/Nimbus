// server/src/core/at_secondary_server.go

package core

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/atsign-foundation/dess/server" // Importing the DESS package
	"server/utils"
)

// AtSecondaryServer struct manages the DESS server instance
type AtSecondaryServer struct {
	config *Config
	logger *utils.Logger
	server *server.AtServer // DESS AtServer type from the actual DESS package
	mu     sync.Mutex       // Mutex to protect start/stop operations
	ctx    context.Context
	cancel context.CancelFunc
}

// NewAtSecondaryServer initializes a new instance of AtSecondaryServer
func NewAtSecondaryServer(config *Config) *AtSecondaryServer {
	ctx, cancel := context.WithCancel(context.Background())
	return &AtSecondaryServer{
		config: config,
		logger: utils.NewLogger(config.LogLevel, "", false, 0), // Ensure logger is initialized with appropriate settings
		ctx:    ctx,
		cancel: cancel,
	}
}

// Start initializes and starts the DESS atProtocol secondary server
func (s *AtSecondaryServer) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.logger.Info("Starting DESS atProtocol secondary server...")

	// Validate configuration before starting
	if err := s.validateConfig(); err != nil {
		s.logger.Error("Configuration validation failed:", err)
		return err
	}

	// Create a new DESS server with the provided configuration
	s.server = server.NewAtServer(&server.Config{
		AtSign:           s.config.AtSign,
		Domain:           s.config.RootDomain,
		Port:             s.config.ServerPort,
		StoragePath:      s.config.StoragePath,
		CommitLogPath:    s.config.CommitLogPath,
		LogLevel:         s.config.LogLevel,
		AuthRequired:     s.config.AuthRequired,
		EncryptionConfig: s.config.EncryptionConfig,
		Secret:           s.config.Secret,
		Email:            s.config.Email,
	})

	// Set up signal handling for graceful shutdown
	s.setupSignalHandler()

	// Attempt to start the DESS server
	if err := s.server.Start(); err != nil {
		s.logger.Error("Error starting DESS atProtocol secondary server:", err)
		return err
	}

	s.logger.Info("DESS atProtocol secondary server started successfully.")
	return nil
}

// Stop gracefully stops the DESS atProtocol secondary server
func (s *AtSecondaryServer) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.logger.Info("Stopping DESS atProtocol secondary server...")

	// Cancel the context to signal all goroutines to terminate
	s.cancel()

	// Attempt to stop the server gracefully
	if err := s.server.Stop(); err != nil {
		s.logger.Error("Error stopping DESS atProtocol secondary server:", err)
		return err
	}

	s.logger.Info("DESS atProtocol secondary server stopped successfully.")
	return nil
}

// setupSignalHandler sets up OS signal handling for graceful shutdown
func (s *AtSecondaryServer) setupSignalHandler() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		s.logger.Info("Shutdown signal received. Initiating graceful shutdown...")

		// Stop the server gracefully
		if err := s.Stop(); err != nil {
			s.logger.Error("Error during graceful shutdown of DESS atProtocol secondary server:", err)
		}

		s.logger.Info("DESS atProtocol secondary server shutdown complete.")
		os.Exit(0)
	}()
}

// MonitorHealth continuously monitors the health of the DESS atProtocol secondary server
func (s *AtSecondaryServer) MonitorHealth() {
	s.logger.Info("Starting health monitoring for DESS atProtocol secondary server...")

	ticker := time.NewTicker(30 * time.Second) // Adjust monitoring interval as needed
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			s.logger.Info("Health monitoring stopped.")
			return
		case <-ticker.C:
			status, err := s.server.HealthCheck()
			if err != nil {
				s.logger.Error(fmt.Sprintf("Health check failed: %v", err))
				continue
			}

			if status.Healthy {
				s.logger.Info(fmt.Sprintf("Server health status: Healthy. Uptime: %v", status.Uptime))
			} else {
				s.logger.Error(fmt.Sprintf("Server health status: Unhealthy. Issues: %v", status.Issues))
			}
		}
	}
}

// Restart handles the restart of the server, typically used during updates or configuration changes
func (s *AtSecondaryServer) Restart() error {
	s.logger.Info("Restarting DESS atProtocol secondary server...")

	if err := s.Stop(); err != nil {
		s.logger.Error("Error stopping DESS atProtocol secondary server during restart:", err)
		return err
	}

	if err := s.Start(); err != nil {
		s.logger.Error("Error restarting DESS atProtocol secondary server:", err)
		return err
	}

	s.logger.Info("DESS atProtocol secondary server restarted successfully.")
	return nil
}

// validateConfig validates the server configuration to ensure it meets requirements
func (s *AtSecondaryServer) validateConfig() error {
	if s.config.AtSign == "" || s.config.RootDomain == "" || s.config.ServerPort == 0 {
		return fmt.Errorf("missing required configuration parameters")
	}
	return nil
}
