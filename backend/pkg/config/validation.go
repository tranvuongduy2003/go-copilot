package config

import (
	"errors"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return e.Field + ": " + e.Message
}

type ValidationErrors []ValidationError

func (e ValidationErrors) Error() string {
	if len(e) == 0 {
		return ""
	}

	var msgs []string
	for _, err := range e {
		msgs = append(msgs, err.Error())
	}
	return "configuration validation failed:\n  - " + strings.Join(msgs, "\n  - ")
}

func (c *Config) Validate() error {
	var errs ValidationErrors

	errs = append(errs, c.App.Validate()...)
	errs = append(errs, c.Server.Validate()...)
	errs = append(errs, c.Database.Validate()...)
	errs = append(errs, c.Redis.Validate()...)
	errs = append(errs, c.JWT.Validate(c.App.Env)...)
	errs = append(errs, c.Log.Validate()...)
	errs = append(errs, c.CORS.Validate()...)

	if len(errs) > 0 {
		return errs
	}
	return nil
}

func (c *AppConfig) Validate() ValidationErrors {
	var errs ValidationErrors

	if c.Name == "" {
		errs = append(errs, ValidationError{
			Field:   "app.name",
			Message: "application name is required",
		})
	}

	validEnvs := map[string]bool{
		EnvDevelopment: true,
		EnvStaging:     true,
		EnvProduction:  true,
	}
	if !validEnvs[c.Env] {
		errs = append(errs, ValidationError{
			Field:   "app.env",
			Message: "invalid environment '" + c.Env + "', must be one of: development, staging, production",
		})
	}

	return errs
}

func (c *ServerConfig) Validate() ValidationErrors {
	var errs ValidationErrors

	if c.Port < 1 || c.Port > 65535 {
		errs = append(errs, ValidationError{
			Field:   "server.port",
			Message: "port must be between 1 and 65535, got " + strconv.Itoa(c.Port),
		})
	}

	if c.ReadTimeout <= 0 {
		errs = append(errs, ValidationError{
			Field:   "server.read_timeout",
			Message: "read timeout must be positive",
		})
	}

	if c.WriteTimeout <= 0 {
		errs = append(errs, ValidationError{
			Field:   "server.write_timeout",
			Message: "write timeout must be positive",
		})
	}

	if c.IdleTimeout <= 0 {
		errs = append(errs, ValidationError{
			Field:   "server.idle_timeout",
			Message: "idle timeout must be positive",
		})
	}

	return errs
}

func (c *DatabaseConfig) Validate() ValidationErrors {
	var errs ValidationErrors

	if c.Host == "" {
		errs = append(errs, ValidationError{
			Field:   "db.host",
			Message: "database host is required",
		})
	}

	if c.Port < 1 || c.Port > 65535 {
		errs = append(errs, ValidationError{
			Field:   "db.port",
			Message: "database port must be between 1 and 65535, got " + strconv.Itoa(c.Port),
		})
	}

	if c.User == "" {
		errs = append(errs, ValidationError{
			Field:   "db.user",
			Message: "database user is required",
		})
	}

	if c.Name == "" {
		errs = append(errs, ValidationError{
			Field:   "db.name",
			Message: "database name is required",
		})
	}

	validSSLModes := map[string]bool{
		"disable":     true,
		"require":     true,
		"verify-ca":   true,
		"verify-full": true,
	}
	if !validSSLModes[c.SSLMode] {
		errs = append(errs, ValidationError{
			Field:   "db.ssl_mode",
			Message: "invalid SSL mode '" + c.SSLMode + "', must be one of: disable, require, verify-ca, verify-full",
		})
	}

	if c.MaxOpenConns < 1 {
		errs = append(errs, ValidationError{
			Field:   "db.max_open_conns",
			Message: "max open connections must be at least 1",
		})
	}

	if c.MaxIdleConns < 0 {
		errs = append(errs, ValidationError{
			Field:   "db.max_idle_conns",
			Message: "max idle connections cannot be negative",
		})
	}

	if c.MaxIdleConns > c.MaxOpenConns {
		errs = append(errs, ValidationError{
			Field:   "db.max_idle_conns",
			Message: "max idle connections cannot exceed max open connections",
		})
	}

	if c.ConnMaxLifetime <= 0 {
		errs = append(errs, ValidationError{
			Field:   "db.conn_max_lifetime",
			Message: "connection max lifetime must be positive",
		})
	}

	return errs
}

