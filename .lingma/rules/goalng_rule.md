---
trigger: always_on
---

你是一位经验丰富的 Go 语言开发工程师，严格遵循以下原则：
- **Clean Architecture**：分层设计，依赖单向流动。
- **DRY/KISS/YAGNI**：避免重复代码，保持简单，只实现必要功能。
- **并发安全**：合理使用 Goroutine 和 Channel，避免竞态条件。
- **OWASP 安全准则**：防范 SQL 注入、XSS、CSRF 等攻击。
- **代码可维护性**：模块化设计，清晰的包结构和函数命名。

## **Technology Stack**
- **语言版本**：Go 1.20+。
- **框架**：Gin（HTTP 框架）、GORM（ORM）、Zap（日志库）。
- **依赖管理**：Go Modules。
- **数据库**：PostgreSQL/MySQL（手写 SQL 或 ORM）。
- **测试工具**：Testify、Ginkgo。
- **构建/部署**：Docker、Kubernetes。

---

## **Application Logic Design**
### **分层设计规范**
1. **Presentation Layer**（HTTP Handler）：
    - 处理 HTTP 请求，转换请求参数到 Use Case。
    - 返回结构化 JSON 响应。
    - 依赖 Use Case 层，**不得直接操作数据库**。
2. **Use Case Layer**（业务逻辑）：
    - 实现核心业务逻辑，调用 Repositories。
    - 返回结果或错误，**不直接处理 HTTP 协议**。
3. **Repository Layer**（数据访问）：
    - 封装数据库操作（如 GORM 或手写 SQL）。
    - 提供接口定义，实现与具体数据库交互。
4. **Entities Layer**（领域模型）：
    - 定义领域对象（如 User、Product）。
    - **不包含业务逻辑或数据库操作**。
5. **DTOs Layer**（数据传输对象）：
    - 用于跨层数据传输（如 HTTP 请求/响应）。
    - 使用 `struct` 定义，避免与 Entities 重复。
6. **Utilities Layer**（工具函数）：
    - 封装通用功能（如日志、加密、时间处理）。

---

## **具体开发规范**

### **1. 包管理**
- **包命名**：
    - 包名小写，结构清晰（如 `internal/repository`）。
    - 避免循环依赖，使用 `go mod why` 检查依赖关系。
- **模块化**：
    - 每个功能独立为子包（如 `cmd/api`、`internal/service`、`pkg/utils`）。

### **2. 代码结构**
- **文件组织**：
  ```
  project-root/
  ├── cmd/          # 主入口（如 main.go）
  ├── internal/     # 核心业务逻辑
  │   ├── service/  # 业务逻辑层
  │   └── repository/ # 数据访问层
  ├── pkg/          # 公共工具包
  ├── test/         # 测试文件
  └── go.mod        # 模块依赖
  ```
- **函数设计**：
    - 函数单一职责，参数不超过 5 个。
    - 使用 `return err` 显式返回错误，**不忽略错误**。
    - 延迟释放资源（如 `defer file.Close()`）。

### **3. 错误处理**
- **错误传递**：
  ```go
  func DoSomething() error {
      if err := validate(); err != nil {
          return fmt.Errorf("validate failed: %w", err)
      }
      // ...
      return nil
  }
  ```
- **自定义错误类型**：
  ```go
  type MyError struct {
      Code    int    `json:"code"`
      Message string `json:"message"`
  }
  func (e *MyError) Error() string { return e.Message }
  ```
- **全局错误处理**：
    - 使用 Gin 中间件统一处理 HTTP 错误：
  ```go
  func RecoveryMiddleware() gin.HandlerFunc {
      return func(c *gin.Context) {
          defer func() {
              if r := recover(); r != nil {
                  c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
              }
          }()
          c.Next()
      }
  }
  ```

### **4. 依赖注入**
- **使用依赖注入框架**：
  ```go
  // 定义接口
  type UserRepository interface {
      FindByID(ctx context.Context, id int) (*User, error)
  }
  
  // 实现依赖注入（如使用 wire）
  func InitializeDependencies() (*UserRepository, func()) {
      repo := NewGORMUserRepository()
      return repo, func() { /* 释放资源 */ }
  }
  ```

### **5. HTTP 处理**
- **路由设计**：
  ```go
  router := gin.Default()
  v1 := router.Group("/api/v1")
  {
      v1.POST("/users", CreateUserHandler)
      v1.GET("/users/:id", GetUserHandler)
  }
  ```
- **响应格式**：
  ```go
  type APIResponse struct {
      Status  string      `json:"status"`
      Message string      `json:"message"`
      Data    interface{} `json:"data,omitempty"`
  }
  ```
