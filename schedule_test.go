package chrono

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"testing"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
)

type customJobMonitor struct {
	jobChan chan JobWatchInterface
}
type customJobSpec struct {
	ID string
}

func (c customJobSpec) GetJobID() string {
	return c.ID
}

func (c customJobSpec) GetJobName() string {
	return "customJobSpec"
}

func (c customJobSpec) GetStartTime() time.Time {
	return time.Now()
}

func (c customJobSpec) GetEndTime() time.Time {
	return time.Now()
}

func (c customJobSpec) GetStatus() gocron.JobStatus {
	return "success"
}

func (c customJobSpec) GetTags() []string {
	return []string{}
}

func (c customJobSpec) Error() error {
	return errors.New("error")
}

func (c customJobMonitor) IncrementJob(id uuid.UUID, name string, tags []string, status gocron.JobStatus) {
	slog.Info("IncrementJob", "JobID", id, "JobName", name, "tags", tags, "status", status)
}

func (c customJobMonitor) RecordJobTiming(startTime, endTime time.Time, id uuid.UUID, name string, tags []string) {
	slog.Info("IncrementJob", "JobID", id, "JobName", name, "tags", tags)
}

func (c customJobMonitor) RecordJobTimingWithStatus(startTime, endTime time.Time, id uuid.UUID, name string, tags []string, status gocron.JobStatus, err error) {
	c.jobChan <- customJobSpec{ID: id.String()}
}

func (c customJobMonitor) Watch() chan JobWatchInterface {
	return c.jobChan
}

