package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// CurlmeMetrics struct with all prom metrics usable into the controller
type CurlmeMetrics struct {
	CmSyncedCount prometheus.Counter
}

// New init all prometheus curlme metrics
func New() *CurlmeMetrics {
	reg := prometheus.NewRegistry()
	cmSyncedCount := promauto.With(reg).NewCounter(prometheus.CounterOpts{
		Name: "curlme_configmap_synced_total",
		Help: "The total number of cm processed by curlme controller",
	})

	return &CurlmeMetrics{
		CmSyncedCount: cmSyncedCount,
	}

}
