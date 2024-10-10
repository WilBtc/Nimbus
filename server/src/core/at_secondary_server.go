//Nimbus/server/src/core/at_secondary_server.go
package core

import (
	"context"
	"fmt"
	"net"
	"sync"

	"github.com/atsign-foundation/at_server/server"
	"server/utils"
)

// AtSecondaryServer manages the DESS server instance and its lifecycle
type AtSecondaryServer struct {
	config          *Config
	logger          *utils.Logger
	server          *server.AtServer
	mu              sync.Mutex
	ctx             context.Context
	cancel          context.CancelFunc
	securityGateway *SecurityGateway
}

// NewAtSecondaryServer initializes a new instance of AtSecondaryServer with provided configuration and security gateway
func NewAtSecondaryServer(config *Config, securityGateway *SecurityGateway) *AtSecondaryServer {
	ctx, cancel := context.WithCancel(context.Background())
	return &AtSecondaryServer{
		config:          config,
		logger:          utils.NewLogger(config.LogLevel, config.LogFilePath, config.RotateLogs, config.RotationInterval),
		ctx:             ctx,
		cancel:          cancel,
		securityGateway: securityGateway,
	}
}

// Start initializes and runs the DESS atProtocol secondary server
func (s *AtSecondaryServer) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.logger.Info("Starting DESS atProtocol secondary server...")

	// Validate configuration before starting
	if err := s.validateConfig(); err != nil {
		s.logger.Error("Configuration validation failed:", err)
		return err
	}

	// Create a new DESS server instance with provided configuration
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

	// Prepare secure listener with SecurityGateway
	address := fmt.Sprintf("%s:%d", s.config.ServerHost, s.config.ServerPort)
	secureListener, err := s.securityGateway.SecureListener(address)
	if err != nil {
		s.logger.Error("Failed to create secure listener:", err)
		return err
	}

	// Run the server in a separate goroutine
	go func() {
		if err := s.server.Serve(secureListener); err != nil {
			s.logger.Error("Error running DESS server:", err)
		}
	}()

	s.logger.Info("DESS atProtocol secondary server started successfully.")
	return nil
}

// Stop gracefully shuts down the DESS atProtocol secondary server
func (s *AtSecondaryServer) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.logger.Info("Stopping DESS atProtocol secondary server...")

	// Cancel the context to signal goroutines to terminate
	s.cancel()

	// Attempt to gracefully stop the server
	if err := s.server.Stop(); err != nil {
		s.logger.Error("Error stopping DESS atProtocol secondary server:", err)
		return err
	}

	s.logger.Info("DESS atProtocol secondary server stopped successfully.")
	return nil
}

// validateConfig ensures that the configuration parameters are correct and complete
func (s *AtSecondaryServer) validateConfig() error {
	if s.config.AtSign == "" || s.config.RootDomain == "" || s.config.ServerPort == 0 {
		return fmt.Errorf("missing required configuration parameters")
	}
	return nil
}
