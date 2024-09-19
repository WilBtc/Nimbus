// src/modules/analytics_engine.go

package modules

import (
	"fmt"
	"log"
	"math"
	"sync"
	"time"

	"github.com/atsign-foundation/dess/server" // Importing DESS server package for integration
)

// AnalyticsEngine handles real-time data processing, analytics, and anomaly detection.
type AnalyticsEngine struct {
	dataMutex        sync.Mutex
	dataStore        []float64          // Stores incoming data for processing
	anomalyThreshold float64            // Threshold for detecting anomalies
	logger           *log.Logger        // Logger for tracking analytics events
	alertHandler     func(string)       // Handler function for sending alerts
	dessServer       *server.AtServer   // Reference to the DESS server for extended security and logging
}

// NewAnalyticsEngine creates a new instance of AnalyticsEngine with DESS integration.
func NewAnalyticsEngine(anomalyThreshold float64, logger *log.Logger, alertHandler func(string), dessServer *server.AtServer) *AnalyticsEngine {
	return &AnalyticsEngine{
		dataStore:        make([]float64, 0),
		anomalyThreshold: anomalyThreshold,
		logger:           logger,
		alertHandler:     alertHandler,
		dessServer:       dessServer, // Integrating DESS server reference
	}
}

// AddData adds new data to the analytics engine for processing.
func (ae *AnalyticsEngine) AddData(data float64) {
	ae.dataMutex.Lock()
	defer ae.dataMutex.Unlock()
	ae.dataStore = append(ae.dataStore, data)
	ae.logger.Printf("Data added: %v\n", data)
}

// ProcessData performs real-time analysis on the collected data.
func (ae *AnalyticsEngine) ProcessData() {
	ae.dataMutex.Lock()
	defer ae.dataMutex.Unlock()

	if len(ae.dataStore) == 0 {
		ae.logger.Println("No data to process.")
		return
	}

	// Calculate mean and standard deviation
	mean := ae.calculateMean()
	stdDev := ae.calculateStdDev(mean)

	// Log and display results
	ae.logger.Printf("Real-time Analytics - Mean: %.2f, Std Dev: %.2f\n", mean, stdDev)

	// Detect anomalies in the processed data
	ae.detectAnomalies(mean, stdDev)

	// Clear data after processing to maintain efficiency
	ae.dataStore = ae.dataStore[:0]
}

// calculateMean computes the mean of the current data set.
func (ae *AnalyticsEngine) calculateMean() float64 {
	sum := 0.0
	for _, value := range ae.dataStore {
		sum += value
	}
	return sum / float64(len(ae.dataStore))
}

// calculateStdDev computes the standard deviation of the current data set.
func (ae *AnalyticsEngine) calculateStdDev(mean float64) float64 {
	varianceSum := 0.0
	for _, value := range ae.dataStore {
		varianceSum += math.Pow(value-mean, 2)
	}
	return math.Sqrt(varianceSum / float64(len(ae.dataStore)))
}

// detectAnomalies checks the data for anomalies based on deviation from the mean.
func (ae *AnalyticsEngine) detectAnomalies(mean, stdDev float64) {
	anomalies := []float64{}
	for _, value := range ae.dataStore {
		if math.Abs(value-mean) > ae.anomalyThreshold*stdDev {
			anomalies = append(anomalies, value)
			ae.logger.Printf("Anomaly detected: %v (Threshold: %.2f)\n", value, ae.anomalyThreshold*stdDev)
			ae.handleAnomaly(fmt.Sprintf("Anomaly detected: %v", value))
		}
	}
	if len(anomalies) == 0 {
		ae.logger.Println("No anomalies detected.")
	}
}

// handleAnomaly handles detected anomalies, integrating DESS alerting mechanisms.
func (ae *AnalyticsEngine) handleAnomaly(message string) {
	if ae.alertHandler != nil {
		ae.alertHandler(message)
	}

	// Log the anomaly detection to DESS if integrated for extended monitoring and alerting
	if ae.dessServer != nil {
		err := ae.dessServer.LogEvent(message) // Hypothetical DESS function for logging
		if err != nil {
			ae.logger.Printf("Failed to log anomaly to DESS: %v", err)
		}
	}
}

// StartProcessing starts periodic data processing.
func (ae *AnalyticsEngine) StartProcessing(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			ae.ProcessData()
		}
	}()
}

// SetAnomalyThreshold allows dynamic adjustment of the anomaly detection threshold.
func (ae *AnalyticsEngine) SetAnomalyThreshold(threshold float64) {
	ae.dataMutex.Lock()
	defer ae.dataMutex.Unlock()
	ae.anomalyThreshold = threshold
	ae.logger.Printf("Anomaly threshold updated to: %.2f", threshold)
}
