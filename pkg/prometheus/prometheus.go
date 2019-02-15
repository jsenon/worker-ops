//Package prometheus build custom metrics based on number of worker node
package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	WorkerNumber = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "worker_total",
			Help: "Number of Workers.",
		},
		[]string{"env", "region"},
	)
)
