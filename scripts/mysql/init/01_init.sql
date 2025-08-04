-- MySQL 初始化腳本
-- 這個腳本會在 MySQL 容器第一次啟動時執行

-- 設定時區
SET time_zone = '+08:00';

-- 建立用戶（如果需要自定義用戶）
-- CREATE USER IF NOT EXISTS 'todolist_user'@'%' IDENTIFIED BY 'your_password';
-- GRANT ALL PRIVILEGES ON todolist_db.* TO 'todolist_user'@'%';

-- 確保資料庫使用 utf8mb4 字符集
ALTER DATABASE todolist_db CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

-- 刷新權限
FLUSH PRIVILEGES;

-- 顯示資料庫資訊
SELECT 'MySQL initialization completed' AS status;
SHOW DATABASES;