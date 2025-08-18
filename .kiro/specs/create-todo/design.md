# 技術設計文件

## 概述
1. 完成todo的新增功能 , client 可以透過API 來新增一筆todo資料
2. 保持clean arch與ddd開發方式 , usecase為CreateTodo
3. 保持TDD的開發模式 , 先完成測試再開發

## 組件和介面

### 1. Domain Layer (領域層)

#### Todo Entity
```go
// internal/domain/entity/todo.go
type Todo struct {
    ID          uint       `json:"id"`
    Title       string     `json:"title"`
    Description *string    `json:"description,omitempty"`
    Status      TodoStatus `json:"status"`
    DueDate     *time.Time `json:"due_date,omitempty"`
    CreatedAt   time.Time  `json:"created_at"`
    UpdatedAt   time.Time  `json:"updated_at"`
}

type TodoStatus string

const (
    StatusPending TodoStatus = "pending"
    StatusDoing   TodoStatus = "doing" 
    StatusDone    TodoStatus = "done"
)
```

#### Todo Repository Interface  
```go
// internal/domain/repository/todo_repository.go
type TodoRepository interface {
    Create(ctx context.Context, todo *entity.Todo) (*entity.Todo, error)
    GetByID(ctx context.Context, id uint) (*entity.Todo, error)
}
```

### 2. Use Case Layer (應用層)

#### CreateTodo UseCase Interface
```go
// internal/usecase/interface/create_todo_usecase.go
type CreateTodoUseCase interface {
    Execute(ctx context.Context, req CreateTodoRequest) (*CreateTodoResponse, error)
}

type CreateTodoRequest struct {
    Title       string     `json:"title" binding:"required"`
    Description *string    `json:"description,omitempty"`
    Status      *string    `json:"status,omitempty"`
    DueDate     *string    `json:"due_date,omitempty"`
}

type CreateTodoResponse struct {
    Todo *entity.Todo `json:"todo"`
}
```

#### CreateTodo UseCase Implementation
```go
// internal/usecase/create_todo_usecase.go
type CreateTodoUseCaseImpl struct {
    todoRepo repository.TodoRepository
    logger   logger.Logger
}

func NewCreateTodoUseCase(todoRepo repository.TodoRepository, logger logger.Logger) CreateTodoUseCase {
    return &CreateTodoUseCaseImpl{
        todoRepo: todoRepo,
        logger:   logger,
    }
}

func (uc *CreateTodoUseCaseImpl) Execute(ctx context.Context, req CreateTodoRequest) (*CreateTodoResponse, error) {
    // 1. 驗證和轉換輸入
    // 2. 建立 Todo entity
    // 3. 呼叫 repository 儲存
    // 4. 回傳結果
}
```

### 3. Infrastructure Layer (基礎設施層)

#### GORM Model for Auto Migration
```go
// internal/infrastructure/database/model/todo.go
type Todo struct {
    ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
    Title       string    `gorm:"type:varchar(80);not null;comment:Todo標題，最多20個中文字符" json:"title"`
    Description *string   `gorm:"type:text;comment:Todo描述，最多100個中文字符" json:"description"`
    Status      string    `gorm:"type:enum('pending','doing','done');default:'pending';not null;comment:Todo狀態;index" json:"status"`
    DueDate     *time.Time `gorm:"type:timestamp;null;comment:到期日期，UTC時間;index" json:"due_date"`
    CreatedAt   time.Time `gorm:"autoCreateTime;index" json:"created_at"`
    UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Todo) TableName() string {
    return "todos"
}
```

#### Todo Repository Implementation
```go
// internal/infrastructure/repository/todo_repository_impl.go
type TodoRepositoryImpl struct {
    db     database.Database
    logger logger.Logger
}

func NewTodoRepository(db database.Database, logger logger.Logger) repository.TodoRepository {
    return &TodoRepositoryImpl{
        db:     db,
        logger: logger,
    }
}

func (r *TodoRepositoryImpl) Create(ctx context.Context, todo *entity.Todo) (*entity.Todo, error) {
    // GORM 實作，包含 entity 和 model 之間的轉換
}
```

### 4. Delivery Layer (交付層)

#### CreateTodo Handler
```go
// internal/delivery/http/handler/todo_handler.go  
type TodoHandler struct {
    createTodoUseCase usecase_interface.CreateTodoUseCase
    logger           logger.Logger
}

func NewTodoHandler(createTodoUseCase usecase_interface.CreateTodoUseCase, logger logger.Logger) *TodoHandler {
    return &TodoHandler{
        createTodoUseCase: createTodoUseCase,
        logger:           logger,
    }
}

func (h *TodoHandler) CreateTodo(c *gin.Context) {
    // 1. 綁定和驗證請求
    // 2. 呼叫 UseCase，傳入 c.Request.Context()
    // 3. 處理回應和錯誤
}
```

### 5. API 契約

#### Request Format
```http
POST /api/v1/todos
Content-Type: application/json

{
    "title": "學習 Go 語言",
    "description": "完成 Go 語言基礎教程",
    "status": "pending",
    "due_date": "2024-12-31T23:59:59Z"
}
```

#### Response Format
```http
HTTP/1.1 201 Created
Content-Type: application/json

{
    "todo": {
        "id": 1,
        "title": "學習 Go 語言", 
        "description": "完成 Go 語言基礎教程",
        "status": "pending",
        "due_date": "2024-12-31T23:59:59Z",
        "created_at": "2024-01-01T10:00:00Z",
        "updated_at": "2024-01-01T10:00:00Z"
    }
}
```

## 實作計劃

### Phase 1: Domain Layer
1. 建立 Todo entity
2. 定義 TodoRepository interface
3. 撰寫 domain layer 的單元測試

### Phase 2: Use Case Layer  
1. 定義 CreateTodoUseCase interface (確保 Execute 方法包含 context.Context)
2. 實作 CreateTodoUseCaseImpl
3. 撰寫 use case 的單元測試

### Phase 3: Infrastructure Layer
1. 建立 GORM Model 和 Auto Migration
2. 實作 TodoRepositoryImpl
3. 撰寫 repository 的整合測試

### Phase 4: Delivery Layer
1. 實作 TodoHandler (在呼叫 UseCase 時傳入 c.Request.Context())
2. 更新 router 配置
3. 撰寫 handler 的單元測試和整合測試

### Phase 5: End-to-End Testing
1. 撰寫 API 整合測試
2. 執行完整的測試套件
3. 驗證所有需求