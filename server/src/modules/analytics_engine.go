// server/src/modules/analytics_engine.go
package modules

import (
	"fmt"
	"log"
	"math"
	"sync"
	"time"

	"github.com/atsign-foundation/at_server/server" // Assuming this is the correct import path for DESS server package
)

// AnalyticsEngine handles real-time data processing, analytics, and anomaly detection.
type AnalyticsEngine struct {
	dataMutex        sync.Mutex        // Mutex for thread-safe access to data
	dataStore        []float64         // Stores incoming data for processing
	anomalyThreshold float64           // Threshold for Z-score anomaly detection
	logger           *log.Logger       // Logger for tracking analytics events
	alertHandler     func(string)      // Handler function for sending alerts
	dessServer       *server.AtServer  // Reference to the DESS server for extended security and logging
}

// NewAnalyticsEngine initializes a new AnalyticsEngine with an anomaly threshold, logger, alert handler, and DESS server integration.
func NewAnalyticsEngine(anomalyThreshold float64, logger *log.Logger, alertHandler func(string), dessServer *server.AtServer) *AnalyticsEngine {
	return &AnalyticsEngine{
		dataStore:        make([]float64, 0),
		anomalyThreshold: anomalyThreshold,
		logger:           logger,
		alertHandler:     alertHandler,
		dessServer:       dessServer,
	}
}

// AddData appends new data to the internal store, ensuring thread-safety with a mutex.
func (ae *AnalyticsEngine) AddData(data float64) {
	ae.dataMutex.Lock()
	defer ae.dataMutex.Unlock()
	ae.dataStore = append(ae.dataStore, data)
	ae.logger.Printf("Data added: %.2f\n", data)
}

// ProcessData processes the stored data to calculate mean and standard deviation, and checks for anomalies.
func (ae *AnalyticsEngine) ProcessData() {
	ae.dataMutex.Lock()
	defer ae.dataMutex.Unlock()

	if len(ae.dataStore) == 0 {
		ae.logger.Println("No data to process.")
		return
	}

	mean := ae.calculateMean()
	stdDev := ae.calculateStdDev(mean)

	ae.logger.Printf("Real-time Analytics - Mean: %.2f, Std Dev: %.2f\n", mean, stdDev)

	// Detect anomalies using Z-score
	ae.detectAnomalies(mean, stdDev)

	// Clear data after processing
	ae.dataStore = ae.dataStore[:0]
}

// calculateMean calculates the average of the data points in the store.
func (ae *AnalyticsEngine) calculateMean() float64 {
	sum := 0.0
	for _, value := range ae.dataStore {
		sum += value
	}
	return sum / float64(len(ae.dataStore))
}

// calculateStdDev calculates the standard deviation of the data points based on the provided mean.
func (ae *AnalyticsEngine) calculateStdDev(mean float64) float64 {
	varianceSum := 0.0
	for _, value := range ae.dataStore {
		varianceSum += math.Pow(value-mean, 2)
	}
	return math.Sqrt(varianceSum / float64(len(ae.dataStore)))
}

// detectAnomalies identifies data points that exceed the anomaly threshold using the Z-score method.
func (ae *AnalyticsEngine) detectAnomalies(mean, stdDev float64) {
	for _, value := range ae.dataStore {
		zScore := math.Abs((value - mean) / stdDev)
		if zScore > ae.anomalyThreshold {
			message := fmt.Sprintf("Anomaly detected: Value %.2f with Z-score %.2f", value, zScore)
			ae.logger.Println(message)
			ae.handleAnomaly(message)
		}
	}
}

// handleAnomaly handles detected anomalies by logging them, sending alerts, and optionally logging events to DESS.
func (ae *AnalyticsEngine) handleAnomaly(message string) {
	// Send an alert if an alert handler is provided
	if ae.alertHandler != nil {
		ae.alertHandler(message)
	}

	// Log the anomaly detection to DESS if integrated
	if ae.dessServer != nil {
		// Hypothetical method to log events in DESS
		ae.dessServer.LogEvent(message)
	}
}

// StartProcessing starts a ticker that processes data at regular intervals.
func (ae *AnalyticsEngine) StartProcessing(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			ae.ProcessData()
		}
	}()
}

// Stop stops the analytics engine by logging the stop event. Cleanup actions can be added if necessary.
func (ae *AnalyticsEngine) Stop() {
	ae.logger.Println("Stopping Analytics Engine...")
	// Implement any cleanup if necessary
}
