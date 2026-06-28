package scheduler

import (
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/piwriw/go-chrono/monitor"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestUUIDEventIDGenerator tests the UUIDEventIDGenerator functionality.
// TestUUIDEventIDGenerator 测试 UUIDEventIDGenerator 的功能。
//
// Test scenarios:
// - Happy path: generate unique UUIDs
// - Boundary conditions: multiple calls should generate different IDs
// - Exception cases: nil jobs spec should still work
func TestUUIDEventIDGenerator(t *testing.T) {
	t.Parallel()

	generator := &monitor.UUIDEventIDGenerator{}

	t.Run("generates valid UUID format", func(t *testing.T) {
		t.Parallel()

		jobSpec := monitor.JobSpec{
			JobID:   "test-jobs-123",
			JobName: "test-jobs",
			Tags:    []string{"tag1", "tag2"},
		}

		eventID := generator.NextID(jobSpec)

		// Verify it's a valid UUID
		// 验证是有效的UUID
		parsedUUID, err := uuid.Parse(eventID)
		require.NoError(t, err, "generated ID should be a valid UUID")
		assert.NotEqual(t, uuid.Nil, parsedUUID, "generated UUID should not be nil")
	})

	t.Run("generates unique IDs for consecutive calls", func(t *testing.T) {
		t.Parallel()

		jobSpec := monitor.JobSpec{
			JobID:   "test-jobs-123",
			JobName: "test-jobs",
		}

		ids := make(map[string]bool)
		iterations := 100

		for range iterations {
			eventID := generator.NextID(jobSpec)
			// Check uniqueness
			// 检查唯一性
			_, exists := ids[eventID]
			assert.False(t, exists, "generated ID should be unique")
			ids[eventID] = true
		}

		assert.Equal(t, iterations, len(ids), "should generate unique IDs")
	})

	t.Run("generates unique IDs for different jobs specs", func(t *testing.T) {
		t.Parallel()

		jobSpecs := []monitor.JobSpec{
			{JobID: "jobs-1", JobName: "jobs-one", Tags: []string{"tag1"}},
			{JobID: "jobs-2", JobName: "jobs-two", Tags: []string{"tag2"}},
			{JobID: "jobs-3", JobName: "jobs-three", Tags: []string{"tag3"}},
		}

		ids := make([]string, len(jobSpecs))
		for i, spec := range jobSpecs {
			ids[i] = generator.NextID(spec)
		}

		// All IDs should be unique
		// 所有ID应该是唯一的
		for i := 0; i < len(ids); i++ {
			for j := i + 1; j < len(ids); j++ {
				assert.NotEqual(t, ids[i], ids[j], "IDs should be unique across different jobs specs")
			}
		}
	})

	t.Run("thread-safe concurrent generation", func(t *testing.T) {
		t.Parallel()

		jobSpec := monitor.JobSpec{
			JobID:   "concurrent-jobs",
			JobName: "concurrent-jobs",
		}

		const goroutines = 50
		const callsPerGoroutine = 20

		var wg sync.WaitGroup
		ids := make(chan string, goroutines*callsPerGoroutine)

		for range goroutines {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < callsPerGoroutine; j++ {
					ids <- generator.NextID(jobSpec)
				}
			}()
		}

		wg.Wait()
		close(ids)

		// Collect and verify uniqueness
		// 收集并验证唯一性
		uniqueIDs := make(map[string]bool)
		for id := range ids {
			_, exists := uniqueIDs[id]
			assert.False(t, exists, "concurrent generation should produce unique IDs")
			uniqueIDs[id] = true
		}

		assert.Equal(t, goroutines*callsPerGoroutine, len(uniqueIDs), "all generated IDs should be unique")
	})

	t.Run("implements monitor.EventIDGenerator interface", func(t *testing.T) {
		t.Parallel()

		var _ monitor.EventIDGenerator = &monitor.UUIDEventIDGenerator{}
		var _ monitor.EventIDGenerator = generator
	})
}

