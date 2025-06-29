package chrono

// 实现一个 Prometheus 监控端点
import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// 用于跟踪正在运行的job
	runningJobs sync.Map // map[string]struct{}
)

func StartPrometheusEndpoint(addr string) {
	if err := validateURLAddr(addr); err != nil {
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
	// Once job metrics
	onceJobTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "once_job_executions_total",
			Help:      "Total number of once job executions",
		},
		[]string{"job_id", "job_name"},
	)
	onceJobFailedTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "once_job_failures_total",
			Help:      "Total number of failed once job executions",
		},
		[]string{"job_id", "job_name", "error"},
	)
	onceJobDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "once_job_duration_seconds",
			Help:      "Duration of once job executions in seconds",
			Buckets:   prometheus.ExponentialBuckets(0.1, 1.5, 5),
		},
		[]string{"job_id", "job_name"},
	)
	onceJobRunning = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "once_job_running",
			Help:      "Number of currently running once jobs",
		},
		[]string{"job_id", "job_name"},
	)
	onceJobStatus = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "once_job_status",
			Help:      "Current status of once jobs (0=inactive, 1=active, 2=failed)",
		},
		[]string{"job_id", "job_name"},
	)

	// Interval job metrics
	intervalJobTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "interval_job_executions_total",
			Help:      "Total number of interval job executions",
		},
		[]string{"job_id", "job_name"},
	)
	intervalJobFailedTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "interval_job_failures_total",
			Help:      "Total number of failed interval job executions",
		},
		[]string{"job_id", "job_name", "error"},
	)
	intervalJobDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "interval_job_duration_seconds",
			Help:      "Duration of interval job executions in seconds",
			Buckets:   prometheus.ExponentialBuckets(0.1, 1.5, 5),
		},
		[]string{"job_id", "job_name"},
	)
	intervalJobRunning = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "interval_job_running",
			Help:      "Number of currently running interval jobs",
		},
		[]string{"job_id", "job_name"},
	)
	intervalJobStatus = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "interval_job_status",
			Help:      "Current status of interval jobs (0=inactive, 1=active, 2=failed)",
		},
		[]string{"job_id", "job_name"},
	)

	// Daily job metrics
	dailyJobTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "daily_job_executions_total",
			Help:      "Total number of daily job executions",
		},
		[]string{"job_id", "job_name"},
	)
	dailyJobFailedTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "daily_job_failures_total",
			Help:      "Total number of failed daily job executions",
		},
		[]string{"job_id", "job_name", "error"},
	)
	dailyJobDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "daily_job_duration_seconds",
			Help:      "Duration of daily job executions in seconds",
			Buckets:   prometheus.ExponentialBuckets(0.1, 1.5, 5),
		},
		[]string{"job_id", "job_name"},
	)
	dailyJobRunning = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "daily_job_running",
			Help:      "Number of currently running daily jobs",
		},
		[]string{"job_id", "job_name"},
	)
	dailyJobStatus = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "daily_job_status",
			Help:      "Current status of daily jobs (0=inactive, 1=active, 2=failed)",
		},
		[]string{"job_id", "job_name"},
	)

	// Weekly job metrics
	weeklyJobTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "weekly_job_executions_total",
			Help:      "Total number of weekly job executions",
		},
		[]string{"job_id", "job_name"},
	)
	weeklyJobFailedTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "weekly_job_failures_total",
			Help:      "Total number of failed weekly job executions",
		},
		[]string{"job_id", "job_name", "error"},
	)
	weeklyJobDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "weekly_job_duration_seconds",
			Help:      "Duration of weekly job executions in seconds",
			Buckets:   prometheus.ExponentialBuckets(0.1, 1.5, 5),
		},
		[]string{"job_id", "job_name"},
	)
	weeklyJobRunning = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "weekly_job_running",
			Help:      "Number of currently running weekly jobs",
		},
		[]string{"job_id", "job_name"},
	)
	weeklyJobStatus = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "weekly_job_status",
			Help:      "Current status of weekly jobs (0=inactive, 1=active, 2=failed)",
		},
		[]string{"job_id", "job_name"},
	)

	// Monthly job metrics
	monthlyJobTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "monthly_job_executions_total",
			Help:      "Total number of monthly job executions",
		},
		[]string{"job_id", "job_name"},
	)
	monthlyJobFailedTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "monthly_job_failures_total",
			Help:      "Total number of failed monthly job executions",
		},
		[]string{"job_id", "job_name", "error"},
	)
	monthlyJobDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "monthly_job_duration_seconds",
			Help:      "Duration of monthly job executions in seconds",
			Buckets:   prometheus.ExponentialBuckets(0.1, 1.5, 5),
		},
		[]string{"job_id", "job_name"},
	)
	monthlyJobRunning = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "monthly_job_running",
			Help:      "Number of currently running monthly jobs",
		},
		[]string{"job_id", "job_name"},
	)
	monthlyJobStatus = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "monthly_job_status",
			Help:      "Current status of monthly jobs (0=inactive, 1=active, 2=failed)",
		},
		[]string{"job_id", "job_name"},
	)

	// Cron job metrics
	cronJobTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "cron_job_executions_total",
			Help:      "Total number of cron job executions",
		},
		[]string{"job_id", "job_name"},
	)
	cronJobFailedTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "cron_job_failures_total",
			Help:      "Total number of failed cron job executions",
		},
		[]string{"job_id", "job_name", "error"},
	)
	cronJobDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "cron_job_duration_seconds",
			Help:      "Duration of cron job executions in seconds",
			Buckets:   prometheus.ExponentialBuckets(0.1, 1.5, 5),
		},
		[]string{"job_id", "job_name"},
	)
	cronJobRunning = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "cron_job_running",
			Help:      "Number of currently running cron jobs",
		},
		[]string{"job_id", "job_name"},
	)
	cronJobStatus = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "cron_job_status",
			Help:      "Current status of cron jobs (0=inactive, 1=active, 2=failed)",
		},
		[]string{"job_id", "job_name"},
	)
)

