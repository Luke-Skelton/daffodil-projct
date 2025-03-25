module github.com/Luke-Skelton/daffodil-projct

go 1.21

require (
	github.com/prometheus/client_golang v1.21.1
	github.com/shirou/gopsutil/v3 v3.24.5
	go.uber.org/zap v1.27.0
)

require go.uber.org/multierr v1.10.0 // indirect

replace github.com/Luke-Skelton/daffodil-projct => ./