package chrono

import (
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestUUIDEventIDGenerator tests the UUIDEventIDGenerator functionality.
// TestUUIDEventIDGenerator 测试 UUIDEventIDGenerator 的功能。
//
// Test scenarios:
// - Happy path: generate unique UUIDs
// - Boundary conditions: multiple calls should generate different IDs
// - Exception cases: nil job spec should still work
func TestUUIDEventIDGenerator(t *testing.T) {
	t.Parallel()

	generator := &UUIDEventIDGenerator{}

	t.Run("generates valid UUID format", func(t *testing.T) {
		t.Parallel()

		jobSpec := JobSpec{
			JobID:   "test-job-123",
			JobName: "test-job",
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

		jobSpec := JobSpec{
			JobID:   "test-job-123",
			JobName: "test-job",
		}

		ids := make(map[string]bool)
		iterations := 100

		for i := 0; i < iterations; i++ {
			eventID := generator.NextID(jobSpec)
			// Check uniqueness
			// 检查唯一性
			_, exists := ids[eventID]
			assert.False(t, exists, "generated ID should be unique")
			ids[eventID] = true
		}

		assert.Equal(t, iterations, len(ids), "should generate unique IDs")
	})

	t.Run("generates unique IDs for different job specs", func(t *testing.T) {
		t.Parallel()

		jobSpecs := []JobSpec{
			{JobID: "job-1", JobName: "job-one", Tags: []string{"tag1"}},
			{JobID: "job-2", JobName: "job-two", Tags: []string{"tag2"}},
			{JobID: "job-3", JobName: "job-three", Tags: []string{"tag3"}},
		}

		ids := make([]string, len(jobSpecs))
		for i, spec := range jobSpecs {
			ids[i] = generator.NextID(spec)
		}

		// All IDs should be unique
		// 所有ID应该是唯一的
		for i := 0; i < len(ids); i++ {
			for j := i + 1; j < len(ids); j++ {
				assert.NotEqual(t, ids[i], ids[j], "IDs should be unique across different job specs")
			}
		}
	})

	t.Run("thread-safe concurrent generation", func(t *testing.T) {
		t.Parallel()

		jobSpec := JobSpec{
			JobID:   "concurrent-job",
			JobName: "concurrent-job",
		}

		const goroutines = 50
		const callsPerGoroutine = 20

		var wg sync.WaitGroup
		ids := make(chan string, goroutines*callsPerGoroutine)

		for i := 0; i < goroutines; i++ {
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

	t.Run("implements EventIDGenerator interface", func(t *testing.T) {
		t.Parallel()

		var _ EventIDGenerator = &UUIDEventIDGenerator{}
		var _ EventIDGenerator = generator
	})
}

// TestTimeEventIDGenerator tests the TimeEventIDGenerator functionality.
// TestTimeEventIDGenerator 测试 TimeEventIDGenerator 的功能。
//
// Test scenarios:
// - Happy path: generate time-based IDs
// - Boundary conditions: custom time format, default time format
// - Exception cases: empty job spec fields
func TestTimeEventIDGenerator(t *testing.T) {
	t.Parallel()

	t.Run("generates ID with default time format", func(t *testing.T) {
		t.Parallel()

		generator := &TimeEventIDGenerator{}
		jobSpec := JobSpec{
			JobID:   "job-123",
			JobName: "test-job",
		}

		eventID := generator.NextID(jobSpec)

		// Format: jobID_jobName_timestamp
		// 格式：jobID_jobName_时间戳
		expectedPrefix := "job-123_test-job_"
		assert.True(t, strings.HasPrefix(eventID, expectedPrefix), "ID should start with jobID and jobName")

		// Extract timestamp part
		// 提取时间戳部分
		timestampPart := strings.TrimPrefix(eventID, expectedPrefix)
		assert.NotEmpty(t, timestampPart, "timestamp part should not be empty")

		// Verify timestamp format (default: 20060102150405)
		// 验证时间戳格式（默认：20060102150405）
		timestampRegex := regexp.MustCompile(`^\d{14}$`)
		assert.True(t, timestampRegex.MatchString(timestampPart), "timestamp should match format YYYYMMDDHHmmss")
	})

	t.Run("generates ID with custom time format", func(t *testing.T) {
		t.Parallel()

		customFormat := "2006-01-02 15:04:05"
		generator := &TimeEventIDGenerator{timeFormat: customFormat}
		jobSpec := JobSpec{
			JobID:   "job-456",
			JobName: "custom-job",
		}

		eventID := generator.NextID(jobSpec)

		// Format: jobID_jobName_customFormat_timestamp
		// 格式：jobID_jobName_自定义格式_时间戳
		parts := strings.Split(eventID, "_")
		assert.GreaterOrEqual(t, len(parts), 3, "ID should have at least 3 parts separated by underscore")
		assert.Equal(t, "job-456", parts[0], "first part should be jobID")
		assert.Equal(t, "custom", parts[1], "second part should be jobName prefix")
	})

	t.Run("generates different IDs for consecutive calls", func(t *testing.T) {
		t.Parallel()

		generator := &TimeEventIDGenerator{}
		jobSpec := JobSpec{
			JobID:   "job-789",
			JobName: "time-test-job",
		}

		id1 := generator.NextID(jobSpec)
		time.Sleep(10 * time.Millisecond) // Ensure different timestamp / 确保不同的时间戳
		id2 := generator.NextID(jobSpec)

		assert.NotEqual(t, id1, id2, "consecutive calls should generate different IDs due to different timestamps")
	})

	t.Run("handles empty job spec fields", func(t *testing.T) {
		t.Parallel()

		generator := &TimeEventIDGenerator{}
		jobSpec := JobSpec{
			JobID:   "",
			JobName: "",
		}

		eventID := generator.NextID(jobSpec)

		// Should still generate ID with empty fields
		// 仍然应该使用空字段生成ID
		expectedPrefix := "__"
		assert.True(t, strings.HasPrefix(eventID, expectedPrefix), "ID should handle empty jobID and jobName")
	})

	t.Run("handles job spec with tags", func(t *testing.T) {
		t.Parallel()

		generator := &TimeEventIDGenerator{}
		jobSpec := JobSpec{
			JobID:   "tagged-job",
			JobName: "tagged-job-name",
			Tags:    []string{"important", "daily"},
		}

		eventID := generator.NextID(jobSpec)

		// Tags should not affect ID generation
		// 标签不应该影响ID生成
		assert.Contains(t, eventID, "tagged-job_tagged-job-name_", "ID should contain jobID and jobName")
	})

	t.Run("thread-safe concurrent generation", func(t *testing.T) {
		t.Parallel()

		generator := &TimeEventIDGenerator{}
		jobSpec := JobSpec{
			JobID:   "concurrent-time-job",
			JobName: "concurrent-time-job",
		}

		const goroutines = 50
		const callsPerGoroutine = 20

		var wg sync.WaitGroup
		ids := make(chan string, goroutines*callsPerGoroutine)

		for i := 0; i < goroutines; i++ {
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

	t.Run("implements EventIDGenerator interface", func(t *testing.T) {
		t.Parallel()

		var _ EventIDGenerator = &TimeEventIDGenerator{}
	})
}

// BenchmarkUUIDEventIDGenerator benchmarks the UUIDEventIDGenerator.
// BenchmarkUUIDEventIDGenerator 对 UUIDEventIDGenerator 进行基准测试。
func BenchmarkUUIDEventIDGenerator(b *testing.B) {
	generator := &UUIDEventIDGenerator{}
	jobSpec := JobSpec{
		JobID:   "benchmark-job",
		JobName: "benchmark-job",
		Tags:    []string{"benchmark"},
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = generator.NextID(jobSpec)
	}
}

// BenchmarkTimeEventIDGenerator benchmarks the TimeEventIDGenerator.
// BenchmarkTimeEventIDGenerator 对 TimeEventIDGenerator 进行基准测试。
func BenchmarkTimeEventIDGenerator(b *testing.B) {
	generator := &TimeEventIDGenerator{timeFormat: "20060102150405"}
	jobSpec := JobSpec{
		JobID:   "benchmark-job",
		JobName: "benchmark-job",
		Tags:    []string{"benchmark"},
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = generator.NextID(jobSpec)
	}
}

// ExampleUUIDEventIDGenerator demonstrates the usage of UUIDEventIDGenerator.
// ExampleUUIDEventIDGenerator 展示 UUIDEventIDGenerator 的使用方法。
func ExampleUUIDEventIDGenerator() {
	generator := &UUIDEventIDGenerator{}
	jobSpec := JobSpec{
		JobID:   "my-job-123",
		JobName: "my-job",
		Tags:    []string{"important"},
	}

	eventID := generator.NextID(jobSpec)
	println(eventID)
	// Output: (a valid UUID string, e.g., 550e8400-e29b-41d4-a716-446655440000)
}

// ExampleTimeEventIDGenerator demonstrates the usage of TimeEventIDGenerator.
// ExampleTimeEventIDGenerator 展示 TimeEventIDGenerator 的使用方法。
func ExampleTimeEventIDGenerator() {
	generator := &TimeEventIDGenerator{timeFormat: "20060102150405"}
	jobSpec := JobSpec{
		JobID:   "my-job-456",
		JobName: "my-job",
		Tags:    []string{"daily"},
	}

	eventID := generator.NextID(jobSpec)
	println(eventID)
	// Output: my-job-456_my-job_20240115143000 (timestamp will vary)
}

// TestEventIDGeneratorInterface tests the EventIDGenerator interface.
// TestEventIDGeneratorInterface 测试 EventIDGenerator 接口。
func TestEventIDGeneratorInterface(t *testing.T) {
	t.Parallel()

	t.Run("UUIDEventIDGenerator implements interface", func(t *testing.T) {
		t.Parallel()

		var generator EventIDGenerator = &UUIDEventIDGenerator{}
		jobSpec := JobSpec{JobID: "test", JobName: "test"}

		eventID := generator.NextID(jobSpec)
		assert.NotEmpty(t, eventID, "should generate non-empty ID")
	})

	t.Run("TimeEventIDGenerator implements interface", func(t *testing.T) {
		t.Parallel()

		var generator EventIDGenerator = &TimeEventIDGenerator{}
		jobSpec := JobSpec{JobID: "test", JobName: "test"}

		eventID := generator.NextID(jobSpec)
		assert.NotEmpty(t, eventID, "should generate non-empty ID")
	})

	t.Run("can use interface polymorphically", func(t *testing.T) {
		t.Parallel()

		generators := []EventIDGenerator{
			&UUIDEventIDGenerator{},
			&TimeEventIDGenerator{timeFormat: "20060102150405"},
		}

		jobSpec := JobSpec{JobID: "poly-test", JobName: "poly-test"}

		for i, gen := range generators {
			eventID := gen.NextID(jobSpec)
			assert.NotEmpty(t, eventID, "generator %d should generate non-empty ID", i)
		}
	})
}