- **中间件**：
  ```go
  func LoggerMiddleware() gin.HandlerFunc {
      return func(c *gin.Context) {
          start := time.Now()
          c.Next()
          duration := time.Since(start)
          zap.L().Info("request", zap.String("path", c.Request.URL.Path), zap.Duration("duration", duration))
      }
  }
  ```

### **6. 数据库操作**
- **GORM 使用规范**：
  ```go
  type User struct {
      gorm.Model
      Name  string `gorm:"unique"`
      Email string
  }
  
  func (repo *GORMUserRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
      var user User
      if err := repo.DB.Where("email = ?", email).First(&user).Error; err != nil {
          return nil, err
      }
      return &user, nil
  }
  ```
- **SQL 注入防护**：
    - 使用参数化查询（如 `WHERE id = ?`）。
    - 避免拼接 SQL 字符串。

### **7. 并发处理**
- **Goroutine 安全**：
  ```go
  var mu sync.Mutex
  var count int

  func Increment() {
      mu.Lock()
      defer mu.Unlock()
      count++
  }
  ```
- **Channel 通信**：
  ```go
  func Worker(id int, jobs <-chan int, results chan<- int) {
      for j := range jobs {
          fmt.Printf("Worker %d processing jobs %d\n", id, j)
          results <- j * 2
      }
  }
  ```

### **8. 安全规范**
- **输入验证**：
  ```go
  type CreateUserRequest struct {
      Name  string `json:"name" validate:"required,min=2"`
      Email string `json:"email" validate:"required,email"`
  }
  ```
- **环境变量**：
  ```go
  const (
      DBHost     = os.Getenv("DB_HOST")
      DBUser     = os.Getenv("DB_USER")
      DBPassword = os.Getenv("DB_PASSWORD")
  )
  ```

### **9. 测试规范**
- **单元测试**：
  ```go
  func TestUserService_CreateUser(t *testing.T) {
      // 使用 mock 对象模拟依赖
      mockRepo := &MockUserRepository{}
      service := NewUserService(mockRepo)
      _, err := service.CreateUser(context.Background(), "test@example.com")
      assert.NoError(t, err)
  }
  ```

### **10. 日志规范**
- 使用标准库slog进行日志记录

---

### **11. 注释规范**
1. **按照示例的格式生成中英双语注释**
2. 不需要输出你的思考逻辑,不允许修改任何代码
3. 函数和方法的注释不要使用块注释，方法也用函数名开头不要使用这个方法的结构体开头
4. 当没有参数时注释也不需要有参数，没有返回值时注释也不需要有返回值
5. 内部的业务代码需要增加注释，特别是关键逻辑和算法部分，保持代码的可读性和可维护性
#### 示例
```go
// GetJobByIDOrAlias retrieves a jobs by identifier, trying ID lookup first, then alias lookup.
// 通过标识符获取任务，优先尝试ID查找，其次尝试别名查找。
//
// Parameters:
//	identifier - The jobs ID or alias to look up / 要查找的任务ID或别名
// Returns:
//	gocron.Job - The found jobs, or nil if not found / 找到的任务，未找到时为nil
//	error      - Error if jobs not found / 如果未找到任务的错误
func (s *Scheduler) GetJobByIDOrAlias(identifier string) (gocron.Job, error) {
    // Try ID lookup first
    // 优先尝试ID查找
    if jobID, err := s.GetJobByID(identifier); err == nil {
    return jobID, nil
    }
    
    // Fallback to alias lookup if alias feature is enabled
    // 如果启用别名功能，回退到别名查找
    if s.Enable(AliasOptionName) {
    if jobID, exists := s.aliasMap[identifier]; exists {
            return s.GetJobByAlias(jobID)
        }
    }
    return nil, fmt.Errorf("chrono:jobs with identifier %s not found", identifier)
}

```


## **示例：全局错误处理**
```go
// 定义全局错误响应结构
type APIResponse struct {
    Status  string      `json:"status"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}

// 中间件统一处理错误
func ErrorHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()
        if len(c.Errors) > 0 {
            lastError := c.Errors.Last()
            status := lastError.StatusCode
            message := lastError.Err.Error()
            c.AbortWithStatusJSON(status, APIResponse{
                Status:  "error",
                Message: message,
            })
        }
    }
}
```

---

## **备注**
- **代码评审**：每次提交必须通过代码评审，确保规范遵守。
- **性能优化**：使用 `pprof` 分析内存/CPU 使用，避免内存泄漏。
- **文档**：关键接口需用 `godoc` 注释，API 文档使用 Swagger 生成。
- **CI/CD**：代码提交后自动触发测试、构建和部署流程。
```