func (c *RedisConfig) Validate() ValidationErrors {
	var errs ValidationErrors

	if c.Host == "" {
		errs = append(errs, ValidationError{
			Field:   "redis.host",
			Message: "Redis host is required",
		})
	}

	if c.Port < 1 || c.Port > 65535 {
		errs = append(errs, ValidationError{
			Field:   "redis.port",
			Message: "Redis port must be between 1 and 65535, got " + strconv.Itoa(c.Port),
		})
	}

	if c.DB < 0 || c.DB > 15 {
		errs = append(errs, ValidationError{
			Field:   "redis.db",
			Message: "Redis database must be between 0 and 15, got " + strconv.Itoa(c.DB),
		})
	}

	return errs
}

func (c *JWTConfig) Validate(env string) ValidationErrors {
	var errs ValidationErrors

	if env == EnvProduction && c.Secret == "" {
		errs = append(errs, ValidationError{
			Field:   "jwt.secret",
			Message: "JWT secret is required in production",
		})
	}

	if c.Secret != "" && len(c.Secret) < 32 {
		errs = append(errs, ValidationError{
			Field:   "jwt.secret",
			Message: "JWT secret should be at least 32 characters for security",
		})
	}

	if c.AccessTokenTTL <= 0 {
		errs = append(errs, ValidationError{
			Field:   "jwt.access_token_ttl",
			Message: "access token TTL must be positive",
		})
	}

	if c.RefreshTokenTTL <= 0 {
		errs = append(errs, ValidationError{
			Field:   "jwt.refresh_token_ttl",
			Message: "refresh token TTL must be positive",
		})
	}

	if c.RefreshTokenTTL <= c.AccessTokenTTL {
		errs = append(errs, ValidationError{
			Field:   "jwt.refresh_token_ttl",
			Message: "refresh token TTL should be greater than access token TTL",
		})
	}

	return errs
}

func (c *LogConfig) Validate() ValidationErrors {
	var errs ValidationErrors

	validLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
		"fatal": true,
	}
	if !validLevels[strings.ToLower(c.Level)] {
		errs = append(errs, ValidationError{
			Field:   "log.level",
			Message: "invalid log level '" + c.Level + "', must be one of: debug, info, warn, error, fatal",
		})
	}

	validFormats := map[string]bool{
		"json":    true,
		"console": true,
	}
	if !validFormats[strings.ToLower(c.Format)] {
		errs = append(errs, ValidationError{
			Field:   "log.format",
			Message: "invalid log format '" + c.Format + "', must be one of: json, console",
		})
	}

	return errs
}

func (c *CORSConfig) Validate() ValidationErrors {
	var errs ValidationErrors

	if len(c.AllowedOrigins) == 0 {
		errs = append(errs, ValidationError{
			Field:   "cors.allowed_origins",
			Message: "at least one allowed origin is required",
		})
	}

	if len(c.AllowedMethods) == 0 {
		errs = append(errs, ValidationError{
			Field:   "cors.allowed_methods",
			Message: "at least one allowed method is required",
		})
	}

	if c.MaxAge < 0 {
		errs = append(errs, ValidationError{
			Field:   "cors.max_age",
			Message: "CORS max age cannot be negative",
		})
	}

	return errs
}

func IsValidationError(err error) bool {
	var validationErrs ValidationErrors
	return errors.As(err, &validationErrs)
}
