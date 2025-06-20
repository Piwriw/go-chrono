#  go-chrono
## 支持功能
### 支持的任务类型
- Cron表达式（Cron Job）：任务可以使用 crontab 格式的时间表达式运行。
- IntervalJob：任务可以每隔 X 秒、分钟、小时、天、周、月运行。
- 每日（DailyJob）：任务可以每隔 X 天在特定的时间运行。
- 每周（WeeklyJob）：任务可以每隔 X 周在特定的星期几和时间运行。
- 每月（MonthlyJob）：任务可以每隔 X 月在特定的日期和时间运行。
- 一次性（OnceJob）：任务可以在特定的时间运行（可以是一次或多次）

### 支持功能
- [x] 允许设置调度器全局配置，优先使用Job自身配置
  - [x] 超时时间
  - [x] Watch监听模式
  - [ ] 钩子函数
- [ ] 允许自定义实现JobClient
- [ ] 设置Schedule最大管控调用任务数量
- [ ] Job分组级别管控
- [ ] 监控模式
  - [ ]  Prometheus端点监控
  - [ ]  Listen Web监控
  - [ ]  可以查询最近X次的执行结果