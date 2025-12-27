package circuitbreaker

import (
	"sync"
)

type Registry struct {
	breakers       map[string]*CircuitBreaker
	mutex          sync.RWMutex
	defaultConfig  func(name string) Config
}

func NewRegistry() *Registry {
	return &Registry{
		breakers:      make(map[string]*CircuitBreaker),
		defaultConfig: DefaultConfig,
	}
}

func NewRegistryWithConfig(defaultConfig func(name string) Config) *Registry {
	return &Registry{
		breakers:      make(map[string]*CircuitBreaker),
		defaultConfig: defaultConfig,
	}
}

func (r *Registry) Get(name string) *CircuitBreaker {
	r.mutex.RLock()
	cb, exists := r.breakers[name]
	r.mutex.RUnlock()

	if exists {
		return cb
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	cb, exists = r.breakers[name]
	if exists {
		return cb
	}

	cb = New(r.defaultConfig(name))
	r.breakers[name] = cb
	return cb
}

func (r *Registry) GetOrCreate(name string, config Config) *CircuitBreaker {
	r.mutex.RLock()
	cb, exists := r.breakers[name]
	r.mutex.RUnlock()

	if exists {
		return cb
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	cb, exists = r.breakers[name]
	if exists {
		return cb
	}

	cb = New(config)
	r.breakers[name] = cb
	return cb
}

func (r *Registry) Register(name string, config Config) *CircuitBreaker {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	cb := New(config)
	r.breakers[name] = cb
	return cb
}

func (r *Registry) Remove(name string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	delete(r.breakers, name)
}

func (r *Registry) All() map[string]*CircuitBreaker {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	result := make(map[string]*CircuitBreaker, len(r.breakers))
	for name, cb := range r.breakers {
		result[name] = cb
	}
	return result
}

func (r *Registry) Stats() map[string]BreakerStats {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	result := make(map[string]BreakerStats, len(r.breakers))
	for name, cb := range r.breakers {
		counts := cb.Counts()
		result[name] = BreakerStats{
			Name:                 name,
			State:                cb.State().String(),
			Requests:             counts.Requests,
			TotalSuccesses:       counts.TotalSuccesses,
			TotalFailures:        counts.TotalFailures,
			ConsecutiveSuccesses: counts.ConsecutiveSuccesses,
			ConsecutiveFailures:  counts.ConsecutiveFailures,
		}
	}
	return result
}

type BreakerStats struct {
	Name                 string `json:"name"`
	State                string `json:"state"`
	Requests             int64  `json:"requests"`
	TotalSuccesses       int64  `json:"total_successes"`
	TotalFailures        int64  `json:"total_failures"`
	ConsecutiveSuccesses int64  `json:"consecutive_successes"`
	ConsecutiveFailures  int64  `json:"consecutive_failures"`
}

var globalRegistry = NewRegistry()

func GetGlobal(name string) *CircuitBreaker {
	return globalRegistry.Get(name)
}

func RegisterGlobal(name string, config Config) *CircuitBreaker {
	return globalRegistry.Register(name, config)
}

func GlobalStats() map[string]BreakerStats {
	return globalRegistry.Stats()
}
