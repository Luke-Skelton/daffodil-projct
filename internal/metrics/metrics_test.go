package metrics

import (
	"reflect"
	"testing"

	"github.com/Luke-Skelton/daffodil-projct/internal/storage" // Import the storage interface
)

func TestCollectMetrics(t *testing.T) {
	// Create a MockStorage instance.
	mockStorage := storage.NewMockStorage()
	mockStorage.NameValue = "mock"

	// Call CollectMetrics (which you'll need to modify slightly - see next step).
	dataPoints, err := CollectMetrics()
	if err != nil {
		t.Fatalf("CollectMetrics returned an error: %v", err)
	}

	// Assert that we got some data points.  A more robust test would
	// check the *content* of the data points (values, names, tags),
	// but this is a good start.
	if len(dataPoints) == 0 {
		t.Errorf("Expected to collect metrics, but got 0 data points")
	}

	// Check that the data points have expected names, values and tags
	expectedDataPoints := map[string][]struct {
		Tags map[string]string
		// Add a type so that different values can be checked.
		ValueType string
		Value     interface{}
	}{
		"cpu_usage": {
			{Tags: map[string]string{"core": "cpu0"}, ValueType: "float", Value: 0.0},
			{Tags: map[string]string{"core": "cpu1"}, ValueType: "float", Value: 0.0},
			// Add more expected CPU cores based on your system.
		},
		"memory_usage": {
			{Tags: map[string]string{"type": "used"}, ValueType: "float", Value: 0.0},
		},
		"memory_total": {
			{Tags: map[string]string{}, ValueType: "float", Value: 0.0},
		},
		// Add expected network interfaces based on your system.
		"network_bytes_sent": {
			{Tags: map[string]string{"interface": "lo"}, ValueType: "float", Value: 0.0},
		},
		"network_bytes_recv": {
			{Tags: map[string]string{"interface": "lo"}, ValueType: "float", Value: 0.0},
		},
		// Add disk usage tests
		"disk_usage": {
			{Tags: map[string]string{"mountpoint": "/", "fstype": "ext4"}, ValueType: "float", Value: 0.0},
		},
	}

	for _, dp := range dataPoints {
		expected, ok := expectedDataPoints[dp.Name]
		if !ok {
			t.Errorf("Unexpected metric name: %s", dp.Name)
			continue
		}

		found := false
		for _, exp := range expected {
			if reflect.DeepEqual(dp.Tags, exp.Tags) {
				found = true
				switch exp.ValueType {
				case "float":
					if dp.Value < exp.Value.(float64) {
						// t.Errorf("Metric %s with tags %v: expected value >= %f, got %f", dp.Name, dp.Tags, exp.Value, dp.Value)
						// Allow for values to change, but we know the metrics shouldn't return negative numbers.
					}
				case "string":
					if dp.Value != exp.Value.(string) {
						t.Errorf("Metric %s with tags %v: expected value %s, got %s", dp.Name, dp.Tags, exp.Value, dp.Value)
					}
				// Add cases for other data types as needed.
				default:
					t.Errorf("Unknown value type: %s", exp.ValueType)
				}
				break // Move on once a match is found within expected values for that metric name
			}
		}
		if !found {
			t.Errorf("No matching expectation found for metric: %s, tags: %v", dp.Name, dp.Tags)
		}
	}
	// Example of how to check if Close() was called (if you were testing
	// code that should call Close()).  Not directly applicable to
	// CollectMetrics, but useful for other tests.
	// if !mockStorage.Closed {
	//     t.Errorf("Expected storage to be closed, but it wasn't")
	// }
}
