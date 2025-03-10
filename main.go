package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/Luke-Skelton/daffodil-projct/internal/metrics" // Import your metrics package
	"github.com/Luke-Skelton/daffodil-projct/internal/storage"
	"github.com/Luke-Skelton/daffodil-projct/internal/storage/prometheus" // Prometheus implementation
	// ... other storage implementations (influxdb, etc.)
)

// Config holds the agent's configuration.
type Config struct {
	StorageType       string
	PrometheusAddress string
	CollectInterval   time.Duration
}

// loadConfig loads the configuration (from flags, environment variables, or a file).
func loadConfig() Config {
	// Use flags for configuration. This is a good, flexible approach.
	storageType := flag.String("storage", "prometheus", "Storage backend (prometheus, mock, etc.)")
	prometheusAddress := flag.String("prometheus-address", "http://localhost:9090", "Prometheus server address")
	collectInterval := flag.Duration("interval", 15*time.Second, "Metric collection interval")

	flag.Parse() // Parse the command-line flags

	return Config{
		StorageType:       *storageType,
		PrometheusAddress: *prometheusAddress,
		CollectInterval:   *collectInterval,
	}
}

// runAgent contains the main agent logic (collecting and writing metrics).
func runAgent(dataStorage storage.Storage, collectInterval time.Duration) {
	for {
		dataPoints, err := metrics.CollectMetrics()
		if err != nil {
			log.Printf("Error collecting metrics: %v", err)
			time.Sleep(10 * time.Second) // Wait before retrying
			continue
		}

		err = dataStorage.WriteDataPoints(dataPoints)
		if err != nil {
			log.Printf("Error writing data to %s: %v", dataStorage.Name(), err)
			// Implement retry logic, error handling (e.g., backoff and retry)
		}

		time.Sleep(collectInterval) // Wait for the next collection interval.
	}
}

func main() {
	config := loadConfig()

	var dataStorage storage.Storage
	var err error

	// Select the storage backend based on configuration.
	switch config.StorageType {
	case "prometheus":
		dataStorage, err = prometheus.NewPrometheusStorage(config.PrometheusAddress)
	// case "influxdb": // Add cases for other storage backends
	//  dataStorage, err = influxdb.NewInfluxDBStorage(...)
	case "mock": // Add a case for the mock storage
		dataStorage = storage.NewMockStorage()

	default:
		log.Fatalf("Invalid storage type: %s", config.StorageType)
	}

	if err != nil {
		log.Fatalf("Error creating storage backend: %v", err)
	}
	defer dataStorage.Close() // Ensure storage is closed on exit.

	fmt.Printf("Starting agent with storage: %s\n", dataStorage.Name()) //Log which storage is used
	runAgent(dataStorage, config.CollectInterval)                       // Pass the storage to runAgent.
}
