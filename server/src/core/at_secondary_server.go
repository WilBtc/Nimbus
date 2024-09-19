// server/src/core/at_secondary_server.go

package core

import (
	"context"
	"fmt"
	"net"
	"sync"

	"github.com/atsign-foundation/at_server/server"
	"server/utils"
)

// AtSecondaryServer struct manages the DESS server instance
type AtSecondaryServer struct {
	config          *Config
	logger          *utils.Logger
	server          *server.AtServer
	mu              sync.Mutex
	ctx             context.Context
	cancel          context.CancelFunc
	securityGateway *SecurityGateway
}

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

	// Set up secure listener using SecurityGateway
	address := fmt.Sprintf("%s:%d", s.config.ServerHost, s.config.ServerPort)
	secureListener, err := s.securityGateway.SecureListener(address)
	if err != nil {
		s.logger.Error("Failed to create secure listener:", err)
		return err
	}

	// Start the server with the secure listener
	go func() {
		if err := s.server.Serve(secureListener); err != nil {
			s.logger.Error("Error running DESS server:", err)
		}
	}()

	s.logger.Info("DESS atProtocol secondary server started successfully.")
	return nil
}

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

func (s *AtSecondaryServer) validateConfig() error {
	if s.config.AtSign == "" || s.config.RootDomain == "" || s.config.ServerPort == 0 {
		return fmt.Errorf("missing required configuration parameters")
	}
	return nil
}
