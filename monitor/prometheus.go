package monitor

// 实现一个 Prometheus 监控端点
import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"sync"

	"github.com/piwriw/go-chrono/pkg/url"

	"github.com/piwriw/go-chrono/jobs"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// 用于跟踪正在运行的job
	runningJobs sync.Map // map[string]struct{}
	registry    = prometheus.NewRegistry()
)

// GetRegistry 获取prometheus的注册器
// GetRegistry returns the prometheus registry.
func GetRegistry() *prometheus.Registry {
	return registry
}

func StartPrometheusEndpoint(addr string) {
	if err := url.ValidateURLAddr(addr); err != nil {
		panic(err)
	}
	http.Handle("/metrics", promhttp.Handler())
	// 打印访问地址
	host := addr
	if strings.HasPrefix(host, ":") {
		host = "localhost" + host
	}
	slog.Info("Prometheus metrics endpoint started", "address", fmt.Sprintf("http://%s/metrics", host))

	go func() {
		if err := http.ListenAndServe(addr, nil); err != nil {
			slog.Error("StartPrometheusEndpoint is failed", slog.Any("err", err))
		}
	}()
}

var (
	// Once jobs metrics
	onceJobTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: Namespace,
			Name:      "once_job_executions_total",
			Help:      "Total number of once jobs executions",
		},
		[]string{"job_id", "job_name"},
	)
	onceJobFailedTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: Namespace,
			Name:      "once_job_failures_total",
			Help:      "Total number of failed once jobs executions",
		},
		[]string{"job_id", "job_name", "error"},
	)
	onceJobDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: Namespace,
			Name:      "once_job_duration_seconds",
			Help:      "Duration of once jobs executions in seconds",
			Buckets:   prometheus.ExponentialBuckets(0.1, 1.5, 5),
		},
		[]string{"job_id", "job_name"},
	)
	onceJobRunning = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: Namespace,
			Name:      "once_job_running",
			Help:      "Number of currently running once jobs",
		},
		[]string{"job_id", "job_name"},
	)
	onceJobStatus = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: Namespace,
			Name:      "once_job_status",
			Help:      "Current status of once jobs (0=inactive, 1=active, 2=failed)",
		},
		[]string{"job_id", "job_name"},
	)

	// Interval jobs metrics
	intervalJobTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: Namespace,
			Name:      "interval_job_executions_total",
			Help:      "Total number of interval jobs executions",
		},
		[]string{"job_id", "job_name"},
	)
	intervalJobFailedTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: Namespace,
			Name:      "interval_job_failures_total",
			Help:      "Total number of failed interval jobs executions",
		},
		[]string{"job_id", "job_name", "error"},
	)
	intervalJobDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: Namespace,
			Name:      "interval_job_duration_seconds",
			Help:      "Duration of interval jobs executions in seconds",
			Buckets:   prometheus.ExponentialBuckets(0.1, 1.5, 5),
		},
		[]string{"job_id", "job_name"},
	)
	intervalJobRunning = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: Namespace,
			Name:      "interval_job_running",
			Help:      "Number of currently running interval jobs",
		},
		[]string{"job_id", "job_name"},
	)
	intervalJobStatus = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: Namespace,
			Name:      "interval_job_status",
			Help:      "Current status of interval jobs (0=inactive, 1=active, 2=failed)",
		},
		[]string{"job_id", "job_name"},
	)

	// Daily jobs metrics
	dailyJobTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: Namespace,
			Name:      "daily_job_executions_total",
			Help:      "Total number of daily jobs executions",
		},
		[]string{"job_id", "job_name"},
	)
	dailyJobFailedTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: Namespace,
			Name:      "daily_job_failures_total",
			Help:      "Total number of failed daily jobs executions",
		},
		[]string{"job_id", "job_name", "error"},
	)
	dailyJobDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: Namespace,
			Name:      "daily_job_duration_seconds",
			Help:      "Duration of daily jobs executions in seconds",
			Buckets:   prometheus.ExponentialBuckets(0.1, 1.5, 5),
		},
		[]string{"job_id", "job_name"},
	)
	dailyJobRunning = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: Namespace,
			Name:      "daily_job_running",
			Help:      "Number of currently running daily jobs",
		},
		[]string{"job_id", "job_name"},
	)
	dailyJobStatus = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: Namespace,
			Name:      "daily_job_status",
			Help:      "Current status of daily jobs (0=inactive, 1=active, 2=failed)",
		},
		[]string{"job_id", "job_name"},
	)

	// Weekly jobs metrics
	weeklyJobTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: Namespace,
			Name:      "weekly_job_executions_total",
			Help:      "Total number of weekly jobs executions",
		},
		[]string{"job_id", "job_name"},
	)
	weeklyJobFailedTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: Namespace,
			Name:      "weekly_job_failures_total",
			Help:      "Total number of failed weekly jobs executions",
		},
		[]string{"job_id", "job_name", "error"},
	)
	weeklyJobDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: Namespace,
			Name:      "weekly_job_duration_seconds",
			Help:      "Duration of weekly jobs executions in seconds",
			Buckets:   prometheus.ExponentialBuckets(0.1, 1.5, 5),
		},
		[]string{"job_id", "job_name"},
	)
	weeklyJobRunning = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: Namespace,
			Name:      "weekly_job_running",
			Help:      "Number of currently running weekly jobs",
		},
		[]string{"job_id", "job_name"},
	)
	weeklyJobStatus = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: Namespace,
			Name:      "weekly_job_status",
			Help:      "Current status of weekly jobs (0=inactive, 1=active, 2=failed)",
		},
		[]string{"job_id", "job_name"},
	)

	// Monthly jobs metrics
	monthlyJobTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: Namespace,
			Name:      "monthly_job_executions_total",
			Help:      "Total number of monthly jobs executions",
		},
		[]string{"job_id", "job_name"},
	)
	monthlyJobFailedTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: Namespace,
			Name:      "monthly_job_failures_total",
			Help:      "Total number of failed monthly jobs executions",
		},
		[]string{"job_id", "job_name", "error"},
	)
	monthlyJobDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: Namespace,
			Name:      "monthly_job_duration_seconds",
			Help:      "Duration of monthly jobs executions in seconds",
			Buckets:   prometheus.ExponentialBuckets(0.1, 1.5, 5),
		},
		[]string{"job_id", "job_name"},
	)
	monthlyJobRunning = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: Namespace,
			Name:      "monthly_job_running",
			Help:      "Number of currently running monthly jobs",
		},
		[]string{"job_id", "job_name"},
	)
	monthlyJobStatus = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: Namespace,
			Name:      "monthly_job_status",
			Help:      "Current status of monthly jobs (0=inactive, 1=active, 2=failed)",
		},
		[]string{"job_id", "job_name"},
	)

	// Cron jobs metrics
	cronJobTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: Namespace,
			Name:      "cron_job_executions_total",
			Help:      "Total number of cron jobs executions",
		},
		[]string{"job_id", "job_name"},
	)
	cronJobFailedTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: Namespace,
			Name:      "cron_job_failures_total",
			Help:      "Total number of failed cron jobs executions",
		},
		[]string{"job_id", "job_name", "error"},
	)
	cronJobDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: Namespace,
			Name:      "cron_job_duration_seconds",
			Help:      "Duration of cron jobs executions in seconds",
			Buckets:   prometheus.ExponentialBuckets(0.1, 1.5, 5),
		},
		[]string{"job_id", "job_name"},
	)
	cronJobRunning = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: Namespace,
			Name:      "cron_job_running",
			Help:      "Number of currently running cron jobs",
		},
		[]string{"job_id", "job_name"},
	)
	cronJobStatus = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: Namespace,
			Name:      "cron_job_status",
			Help:      "Current status of cron jobs (0=inactive, 1=active, 2=failed)",
		},
		[]string{"job_id", "job_name"},
	)
)