func TestCustomJobMonitor(t *testing.T) {
	scheduler, err := NewScheduler(context.TODO(), customJobMonitor{jobChan: make(chan JobWatchInterface)})
	if err != nil {
		t.Fatal(err)
	}
	// 添加一个 Cron 任务
	task2 := func() error {
		fmt.Println("Task2 executed with parameters:")
		return nil
	}
	if err != nil {
		t.Fatal(err)
	}
	scheduler.Start()
	intervalJob2 := NewIntervalJob(time.Second * 20).
		JobID("550e8400-e29b-41d4-a716-446655440000").
		Names("TestTwoJob").
		Task(task2).Watch(func(event JobWatchInterface) {
		fmt.Println("StartTime", event.GetStartTime().Format("2006-04-02 15-04-05"),
			"EndTime", event.GetEndTime().Format("2006-04-02 15-04-05"),
			"Duration", event.GetEndTime().Sub(event.GetStartTime()))
	})

	job, err := scheduler.AddIntervalJob(intervalJob2)
	if err != nil {
		t.Fatal(err)
	}
	nextRun, err := job.NextRun()
	go scheduler.Watch()
	t.Log("First Task", job.ID(), "TASK NAME", job.Name(), "nextRunTime", nextRun.Format("2006-01-02 15:04:05"))
	// block until you are ready to shut down
	select {
	case <-time.After(time.Minute * 10):
	}
}
func TestTwoJob(t *testing.T) {
	scheduler, err := NewScheduler(nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	// 添加一个 Cron 任务
	task2 := func() error {
		fmt.Println("Task2 executed with parameters:")
		return nil
	}
	if err != nil {
		t.Fatal(err)
	}
	scheduler.Start()
	intervalJob2 := NewIntervalJob(time.Second*20).
		Names("TestTwoJob").
		Task(task2, 12).Watch(func(event JobWatchInterface) {
		fmt.Println("StartTime", event.GetStartTime().Format("2006-04-02 15-04-05"),
			"EndTime", event.GetEndTime().Format("2006-04-02 15-04-05"),
			"Duration", event.GetEndTime().Sub(event.GetStartTime()))
	})

	job, err := scheduler.AddIntervalJob(intervalJob2)
	if err != nil {
		t.Fatal(err)
	}
	nextRun, err := job.NextRun()
	go scheduler.Watch()
	t.Log("First Task", job.ID(), "TASK NAME", job.Name(), "nextRunTime", nextRun.Format("2006-01-02 15:04:05"))
	// block until you are ready to shut down
	select {
	case <-time.After(time.Minute * 10):
	}
}

func TestTimeOutJobPanic(t *testing.T) {
	scheduler, err := NewScheduler(nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	// 添加一个 Cron 任务
	task := func(a, b int) error {
		panic("panic")
	}
	name := "Joohwan"
	intervalJob := NewIntervalJob(time.Second*20).
		Names("TestTimeOutJobPanic").
		AfterJobRunsWithPanic(func(jobID uuid.UUID, jobName string, recoverData any) {
			fmt.Println("watchFunc", err, name)
		}).
		Task(task, 1, 2).Watch(func(event JobWatchInterface) {
		fmt.Println("StartTime", event.GetStartTime().Format("2006-04-02 15-04-05"),
			"EndTime", event.GetEndTime().Format("2006-04-02 15-04-05"),
			"Duration", event.GetEndTime().Sub(event.GetStartTime()))
		fmt.Println("watchFunc", event, name)
	})

	job, err := scheduler.AddIntervalJob(intervalJob)

	if err != nil {
		t.Fatal(err)
	}
	scheduler.Start()
	nextRun, err := job.NextRun()
	go scheduler.Watch()
	t.Log("First Task", job.ID(), "TASK NAME", job.Name(), "nextRunTime", nextRun.Format("2006-01-02 15:04:05"))
	// block until you are ready to shut down
	select {
	case <-time.After(time.Minute * 10):
	}
}
func TestTimeOutJob(t *testing.T) {
	scheduler, err := NewScheduler(nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	// 添加一个 Cron 任务
	task := func(a, b int) error {
		fmt.Println("Task executed with parameters:", a, b)
		return nil
	}
	name := "Joohwan"
	intervalJob := NewIntervalJob(time.Second*20).
		Names("TestTimeOutJob").
		Timeout(time.Second*80).
		Task(task, 1, 2).
		Watch(func(event JobWatchInterface) {
			// fmt.Println("StartTime", event.StartTime.Format("2006-01-02 15-04-05"),
			// 	"EndTime", event.EndTime.Format("2006-01-02 15-04-05"),
			// 	"Duration", event.EndTime.Sub(event.StartTime))
			fmt.Println("watchFunc", event, name)
		})

	job, err := scheduler.AddIntervalJob(intervalJob)

	if err != nil {
		t.Fatal(err)
	}
	scheduler.Start()
	nextRun, err := job.NextRun()
	go scheduler.Watch()
	t.Log("First Task", job.ID(), "TASK NAME", job.Name(), "nextRunTime", nextRun.Format("2006-01-02 15:04:05"))
	// block until you are ready to shut down
	select {
	case <-time.After(time.Minute * 10):
	}
}

func TestWithNoTimeOutJob(t *testing.T) {
	scheduler, err := NewScheduler(nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	// 添加一个 Cron 任务
	task := func(a, b int) error {
		fmt.Println("Task executed with parameters:", a, b)
		return nil
	}
	intervalJob := NewIntervalJob(time.Second*20).
		Names("TestWithNoTimeOutJob").
		Timeout(time.Second*80).
		Task(task, 1, 2).Watch(func(event JobWatchInterface) {
		fmt.Println("StartTime", event.GetStartTime().Format("2006-04-02 15-04-05"),
			"EndTime", event.GetEndTime().Format("2006-04-02 15-04-05"),
			"Duration", event.GetEndTime().Sub(event.GetStartTime()))
	})

	job, err := scheduler.AddIntervalJob(intervalJob)

	if err != nil {
		t.Fatal(err)
	}
	scheduler.Start()
	nextRun, err := job.NextRun()
	go scheduler.Watch()
	t.Log("First Task", job.ID(), "TASK NAME", job.Name(), "nextRunTime", nextRun.Format("2006-01-02 15:04:05"))
	// block until you are ready to shut down
	select {
	case <-time.After(time.Minute * 10):
	}
}
func TestValidateTimeOutJob(t *testing.T) {
	scheduler, err := NewScheduler(nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	// 添加一个 Cron 任务
	task := func(a, b int) error {
		fmt.Println("Task executed with parameters:", a, b)
		return nil
	}
	intervalJob := NewIntervalJob(time.Second*20).
		Names("TestValidateTimeOutJob").
		Timeout(-1).
		Task(task, 1, 2).Watch(func(event JobWatchInterface) {
		fmt.Println("StartTime", event.GetStartTime().Format("2006-04-02 15-04-05"),
			"EndTime", event.GetEndTime().Format("2006-04-02 15-04-05"),
			"Duration", event.GetEndTime().Sub(event.GetStartTime()))
	})

	job, err := scheduler.AddIntervalJob(intervalJob)

	if err != nil {
		t.Fatal(err)
	}
	scheduler.Start()
	nextRun, err := job.NextRun()
	go scheduler.Watch()
	t.Log("First Task", job.ID(), "TASK NAME", job.Name(), "nextRunTime", nextRun.Format("2006-01-02 15:04:05"))
	// block until you are ready to shut down
	select {
	case <-time.After(time.Minute * 10):
	}
}
func TestWatchJob(t *testing.T) {
	scheduler, err := NewScheduler(nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	// 添加一个 Cron 任务
	task := func(a, b int) error {
		fmt.Println("Task executed with parameters:", a, b)
		return nil
	}
	name := "Joohwan"
	cronJob := NewCronJob(DayTimeToCron(time.Now().Add(time.Minute*1))).
		Names("TestWatchJob").
		Task(task, 1, 2).Watch(func(event JobWatchInterface) {
		fmt.Println("watchFunc", event, name)
	})

	job, err := scheduler.AddCronJob(cronJob)

	if err != nil {
		t.Fatal(err)
	}
	scheduler.Start()
	nextRun, err := job.NextRun()
	go scheduler.Watch()
	t.Log("First Task", job.ID(), "TASK NAME", job.Name(), "nextRunTime", nextRun.Format("2006-01-02 15:04:05"))
	// block until you are ready to shut down
	select {
	case <-time.After(time.Minute * 10):
	}
}

func TestMonthlyJob(t *testing.T) {
	scheduler, err := NewScheduler(nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	task := func(a, b int) error {
		fmt.Println("Task executed with parameters:", a, b)
		return nil
	}
	monthlyJob := NewMonthJob(1, gocron.NewDaysOfTheMonth(23), gocron.NewAtTimes(gocron.NewAtTime(11, 43, 3))).
		Names("TestMonthlyJob").
		Task(task, 1, 2)

	job, err := scheduler.AddMonthlyJob(monthlyJob)
	if err != nil {
		t.Fatal(err)
	}
	scheduler.Start()
	nextRun, err := job.NextRun()
	t.Log("First Task", job.ID(), "TASK NAME", job.Name(), "nextRunTime", nextRun.Format("2006-01-02 15:04:05"))
	// block until you are ready to shut down
	select {
	case <-time.After(time.Minute * 10):
	}
}

func TestWeeklyJob(t *testing.T) {
	scheduler, err := NewScheduler(nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	task := func(a, b int) error {
		fmt.Println("Task executed with parameters:", a, b)
		return nil
	}
	weeklyJob := NewWeeklyJob(1, gocron.NewWeekdays(time.Thursday), gocron.NewAtTimes(gocron.NewAtTime(11, 34, 3))).
		Names("TestWeeklyJob").
		Task(task, 1, 2)

	job, err := scheduler.AddWeeklyJob(weeklyJob)
	if err != nil {
		t.Fatal(err)
	}
	scheduler.Start()
	nextRun, err := job.NextRun()
	t.Log("First Task", job.ID(), "TASK NAME", job.Name(), "nextRunTime", nextRun.Format("2006-01-02 15:04:05"))
	// block until you are ready to shut down
	select {
	case <-time.After(time.Minute * 10):
	}
}
func TestDailyJob(t *testing.T) {
	scheduler, err := NewScheduler(nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	task := func(a, b int) error {
		fmt.Println("Task executed with parameters:", a, b)
		return nil
	}
	dailyJob := NewDailyJob(1, gocron.NewAtTimes(gocron.NewAtTime(19, 40, 3))).
		Names("TestDailyJob").
		Task(task, 1, 2)

	job, err := scheduler.AddDailyJob(dailyJob)
	if err != nil {
		t.Fatal(err)
	}
	scheduler.Start()
	nextRun, err := job.NextRun()
	t.Log("First Task", job.ID(), "TASK NAME", job.Name(), "nextRunTime", nextRun.Format("2006-01-02 15:04:05"))
	// block until you are ready to shut down
	select {
	case <-time.After(time.Minute * 10):
	}
}

func add(a, b int) {
	fmt.Println(a + b)
}

func countdis(a, b int) {
	fmt.Println(a - b)
}

func TestIntervalJob(t *testing.T) {
	scheduler, err := NewScheduler(nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	task := func(a, b int) error {
		fmt.Println("Task executed with parameters:", a, b)
		return nil
	}
	taskErr := func(a, b int) error {
		fmt.Println("Task executed with parameters:", a, b)
		return errors.New("error")
	}
	// 定义 funcName
	a := 1
	b := 2
	// 定义 watchFunc
	watchFunc := func(event JobWatchInterface) {
		fmt.Println("watchFunc", event)
		add(a, b)
	}
	watchFunc2 := func(event JobWatchInterface) {
		fmt.Println("watchFunc2", event)
		countdis(a, b)
	}
	intervalJob := NewIntervalJob(15*time.Second).
		Names("TestDurationJob").
		Watch(watchFunc).
		Task(task, 1, 2)

	job, err := scheduler.AddIntervalJob(intervalJob)
	if err != nil {
		t.Fatal(err)
	}

	intervalJob2 := NewIntervalJob(15*time.Second).
		Names("TestDurationJob2").
		Watch(watchFunc2).
		Task(taskErr, 3, 5)

	_, err = scheduler.AddIntervalJob(intervalJob2)
	if err != nil {
		t.Fatal(err)
	}

	scheduler.Start()
	go scheduler.Watch()
	nextRun, err := job.NextRun()
	t.Log("First Task", job.ID(), "TASK NAME", job.Name(), "nextRunTime", nextRun.Format("2006-01-02 15:04:05"))
	// block until you are ready to shut down
	select {
	case <-time.After(time.Minute * 10):
	}
}

func TestIntervalJobRe(t *testing.T) {
	scheduler, err := NewScheduler(nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	task := func() error {
		fmt.Println("Task executed start")
		time.Sleep(30 * time.Second)
		fmt.Println("Task executed done")
		return nil
	}

	onceJob := NewOnceJob(time.Now().Add(time.Second * 5)).
		Names("TestOnceJob").
		Task(task)
	job, err := scheduler.AddOnceJob(onceJob)
	if err != nil {
		t.Fatal(err)
	}
	scheduler.Start()
	go func() {
		time.Sleep(10 * time.Second)
		if err := scheduler.RemoveJob(job.ID().String()); err != nil {
			t.Fatal(err)
		}
	}()
	go func() {
		for {
			time.Sleep(5 * time.Second)
			jobs, err := scheduler.GetJobs()
			if err != nil {
				t.Fatal(err)
			}
			for _, job := range jobs {
				t.Log("job", job.ID(), "TASK NAME", job.Name())
			}
		}
	}()
	select {}
}

func TestOnceJob(t *testing.T) {
	scheduler, err := NewScheduler(nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	// 添加一个 Cron 任务
	task := func(a, b int) error {
		fmt.Println("Task executed with parameters:", a, b)
		return nil
	}
	onceJob := NewOnceJob(time.Now().Add(time.Minute)).
		Names("TestOnceJob").
		Task(task, 1, 2)

	job, err := scheduler.AddOnceJob(onceJob)

	if err != nil {
		t.Fatal(err)
	}
	scheduler.Start()
	nextRun, err := job.NextRun()
	t.Log("First Task", job.ID(), "TASK NAME", job.Name(), "nextRunTime", nextRun.Format("2006-01-02 15:04:05"))
	// block until you are ready to shut down
	select {
	case <-time.After(time.Minute * 10):
	}
}

func TestMonitor(t *testing.T) {
	scheduler, err := NewScheduler(nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	// 添加一个 Cron 任务
	task := func(a, b int) error {
		fmt.Println("Task executed with parameters:", a, b)
		return nil
	}
	cronJob := NewCronJob(DayTimeToCron(time.Now().Add(time.Minute*1))).
		Names("TestMonitor").
		Task(task, 1, 2)

	job, err := scheduler.AddCronJob(cronJob)

	if err != nil {
		t.Fatal(err)
	}
	scheduler.Start()
	nextRun, err := job.NextRun()
	go scheduler.Watch()
	t.Log("First Task", job.ID(), "TASK NAME", job.Name(), "nextRunTime", nextRun.Format("2006-01-02 15:04:05"))
	// block until you are ready to shut down
	select {
	case <-time.After(time.Minute * 10):
	}
}

func TestDefaultHooks(t *testing.T) {
	scheduler, err := NewScheduler(nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	// 添加一个 Cron 任务
	task := func(a, b int) error {
		fmt.Println("Task executed with parameters:", a, b)
		return errors.New("some error")
	}
	cronJob := NewCronJob(
		DayTimeToCron(time.Now().Add(time.Minute*1))).
		Task(task, 1, 2).
		DefaultHooks()

	job, err := scheduler.AddCronJob(cronJob)

	if err != nil {
		t.Fatal(err)
	}
	scheduler.Start()
	nextRun, err := job.NextRun()
	t.Log("First Task", job.ID(), "TASK NAME", job.Name(), "nextRunTime", nextRun.Format("2006-01-02 15:04:05"))
	// block until you are ready to shut down
	select {
	case <-time.After(time.Minute * 10):
	}
}

func TestDayTimeToCron(t *testing.T) {
	scheduler, err := NewScheduler(nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	// 添加一个 Cron 任务
	// Task with a parameter using a closure
	task := func(a, b int) {
		jobs, err := scheduler.GetJobs()
		if err != nil {
			fmt.Println(err)
		}
		for _, job := range jobs {
			// 获取当前任务信息
			jobID := job.ID()
			jobName := job.Name()
			nextRun, err := job.NextRun()
			if err != nil {
				t.Fatal(err)
			}
			t.Log("TASKID", job.ID(), "TASK NAME", job.Name(), "nextRunTime", nextRun.Format("2006-01-02 15:04:05"))
			lastRun, err := job.LastRun()
			if err != nil {
				t.Fatal(err)
			}
			t.Log("TASKID", job.ID(), "TASK NAME", job.Name(), "lastRunTime", lastRun.Format("2006-01-02 15:04:05"))
			// 使用 Monitor 记录任务信息
			scheduler.monitor.IncrementJob(jobID, jobName, nil, gocron.Success)

			// 模拟任务执行
			startTime := time.Now()
			time.Sleep(2 * time.Second) // 模拟任务运行 2 秒
			endTime := time.Now()

			// 记录任务执行时间
			scheduler.monitor.RecordJobTiming(startTime, endTime, jobID, jobName, nil)
		}
		fmt.Println("Task executed with parameters:", a, b)

	}
	cronJob := NewCronJob(DayTimeToCron(time.Now().Add(time.Minute*1))).
		Task(task, 1, 2).
		AfterJobRuns(func(jobID uuid.UUID, jobName string) {
			fmt.Println("AfterJobRuns")
		}).
		BeforeJobRuns(func(jobID uuid.UUID, jobName string) {
			fmt.Println("BeforeJobRuns")
		}).AfterJobRunsWithError(func(jobID uuid.UUID, jobName string, err error) {
		fmt.Println("AfterJobRuns")
	})

	job, err := scheduler.AddCronJob(cronJob)

	if err != nil {
		t.Fatal(err)
	}
	scheduler.Start()
	nextRun, err := job.NextRun()
	t.Log("First Task", job.ID(), "TASK NAME", job.Name(), "nextRunTime", nextRun.Format("2006-01-02 15:04:05"))
	// block until you are ready to shut down
	select {
	case <-time.After(time.Minute * 10):
	}
}
