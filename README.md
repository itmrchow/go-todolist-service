# go-todolist-service
todolist 基於 golang 與 clean architecture 的實作.

## 規格驅動開發
[參考](https://github.com/gotalab/claude-code-spec)

## tech stack
### language
- golang
### database
- mysql
### packages
- api: gin
- log: zerolog
- orm: gorm
- config: viper
- test: testify
- mock: 

## features
### task
- task usecase
  - 新增task
  - 查詢task
  - 更新task狀態
  - 更新task內容
  - 刪除task

- schema
| name        | type   | pk  | 必填 | description      |
| ----------- | ------ | --- | ---- | ---------------- |
| id          | int    | Y   | Y    | task id          |
| title       | string | N   | Y    | task title       |
| description | string | N   | Y    | task description |
| status      | int    | N   | Y    | task status 1:todo , 2:doing, 3:done |
| created_at  | time   | N   | Y    | task create time |
| updated_at  | time   | N   | Y    | task update time |
| deleted_at  | time   | N   | Y    | task update time |
