# CreateTodo 功能實作計劃

## 概述
基於測試驅動開發 (TDD) 方法實作 CreateTodo 功能。此計劃遵循 Clean Architecture 和 DDD 原則，逐步建立完整的 Todo 新增功能，確保代碼品質和可維護性。

## 實作任務清單

### 階段 1: Domain Layer (領域層) - Inside-Out 開發起點 ✅

- [x] **1. 建立 Todo Entity 和 Value Objects**
  - ✅ 在 `internal/domain/entity/` 建立 `todo.go`
  - ✅ 實作 `Todo` 結構體和 `TodoStatus` 列舉
  - ✅ 先撰寫 entity 的單元測試：測試結構體初始化、JSON 序列化
  - ✅ 加入軟刪除功能 (`DeletedAt`, `IsDeleted()`, `Delete()`, `Restore()`)
  - ✅ 撰寫完整的驗證測試：長度限制、空值處理、邊界測試
  - ✅ 確保所有 domain 測試通過 (8個測試案例全部通過)
  - _需求映射: title (≤20中文字), description (≤100中文字), status, due_date_

- [x] **2. 定義 Repository Interface**
  - ✅ 在 `internal/domain/repository/` 建立 `todo_repository.go`
  - ✅ 定義 `TodoRepository` 介面，包含完整的 CRUD 操作
  - ✅ 支援軟刪除、分頁查詢、過濾條件等進階功能
  - ✅ 確保方法簽名正確包含 `context.Context` 參數
  - ✅ 撰寫 interface 的文件註解
  - _需求映射: 資料持久化需求_

### 階段 2: Infrastructure Layer (基礎設施層) - 持久化實作 ✅

- [x] **3. 建立 GORM Model 和 Auto Migration**
  - ✅ 在 `internal/infrastructure/database/model/` 建立 `todo.go`
  - ✅ 實作 GORM model 結構體，包含正確的 gorm tags
  - ✅ 建立 entity 和 model 之間的轉換函數 (`EntityToModel`, `ModelToEntity`)
  - ✅ 撰寫轉換函數的單元測試 (`todo_test.go`)
  - [ ] 在資料庫服務中加入 Todo model 的 auto migration (待整合到 main.go)
  - _需求映射: 資料庫結構設計_

- [x] **4. 實作 TodoRepository Implementation**
  - ✅ 先撰寫 Repository 的整合測試：
    - ✅ 使用 SQLite 進行測試
    - ✅ 測試 Create 方法的正常流程
    - ✅ 測試資料庫連線錯誤處理
    - ✅ 測試重複 ID 處理
  - ✅ 在 `internal/infrastructure/repository/` 建立 `todo_repository_impl.go`
  - ✅ 實作 `TodoRepositoryImpl` 結構體
  - ✅ 實作 GORM 資料庫操作，包含 entity/model 轉換
  - ✅ 實作適當的錯誤處理和日誌記錄
  - ✅ 確保所有 Repository 測試通過 (`todo_repository_impl_test.go`)
  - _需求映射: 資料持久化實作_

### 階段 3: Use Case Layer (應用層) - 業務邏輯 ✅

- [x] **5. 實作 CreateTodo UseCase Interface**
  - ✅ 在 `internal/domain/usecase/` 建立 `todo_uc.go` (架構調整)
  - ✅ 定義 `TodoUseCase` 介面和 `CreateTodoRequest`/`CreateTodoResponse` 結構體
  - ✅ 確保 `CreateTodo` 方法正確包含 `context.Context` 參數
  - ✅ 定義完整的 request/response 結構，支援所有欄位
  - _需求映射: API 契約定義_

- [x] **6. 實作 CreateTodo UseCase Implementation**
  - ✅ 先撰寫 UseCase 的單元測試 (`todo_uc_impl_test.go`)：
    - ✅ 正常流程測試（成功建立 Todo）
    - ✅ 驗證錯誤測試（entity 建立失敗）
    - ✅ Repository 錯誤處理測試
  - ✅ 在 `internal/domain/usecase/` 建立 `todo_uc_impl.go`
  - ✅ 實作 `todoUseCaseImpl` 結構體
  - ✅ 實作完整的業務邏輯：輸入驗證、entity 建立、repository 呼叫
  - ✅ 使用 gomock 模擬 repository 依賴
  - ✅ 確保所有 UseCase 測試通過 (3個測試案例)
  - _需求映射: 業務邏輯驗證_

### 階段 4: Delivery Layer (交付層) - HTTP Handler

- [ ] **7. 實作 CreateTodo Handler**
  - 先撰寫 Handler 的單元測試：
    - 正常請求處理測試
    - JSON 綁定錯誤測試
    - 驗證錯誤回應測試
    - UseCase 錯誤處理測試
  - 在 `internal/delivery/http/handler/` 建立 `todo_handler.go`
  - 實作 `TodoHandler` 結構體和 `CreateTodo` 方法
  - 實作請求綁定、驗證、UseCase 呼叫、回應處理
  - 確保正確傳入 `c.Request.Context()` 到 UseCase
  - 實作適當的 HTTP 狀態碼和錯誤回應
  - 使用 testify/mock 模擬 UseCase 依賴
  - 確保所有 Handler 測試通過
  - _需求映射: HTTP API 實作_