// TestTimeEventIDGenerator tests the TimeEventIDGenerator functionality.
// TestTimeEventIDGenerator 测试 TimeEventIDGenerator 的功能。
//
// Test scenarios:
// - Happy path: generate time-based IDs
// - Boundary conditions: custom time format, default time format
// - Exception cases: empty jobs spec fields
func TestTimeEventIDGenerator(t *testing.T) {
	t.Parallel()

	t.Run("generates ID with default time format", func(t *testing.T) {
		t.Parallel()

		generator := &monitor.TimeEventIDGenerator{}
		jobSpec := monitor.JobSpec{
			JobID:   "jobs-123",
			JobName: "test-jobs",
		}

		eventID := generator.NextID(jobSpec)

		// Format: jobID_jobName_timestamp
		// 格式：jobID_jobName_时间戳
		expectedPrefix := "jobs-123_test-job_"
		assert.True(t, strings.HasPrefix(eventID, expectedPrefix), "ID should start with jobID and jobName")

		// Extract timestamp part
		// 提取时间戳部分
		timestampPart := strings.TrimPrefix(eventID, expectedPrefix)
		assert.NotEmpty(t, timestampPart, "timestamp part should not be empty")

		// Verify timestamp format (default: 20060102150405.000000000), with optional counter suffix for concurrency safety
		// 验证时间戳格式（默认：20060102150405.000000000），并允许为了并发安全追加可选的计数器后缀
		timestampRegex := regexp.MustCompile(`^\d{14}\.\d{9}(_\d+)?$`)
		assert.True(t, timestampRegex.MatchString(timestampPart), "timestamp should match format YYYYMMDDHHmmss.nnnnnnnnn with optional counter suffix")
	})

	t.Run("generates ID with custom time format", func(t *testing.T) {
		t.Parallel()

		customFormat := "2006-01-02 15:04:05"
		generator := monitor.NewTimeEventIDGenerator(customFormat)
		jobSpec := monitor.JobSpec{
			JobID:   "jobs-456",
			JobName: "custom-jobs",
		}

		eventID := generator.NextID(jobSpec)

		// Format: jobID_jobName_customFormat_timestamp
		// 格式：jobID_jobName_自定义格式_时间戳
		parts := strings.Split(eventID, "_")
		assert.GreaterOrEqual(t, len(parts), 3, "ID should have at least 3 parts separated by underscore")
		assert.Equal(t, "jobs-456", parts[0], "first part should be jobID")
		assert.Equal(t, "custom-job", parts[1], "second part should be jobName prefix")
	})

	t.Run("generates different IDs for consecutive calls", func(t *testing.T) {
		t.Parallel()

		generator := &monitor.TimeEventIDGenerator{}
		jobSpec := monitor.JobSpec{
			JobID:   "jobs-789",
			JobName: "time-test-jobs",
		}

		id1 := generator.NextID(jobSpec)
		time.Sleep(10 * time.Millisecond) // Ensure different timestamp / 确保不同的时间戳
		id2 := generator.NextID(jobSpec)

		assert.NotEqual(t, id1, id2, "consecutive calls should generate different IDs due to different timestamps")
	})

	t.Run("handles empty jobs spec fields", func(t *testing.T) {
		t.Parallel()

		generator := &monitor.TimeEventIDGenerator{}
		jobSpec := monitor.JobSpec{
			JobID:   "",
			JobName: "",
		}

		eventID := generator.NextID(jobSpec)

		// Should still generate ID with empty fields
		// 仍然应该使用空字段生成ID
		expectedPrefix := "__"
		assert.True(t, strings.HasPrefix(eventID, expectedPrefix), "ID should handle empty jobID and jobName")
	})

	t.Run("handles jobs spec with tags", func(t *testing.T) {
		t.Parallel()

		generator := &monitor.TimeEventIDGenerator{}
		jobSpec := monitor.JobSpec{
			JobID:   "tagged-jobs",
			JobName: "tagged-jobs-name",
			Tags:    []string{"important", "daily"},
		}

		eventID := generator.NextID(jobSpec)

		// Tags should not affect ID generation
		// 标签不应该影响ID生成
		assert.Contains(t, eventID, "tagged-job_tagged-jobs-name_", "ID should contain jobID and jobName")
	})

	t.Run("thread-safe concurrent generation", func(t *testing.T) {
		t.Parallel()

		generator := &monitor.TimeEventIDGenerator{}
		jobSpec := monitor.JobSpec{
			JobID:   "concurrent-time-jobs",
			JobName: "concurrent-time-jobs",
		}

		const goroutines = 50
		const callsPerGoroutine = 20

		var wg sync.WaitGroup
		ids := make(chan string, goroutines*callsPerGoroutine)

		for range goroutines {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < callsPerGoroutine; j++ {
					ids <- generator.NextID(jobSpec)
				}
			}()
		}

		wg.Wait()
		close(ids)

		// All IDs should be unique (due to different timestamps)
		// 所有ID应该是唯一的（由于不同的时间戳）
		uniqueIDs := make(map[string]bool)
		for id := range ids {
			uniqueIDs[id] = true
		}

		assert.Equal(t, goroutines*callsPerGoroutine, len(uniqueIDs), "concurrent generation should produce unique IDs")
	})

	t.Run("implements monitor.EventIDGenerator interface", func(t *testing.T) {
		t.Parallel()

		var _ monitor.EventIDGenerator = &monitor.TimeEventIDGenerator{}
	})
}

