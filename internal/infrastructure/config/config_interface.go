package config

// Config interface defines methods for loading and accessing configuration settings.
type Config interface {
	LoadConfig() error
	GetDatabaseConfig() *DatabaseConfig
	GetAPIServerConfig() *APIServerConfig
	GetLogConfig() *LogConfig
}

// DatabaseConfig 資料庫設定值
type DatabaseConfig struct {
	URLSuffix string // 資料庫連接字串後綴
	Account   string // 資料庫帳號
	Password  string // 資料庫密碼
	Host      string // 資料庫主機
	Port      int    // 資料庫端口
	Name      string // 資料庫名稱
}

// APIServerConfig API 服務設定值
type APIServerConfig struct {
	ServerPort int // API 服務端口
}

// LogConfig 日誌設定值
type LogConfig struct {
	Level string // 日誌級別
}