- [ ] **8. 更新 Router 配置**
  - 在 `internal/infrastructure/router/router_impl.go` 的 `RegisterV1Routes` 中新增路由
  - 新增 `POST /api/v1/todos` 路由到 `CreateTodo` handler
  - 撰寫路由註冊的測試
  - 確保路由正確對應到 handler 方法
  - _需求映射: API 路由配置_

### 階段 5: System Integration (系統整合)

- [ ] **9. 建立依賴注入容器**
  - 在 `internal/di/` 建立 `container.go`
  - 實作 Todo 相關服務的依賴注入：
    - TodoRepository 實例化
    - TodoUseCase 實例化  
    - TodoHandler 實例化
  - 更新 main.go 或現有的 DI 系統
  - 🔥 **重要**: 在初始化過程中加入 Todo model 的 auto migration
  - 確保所有依賴正確注入和初始化
  - _需求映射: 系統整合_

### 階段 6: End-to-End Testing (端到端測試)

- [ ] **10. 建立 API 整合測試**
  - 在 `tests/integration/` 建立 `create_todo_test.go`
  - 撰寫完整的 API 整合測試：
    - 測試完整的 HTTP 請求/回應流程
    - 測試資料庫實際寫入和讀取
    - 測試各種錯誤情況的完整處理
  - 使用真實的資料庫環境（Docker 或測試資料庫）
  - 測試案例涵蓋 `.kiro/specs/create-todo/api.http` 中的所有情境
  - 確保所有整合測試通過
  - _需求映射: 完整功能驗證_

### 階段 7: 最終驗證和優化

- [ ] **11. 執行完整測試套件**
  - 執行所有單元測試：`go test ./internal/...`
  - 執行所有整合測試：`go test ./tests/...`
  - 確保測試覆蓋率 ≥80%
  - 檢查是否有測試失敗或不穩定的測試
  - _需求映射: 品質保證_

- [ ] **12. 代碼品質檢查**
  - 執行 `go fmt ./...` 確保代碼格式符合規範
  - 執行 `go vet ./...` 檢查代碼問題
  - 執行 `go mod tidy` 清理依賴
  - 檢查 import 循環依賴問題
  - 確保所有 public function 都有適當的註解
  - _需求映射: 代碼品質標準_

- [ ] **13. 功能驗證**
  - 使用 `.kiro/specs/create-todo/api.http` 手動測試所有案例
  - 驗證錯誤回應格式正確
  - 驗證成功回應包含所有必要欄位
  - 確認中文字符長度限制正確運作
  - 驗證日期格式處理正確
  - _需求映射: 功能需求完整驗證_

## 任務執行指導原則

### 測試驅動開發 (TDD) 流程
1. **紅色階段**: 先寫測試，測試必然失敗
2. **綠色階段**: 寫最少的代碼讓測試通過  
3. **重構階段**: 改進代碼品質，保持測試通過

### Clean Architecture 層級依賴
- Domain Layer (不依賴任何層級)
- Use Case Layer (依賴 Domain Layer)
- Infrastructure Layer (依賴 Domain Layer)
- Delivery Layer (依賴 Use Case Layer)

### 開發順序建議 (Inside-Out TDD)
1. **Domain → Infrastructure → Use Case → Delivery → Integration → E2E Testing**
2. 從核心業務邏輯開始，逐步向外層擴展
3. 每個階段都必須先完成測試，再進行實作
4. 使用 mock 隔離外部依賴進行單元測試
5. 最後階段進行完整的整合測試和端到端驗證

### 代碼品質要求
- 遵循 Go 官方代碼規範和最佳實踐
- 保持 Clean Architecture 邊界清晰，避免循環依賴
- 每個公開函數都應有適當的測試覆蓋
- 錯誤處理要完整和一致
- 使用 context.Context 進行請求追蹤和取消

### 測試策略
- **單元測試**: 測試個別函數和方法的邏輯，使用 mock 隔離依賴
- **整合測試**: 測試組件間的交互，使用真實的資料庫
- **端到端測試**: 測試完整的 HTTP API 流程
- **測試覆蓋率**: 目標 ≥80% 代碼覆蓋率

### 驗證規則實作
- **Title**: 必填，1-20 個中文字符，使用 `[]rune` 計算長度
- **Description**: 可選，最多 100 個中文字符，使用 `[]rune` 計算長度
- **Status**: 可選，預設 "pending"，只接受 "pending", "doing", "done"
- **DueDate**: 可選，必須是有效的 RFC3339 格式 UTC 時間

### 錯誤處理標準
- 使用自定義錯誤類型包含 code, message, details
- HTTP 狀態碼對應：400 (驗證錯誤), 500 (系統錯誤), 201 (成功建立)
- 所有錯誤都要記錄到日誌系統
- 錯誤回應格式必須一致且結構化