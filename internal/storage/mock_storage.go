package storage

// MockStorage is a mock implementation of the Storage interface for testing.
type MockStorage struct {
	DataPoints []DataPoint // Store received data points here.
	Closed     bool        // Track whether Close() was called.
	NameValue  string
}

// NewMockStorage creates a new MockStorage instance.
func NewMockStorage() *MockStorage {
	return &MockStorage{}
}

// WriteDataPoints stores the data points in the MockStorage's DataPoints slice.
func (m *MockStorage) WriteDataPoints(dataPoints []DataPoint) error {
	m.DataPoints = append(m.DataPoints, dataPoints...) // Append to the slice.
	return nil                                         // Simulate successful write.
}

// Close sets the Closed flag to true.
func (m *MockStorage) Close() error {
	m.Closed = true
	return nil
}

func (m *MockStorage) Name() string {
	return m.NameValue
}
