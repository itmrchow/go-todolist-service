# Docker 開發環境設定

本專案使用 Docker Compose 來提供本地開發所需的外部依賴服務。

## 前置需求

- Docker
- Docker Compose

## 快速開始

### 1. 環境變數設定

複製環境變數範例檔案：
```bash
cp .env.example .env
```

編輯 `.env` 文件，設定適當的資料庫密碼和其他配置。

### 2. 啟動服務

啟動 MySQL 資料庫：
```bash
docker-compose up -d
```

查看服務狀態：
```bash
docker-compose ps
```

### 3. 停止服務

停止所有服務：
```bash
docker-compose down
```

停止服務並刪除資料：
```bash
docker-compose down -v
```

## 服務說明

### MySQL 資料庫

- **容器名稱**: `todolist-mysql`
- **連接埠**: `3306` (可透過 `DB_PORT` 環境變數修改)
- **預設資料庫**: `todolist_db`
- **資料持久化**: 使用 Docker volume `mysql_data`
- **健康檢查**: 自動檢查資料庫連接狀態

### 連接資料庫

應用程式會自動讀取 `.env` 文件中的環境變數來連接資料庫。

如需直接連接資料庫進行除錯：
```bash
docker exec -it todolist-mysql mysql -u root -p
```

## 環境變數說明

| 變數名稱 | 說明 | 預設值 |
|---------|------|--------|
| `DB_HOST` | 資料庫主機 | `localhost` |
| `DB_PORT` | 資料庫連接埠 | `3306` |
| `DB_NAME` | 資料庫名稱 | `todolist_db` |
| `DB_ACCOUNT` | 資料庫用戶名 | `root` |
| `DB_PASSWORD` | 資料庫密碼 | 需要設定 |

## 注意事項

- 首次啟動時，MySQL 會執行 `scripts/mysql/init/` 目錄下的初始化腳本
- 資料會持久化在 Docker volume 中，除非使用 `-v` 參數刪除
- 如果遇到連接問題，請確認防火牆設定和連接埠是否被占用