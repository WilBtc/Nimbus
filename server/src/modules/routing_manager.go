// src/modules/routing_manager.go

package modules

import (
	"fmt"
	"log"
	"net"
	"sort"
	"sync"
	"time"

	"github.com/atsign-foundation/dess/server" // Importing DESS server package for integration
)

// RoutingManager manages data routing between devices and external systems, with support for prioritized routing.
type RoutingManager struct {
	routes      map[string]net.Conn    // Maps routes to device connections
	priorityMap map[string]int         // Priority map for routing critical data first
	routeMutex  sync.RWMutex           // Read-write mutex for managing routes and priorities
	logger      *log.Logger            // Logger for tracking routing events
	dessServer  *server.AtServer       // Reference to the DESS server for extended security and routing management
}

// NewRoutingManager creates a new RoutingManager instance with DESS integration.
func NewRoutingManager(logger *log.Logger, dessServer *server.AtServer) *RoutingManager {
	return &RoutingManager{
		routes:      make(map[string]net.Conn),
		priorityMap: make(map[string]int),
		logger:      logger,
		dessServer:  dessServer, // Integrating DESS server reference
	}
}

// AddRoute adds a new route to the routing manager with a specified priority.
func (rm *RoutingManager) AddRoute(deviceID string, conn net.Conn, priority int) {
	rm.routeMutex.Lock()
	defer rm.routeMutex.Unlock()

	// Ensure the device is authenticated with DESS before adding the route
	if !rm.dessServer.IsAuthenticated(deviceID) {
		rm.logger.Printf("Attempt to add route for unauthenticated device %s denied.\n", deviceID)
		return
	}

	rm.routes[deviceID] = conn
	rm.priorityMap[deviceID] = priority
	rm.logger.Printf("Route added for device %s with priority %d\n", deviceID, priority)
}

// RemoveRoute removes a route from the routing manager.
func (rm *RoutingManager) RemoveRoute(deviceID string) {
	rm.routeMutex.Lock()
	defer rm.routeMutex.Unlock()
	delete(rm.routes, deviceID)
	delete(rm.priorityMap, deviceID)
	rm.logger.Printf("Route removed for device %s\n", deviceID)
}

// RouteData routes data to the specified device based on priority.
func (rm *RoutingManager) RouteData(deviceID string, data []byte) {
	rm.routeMutex.RLock()
	defer rm.routeMutex.RUnlock()

	conn, exists := rm.routes[deviceID]
	if !exists {
		rm.logger.Printf("No route found for device %s\n", deviceID)
		return
	}

	// Log and route data based on priority
	priority := rm.priorityMap[deviceID]
	rm.logger.Printf("Routing data for device %s with priority %d\n", deviceID, priority)

	// Attempt to send data to the corresponding connection
	_, err := conn.Write(data)
	if err != nil {
		rm.logger.Printf("Error routing data to device %s: %v\n", deviceID, err)
	}
}

// RouteDataByPriority routes data to devices based on their priority, ensuring higher priority data is sent first.
func (rm *RoutingManager) RouteDataByPriority(data map[string][]byte) {
	rm.routeMutex.RLock()
	defer rm.routeMutex.RUnlock()

	// Sort device IDs by their priority
	sortedDevices := rm.sortDevicesByPriority()

	for _, deviceID := range sortedDevices {
		conn, exists := rm.routes[deviceID]
		if !exists {
			rm.logger.Printf("No route found for device %s\n", deviceID)
			continue
		}

		priority := rm.priorityMap[deviceID]
		rm.logger.Printf("Routing data for device %s with priority %d\n", deviceID, priority)

		// Send data to the corresponding connection
		_, err := conn.Write(data[deviceID])
		if err != nil {
			rm.logger.Printf("Error routing data to device %s: %v\n", deviceID, err)
		}
	}
}

// sortDevicesByPriority returns a slice of device IDs sorted by priority.
func (rm *RoutingManager) sortDevicesByPriority() []string {
	devices := make([]string, 0, len(rm.priorityMap))
	for deviceID := range rm.priorityMap {
		devices = append(devices, deviceID)
	}

	sort.Slice(devices, func(i, j int) bool {
		// Sort devices in descending order of priority
		return rm.priorityMap[devices[i]] > rm.priorityMap[devices[j]]
	})

	return devices
}

// MonitorRoutes periodically checks the health of all routes and ensures connections remain active.
func (rm *RoutingManager) MonitorRoutes(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			rm.checkRoutesHealth()
		}
	}()
}

// checkRoutesHealth inspects the health of each route, ensuring connections are still valid.
func (rm *RoutingManager) checkRoutesHealth() {
	rm.routeMutex.RLock()
	defer rm.routeMutex.RUnlock()

	for deviceID, conn := range rm.routes {
		if conn == nil || rm.isConnectionClosed(conn) {
			rm.logger.Printf("Route for device %s is inactive, removing route.\n", deviceID)
			rm.RemoveRoute(deviceID)
		} else {
			rm.logger.Printf("Route for device %s is active.\n", deviceID)
		}
	}
}

// isConnectionClosed checks if a given connection is closed.
func (rm *RoutingManager) isConnectionClosed(conn net.Conn) bool {
	// Placeholder for connection health check logic
	// Replace with actual implementation as needed
	_, err := conn.Write([]byte{})
	return err != nil
}