// RecordJobExecution Unified entry for recording job execution
func RecordJobExecution(jobType JobType, jobID, jobName string, durationSeconds float64, success bool, err error) {
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
	case JobTypeOnce:
		total = onceJobTotal
		failed = onceJobFailedTotal
		duration = onceJobDuration
		status = onceJobStatus
	case JobInterval:
		total = intervalJobTotal
		failed = intervalJobFailedTotal
		duration = intervalJobDuration
		status = intervalJobStatus
	case JobTypeDaily:
		total = dailyJobTotal
		failed = dailyJobFailedTotal
		duration = dailyJobDuration
		status = dailyJobStatus
	case JobTypeWeekly:
		total = weeklyJobTotal
		failed = weeklyJobFailedTotal
		duration = weeklyJobDuration
		status = weeklyJobStatus
	case JobTypeMonthly:
		total = monthlyJobTotal
		failed = monthlyJobFailedTotal
		duration = monthlyJobDuration
		status = monthlyJobStatus
	case JobTypeCron:
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
func IncJobRunning(jobType JobType, jobID, jobName string) bool {
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
	case JobTypeOnce:
		onceJobRunning.With(labels).Inc()
	case JobInterval:
		intervalJobRunning.With(labels).Inc()
	case JobTypeDaily:
		dailyJobRunning.With(labels).Inc()
	case JobTypeWeekly:
		weeklyJobRunning.With(labels).Inc()
	case JobTypeMonthly:
		monthlyJobRunning.With(labels).Inc()
	case JobTypeCron:
		cronJobRunning.With(labels).Inc()
	}
	return true
}

// DecJobRunning Unified entry for decrementing running jobs
func DecJobRunning(jobType JobType, jobID, jobName string) {
	// 从运行map中移除
	runningJobs.Delete(jobID)

	labels := prometheus.Labels{
		"job_id":   jobID,
		"job_name": jobName,
	}
	switch jobType {
	case JobTypeOnce:
		onceJobRunning.With(labels).Dec()
	case JobInterval:
		intervalJobRunning.With(labels).Dec()
	case JobTypeDaily:
		dailyJobRunning.With(labels).Dec()
	case JobTypeWeekly:
		weeklyJobRunning.With(labels).Dec()
	case JobTypeMonthly:
		monthlyJobRunning.With(labels).Dec()
	case JobTypeCron:
		cronJobRunning.With(labels).Dec()
	}
}

// In init(), register all metrics
func init() {
	runningJobs = sync.Map{}
	// Once
	prometheus.MustRegister(onceJobTotal)
	prometheus.MustRegister(onceJobFailedTotal)
	prometheus.MustRegister(onceJobDuration)
	prometheus.MustRegister(onceJobRunning)
	prometheus.MustRegister(onceJobStatus)
	// Interval
	prometheus.MustRegister(intervalJobTotal)
	prometheus.MustRegister(intervalJobFailedTotal)
	prometheus.MustRegister(intervalJobDuration)
	prometheus.MustRegister(intervalJobRunning)
	prometheus.MustRegister(intervalJobStatus)
	// Daily
	prometheus.MustRegister(dailyJobTotal)
	prometheus.MustRegister(dailyJobFailedTotal)
	prometheus.MustRegister(dailyJobDuration)
	prometheus.MustRegister(dailyJobRunning)
	prometheus.MustRegister(dailyJobStatus)
	// Weekly
	prometheus.MustRegister(weeklyJobTotal)
	prometheus.MustRegister(weeklyJobFailedTotal)
	prometheus.MustRegister(weeklyJobDuration)
	prometheus.MustRegister(weeklyJobRunning)
	prometheus.MustRegister(weeklyJobStatus)
	// Monthly
	prometheus.MustRegister(monthlyJobTotal)
	prometheus.MustRegister(monthlyJobFailedTotal)
	prometheus.MustRegister(monthlyJobDuration)
	prometheus.MustRegister(monthlyJobRunning)
	prometheus.MustRegister(monthlyJobStatus)
	// Cron
	prometheus.MustRegister(cronJobTotal)
	prometheus.MustRegister(cronJobFailedTotal)
	prometheus.MustRegister(cronJobDuration)
	prometheus.MustRegister(cronJobRunning)
	prometheus.MustRegister(cronJobStatus)
}
