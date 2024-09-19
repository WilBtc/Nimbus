package core

import (
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

// Config struct holds the configuration parameters for the Nimbus Edge Server
type Config struct {
	AtSign           string // atSign used for DESS server communication
	RootDomain       string // Fully Qualified Domain Name for the secondary server
	Namespace        string
	ServerHost       string
	ServerPort       int
	StoragePath      string
	CommitLogPath    string
	LogLevel         string
	AuthRequired     bool
	EncryptionConfig string // Placeholder for specific encryption settings, adjust based on DESS requirements
	SecurityLevel    int
	EdgeAnalytics    bool
	Secret           string // Secret used for cram authentication to the secondary
	Email            string // Email address for SSL certificate management
}

// LoadConfig loads configuration from environment variables with defaults
func LoadConfig() (*Config, error) {
	config := &Config{
		AtSign:           getEnv("ATSIGN", ""),
		RootDomain:       getEnv("ROOT_DOMAIN", "root.atsign.org"),
		Namespace:        getEnv("NAMESPACE", "nimbus"),
		ServerHost:       getEnv("SERVER_HOST", "0.0.0.0"),
		ServerPort:       getEnvAsInt("SERVER_PORT", 6464),
		StoragePath:      getEnv("STORAGE_PATH", "./storage"),
		CommitLogPath:    getEnv("COMMIT_LOG_PATH", "./commitLog"),
		LogLevel:         getEnv("LOG_LEVEL", "INFO"),
		AuthRequired:     getEnvAsBool("AUTH_REQUIRED", true),
		EncryptionConfig: getEnv("ENCRYPTION_CONFIG", "default_encryption"),
		SecurityLevel:    getEnvAsInt("SECURITY_LEVEL", 2),
		EdgeAnalytics:    getEnvAsBool("EDGE_ANALYTICS", true),
		Secret:           getEnvSecure("CRAM_SECRET", ""), // Securely handle secrets
		Email:            getEnv("EMAIL", ""),             // Email for SSL certificate requests
	}

	// Perform validation on loaded configuration
	if err := validateConfig(config); err != nil {
		return nil, err
	}

	return config, nil
}

// validateConfig performs validation on the loaded configuration
func validateConfig(config *Config) error {
	// Ensure all required DESS-specific configurations are set
	if config.AtSign == "" {
		return errors.New("ATSIGN is required but not set in the environment variables")
	}
	if config.Secret == "" {
		return errors.New("SECRET is required for authentication but not set in the environment variables")
	}
	if config.Email == "" {
		return errors.New("EMAIL is required for SSL management but not set in the environment variables")
	}

	// Validate email format
	if !isValidEmail(config.Email) {
		return fmt.Errorf("invalid email format: %s", config.Email)
	}

	// Check if essential paths exist or are writable
	if err := ensurePathExists(config.StoragePath); err != nil {
		return fmt.Errorf("storage path error: %w", err)
	}

	if err := ensurePathExists(config.CommitLogPath); err != nil {
		return fmt.Errorf("commit log path error: %w", err)
	}

	// Validate server port range
	if config.ServerPort < 1024 || config.ServerPort > 65535 {
		return fmt.Errorf("invalid SERVER_PORT: %d. Must be between 1024 and 65535", config.ServerPort)
	}

	return nil
}

// ensurePathExists checks if a directory exists and creates it if it does not
func ensurePathExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("failed to create path %s: %w", path, err)
		}
	}
	return nil
}

// getEnv retrieves environment variables or returns a default value if not set
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// getEnvSecure retrieves sensitive environment variables securely
func getEnvSecure(key, defaultValue string) string {
	value := getEnv(key, defaultValue)
	if value == "" && defaultValue == "" {
		log.Printf("Warning: %s environment variable is empty or not set", key)
	}
	return value
}

// getEnvAsInt retrieves an environment variable as an integer, with a fallback to a default value
func getEnvAsInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultVal
}

// getEnvAsBool retrieves an environment variable as a boolean, with a fallback to a default value
func getEnvAsBool(name string, defaultVal bool) bool {
	valStr := getEnv(name, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}
	return defaultVal
}

// isValidEmail validates the email address format using regex
func isValidEmail(email string) bool {
	// Basic email validation regex pattern
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}
