// src/modules/access_control.go

package modules

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/atsign-foundation/dess/server" // Importing DESS server package
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
	permissions    map[string]Role      // Stores access permissions for devices with associated roles
	permMutex      sync.RWMutex         // Read-write mutex for handling concurrent access to permissions
	logger         *log.Logger          // Logger for recording access control activities
	auditLog       []AccessEvent        // Stores audit logs for access events
	auditMutex     sync.Mutex           // Mutex for concurrent access to audit logs
	dessServer     *server.AtServer     // Reference to the DESS server for authentication and access control
}

// AccessEvent represents an access attempt or change
type AccessEvent struct {
	DeviceID  string
	Timestamp time.Time
	Action    string
	Role      Role
}

// NewAccessControl creates a new AccessControl instance integrated with DESS.
func NewAccessControl(logger *log.Logger, dessServer *server.AtServer) *AccessControl {
	return &AccessControl{
		permissions: make(map[string]Role),
		logger:      logger,
		auditLog:    make([]AccessEvent, 0),
		dessServer:  dessServer, // Integrate DESS server instance
	}
}

// GrantAccess grants access to a specific device ID with a specified role.
func (ac *AccessControl) GrantAccess(deviceID string, role Role) {
	ac.permMutex.Lock()
	defer ac.permMutex.Unlock()
	ac.permissions[deviceID] = role
	ac.logger.Printf("Access granted to device %s with role %s\n", deviceID, role)
	ac.logAccessEvent(deviceID, "granted access", role)
}

// RevokeAccess revokes access for a specific device ID.
func (ac *AccessControl) RevokeAccess(deviceID string) {
	ac.permMutex.Lock()
	defer ac.permMutex.Unlock()
	role := ac.permissions[deviceID]
	delete(ac.permissions, deviceID)
	ac.logger.Printf("Access revoked for device %s\n", deviceID)
	ac.logAccessEvent(deviceID, "revoked access", role)
}

// CheckAccess checks if a device has permission to access the server based on role.
// Validates against DESS server's authentication.
func (ac *AccessControl) CheckAccess(deviceID string, requiredRole Role) bool {
	ac.permMutex.RLock()
	defer ac.permMutex.RUnlock()

	// Check if the device is authenticated with the DESS server
	if !ac.dessServer.IsAuthenticated(deviceID) {
		ac.logger.Printf("Access denied for unauthenticated device %s\n", deviceID)
		ac.logAccessEvent(deviceID, "unauthenticated access denied", NoAccess)
		return false
	}

	// Check permissions against assigned role
	role, exists := ac.permissions[deviceID]
	if !exists || role == NoAccess {
		ac.logger.Printf("Access denied for device %s\n", deviceID)
		ac.logAccessEvent(deviceID, "access denied", NoAccess)
		return false
	}

	if !ac.roleSufficient(role, requiredRole) {
		ac.logger.Printf("Insufficient permissions for device %s. Required role: %s, current role: %s\n", deviceID, requiredRole, role)
		ac.logAccessEvent(deviceID, "insufficient permissions", role)
		return false
	}

	ac.logger.Printf("Access granted to device %s with role %s\n", deviceID, role)
	ac.logAccessEvent(deviceID, "access granted", role)
	return true
}

// LogUnauthorizedAccess logs an attempt to access the server without permission.
func (ac *AccessControl) LogUnauthorizedAccess(deviceID string) {
	ac.logger.Printf("Unauthorized access attempt by device %s\n", deviceID)
	fmt.Printf("Unauthorized access attempt detected for device: %s\n", deviceID)
	ac.logAccessEvent(deviceID, "unauthorized access", NoAccess)
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

// logAccessEvent logs the access event to the audit log
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
	ac.logger.Printf("Audit log: %v", event)
}

// GetAuditLog returns a copy of the audit log for review
func (ac *AccessControl) GetAuditLog() []AccessEvent {
	ac.auditMutex.Lock()
	defer ac.auditMutex.Unlock()
	return append([]AccessEvent(nil), ac.auditLog...) // Returns a copy of the audit log
}
