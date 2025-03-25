package logs

import (
    "bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"your_agent/internal/storage"
)
var (
	Logger *zap.Logger
)
func init() {
	// Configure Zap logger.  This is a basic configuration; you can customize it.
	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	var err error //Declare err to be used later.
	Logger, err = config.Build()
	if err != nil {
		panic(err) // Fail fast if logger can't be initialized.
	}
}
func ProcessLogFile(filePath string, dataStorage storage.Storage) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening log file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Parse the log line (this is VERY basic; you'll need a more robust parser).
		fields := strings.Fields(line) // Split by spaces (very simplistic parsing)
		if len(fields) < 2 {
			continue // Skip lines that don't have at least two fields.
		}
		// Simple example. You will change the parsing depending on your log format.
		timestamp, err := time.Parse(time.RFC3339, fields[0]) // Parse the timestamp
		if err != nil{
			continue
		}
		message := strings.Join(fields[1:], " ")               // The rest is the message


		dataPoints := []storage.DataPoint{
			{
			Name: "log_entry",
			Timestamp: timestamp.UnixMilli(),
			Value: 1, // Use value 1 to just "count" log lines
			Tags: map[string]string{
				"level":   "info",  // Add level (you'd parse this from the log line)
				"message": message, // Store the entire message as a tag
			},

			},
		}
		if err := dataStorage.WriteDataPoints(dataPoints); err != nil{
			Logger.Error("Failed to write datapoints", zap.Error(err)) //Log error.
			return fmt.Errorf("error writing to data storage : %w", err)
		}

		// Log using Zap (structured logging).
		Logger.Info(message,
			zap.String("level", "info"),  // Example fields
			zap.Time("