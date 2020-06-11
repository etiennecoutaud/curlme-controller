package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type CurlmeMetrics struct {
	CmSyncedCount prometheus.Counter
}

func New() *CurlmeMetrics {
	cmSyncedCount := promauto.NewCounter(prometheus.CounterOpts{
		Name: "curlme_configmap_synced_total",
		Help: "The total number of cm processed by curlme controller",
	})

	return &CurlmeMetrics{
		CmSyncedCount: cmSyncedCount,
	}

}

