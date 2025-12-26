package config

import (
	"sync"
	"sync/atomic"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type ReloadableConfig struct {
	Log  LogConfig  `mapstructure:"log"`
	CORS CORSConfig `mapstructure:"cors"`
}

type ReloadCallback func(old, new *ReloadableConfig)

type HotReloader struct {
	v         *viper.Viper
	config    atomic.Pointer[ReloadableConfig]
	callbacks []ReloadCallback
	mu        sync.RWMutex
	watching  bool
	stopCh    chan struct{}
}

func NewHotReloader(v *viper.Viper) *HotReloader {
	hr := &HotReloader{
		v:      v,
		stopCh: make(chan struct{}),
	}

	cfg := hr.loadReloadableConfig()
	hr.config.Store(cfg)

	return hr
}

func (hr *HotReloader) GetReloadable() *ReloadableConfig {
	return hr.config.Load()
}

func (hr *HotReloader) OnReload(callback ReloadCallback) {
	hr.mu.Lock()
	defer hr.mu.Unlock()
	hr.callbacks = append(hr.callbacks, callback)
}

func (hr *HotReloader) Watch() error {
	hr.mu.Lock()
	if hr.watching {
		hr.mu.Unlock()
		return nil
	}
	hr.watching = true
	hr.mu.Unlock()

	hr.v.OnConfigChange(func(e fsnotify.Event) {
		hr.reload()
	})

	hr.v.WatchConfig()
	return nil
}

func (hr *HotReloader) Stop() {
	hr.mu.Lock()
	defer hr.mu.Unlock()

	if !hr.watching {
		return
	}

	hr.watching = false
	close(hr.stopCh)
	hr.stopCh = make(chan struct{})
}

func (hr *HotReloader) Reload() error {
	if err := hr.v.ReadInConfig(); err != nil {
		return err
	}
	hr.reload()
	return nil
}

func (hr *HotReloader) reload() {
	oldConfig := hr.config.Load()
	newConfig := hr.loadReloadableConfig()

	hr.config.Store(newConfig)

	hr.mu.RLock()
	callbacks := make([]ReloadCallback, len(hr.callbacks))
	copy(callbacks, hr.callbacks)
	hr.mu.RUnlock()

	for _, cb := range callbacks {
		cb(oldConfig, newConfig)
	}
}

func (hr *HotReloader) loadReloadableConfig() *ReloadableConfig {
	cfg := &ReloadableConfig{}

	cfg.Log.Level = hr.v.GetString("log.level")
	cfg.Log.Format = hr.v.GetString("log.format")

	cfg.CORS.AllowedOrigins = hr.v.GetStringSlice("cors.allowed_origins")
	cfg.CORS.AllowedMethods = hr.v.GetStringSlice("cors.allowed_methods")
	cfg.CORS.AllowedHeaders = hr.v.GetStringSlice("cors.allowed_headers")
	cfg.CORS.MaxAge = hr.v.GetInt("cors.max_age")

	return cfg
}

func (hr *HotReloader) IsWatching() bool {
	hr.mu.RLock()
	defer hr.mu.RUnlock()
	return hr.watching
}
