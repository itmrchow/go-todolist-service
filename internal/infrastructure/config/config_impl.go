package config

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

var _ Config = &ConfigImpl{}

// ConfigImpl implements the Config interface.
type ConfigImpl struct {
}

func (c *ConfigImpl) LoadConfig() error {
	viper.AutomaticEnv()
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	
	// 簡化路徑設定，優先從根目錄載入
	viper.AddConfigPath(".")                                    // 根目錄（main 執行）
	viper.AddConfigPath("../../../")                           // 測試執行時的相對路徑
	viper.AddConfigPath("./internal/infrastructure/config")     // 備用路徑

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	log.Info().Str("module", "config").Msgf("config init success")
	return nil
}

func (c *ConfigImpl) GetDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		URLSuffix: viper.GetString("DB_URL_SUFFIX"),
		Account:   viper.GetString("DB_ACCOUNT"),
		Password:  viper.GetString("DB_PASSWORD"),
		Host:      viper.GetString("DB_HOST"),
		Port:      viper.GetInt("DB_PORT"),
		Name:      viper.GetString("DB_NAME"),
	}
}

func (c *ConfigImpl) GetAPIServerConfig() *APIServerConfig {
	return &APIServerConfig{
		ServerPort: viper.GetInt("SERVER_PORT"),
	}
}

func (c *ConfigImpl) GetLogConfig() *LogConfig {
	return &LogConfig{
		Level: viper.GetString("LOG_LEVEL"),
	}
}
