#  go-chrono
## 支持功能
### 支持的任务类型
- Cron表达式（Cron Job）：任务可以使用 crontab 格式的时间表达式运行
- IntervalJob：任务可以每隔 X 秒、分钟、小时、天、周、月运行
- 每日（DailyJob）：任务可以每隔 X 天在特定的时间运行
- 每周（WeeklyJob）：任务可以每隔 X 周在特定的星期几和时间运行
- 每月（MonthlyJob）：任务可以每隔 X 月在特定的日期和时间运行
- 一次性（OnceJob）：任务可以在特定的时间运行（可以是一次或多次）

### 支持功能
- [x] 允许设置调度器全局配置，优先使用Job自身配置
  - [x] 超时时间
  - [x] Watch监听模式
  - [ ] 钩子函数
  - [ ] 重试策略(是否至少间隔多大?，防止溢出风险)
- [ ] 允许自定义实现JobClient
- [ ] 统一Time.Format格式
- [x] 设置Schedule最大管控调用任务数量
- [ ] 自定义Logger
- [x] Job标签级别管控
  - 缺少移除func
- [x] 优雅终止关闭
- [x] 监控模式
  - [ ]  Prometheus端点监控
  - [x]  Listen Web监控
  - [x]  可以查询最近X次的执行结果
     - [ ]  支持查询指定时间范围的执行结果
     - [ ]  支持查询指定标签的执行结果
     - [ ]  支持查询指定任务的执行结果
     - [ ]  支持查询指定任务的执行结果
     - [x]  新增EventID支持
        - 允许自定义实现Event ID生成器
        - 默认使用JobID+JobName+时间戳
        - 内置支持生成器：
          1. UUID
          2. JobID+JobName+时间戳