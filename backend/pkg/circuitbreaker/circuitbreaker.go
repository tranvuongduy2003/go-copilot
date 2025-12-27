package circuitbreaker

import (
	"context"
	"errors"
	"sync"
	"time"
)

type State int

const (
	StateClosed State = iota
	StateOpen
	StateHalfOpen
)

func (s State) String() string {
	switch s {
	case StateClosed:
		return "closed"
	case StateOpen:
		return "open"
	case StateHalfOpen:
		return "half-open"
	default:
		return "unknown"
	}
}

var (
	ErrCircuitOpen     = errors.New("circuit breaker is open")
	ErrTooManyRequests = errors.New("too many requests in half-open state")
)

type Config struct {
	Name                   string
	FailureThreshold       int
	SuccessThreshold       int
	Timeout                time.Duration
	MaxHalfOpenRequests    int
	OnStateChange          func(name string, from State, to State)
	IsSuccessful           func(err error) bool
}

func DefaultConfig(name string) Config {
	return Config{
		Name:                name,
		FailureThreshold:    5,
		SuccessThreshold:    2,
		Timeout:             30 * time.Second,
		MaxHalfOpenRequests: 1,
		OnStateChange:       nil,
		IsSuccessful:        defaultIsSuccessful,
	}
}

func defaultIsSuccessful(err error) bool {
	return err == nil
}

type Counts struct {
	Requests             int64
	TotalSuccesses       int64
	TotalFailures        int64
	ConsecutiveSuccesses int64
	ConsecutiveFailures  int64
}

func (c *Counts) onRequest() {
	c.Requests++
}

func (c *Counts) onSuccess() {
	c.TotalSuccesses++
	c.ConsecutiveSuccesses++
	c.ConsecutiveFailures = 0
}

func (c *Counts) onFailure() {
	c.TotalFailures++
	c.ConsecutiveFailures++
	c.ConsecutiveSuccesses = 0
}

func (c *Counts) clear() {
	c.Requests = 0
	c.TotalSuccesses = 0
	c.TotalFailures = 0
	c.ConsecutiveSuccesses = 0
	c.ConsecutiveFailures = 0
}

type CircuitBreaker struct {
	name                string
	failureThreshold    int
	successThreshold    int
	timeout             time.Duration
	maxHalfOpenRequests int
	onStateChange       func(name string, from State, to State)
	isSuccessful        func(err error) bool

	mutex      sync.RWMutex
	state      State
	counts     Counts
	expiry     time.Time
	generation uint64
}

func New(config Config) *CircuitBreaker {
	if config.FailureThreshold <= 0 {
		config.FailureThreshold = 5
	}
	if config.SuccessThreshold <= 0 {
		config.SuccessThreshold = 2
	}
	if config.Timeout <= 0 {
		config.Timeout = 30 * time.Second
	}
	if config.MaxHalfOpenRequests <= 0 {
		config.MaxHalfOpenRequests = 1
	}
	if config.IsSuccessful == nil {
		config.IsSuccessful = defaultIsSuccessful
	}

	return &CircuitBreaker{
		name:                config.Name,
		failureThreshold:    config.FailureThreshold,
		successThreshold:    config.SuccessThreshold,
		timeout:             config.Timeout,
		maxHalfOpenRequests: config.MaxHalfOpenRequests,
		onStateChange:       config.OnStateChange,
		isSuccessful:        config.IsSuccessful,
		state:               StateClosed,
	}
}

func (cb *CircuitBreaker) Name() string {
	return cb.name
}

func (cb *CircuitBreaker) State() State {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()

	now := time.Now()
	state, _ := cb.currentState(now)
	return state
}

func (cb *CircuitBreaker) Counts() Counts {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()
	return cb.counts
}

func (cb *CircuitBreaker) Execute(request func() (interface{}, error)) (interface{}, error) {
	return cb.ExecuteContext(context.Background(), func(ctx context.Context) (interface{}, error) {
		return request()
	})
}

func (cb *CircuitBreaker) ExecuteContext(ctx context.Context, request func(ctx context.Context) (interface{}, error)) (interface{}, error) {
	generation, err := cb.beforeRequest()
	if err != nil {
		return nil, err
	}

	defer func() {
		panicErr := recover()
		if panicErr != nil {
			cb.afterRequest(generation, false)
			panic(panicErr)
		}
	}()

	result, err := request(ctx)
	cb.afterRequest(generation, cb.isSuccessful(err))
	return result, err
}

func (cb *CircuitBreaker) beforeRequest() (uint64, error) {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	now := time.Now()
	state, generation := cb.currentState(now)

	if state == StateOpen {
		return generation, ErrCircuitOpen
	}

	if state == StateHalfOpen && cb.counts.Requests >= int64(cb.maxHalfOpenRequests) {
		return generation, ErrTooManyRequests
	}

	cb.counts.onRequest()
	return generation, nil
}

func (cb *CircuitBreaker) afterRequest(beforeGeneration uint64, success bool) {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	now := time.Now()
	state, generation := cb.currentState(now)
	if generation != beforeGeneration {
		return
	}

	if success {
		cb.onSuccess(state, now)
	} else {
		cb.onFailure(state, now)
	}
}

func (cb *CircuitBreaker) onSuccess(state State, now time.Time) {
	switch state {
	case StateClosed:
		cb.counts.onSuccess()
	case StateHalfOpen:
		cb.counts.onSuccess()
		if cb.counts.ConsecutiveSuccesses >= int64(cb.successThreshold) {
			cb.setState(StateClosed, now)
		}
	}
}

func (cb *CircuitBreaker) onFailure(state State, now time.Time) {
	switch state {
	case StateClosed:
		cb.counts.onFailure()
		if cb.counts.ConsecutiveFailures >= int64(cb.failureThreshold) {
			cb.setState(StateOpen, now)
		}
	case StateHalfOpen:
		cb.setState(StateOpen, now)
	}
}

func (cb *CircuitBreaker) currentState(now time.Time) (State, uint64) {
	switch cb.state {
	case StateClosed:
		return StateClosed, cb.generation
	case StateOpen:
		if cb.expiry.Before(now) {
			cb.setState(StateHalfOpen, now)
		}
	}
	return cb.state, cb.generation
}

func (cb *CircuitBreaker) setState(state State, now time.Time) {
	if cb.state == state {
		return
	}

	previousState := cb.state
	cb.state = state
	cb.generation++
	cb.counts.clear()

	switch state {
	case StateClosed:
		cb.expiry = time.Time{}
	case StateOpen:
		cb.expiry = now.Add(cb.timeout)
	case StateHalfOpen:
		cb.expiry = time.Time{}
	}

	if cb.onStateChange != nil {
		go cb.onStateChange(cb.name, previousState, state)
	}
}

type ExecuteFunc[T any] func(ctx context.Context) (T, error)

func Execute[T any](cb *CircuitBreaker, ctx context.Context, fn ExecuteFunc[T]) (T, error) {
	result, err := cb.ExecuteContext(ctx, func(ctx context.Context) (interface{}, error) {
		return fn(ctx)
	})
	if err != nil {
		var zero T
		return zero, err
	}
	return result.(T), nil
}
