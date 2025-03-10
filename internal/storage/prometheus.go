package prometheus

import (
	"context"
	"fmt"
	"time"

	"github.com/Luke-Skelton/daffodil-projct/internal/storage" // Import the storage interface

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

// PrometheusStorage implements the Storage interface for Prometheus.
type PrometheusStorage struct {
	client v1.API
	addr   string // Prometheus server address
}

// NewPrometheusStorage creates a new Prometheus storage backend.
func NewPrometheusStorage(addr string) (storage.Storage, error) {
	client, err := api.NewClient(api.Config{
		Address: addr,
	})
	if err != nil {
		return nil, fmt.Errorf("error creating Prometheus client: %w", err)
	}

	v1api := v1.NewAPI(client)
	return &PrometheusStorage{client: v1api, addr: addr}, nil
}

// WriteDataPoints writes data points to Prometheus.
func (p *PrometheusStorage) WriteDataPoints(dataPoints []storage.DataPoint) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // Add a timeout
	defer cancel()

	for _, dp := range dataPoints {
		// Convert your DataPoint to a Prometheus sample.
		sample := model.Sample{
			Metric:    model.Metric{}, // Start with an empty metric
			Timestamp: model.Time(dp.Timestamp),
			Value:     model.SampleValue(dp.Value),
		}

		// Add the metric name as a label.  This is crucial!
		sample.Metric[model.MetricNameLabel] = model.LabelValue(dp.Name)

		// Add tags as labels.
		for k, v := range dp.Tags {
			sample.Metric[model.LabelName(k)] = model.LabelValue(v)
		}
		var samples model.Samples
		samples = append(samples, &sample)
		err := p.client.InsertSamples(ctx, samples, v1.InsertOptions{})

		if err != nil {
			return fmt.Errorf("error writing data point to Prometheus: %w", err)
		}
	}

	return nil
}

// Close closes the Prometheus storage backend (no-op in this simple example).
func (p *PrometheusStorage) Close() error {
	//  In a more complex scenario, you might need to close connections here.
	return nil
}

// Name returns the name of the storage backend.
func (p *PrometheusStorage) Name() string {
	return "prometheus"
}