// BenchmarkUUIDEventIDGenerator benchmarks the UUIDEventIDGenerator.
// BenchmarkUUIDEventIDGenerator 对 UUIDEventIDGenerator 进行基准测试。
func BenchmarkUUIDEventIDGenerator(b *testing.B) {
	generator := &monitor.UUIDEventIDGenerator{}
	jobSpec := monitor.JobSpec{
		JobID:   "benchmark-jobs",
		JobName: "benchmark-jobs",
		Tags:    []string{"benchmark"},
	}

	b.ResetTimer()
	b.ReportAllocs()

	for range b.N {
		_ = generator.NextID(jobSpec)
	}
}

// BenchmarkTimeEventIDGenerator benchmarks the TimeEventIDGenerator.
// BenchmarkTimeEventIDGenerator 对 TimeEventIDGenerator 进行基准测试。
func BenchmarkTimeEventIDGenerator(b *testing.B) {
	generator := monitor.NewTimeEventIDGenerator("20060102150405")
	jobSpec := monitor.JobSpec{
		JobID:   "benchmark-jobs",
		JobName: "benchmark-jobs",
		Tags:    []string{"benchmark"},
	}

	b.ResetTimer()
	b.ReportAllocs()

	for range b.N {
		_ = generator.NextID(jobSpec)
	}
}

// TestUUIDEventIDGenerator_Demo demonstrates the usage of UUIDEventIDGenerator.
// TestUUIDEventIDGenerator_Demo 展示 UUIDEventIDGenerator 的使用方法。
func TestUUIDEventIDGenerator_Demo(t *testing.T) {
	generator := &monitor.UUIDEventIDGenerator{}
	jobSpec := monitor.JobSpec{
		JobID:   "my-jobs-123",
		JobName: "my-jobs",
		Tags:    []string{"important"},
	}

	eventID := generator.NextID(jobSpec)
	assert.NotEmpty(t, eventID, "UUID event ID should not be empty")
}

// TestTimeEventIDGenerator_Demo demonstrates the usage of TimeEventIDGenerator.
// TestTimeEventIDGenerator_Demo 展示 TimeEventIDGenerator 的使用方法。
func TestTimeEventIDGenerator_Demo(t *testing.T) {
	generator := monitor.NewTimeEventIDGenerator("20060102150405")
	jobSpec := monitor.JobSpec{
		JobID:   "my-jobs-456",
		JobName: "my-jobs",
		Tags:    []string{"daily"},
	}

	eventID := generator.NextID(jobSpec)
	assert.NotEmpty(t, eventID, "time-based event ID should not be empty")
}

// TestEventIDGeneratorInterface tests the monitor.EventIDGenerator interface.
// TestEventIDGeneratorInterface 测试 monitor.EventIDGenerator 接口。
func TestEventIDGeneratorInterface(t *testing.T) {
	t.Parallel()

	t.Run("UUIDEventIDGenerator implements interface", func(t *testing.T) {
		t.Parallel()

		var generator monitor.EventIDGenerator = &monitor.UUIDEventIDGenerator{}
		jobSpec := monitor.JobSpec{JobID: "test", JobName: "test"}

		eventID := generator.NextID(jobSpec)
		assert.NotEmpty(t, eventID, "should generate non-empty ID")
	})

	t.Run("TimeEventIDGenerator implements interface", func(t *testing.T) {
		t.Parallel()

		var generator monitor.EventIDGenerator = &monitor.TimeEventIDGenerator{}
		jobSpec := monitor.JobSpec{JobID: "test", JobName: "test"}

		eventID := generator.NextID(jobSpec)
		assert.NotEmpty(t, eventID, "should generate non-empty ID")
	})

	t.Run("can use interface polymorphically", func(t *testing.T) {
		t.Parallel()

		generators := []monitor.EventIDGenerator{
			&monitor.UUIDEventIDGenerator{},
			monitor.NewTimeEventIDGenerator("20060102150405"),
		}

		jobSpec := monitor.JobSpec{JobID: "poly-test", JobName: "poly-test"}

		for i, gen := range generators {
			eventID := gen.NextID(jobSpec)
			assert.NotEmpty(t, eventID, "generator %d should generate non-empty ID", i)
		}
	})
}
