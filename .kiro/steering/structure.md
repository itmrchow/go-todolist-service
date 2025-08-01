# Project Structure

## Current State
專案目前處於初始階段，僅包含基本文件：
- `README.md` - 專案說明文件
- `CLAUDE.md` - Claude 開發規範
- `.kiro/` - 規格驅動開發配置

## Planned Directory Organization

### Root Directory Structure
```
go-todolist-service/
├── README.md                 # 專案說明文件
├── CLAUDE.md                 # Claude 開發規範
├── go.mod                    # Go modules 依賴管理
├── go.sum                    # 依賴版本鎖定
├── main.go                   # 應用程式入口點
├── .env.example              # 環境變數範例
├── .gitignore                # Git 忽略文件規則
├── docker-compose.yml        # Docker 容器配置
├── Dockerfile                # Docker 映像檔定義
│
├── .kiro/                    # 規格驅動開發
│   ├── steering/             # 專案指導文件
│   └── specs/                # 功能規格文件
│
├── cmd/                      # 應用程式入口點
│   ├── server/               # HTTP 服務器
│   │   └── main.go
│   └── migrate/              # 資料庫遷移工具
│       └── main.go
│
├── internal/                 # 私有應用程式代碼
│   ├── entity/               # 業務實體
│   │   └── task.go
│   ├── usecase/              # 業務邏輯
│   │   ├── interface.go
│   │   └── task.go
│   ├── repository/           # 資料存取層
│   │   ├── interface.go
│   │   └── mysql/
│   │       └── task.go
│   ├── delivery/             # 交付層
│   │   └── http/
│   │       ├── handler/
│   │       │   └── task.go
│   │       ├── middleware/
│   │       │   └── cors.go
│   │       └── router.go
│   └── config/               # 配置管理
│       └── config.go
│
├── pkg/                      # 可重用的函式庫
│   ├── database/             # 資料庫連接
│   │   └── mysql.go
│   ├── logger/               # 日誌配置
│   │   └── zerolog.go
│   └── validator/            # 資料驗證
│       └── validator.go
│
├── tests/                    # 測試文件
│   ├── integration/          # 整合測試
│   ├── unit/                 # 單元測試
│   └── mocks/                # 測試模擬
│
├── docs/                     # 文件
│   ├── api/                  # API 文件
│   └── architecture/         # 架構文件
│
└── scripts/                  # 構建和部署腳本
    ├── build.sh
    └── deploy.sh
```

## Code Organization Patterns

### Clean Architecture Layers

#### 1. Entities (`internal/entity/`)
核心業務實體，不依賴於任何外部層
```go
// internal/entity/task.go
type Task struct {
    ID          int       `json:"id"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
    Status      int       `json:"status"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
    DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}
```

#### 2. Use Cases (`internal/usecase/`)
業務邏輯層，定義應用程式的用例
```go
// internal/usecase/interface.go
type TaskUseCase interface {
    Create(task *entity.Task) error
    GetByID(id int) (*entity.Task, error)
    UpdateStatus(id int, status int) error
    UpdateContent(id int, title, description string) error
    Delete(id int) error
}
```

#### 3. Repository (`internal/repository/`)
資料存取抽象層
```go
// internal/repository/interface.go
type TaskRepository interface {
    Create(task *entity.Task) error
    GetByID(id int) (*entity.Task, error)
    Update(task *entity.Task) error
    Delete(id int) error
}
```

#### 4. Delivery (`internal/delivery/`)
交付機制 (HTTP handlers, gRPC, etc.)

## File Naming Conventions

### General Rules
- **Package names**: 小寫，單數形式 (e.g., `task`, `user`)
- **File names**: 小寫，使用底線分隔 (e.g., `task_handler.go`)
- **Interface files**: `interface.go` 或 `interfaces.go`
- **Test files**: `*_test.go` 後綴

### Specific Patterns
- **Entities**: 單數名詞 (e.g., `task.go`, `user.go`)
- **Handlers**: `*_handler.go` (e.g., `task_handler.go`)
- **Repositories**: `*_repository.go` 或依資料庫類型分資料夾
- **Use Cases**: `*_usecase.go` 或 `*.go` 直接命名

## Import Organization

### Import Order
1. Standard library packages
2. Third-party packages
3. Internal packages (relative imports)

### Example
```go
import (
    // Standard library
    "context"
    "fmt"
    "net/http"
    
    // Third-party
    "github.com/gin-gonic/gin"
    "github.com/rs/zerolog"
    "gorm.io/gorm"
    
    // Internal
    "github.com/username/go-todolist-service/internal/entity"
    "github.com/username/go-todolist-service/internal/usecase"
)
```

## Key Architectural Principles

### 1. Dependency Inversion
- 高層模組不依賴低層模組
- 兩者都依賴於抽象
- 抽象不依賴於細節

### 2. Single Responsibility
- 每個包、類型、函數都有單一職責
- 避免過度複雜的結構

### 3. Interface Segregation
- 介面應該小而專一
- 客戶端不應該依賴它們不使用的介面

### 4. Clean Architecture Boundaries
- 業務邏輯不依賴於框架
- 資料庫和 Web 是外部細節
- 依賴方向始終向內指向業務邏輯

## Configuration Management
- **使用 Viper 進行配置管理**
- 支援多種配置格式 (YAML, JSON, TOML)
- 環境變數自動綁定和覆蓋
- 分環境配置文件 (development, staging, production)
- `.env` 文件用於本地開發
- 配置驗證和預設值設定
- 敏感資訊不提交到版本控制
- 配置熱重載支援 (開發環境)