package scheduler

import (
	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
)

// GocronJobWrapper adapts gocron.Job to pkg.Job interface
// GocronJobWrapper 将 gocron.Job 适配到 pkg.Job 接口
type GocronJobWrapper struct {
	job gocron.Job
}

// NewGocronJobWrapper creates a new wrapper for a gocron.Job
// NewGocronJobWrapper 为 gocron.Job 创建一个新的包装器
//
// Parameters:
//
//	jobs - The gocron.Job to wrap / 要包装的 gocron.Job
//
// Returns:
//
//	*GocronJobWrapper - The wrapper / 包装器
func NewGocronJobWrapper(job gocron.Job) *GocronJobWrapper {
	return &GocronJobWrapper{job: job}
}

// ID returns the jobs ID.
// ID 返回任务 ID。
//
// Returns:
//
//	uuid.UUID - The jobs ID / 任务 ID
func (w *GocronJobWrapper) ID() uuid.UUID {
	return w.job.ID()
}

// Name returns the jobs name.
// Name 返回任务名称。
//
// Returns:
//
//	string - The jobs name / 任务名称
func (w *GocronJobWrapper) Name() string {
	return w.job.Name()
}

// GetJobID returns the jobs ID as string.
// GetJobID 返回任务ID字符串。
//
// Returns:
//
//	string - The jobs ID / 任务ID
func (w *GocronJobWrapper) GetJobID() string {
	return w.job.ID().String()
}

// GetJobName returns the jobs name.
// GetJobName 返回任务名称。
//
// Returns:
//
//	string - The jobs name / 任务名称
func (w *GocronJobWrapper) GetJobName() string {
	return w.job.Name()
}
