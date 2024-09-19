// server/src/modules/analytics_engine.go

package modules

import (
	"fmt"
	"log"
	"math"
	"sync"
	"time"

	"github.com/atsign-foundation/at_server/server"
)

// AnalyticsEngine handles real-time data processing, analytics, and anomaly detection.
type AnalyticsEngine struct {
	dataMutex        sync.Mutex
	dataStore        []float64          // Stores incoming data for processing
	anomalyThreshold float64            // Threshold for Z-score anomaly detection
	logger           *log.Logger        // Logger for tracking analytics events
	alertHandler     func(string)       // Handler function for sending alerts
	dessServer       *server.AtServer   // Reference to the DESS server for extended security and logging
}

func NewAnalyticsEngine(anomalyThreshold float64, logger *log.Logger, alertHandler func(string), dessServer *server.AtServer) *AnalyticsEngine {
	return &AnalyticsEngine{
		dataStore:        make([]float64, 0),
		anomalyThreshold: anomalyThreshold,
		logger:           logger,
		alertHandler:     alertHandler,
		dessServer:       dessServer,
	}
}

func (ae *AnalyticsEngine) AddData(data float64) {
	ae.dataMutex.Lock()
	defer ae.dataMutex.Unlock()
	ae.dataStore = append(ae.dataStore, data)
	ae.logger.Printf("Data added: %v\n", data)
}

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

func (ae *AnalyticsEngine) calculateMean() float64 {
	sum := 0.0
	for _, value := range ae.dataStore {
		sum += value
	}
	return sum / float64(len(ae.dataStore))
}

func (ae *AnalyticsEngine) calculateStdDev(mean float64) float64 {
	varianceSum := 0.0
	for _, value := range ae.dataStore {
		varianceSum += math.Pow(value-mean, 2)
	}
	return math.Sqrt(varianceSum / float64(len(ae.dataStore)))
}

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

func (ae *AnalyticsEngine) handleAnomaly(message string) {
	if ae.alertHandler != nil {
		ae.alertHandler(message)
	}

	// Log the anomaly detection to DESS if integrated
	if ae.dessServer != nil {
		// Hypothetical method to log events in DESS
		ae.dessServer.LogEvent(message)
	}
}

func (ae *AnalyticsEngine) StartProcessing(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			ae.ProcessData()
		}
	}()
}

func (ae *AnalyticsEngine) Stop() {
	ae.logger.Println("Stopping Analytics Engine...")
	// Implement any cleanup if necessary
}