// RecordJobExecution Unified entry for recording jobs execution
func RecordJobExecution(jobType jobs.JobType, jobID, jobName string, durationSeconds float64, success bool, err error) {
	labels := prometheus.Labels{
		"job_id":   jobID,
		"job_name": jobName,
	}
	var (
		total    *prometheus.CounterVec
		failed   *prometheus.CounterVec
		duration *prometheus.HistogramVec
		status   *prometheus.GaugeVec
	)

	switch jobType {
	case jobs.JobTypeOnce:
		total = onceJobTotal
		failed = onceJobFailedTotal
		duration = onceJobDuration
		status = onceJobStatus
	case jobs.JobInterval:
		total = intervalJobTotal
		failed = intervalJobFailedTotal
		duration = intervalJobDuration
		status = intervalJobStatus
	case jobs.JobTypeDaily:
		total = dailyJobTotal
		failed = dailyJobFailedTotal
		duration = dailyJobDuration
		status = dailyJobStatus
	case jobs.JobTypeWeekly:
		total = weeklyJobTotal
		failed = weeklyJobFailedTotal
		duration = weeklyJobDuration
		status = weeklyJobStatus
	case jobs.JobTypeMonthly:
		total = monthlyJobTotal
		failed = monthlyJobFailedTotal
		duration = monthlyJobDuration
		status = monthlyJobStatus
	case jobs.JobTypeCron:
		total = cronJobTotal
		failed = cronJobFailedTotal
		duration = cronJobDuration
		status = cronJobStatus
	default:
		return // unknown type, do nothing
	}
	total.With(labels).Inc()
	duration.With(labels).Observe(durationSeconds)
	if !success && err != nil {
		errorLabels := prometheus.Labels{
			"job_id":   jobID,
			"job_name": jobName,
			"error":    err.Error(),
		}
		failed.With(errorLabels).Inc()
		status.With(labels).Set(0)
	} else {
		status.With(labels).Set(1)
	}
}

