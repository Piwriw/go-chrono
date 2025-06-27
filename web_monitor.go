package chrono

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"
)

type WebMonitor struct {
	scheduler       *Scheduler
	scheduleMonitor SchedulerMonitor
	addr            string
	server          *http.Server
	mu              sync.Mutex
}

func NewWebMonitor(s *Scheduler, addr string) *WebMonitor {
	return &WebMonitor{
		scheduler:       s,
		addr:            addr,
		scheduleMonitor: s.monitor,
	}
}

func (wm *WebMonitor) Start() error {
	mux := http.NewServeMux()
	endpoints := []string{"/healthz", "/jobs"}

	mux.HandleFunc("/healthz", wm.handleHealthz)
	mux.HandleFunc("/jobs", wm.handleJobs)

	if err := validateURLAddr(wm.addr); err != nil {
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
			http.Error(w, "failed to get job last and next", http.StatusInternalServerError)
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
		if wm.scheduler.Enable(AliasOptionName) {
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
