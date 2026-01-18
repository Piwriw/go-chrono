package monitor

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/piwriw/go-chrono/common"
	"github.com/piwriw/go-chrono/pkg/url"
)

type WebMonitor struct {
	scheduler       common.SchedulerInterface
	scheduleMonitor SchedulerMonitor
	addr            string
	server          *http.Server
	mu              sync.Mutex
}

func NewWebMonitor(s common.SchedulerInterface, m SchedulerMonitor, addr string) *WebMonitor {
	return &WebMonitor{
		scheduler:       s,
		scheduleMonitor: m,
		addr:            addr,
	}
}

func (wm *WebMonitor) Start() error {
	mux := http.NewServeMux()
	endpoints := []string{"/healthz", "/jobs", "/jobs/{job_id}/retries"}

	mux.HandleFunc("/healthz", wm.handleHealthz)
	mux.HandleFunc("/jobs", wm.handleJobs)
	mux.HandleFunc("/jobs/", wm.handleJobRetries)

	if err := url.ValidateURLAddr(wm.addr); err != nil {
		return fmt.Errorf("invalid address: %v", err)
	}
	wm.server = &http.Server{
		Addr:    wm.addr,
		Handler: mux,
	}
	// 打印接口地址
	host := wm.addr
	if strings.HasPrefix(host, ":") {
		host = "localhost" + host
	}
	slog.Info("chrono:web monitor started", "address", wm.addr)
	for _, ep := range endpoints {
		slog.Info(fmt.Sprintf("  http://%s%s", host, ep))
	}
	go func() {
		if err := wm.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Info("WebMonitor Listen is error", slog.Any("err", err))
		}
	}()
	return nil
}

func (wm *WebMonitor) handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(`{"status":"ok"}`))
	if err != nil {
		slog.Error("write healthz response error", slog.Any("err", err))
		return
	}
}

type JobMonitorSpec struct {
	ID      string      `json:"id"`
	Name    string      `json:"name"`
	Alias   *string     `json:"alias"`
	LastRun *time.Time  `json:"last_run"`
	NextRun *time.Time  `json:"next_run"`
	Events  []*JobEvent `json:"events"`
}

func (j JobMonitorSpec) MarshalJSON() ([]byte, error) {
	type Alias JobMonitorSpec

	var lastRunStr, nextRunStr string

	if j.LastRun != nil {
		lastRunStr = j.LastRun.Format(time.DateTime)
	}

	if j.NextRun != nil {
		nextRunStr = j.NextRun.Format(time.DateTime)
	}

	return json.Marshal(&struct {
		Alias
		LastRun string `json:"last_run"`
		NextRun string `json:"next_run"`
	}{
		Alias:   (Alias)(j),
		LastRun: lastRunStr,
		NextRun: nextRunStr,
	})
}

func (wm *WebMonitor) handleJobs(w http.ResponseWriter, r *http.Request) {
	jobs, err := wm.scheduler.GetJobs()
	if err != nil {
		http.Error(w, "failed to get jobs", http.StatusInternalServerError)
		return
	}
	resJobs := make([]JobMonitorSpec, 0, len(jobs))
	for _, job := range jobs {
		last, next, err := wm.scheduler.GetJobLastAndNextByID(job.ID().String())
		if err != nil {
			http.Error(w, "failed to get jobs last and next", http.StatusInternalServerError)
		}
		events := wm.scheduleMonitor.GetJobEvents(job.ID().String())
		spec := JobMonitorSpec{
			ID:     job.ID().String(),
			Name:   job.Name(),
			Events: events,
		}
		if !last.IsZero() {
			spec.LastRun = &last
		}
		if !next.IsZero() {
			spec.NextRun = &next
		}
		if wm.scheduler.Enable(common.AliasOptionName) {
			alias, err := wm.scheduler.GetAlias(job.ID().String())
			if err != nil {
				http.Error(w, "failed to get alias", http.StatusInternalServerError)
			}
			spec.Alias = &alias
		}
		resJobs = append(resJobs, spec)
	}
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(resJobs); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

// handleJobRetries handles requests to get retry history for a specific jobs.
// handleJobRetries 处理获取特定任务重试历史的请求。
//
// URL pattern: /jobs/{job_id}/retries
//
// Parameters:
//
//	w - HTTP response writer / HTTP 响应写入器
//	r - HTTP request / HTTP 请求
func (wm *WebMonitor) handleJobRetries(w http.ResponseWriter, r *http.Request) {
	// Extract jobs ID from URL path
	// 从 URL 路径中提取任务 ID
	// Expected format: /jobs/{job_id}/retries
	path := r.URL.Path
	// Remove /jobs/ prefix and /retries suffix
	// 移除 /jobs/ 前缀和 /retries 后缀
	parts := strings.Split(strings.TrimPrefix(path, "/jobs/"), "/")
	if len(parts) < 2 || parts[1] != "retries" {
		http.Error(w, "invalid URL format, expected /jobs/{job_id}/retries", http.StatusBadRequest)
		return
	}

	jobID := parts[0]
	if jobID == "" {
		http.Error(w, "job_id is required", http.StatusBadRequest)
		return
	}

	// Get retry history from monitor
	// 从监控器获取重试历史
	history := wm.scheduleMonitor.GetRetryHistory(jobID)

	// Return JSON response
	// 返回 JSON 响应
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(history); err != nil {
		http.Error(w, "failed to encode retry history", http.StatusInternalServerError)
		slog.Error("failed to encode retry history", "error", err)
	}
}
