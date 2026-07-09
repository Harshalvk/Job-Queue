package jobqueue

import "github.com/prometheus/client_golang/prometheus"

var (
	JobsProcessed = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "jobqueue_jobs_processed_total",
			Help: "Total number of jobs processed, labeled by type and outcome.",
		},
		[]string{"type", "outcome"},
	)

	JobDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "jobqueue_job_duration_seconds",
			Help: "Duration of job handler exectuion in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"type"},
	)

	QueueDepth = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "jobqueue_pending_depth",
			Help: "Current number of jobs waiting in the pending queue.",
		},
	)
)

func init() {
	prometheus.MustRegister(JobsProcessed, JobDuration, QueueDepth)
}