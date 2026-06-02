package config

import (
	"log"
	"strings"
	"time"

	"github.com/spf13/viper"
)

var (
	cfg    *viper.Viper
	config *Config
)

type IConfig interface {
	GetString(key string) string
}

type Config struct {
	Http struct {
		Port            int
		ShutdownTimeOut time.Duration
	}

	Database struct {
		Host     string
		Port     int
		User     string
		Password string
		Name     string
	}

	Auth struct {
		Enabled      bool
		TOTPRequired bool `mapstructure:"totp_required"`
	}

	Firebase struct {
		ProjectID string `mapstructure:"project_id"`
	}

	Ollama struct {
		Model   string
		BaseURL string `mapstructure:"base_url"`
	}
}

func Get() *Config {
	if config != nil {
		return config
	}
	newConfig()
	config = &Config{}
	if err := cfg.Unmarshal(config); err != nil {
		log.Printf("Failed to unmarshal config: %v", err)
	}
	return config
}

func newConfig() IConfig {
	if cfg != nil {
		return cfg
	}
	cfg = viper.New()
	cfg.SetDefault("http.port", 8080)
	cfg.SetDefault("http.shutdownTimeOut", 10*time.Second)
	cfg.SetDefault("database.host", "localhost")
	cfg.SetDefault("database.port", 5432)
	cfg.SetDefault("database.user", "postgres")
	cfg.SetDefault("database.password", "password")
	cfg.SetDefault("database.name", "support_copilot")
	cfg.SetDefault("auth.enabled", false)
	cfg.SetDefault("auth.totp_required", false)
	cfg.SetDefault("firebase.project_id", "")
	cfg.SetDefault("ollama.model", "llama3.2")
	cfg.SetDefault("ollama.base_url", "http://localhost:11434")
	cfg.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	cfg.AutomaticEnv()
	cfg.SetConfigName("config")
	cfg.SetConfigType("yaml")
	cfg.AddConfigPath(".")
	cfg.AddConfigPath("./config")

	if err := cfg.ReadInConfig(); err != nil {
		log.Printf("Failed to read config: %v", err)
	}

	return cfg
}
