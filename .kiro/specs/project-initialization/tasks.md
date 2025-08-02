d# 實作計劃

## 概述
基於測試驅動開發 (TDD) 方法的 Go Todolist Service 專案初始化實作計劃。此計劃將逐步建立完整的專案基礎架構，每個任務都包含測試優先的方法，確保代碼品質和可維護性。

## 實作任務清單

### 階段 1: 專案基礎建立

- [x] **1. 建立專案結構和依賴管理**
  - 建立符合 Clean Architecture 的目錄結構：`internal/`, `pkg/`, `cmd/`, `configs/`, `tests/`
  - 設定 `go.mod` 文件並初始化必要的第三方依賴：Gin, GORM, Viper, Zerolog
  - 建立 `.gitignore` 文件排除不必要的檔案
  - 設定測試框架和基本測試目錄結構
  - _需求映射: REQ-1_

- [x] **2. 建立基礎設定架構**
  - 預設系統已讀取.env所設定的環境變數
  - 在 `internal/config/` 建立配置管理結構體和介面定義
  - 先撰寫配置管理的測試案例：測試環境變數讀取、YAML 檔案解析、預設值設定
  - 實作 `Config`, `DatabaseConfig`, `ServerConfig`, `LogConfig` 結構體
  - 確保所有測試通過
  - _需求映射: REQ-2.1, REQ-2.2_

### 階段 2: 配置管理系統

- [x] **3. 實作 Viper 配置管理服務**
  - 先撰寫 `Config` 介面的測試案例
  - 實作 `Config` 的具體實現，支援 YAML 和環境變數
  - 建立 `configs/config.yaml`
  - 實作環境變數自動綁定和覆蓋機制
  - 驗證所有配置載入測試通過
  - _需求映射: REQ-2.1, REQ-2.2_

### 階段 3: 日誌系統

- [x] **5. 實作 Zerolog 日誌系統**
  - 在 `internal/infrastructure/logger/` 建立日誌服務結構
  - 先撰寫日誌服務的測試：測試不同日誌等級、格式化、輸出目標
  - 實作 `Logger` 介面和具體實現 `LoggerImpl`
  - 支援 JSON 格式、可配置的日誌等級
  - 整合配置系統，讓日誌設定可透過配置檔案和環境變數控制
  - _需求映射: REQ-3.1_

### 階段 4: 資料庫連接

- [x] **7. 實作資料庫連接管理**
  - 在 `internal/infrastructure/database/` 建立資料庫服務結構
  - 撰寫 DSN 產生和配置驗證的單元測試（避免測試外部依賴）
  - 實作 `Database` 介面，支援 GORM MySQL 連接
  - 整合配置系統，從環境變數讀取資料庫設定
  - 建立 `docker-compose.yml` 提供 MySQL 外部依賴
  - 確保 `.env` 文件包含資料庫相關環境變數
  - 更新 main.go 整合資料庫服務
  - _需求映射: REQ-4.1, REQ-4.3_

- [x] **8. 建立資料庫遷移基礎**
  - 在 Database 介面中加入 Migrate 方法
  - 實作基礎的遷移管理功能，使用 GORM AutoMigrate
  - 建立遷移介面和基礎結構，為未來的業務實體做準備
  - 支援多個模型的批量遷移
  - _需求映射: REQ-4.2_

### 階段 5: HTTP 服務框架

- [x] **9. 實作 Gin HTTP 服務器基礎**
  - 在 `internal/infrastructure/server/` 建立服務器管理結構
  - 實作 `ServerService` 介面和基礎的 HTTP 服務器功能
  - 支援從環境變數讀取埠號設定
  - 實作優雅關機機制（context-based shutdown）
  - 整合 Router 和配置系統
  - 更新 main.go 完整的服務器啟動流程
  - _需求映射: REQ-5.1, REQ-5.2_

- [x] **10. 建立路由群組和中間件架構**
  - 撰寫路由群組設定的測試（使用 testify suite）
  - 在 `internal/infrastructure/router/` 實作路由管理
  - 在 `internal/delivery/http/handler/` 實作路由 handler
  - 在 `internal/delivery/http/middleware/` 實作中間件（符合 Clean Architecture）
  - 建立 `/api/v1` 路由群組結構
  - 整合錯誤處理中間件和 CORS 中間件
  - 實作基礎的 `/health` 和 `/version` 端點
  - _需求映射: REQ-5.3_

### 階段 6: 應用程式整合

- [x] **11. 實作 main.go 應用程式入口**
  - 在根目錄建立 `main.go`，整合所有系統組件
  - 實作依賴注入和初始化順序：配置 → 日誌 → 資料庫 → 路由 → 服務器
  - 添加應用程式版本資訊（在 `/version` 端點）
  - 確保所有組件能正確協作和錯誤處理
  - 實作優雅關機和錯誤記錄機制
  - _需求映射: REQ-1_

- [x] **12. 建立錯誤處理和中間件系統**
  - 在 `internal/delivery/http/middleware/` 實作錯誤處理中間件（符合 Clean Architecture）
  - 實作 `AppError` 自定義錯誤類型，支援 code、message、details 欄位
  - 支援結構化錯誤回應和適當的 HTTP 狀態碼處理
  - 實作 CORS 中間件和錯誤處理中間件
  - 整合到路由系統，提供全域錯誤處理
  - _需求映射: REQ-5.1_

### 階段 7: 文件和最終驗證

**註：整合測試和端到端測試已移至業務邏輯開發階段，專案初始化階段專注於基礎架構完整性**

- [x] **13. 建立開發環境配置文件**
  - ✅ 建立完整的 `internal/infrastructure/config/config.yaml`
  - ✅ 更新 `.env.example` 包含所有必要的環境變數（服務器、資料庫、日誌）
  - ✅ 撰寫基本的 README 使用說明（技術棧、功能規劃、資料庫結構）
  - ✅ 提供 docker-compose.yml 開發環境支援
  - _需求映射: REQ-2.2_

- [x] **14. 最終整合驗證和優化**
  - ✅ 執行完整的測試套件，所有測試通過（config, database, logger, router）
  - ✅ 驗證應用程式編譯成功（go build 無錯誤）
  - ✅ 檢查代碼品質和 Go 官方規範遵循（gofmt, go vet 通過）
  - ✅ 依賴管理優化完成（go mod tidy）
  - ✅ 專案已準備進入下一階段的業務邏輯開發
  - _需求映射: 所有需求的最終驗證_

## 任務執行指導原則

### 測試驅動開發 (TDD) 流程
1. **紅色階段**: 先寫測試，測試必然失敗
2. **綠色階段**: 寫最少的代碼讓測試通過
3. **重構階段**: 改進代碼品質，保持測試通過

### 依賴順序
- 配置系統 → 日誌系統 → 資料庫連接 → HTTP 服務器 → 應用程式整合
- 每個階段都必須先完成測試，再進行實作
- 後續任務可以依賴前面任務的輸出

### 代碼品質要求
- 遵循 Go 官方代碼規範和最佳實踐
- 保持 Clean Architecture 邊界清晰
- 每個公開函數都應有適當的測試覆蓋
- 錯誤處理要完整和一致

### 測試策略
- **單元測試**: 測試個別函數和方法的邏輯
- **整合測試**: 測試組件間的交互和配置
- **端到端測試**: 測試完整的應用程式流程
- **測試覆蓋率**: 目標 ≥80% 代碼覆蓋率