package main

import (
	"flag"
	"log"
	"time"

	// --- Your Internal Packages ---
	// Make sure this module path matches the one in your go.mod file!
	"github.com/Luke-Skelton/daffodil_projct/internal/metrics"
	"github.com/Luke-Skelton/daffodil_projct/internal/storage"
	"github.com/Luke-Skelton/daffodil_projct/internal/storage/prometheus"
	// If you create other storage backends, import them here.
)

// Config holds the agent's configuration.
type Config struct {
	StorageType       string        // e.g., "prometheus", "mock"
	PrometheusAddress string        // Address for the Prometheus backend
	CollectInterval   time.Duration // How often to collect metrics
	// Add other configuration options here as needed (e.g., log file paths)
}

// loadConfig loads the configuration using command-line flags.
// You could expand this later to read from environment variables or a config file.
func loadConfig() Config {
	// Define flags
	storageType := flag.String("storage", "prometheus", "Storage backend type (e.g., 'prometheus', 'mock')")
	prometheusAddress := flag.String("prometheus-address", "http://localhost:9090", "Prometheus server address (if using Prometheus storage)")
	collectInterval := flag.Duration("interval", 15*time.Second, "Metric collection interval (e.g., '10s', '1m')")

	// Parse the flags from command-line arguments
	flag.Parse()

	// Return the configuration struct populated from the flags
	return Config{
		StorageType:       *storageType,
		PrometheusAddress: *prometheusAddress,
		CollectInterval:   *collectInterval,
	}
}

// runAgent contains the main loop for collecting and writing metrics.
// It accepts the storage implementation and collection interval as arguments.
func runAgent(dataStorage storage.Storage, collectInterval time.Duration) {
	// Create a ticker to control the collection interval
	ticker := time.NewTicker(collectInterval)
	defer ticker.Stop() // Ensure the ticker is stopped when the function exits

	log.Printf("Agent started. Collecting metrics every %v", collectInterval)

	// Main agent loop
	for range ticker.C { // This loop waits for the ticker to fire
		log.Println("Collecting metrics...")
		dataPoints, err := metrics.CollectMetrics()
		if err != nil {
			// Log the error but continue running the agent
			log.Printf("ERROR collecting metrics: %v", err)
			continue // Skip this iteration on collection error
		}

		if len(dataPoints) > 0 {
			log.Printf("Attempting to write %d data points to %s storage", len(dataPoints), dataStorage.Name())
			err = dataStorage.WriteDataPoints(dataPoints)
			if err != nil {
				// Log the error but continue running the agent
				// Implement more robust retry logic here if needed (e.g., exponential backoff)
				log.Printf("ERROR writing data to %s: %v", dataStorage.Name(), err)
				// Consider adding metrics *about* write errors here!
			} else {
				log.Printf("Successfully wrote %d data points to %s storage", len(dataPoints), dataStorage.Name())
			}
		} else {
			log.Println("No metrics collected in this interval.")
		}
	}
}

func main() {
	// 1. Load configuration
	config := loadConfig()

	// 2. Declare storage variables
	var dataStorage storage.Storage
	var err error // To capture errors during storage creation

	// 3. Select and create the storage backend based on configuration
	log.Printf("Selected storage type: %s", config.StorageType)
	switch config.StorageType {
	case "prometheus":
		log.Printf("Initializing Prometheus storage backend (Address: %s)", config.PrometheusAddress)
		dataStorage, err = prometheus.NewPrometheusStorage(config.PrometheusAddress)
		if err != nil {
			log.Fatalf("FATAL: Error creating Prometheus storage backend: %v", err)
		}
	case "mock":
		log.Println("Initializing Mock storage backend")
		// Assuming NewMockStorage is in the storage package and returns *MockStorage which implements Storage
		dataStorage = storage.NewMockStorage()
		// The mock storage doesn't typically have creation errors in this simple form
		// If NewMockStorage could fail, you'd handle 'err' here too.
	// --- Add cases for other storage backends here ---
	// case "influxdb":
	//  dataStorage, err = influxdb.NewInfluxDBStorage(...)
	//  if err != nil {
	//      log.Fatalf("FATAL: Error creating InfluxDB storage backend: %v", err)
	//  }
	default:
		log.Fatalf("FATAL: Invalid storage type specified in configuration: '%s'", config.StorageType)
	}

	// 4. Ensure storage is closed gracefully on exit
	// The defer statement schedules the Close() call to run when the main function returns.
	defer func() {
		log.Println("Shutting down agent, closing storage...")
		if err := dataStorage.Close(); err != nil {
			log.Printf("ERROR closing storage backend %s: %v", dataStorage.Name(), err)
		} else {
			log.Printf("Storage backend %s closed successfully.", dataStorage.Name())
		}
	}()

	// 5. Start the main agent loop
	runAgent(dataStorage, config.CollectInterval)
}
