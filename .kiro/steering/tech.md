# Technology Stack

## Architecture
**Clean Architecture** - 採用分層架構模式，確保業務邏輯與外部依賴的分離，提高代碼的可測試性和可維護性。

### Architecture Layers (預期實作)
- **Entities**: 核心業務實體 (Task)
- **Use Cases**: 業務邏輯層 (Task CRUD operations)
- **Interface Adapters**: 控制器和資料庫適配器
- **Frameworks & Drivers**: Web 框架和資料庫驅動

## Backend Technology

### Programming Language
- **Golang** - 主要開發語言，提供高效能和併發處理能力

### Web Framework
- **Gin** - 輕量級 HTTP web 框架
  - 高效能路由處理
  - 中間件支援
  - JSON 綁定和驗證

### Database
- **MySQL** - 關聯式資料庫
  - 可靠的 ACID 特性
  - 成熟的生態系統
  - 良好的效能表現

### ORM
- **GORM** - Golang ORM 函式庫
  - 自動遷移功能
  - 關聯處理
  - 查詢建構器

### Logging
- **Zerolog** - 結構化日誌函式庫
  - 高效能日誌處理
  - JSON 格式輸出
  - 日誌等級控制

### Configuration Management
- **Viper** - 配置管理函式庫
  - 支援多種配置格式 (JSON, YAML, TOML, etc.)
  - 環境變數自動綁定
  - 配置熱重載功能
  - 配置驗證和預設值處理

## Development Environment

### Required Tools
- **Go 1.19+** - Golang 執行環境
- **MySQL 8.0+** - 資料庫服務
- **Git** - 版本控制

### Common Commands (預期)
```bash
# 啟動開發服務器
go run main.go

# 執行測試
go test ./...

# 建置專案
go build -o bin/todolist-service

# 安裝依賴
go mod tidy

# 資料庫遷移
go run cmd/migrate/main.go
```

## Configuration Management (使用 Viper)

### Configuration Files
- `config.yaml` - 主要配置文件
- `config.development.yaml` - 開發環境配置
- `config.production.yaml` - 生產環境配置
- `.env` - 環境變數文件 (本地開發用)

### Environment Variables

#### Database Configuration
- `DB_HOST` - 資料庫主機地址 (default: localhost)
- `DB_PORT` - 資料庫埠號 (default: 3306)
- `DB_NAME` - 資料庫名稱
- `DB_USER` - 資料庫使用者名稱
- `DB_PASSWORD` - 資料庫密碼

#### Application Configuration
- `APP_PORT` - 應用程式埠號 (default: 8080)
- `APP_ENV` - 執行環境 (development/staging/production)
- `CONFIG_PATH` - 配置文件路徑 (default: ./configs)

#### Logging Configuration
- `LOG_LEVEL` - 日誌等級 (debug/info/warn/error)
- `LOG_FORMAT` - 日誌格式 (json/text)
- `LOG_OUTPUT` - 日誌輸出 (stdout/file)

## Port Configuration

### Standard Ports
- **8080** - 主要 API 服務埠
- **3306** - MySQL 資料庫埠

## Data Schema

### Task Entity
| Field | Type | Primary Key | Required | Description |
|-------|------|-------------|----------|-------------|
| id | int | ✓ | ✓ | 任務唯一識別碼 |
| title | string | ✗ | ✓ | 任務標題 |
| description | string | ✗ | ✓ | 任務描述 |
| status | int | ✗ | ✓ | 任務狀態 (1:todo, 2:doing, 3:done) |
| created_at | time | ✗ | ✓ | 建立時間 |
| updated_at | time | ✗ | ✓ | 更新時間 |
| deleted_at | time | ✗ | ✓ | 軟刪除時間 |

## Development Principles
- **規格驅動開發**: 遵循 spec-driven development 流程
- **測試驅動開發**: 優先撰寫測試案例
- **代碼品質**: 遵循 Go 官方代碼規範
- **文件驅動**: 完整的 API 文件和使用說明