// IncJobRunning Unified entry for incrementing running jobs
func IncJobRunning(jobType jobs.JobType, jobID, jobName string) bool {
	// 检查是否已经在运行
	if _, loaded := runningJobs.LoadOrStore(jobID, struct{}{}); loaded {
		slog.Warn("Job is already running", "job_id", jobID, "job_name", jobName)
		return false
	}

	labels := prometheus.Labels{
		"job_id":   jobID,
		"job_name": jobName,
	}
	switch jobType {
	case jobs.JobTypeOnce:
		onceJobRunning.With(labels).Inc()
	case jobs.JobInterval:
		intervalJobRunning.With(labels).Inc()
	case jobs.JobTypeDaily:
		dailyJobRunning.With(labels).Inc()
	case jobs.JobTypeWeekly:
		weeklyJobRunning.With(labels).Inc()
	case jobs.JobTypeMonthly:
		monthlyJobRunning.With(labels).Inc()
	case jobs.JobTypeCron:
		cronJobRunning.With(labels).Inc()
	}
	return true
}

// DecJobRunning Unified entry for decrementing running jobs
func DecJobRunning(jobType jobs.JobType, jobID, jobName string) {
	// 从运行map中移除
	runningJobs.Delete(jobID)

	labels := prometheus.Labels{
		"job_id":   jobID,
		"job_name": jobName,
	}
	switch jobType {
	case jobs.JobTypeOnce:
		onceJobRunning.With(labels).Dec()
	case jobs.JobInterval:
		intervalJobRunning.With(labels).Dec()
	case jobs.JobTypeDaily:
		dailyJobRunning.With(labels).Dec()
	case jobs.JobTypeWeekly:
		weeklyJobRunning.With(labels).Dec()
	case jobs.JobTypeMonthly:
		monthlyJobRunning.With(labels).Dec()
	case jobs.JobTypeCron:
		cronJobRunning.With(labels).Dec()
	}
}

// In init(), register all metrics
func init() {
	runningJobs = sync.Map{}
	// Once
	registry.MustRegister(onceJobTotal)
	registry.MustRegister(onceJobFailedTotal)
	registry.MustRegister(onceJobDuration)
	registry.MustRegister(onceJobRunning)
	registry.MustRegister(onceJobStatus)
	// Interval
	registry.MustRegister(intervalJobTotal)
	registry.MustRegister(intervalJobFailedTotal)
	registry.MustRegister(intervalJobDuration)
	registry.MustRegister(intervalJobRunning)
	registry.MustRegister(intervalJobStatus)
	// Daily
	registry.MustRegister(dailyJobTotal)
	registry.MustRegister(dailyJobFailedTotal)
	registry.MustRegister(dailyJobDuration)
	registry.MustRegister(dailyJobRunning)
	registry.MustRegister(dailyJobStatus)
	// Weekly
	registry.MustRegister(weeklyJobTotal)
	registry.MustRegister(weeklyJobFailedTotal)
	registry.MustRegister(weeklyJobDuration)
	registry.MustRegister(weeklyJobRunning)
	registry.MustRegister(weeklyJobStatus)
	// Monthly
	registry.MustRegister(monthlyJobTotal)
	registry.MustRegister(monthlyJobFailedTotal)
	registry.MustRegister(monthlyJobDuration)
	registry.MustRegister(monthlyJobRunning)
	registry.MustRegister(monthlyJobStatus)
	// Cron
	registry.MustRegister(cronJobTotal)
	registry.MustRegister(cronJobFailedTotal)
	registry.MustRegister(cronJobDuration)
	registry.MustRegister(cronJobRunning)
	registry.MustRegister(cronJobStatus)
}
