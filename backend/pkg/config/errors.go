package config

import "errors"

var (
	ErrReadConfigFile   = errors.New("failed to read config file")
	ErrUnmarshalConfig  = errors.New("failed to unmarshal config")
	ErrValidationFailed = errors.New("config validation failed")
	ErrLoadFailed       = errors.New("failed to load config")
)

type ConfigError struct {
	Op  string
	Err error
}

func (e *ConfigError) Error() string {
	if e.Err != nil {
		return e.Op + ": " + e.Err.Error()
	}
	return e.Op
}

func (e *ConfigError) Unwrap() error {
	return e.Err
}
