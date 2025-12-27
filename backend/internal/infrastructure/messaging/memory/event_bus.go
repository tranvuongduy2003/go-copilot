package memory

import (
	"context"
	"reflect"
	"sync"
	"time"

	"github.com/tranvuongduy2003/go-copilot/internal/domain/shared"
	"github.com/tranvuongduy2003/go-copilot/pkg/logger"
)

type EventBusOption func(*eventBusOptions)

type eventBusOptions struct {
	async          bool
	workerPoolSize int
	handlerTimeout time.Duration
}

func WithAsync(async bool) EventBusOption {
	return func(o *eventBusOptions) {
		o.async = async
	}
}

func WithWorkerPoolSize(size int) EventBusOption {
	return func(o *eventBusOptions) {
		o.workerPoolSize = size
	}
}

func WithHandlerTimeout(timeout time.Duration) EventBusOption {
	return func(o *eventBusOptions) {
		o.handlerTimeout = timeout
	}
}

type eventJob struct {
	ctx     context.Context
	event   shared.DomainEvent
	handler shared.EventHandler
}

type InMemoryEventBus struct {
	handlers       map[string][]shared.EventHandler
	mu             sync.RWMutex
	logger         logger.Logger
	options        *eventBusOptions
	jobChan        chan eventJob
	wg             sync.WaitGroup
	shutdownChan   chan struct{}
	shutdownOnce   sync.Once
}

func NewInMemoryEventBus(log logger.Logger, opts ...EventBusOption) *InMemoryEventBus {
	options := &eventBusOptions{
		async:          false,
		workerPoolSize: 10,
		handlerTimeout: 30 * time.Second,
	}

	for _, opt := range opts {
		opt(options)
	}

	bus := &InMemoryEventBus{
		handlers:     make(map[string][]shared.EventHandler),
		logger:       log,
		options:      options,
		shutdownChan: make(chan struct{}),
	}

	if options.async {
		bus.jobChan = make(chan eventJob, options.workerPoolSize*10)
		bus.startWorkers()
	}

	return bus
}

func (b *InMemoryEventBus) startWorkers() {
	for i := 0; i < b.options.workerPoolSize; i++ {
		b.wg.Add(1)
		go b.worker()
	}
}

func (b *InMemoryEventBus) worker() {
	defer b.wg.Done()

	for {
		select {
		case <-b.shutdownChan:
			return
		case job, ok := <-b.jobChan:
			if !ok {
				return
			}
			b.executeHandler(job.ctx, job.event, job.handler)
		}
	}
}

func (b *InMemoryEventBus) executeHandler(ctx context.Context, event shared.DomainEvent, handler shared.EventHandler) {
	defer func() {
		if r := recover(); r != nil {
			b.logger.Error("event handler panicked",
				logger.String("event_type", event.EventType()),
				logger.Any("panic", r),
			)
		}
	}()

	handlerCtx := ctx
	if b.options.handlerTimeout > 0 {
		var cancel context.CancelFunc
		handlerCtx, cancel = context.WithTimeout(ctx, b.options.handlerTimeout)
		defer cancel()
	}

	if err := handler(handlerCtx, event); err != nil {
		b.logger.Error("event handler failed",
			logger.String("event_type", event.EventType()),
			logger.Err(err),
		)
	}
}

func (b *InMemoryEventBus) Publish(ctx context.Context, events ...shared.DomainEvent) error {
	for _, event := range events {
		b.mu.RLock()
		handlers, ok := b.handlers[event.EventType()]
		b.mu.RUnlock()

		if !ok || len(handlers) == 0 {
			b.logger.Debug("no handlers registered for event",
				logger.String("event_type", event.EventType()),
			)
			continue
		}

		for _, handler := range handlers {
			if b.options.async {
				select {
				case b.jobChan <- eventJob{ctx: ctx, event: event, handler: handler}:
				case <-b.shutdownChan:
					return nil
				default:
					b.logger.Warn("event job channel full, executing synchronously",
						logger.String("event_type", event.EventType()),
					)
					b.executeHandler(ctx, event, handler)
				}
			} else {
				b.executeHandler(ctx, event, handler)
			}
		}
	}

	return nil
}

func (b *InMemoryEventBus) Subscribe(eventType string, handler shared.EventHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.handlers[eventType] = append(b.handlers[eventType], handler)
	b.logger.Debug("handler subscribed",
		logger.String("event_type", eventType),
	)
}

func (b *InMemoryEventBus) Unsubscribe(eventType string, handler shared.EventHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()

	handlers, ok := b.handlers[eventType]
	if !ok {
		return
	}

	handlerPtr := reflect.ValueOf(handler).Pointer()
	for i, h := range handlers {
		if reflect.ValueOf(h).Pointer() == handlerPtr {
			b.handlers[eventType] = append(handlers[:i], handlers[i+1:]...)
			b.logger.Debug("handler unsubscribed",
				logger.String("event_type", eventType),
			)
			return
		}
	}
}

func (b *InMemoryEventBus) Shutdown(ctx context.Context) error {
	b.shutdownOnce.Do(func() {
		close(b.shutdownChan)

		if b.options.async && b.jobChan != nil {
			close(b.jobChan)
		}
	})

	done := make(chan struct{})
	go func() {
		b.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		b.logger.Info("event bus shutdown completed")
		return nil
	case <-ctx.Done():
		b.logger.Warn("event bus shutdown timed out")
		return ctx.Err()
	}
}

func (b *InMemoryEventBus) HandlerCount(eventType string) int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return len(b.handlers[eventType])
}
