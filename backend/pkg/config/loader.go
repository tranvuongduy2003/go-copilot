package config

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

func Load(configPath string) (*Config, error) {
	v := viper.New()

	setDefaults(v)
	bindEnvVars(v)

	if configPath != "" {
		v.SetConfigFile(configPath)
		if err := v.ReadInConfig(); err != nil {
			return nil, &ConfigError{Op: "read config file", Err: err}
		}
	} else {
		v.SetConfigName("config")
		v.SetConfigType("yaml")
		v.AddConfigPath(".")
		v.AddConfigPath("./config")
		v.AddConfigPath("/etc/app")

		_ = v.ReadInConfig()
	}

	env := v.GetString("app.env")
	if env != "" {
		v.SetConfigName("config." + env)
		_ = v.MergeInConfig()
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, &ConfigError{Op: "unmarshal config", Err: err}
	}

	if err := cfg.Validate(); err != nil {
		return nil, &ConfigError{Op: "validate config", Err: err}
	}

	return &cfg, nil
}

func LoadFromEnv() (*Config, error) {
	return Load("")
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("app.name", "go-copilot")
	v.SetDefault("app.env", EnvDevelopment)
	v.SetDefault("app.debug", false)

	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.read_timeout", 30*time.Second)
	v.SetDefault("server.write_timeout", 30*time.Second)
	v.SetDefault("server.idle_timeout", 60*time.Second)

	v.SetDefault("db.host", "localhost")
	v.SetDefault("db.port", 5432)
	v.SetDefault("db.user", "postgres")
	v.SetDefault("db.password", "")
	v.SetDefault("db.name", "app")
	v.SetDefault("db.ssl_mode", "disable")
	v.SetDefault("db.max_open_conns", 25)
	v.SetDefault("db.max_idle_conns", 5)
	v.SetDefault("db.conn_max_lifetime", 5*time.Minute)

	v.SetDefault("redis.host", "localhost")
	v.SetDefault("redis.port", 6379)
	v.SetDefault("redis.password", "")
	v.SetDefault("redis.db", 0)

	v.SetDefault("jwt.secret", "")
	v.SetDefault("jwt.access_token_ttl", 15*time.Minute)
	v.SetDefault("jwt.refresh_token_ttl", 7*24*time.Hour)

	v.SetDefault("log.level", "info")
	v.SetDefault("log.format", "json")

	v.SetDefault("cors.allowed_origins", []string{"*"})
	v.SetDefault("cors.allowed_methods", []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"})
	v.SetDefault("cors.allowed_headers", []string{"Authorization", "Content-Type", "X-Request-ID"})
	v.SetDefault("cors.max_age", 86400)
}

func bindEnvVars(v *viper.Viper) {
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	envBindings := map[string]string{
		"app.name":  "APP_NAME",
		"app.env":   "APP_ENV",
		"app.debug": "APP_DEBUG",

		"server.host":          "SERVER_HOST",
		"server.port":          "SERVER_PORT",
		"server.read_timeout":  "SERVER_READ_TIMEOUT",
		"server.write_timeout": "SERVER_WRITE_TIMEOUT",
		"server.idle_timeout":  "SERVER_IDLE_TIMEOUT",

		"db.host":              "DB_HOST",
		"db.port":              "DB_PORT",
		"db.user":              "DB_USER",
		"db.password":          "DB_PASSWORD",
		"db.name":              "DB_NAME",
		"db.ssl_mode":          "DB_SSL_MODE",
		"db.max_open_conns":    "DB_MAX_OPEN_CONNS",
		"db.max_idle_conns":    "DB_MAX_IDLE_CONNS",
		"db.conn_max_lifetime": "DB_CONN_MAX_LIFETIME",

		"redis.host":     "REDIS_HOST",
		"redis.port":     "REDIS_PORT",
		"redis.password": "REDIS_PASSWORD",
		"redis.db":       "REDIS_DB",

		"jwt.secret":            "JWT_SECRET",
		"jwt.access_token_ttl":  "JWT_ACCESS_TOKEN_TTL",
		"jwt.refresh_token_ttl": "JWT_REFRESH_TOKEN_TTL",

		"log.level":  "LOG_LEVEL",
		"log.format": "LOG_FORMAT",

		"cors.allowed_origins": "CORS_ALLOWED_ORIGINS",
		"cors.allowed_methods": "CORS_ALLOWED_METHODS",
		"cors.allowed_headers": "CORS_ALLOWED_HEADERS",
		"cors.max_age":         "CORS_MAX_AGE",
	}

	for key, envVar := range envBindings {
		_ = v.BindEnv(key, envVar)
	}
}

func MustLoad(configPath string) *Config {
	cfg, err := Load(configPath)
	if err != nil {
		panic("failed to load config: " + err.Error())
	}
	return cfg
}

type LoadResult struct {
	Config      *Config
	HotReloader *HotReloader
}

func LoadWithHotReload(configPath string) (*LoadResult, error) {
	v := viper.New()

	setDefaults(v)
	bindEnvVars(v)

	hasConfigFile := false

	if configPath != "" {
		v.SetConfigFile(configPath)
		if err := v.ReadInConfig(); err != nil {
			return nil, &ConfigError{Op: "read config file", Err: err}
		}
		hasConfigFile = true
	} else {
		v.SetConfigName("config")
		v.SetConfigType("yaml")
		v.AddConfigPath(".")
		v.AddConfigPath("./config")
		v.AddConfigPath("/etc/app")

		if err := v.ReadInConfig(); err == nil {
			hasConfigFile = true
		}
	}

	env := v.GetString("app.env")
	if env != "" {
		v.SetConfigName("config." + env)
		_ = v.MergeInConfig()
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, &ConfigError{Op: "unmarshal config", Err: err}
	}

	if err := cfg.Validate(); err != nil {
		return nil, &ConfigError{Op: "validate config", Err: err}
	}

	result := &LoadResult{
		Config: &cfg,
	}

	if hasConfigFile {
		result.HotReloader = NewHotReloader(v)
	}

	return result, nil
}

func MustLoadWithHotReload(configPath string) *LoadResult {
	result, err := LoadWithHotReload(configPath)
	if err != nil {
		panic("failed to load config: " + err.Error())
	}
	return result
}
