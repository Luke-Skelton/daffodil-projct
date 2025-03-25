package storage

// DataPoint represents a single metric data point.
type DataPoint struct {
	Name      string
	Timestamp int64 // Unix timestamp in milliseconds (or seconds, choose a consistent unit)
	Value     float64
	Tags      map[string]string
}

// Storage is the interface that all data storage backends must implement.
type Storage interface {
	WriteDataPoints(dataPoints []DataPoint) error
	Close() error // Important for releasing resources (connections, etc.)
	Name() string // Useful for identifying the storage in logs/configuration
}
