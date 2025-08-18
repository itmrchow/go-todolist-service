# UpdateTodo 功能實作計劃

## 概述
基於測試驅動開發 (TDD) 方法實作 UpdateTodo 功能。此計劃遵循 Clean Architecture 和 DDD 原則，逐步建立完整的 Todo 新增功能，確保代碼品質和可維護性。

## 開發順序

1. API router
2. repository
3. usecase

## 實作任務清單

### Phase 1: API Router 層實作
- [ ] 1.1 在 router 中新增 PUT /todos/:id 端點
- [ ] 1.2 定義 UpdateTodo 請求和回應結構
- [ ] 1.3 實作 UpdateTodo handler 方法 (先用 mock 或簡單實作)
- [ ] 1.4 加入請求驗證和錯誤處理
- [ ] 1.5 編寫 router 層的測試

### Phase 2: Repository 層完善
- [ ] 2.1 確認 repository Update 方法實作完整性
- [ ] 2.2 驗證 Update 方法的測試涵蓋度
- [ ] 2.3 如需要，補充額外的邊界情況測試
- [ ] 2.4 確保與 API 層的整合需求一致

### Phase 3: UseCase 層實作
- [ ] 3.1 定義 UpdateTodoRequest 和 UpdateTodoResponse 結構
- [ ] 3.2 在 TodoUseCase 接口中取消註釋 UpdateTodo 方法
- [ ] 3.3 在 todo_uc_impl.go 中實作 UpdateTodo 方法
- [ ] 3.4 編寫 UpdateTodo 的單元測試
- [ ] 3.5 整合 router 和 repository 層

### Phase 4: 整合測試與驗證
- [ ] 4.1 編寫端到端測試驗證完整流程
- [ ] 4.2 測試各種錯誤情境 (404, 400, 500)
- [ ] 4.3 驗證資料驗證規則
- [ ] 4.4 確認 UTC 時間格式處理正確

## 技術細節

### 驗證規則
- title: 必填，最大 20 個中文字符
- description: 可選，最大 100 個中文字符  
- status: 必須是 "pending", "doing", "done" 之一
- due_date: 可選，必須是未來時間，格式 RFC3339

### 錯誤處理
- 400: 驗證錯誤 (invalid request)
- 404: Todo 不存在 (not found)
- 500: 內部錯誤 (internal error)

