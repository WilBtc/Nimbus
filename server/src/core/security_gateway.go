// server/src/core/security_gateway.go

package core

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"time"

	"server/utils"
)

// SecurityGateway manages security features such as encryption, traffic inspection, and anomaly detection
type SecurityGateway struct {
	config           *Config
	logger           *utils.Logger
	tlsConfig        *tls.Config
	idsEnabled       bool // Flag to enable or disable IDS/IPS
	ctx              context.Context
	cancel           context.CancelFunc
}

// NewSecurityGateway initializes a new SecurityGateway with the given configuration
func NewSecurityGateway(config *Config) *SecurityGateway {
	ctx, cancel := context.WithCancel(context.Background())
	return &SecurityGateway{
		config: config,
		logger: utils.NewLogger(config.LogLevel, config.LogFilePath, config.RotateLogs, config.RotationInterval),
		ctx:    ctx,
		cancel: cancel,
	}
}

// Initialize sets up the Security Gateway with necessary security measures
func (s *SecurityGateway) Initialize() error {
	s.logger.Info("Initializing Security Gateway...")

	// Validate security configuration
	if err := s.validateConfig(); err != nil {
		s.logger.Error("Security configuration validation failed:", err)
		return err
	}

	// Initialize TLS configuration for secure communication
	if err := s.initTLSConfig(); err != nil {
		s.logger.Error("Failed to initialize TLS configuration:", err)
		return err
	}
	s.logger.Info("TLS configuration initialized successfully.")

	// Enable advanced security features if security level requires it
	if s.config.SecurityLevel > 1 {
		s.logger.Info("Advanced security enabled: Traffic inspection and anomaly detection activated.")
		if err := s.setupIntrusionDetection(); err != nil {
			s.logger.Error("Error setting up IDS/IPS:", err)
			return err
		}
		s.idsEnabled = true
	}

	s.logger.Info("Security Gateway initialized successfully.")
	return nil
}

// validateConfig checks that the security configuration is valid
func (s *SecurityGateway) validateConfig() error {
	if s.config.SecurityLevel < 1 {
		return fmt.Errorf("invalid security level, must be >= 1")
	}
	if s.config.SSLCertPath == "" || s.config.SSLKeyPath == "" {
		return fmt.Errorf("SSL certificate and key paths are required but not set")
	}
	return nil
}

// initTLSConfig initializes the TLS configuration for secure communication
func (s *SecurityGateway) initTLSConfig() error {
	cert, err := tls.LoadX509KeyPair(s.config.SSLCertPath, s.config.SSLKeyPath)
	if err != nil {
		return fmt.Errorf("failed to load SSL certificates: %w", err)
	}

	s.tlsConfig = &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12, // Ensuring strong TLS encryption
	}

	return nil
}

// setupIntrusionDetection sets up the intrusion detection and prevention system (IDS/IPS)
func (s *SecurityGateway) setupIntrusionDetection() error {
	s.logger.Info("Setting up Intrusion Detection and Prevention System (IDS/IPS)...")
	// Placeholder for integrating with actual IDS solutions like Suricata or Snort
	return nil
}

// MonitorTraffic inspects traffic to detect anomalies or malicious activities
func (s *SecurityGateway) MonitorTraffic() {
	s.logger.Info("Starting traffic monitoring...")

	// Monitoring loop with graceful shutdown using context
	ticker := time.NewTicker(10 * time.Second) // Adjust monitoring interval as needed
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			s.logger.Info("Traffic monitoring stopped.")
			return
		case <-ticker.C:
			if s.idsEnabled {
				anomalies := s.checkForAnomalies()
				if len(anomalies) > 0 {
					s.logger.Warn(fmt.Sprintf("Anomalies detected: %v", anomalies))
					s.takeActionOnAnomalies(anomalies)
				} else {
					s.logger.Info("No anomalies detected during traffic monitoring.")
				}
			}
		}
	}
}

// checkForAnomalies performs a mock check for traffic anomalies (placeholder)
func (s *SecurityGateway) checkForAnomalies() []string {
	// Placeholder function, replace with actual anomaly detection logic
	// Simulate detection of traffic anomalies for demonstration
	return []string{}
}

// takeActionOnAnomalies handles actions based on detected traffic anomalies
func (s *SecurityGateway) takeActionOnAnomalies(anomalies []string) {
	for _, anomaly := range anomalies {
		// Example action: Log, alert, or initiate protective measures
		s.logger.Error(fmt.Sprintf("Taking action on detected anomaly: %s", anomaly))
		// Additional actions, such as blocking IPs, can be implemented here
	}
}

// SecureListener returns a secure listener using TLS for encrypted communication
func (s *SecurityGateway) SecureListener(address string) (net.Listener, error) {
	listener, err := tls.Listen("tcp", address, s.tlsConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create TLS listener: %w", err)
	}
	return listener, nil
}

// Stop gracefully stops the traffic monitoring and shuts down security processes
func (s *SecurityGateway) Stop() {
	s.logger.Info("Stopping Security Gateway...")
	s.cancel() // Cancel the context to stop all monitoring activities
	s.logger.Info("Security Gateway stopped successfully.")
}
