// server/src/modules/access_control.go
package modules

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/atsign-foundation/at_server/server" // Assuming this is the correct import path for DESS server package
)

// Role represents different levels of access permissions
type Role string

const (
	AdminRole   Role = "admin"
	UserRole    Role = "user"
	GuestRole   Role = "guest"
	NoAccess    Role = "no_access"
)

// AccessControl manages permissions, roles, and secure access to the server.
type AccessControl struct {
	permissions    sync.Map         // Stores access permissions for devices with associated roles
	logger         *log.Logger      // Logger for recording access control activities
	auditLog       []AccessEvent    // Stores audit logs for access events
	auditMutex     sync.Mutex       // Mutex for concurrent access to audit logs
	dessServer     *server.AtServer // Reference to the DESS server for authentication and access control
}

// AccessEvent represents an access attempt or change
type AccessEvent struct {
	DeviceID  string
	Timestamp time.Time
	Action    string
	Role      Role
}

// NewAccessControl initializes a new AccessControl instance with logging and DESS server integration
func NewAccessControl(logger *log.Logger, dessServer *server.AtServer) *AccessControl {
	return &AccessControl{
		logger:     logger,
		auditLog:   make([]AccessEvent, 0),
		dessServer: dessServer,
	}
}

// GrantAccess grants access to a specific device ID with a specified role
func (ac *AccessControl) GrantAccess(deviceID string, role Role) {
	ac.permissions.Store(deviceID, role)
	ac.logger.Printf("Access granted to device %s with role %s\n", deviceID, role)
	ac.logAccessEvent(deviceID, "granted access", role)
}

// RevokeAccess revokes access for a specific device ID
func (ac *AccessControl) RevokeAccess(deviceID string) {
	ac.permissions.Delete(deviceID)
	ac.logger.Printf("Access revoked for device %s\n", deviceID)
	ac.logAccessEvent(deviceID, "revoked access", NoAccess)
}

// CheckAccess checks if a device has permission to access the server based on the required role
func (ac *AccessControl) CheckAccess(deviceID string, requiredRole Role) bool {
	// Check if the device is authenticated with the DESS server
	if !ac.dessServer.IsAuthenticated(deviceID) {
		ac.logger.Printf("Access denied for unauthenticated device %s\n", deviceID)
		ac.logAccessEvent(deviceID, "unauthenticated access denied", NoAccess)
		return false
	}

	// Retrieve the role from the permissions map
	value, ok := ac.permissions.Load(deviceID)
	if !ok {
		ac.logger.Printf("Access denied for device %s: No role assigned\n", deviceID)
		ac.logAccessEvent(deviceID, "access denied", NoAccess)
		return false
	}
	role := value.(Role)

	if !ac.roleSufficient(role, requiredRole) {
		ac.logger.Printf("Insufficient permissions for device %s. Required role: %s, current role: %s\n", deviceID, requiredRole, role)
		ac.logAccessEvent(deviceID, "insufficient permissions", role)
		return false
	}

	ac.logger.Printf("Access granted to device %s with role %s\n", deviceID, role)
	ac.logAccessEvent(deviceID, "access granted", role)
	return true
}

// roleSufficient checks if the current role has sufficient permissions compared to the required role
func (ac *AccessControl) roleSufficient(currentRole, requiredRole Role) bool {
	roles := map[Role]int{
		NoAccess:  0,
		GuestRole: 1,
		UserRole:  2,
		AdminRole: 3,
	}
	return roles[currentRole] >= roles[requiredRole]
}

// logAccessEvent records access events to the audit log
func (ac *AccessControl) logAccessEvent(deviceID, action string, role Role) {
	ac.auditMutex.Lock()
	defer ac.auditMutex.Unlock()

	event := AccessEvent{
		DeviceID:  deviceID,
		Timestamp: time.Now(),
		Action:    action,
		Role:      role,
	}

	ac.auditLog = append(ac.auditLog, event)
	ac.logger.Printf("Audit log recorded: %v", event)
}

// GetAuditLog returns a copy of the audit log for review
func (ac *AccessControl) GetAuditLog() []AccessEvent {
	ac.auditMutex.Lock()
	defer ac.auditMutex.Unlock()
	return append([]AccessEvent(nil), ac.auditLog...) // Returns a copy of the audit log
}
