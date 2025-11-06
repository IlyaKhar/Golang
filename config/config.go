package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Port        int
	PprofPort   int
	SQLiteDSN   string
	LogLevel    string
	JWTSecret   string
	JWTIssuer   string
	JWTAudience string
}

func Load() (Config, error) {
	v := viper.New()

	// .env (не обязателен)
	v.SetConfigFile(".env")
	_ = v.ReadInConfig()

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// defaults
	v.SetDefault("PORT", 8080)
	v.SetDefault("PPROF_PORT", 6060)
	v.SetDefault("SQLITE_DSN", "file:./app.db?cache=shared&mode=rwc")
	v.SetDefault("LOG_LEVEL", "dev")
	v.SetDefault("JWT_SECRET", "super-secret")
	v.SetDefault("JWT_ISSUER", "day2")
	v.SetDefault("JWT_AUDIENCE", "day2-clients")

	cfg := Config{
		Port:        v.GetInt("PORT"),
		PprofPort:   v.GetInt("PPROF_PORT"),
		SQLiteDSN:   v.GetString("SQLITE_DSN"),
		LogLevel:    v.GetString("LOG_LEVEL"),
		JWTSecret:   v.GetString("JWT_SECRET"),
		JWTIssuer:   v.GetString("JWT_ISSUER"),
		JWTAudience: v.GetString("JWT_AUDIENCE"),
	}

	if cfg.Port <= 0 || cfg.Port > 65535 {
		return Config{}, fmt.Errorf("invalid PORT: %d", cfg.Port)
	}
	if cfg.PprofPort <= 0 || cfg.PprofPort > 65535 {
		return Config{}, fmt.Errorf("invalid PPROF_PORT: %d", cfg.PprofPort)
	}
	if cfg.JWTSecret == "" {
		return Config{}, fmt.Errorf("JWT_SECRET is required")
	}

	return cfg, nil
}
