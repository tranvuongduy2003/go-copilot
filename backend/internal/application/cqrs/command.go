package cqrs

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	"github.com/tranvuongduy2003/go-copilot/pkg/logger"
)

type Command interface{}

type CommandHandler[C Command, R any] interface {
	Handle(ctx context.Context, command C) (R, error)
}

type CommandBus interface {
	Dispatch(ctx context.Context, command Command) (interface{}, error)
	Register(commandType Command, handler interface{})
	Use(middleware ...CommandMiddleware)
}

type InMemoryCommandBus struct {
	handlers    map[reflect.Type]interface{}
	mutex       sync.RWMutex
	logger      logger.Logger
	middlewares []CommandMiddleware
}

func NewInMemoryCommandBus(log logger.Logger) *InMemoryCommandBus {
	return &InMemoryCommandBus{
		handlers:    make(map[reflect.Type]interface{}),
		logger:      log,
		middlewares: make([]CommandMiddleware, 0),
	}
}

func (bus *InMemoryCommandBus) Use(middlewares ...CommandMiddleware) {
	bus.mutex.Lock()
	defer bus.mutex.Unlock()
	bus.middlewares = append(bus.middlewares, middlewares...)
}

func (bus *InMemoryCommandBus) Register(commandType Command, handler interface{}) {
	bus.mutex.Lock()
	defer bus.mutex.Unlock()

	cmdType := reflect.TypeOf(commandType)
	bus.handlers[cmdType] = handler

	bus.logger.Debug("registered command handler",
		logger.String("command_type", cmdType.String()),
	)
}

func (bus *InMemoryCommandBus) Dispatch(ctx context.Context, command Command) (interface{}, error) {
	coreDispatcher := func(dispatchContext context.Context, dispatchCommand Command) (interface{}, error) {
		return bus.executeHandler(dispatchContext, dispatchCommand)
	}

	bus.mutex.RLock()
	middlewareChain := make([]CommandMiddleware, len(bus.middlewares))
	copy(middlewareChain, bus.middlewares)
	bus.mutex.RUnlock()

	finalDispatcher := ChainCommandMiddleware(middlewareChain...)(coreDispatcher)
	return finalDispatcher(ctx, command)
}

func (bus *InMemoryCommandBus) executeHandler(ctx context.Context, command Command) (interface{}, error) {
	bus.mutex.RLock()
	cmdType := reflect.TypeOf(command)
	handler, exists := bus.handlers[cmdType]
	bus.mutex.RUnlock()

	if !exists {
		return nil, fmt.Errorf("no handler registered for command type: %s", cmdType.String())
	}

	bus.logger.Debug("dispatching command",
		logger.String("command_type", cmdType.String()),
	)

	handlerValue := reflect.ValueOf(handler)
	handleMethod := handlerValue.MethodByName("Handle")

	if !handleMethod.IsValid() {
		return nil, fmt.Errorf("handler does not have Handle method for command type: %s", cmdType.String())
	}

	args := []reflect.Value{
		reflect.ValueOf(ctx),
		reflect.ValueOf(command),
	}

	results := handleMethod.Call(args)

	if len(results) == 2 {
		var result interface{}
		if canBeNil(results[0]) {
			if !results[0].IsNil() {
				result = results[0].Interface()
			}
		} else {
			result = results[0].Interface()
		}

		var err error
		if canBeNil(results[1]) && !results[1].IsNil() {
			err = results[1].Interface().(error)
		}

		return result, err
	}

	if len(results) == 1 {
		if canBeNil(results[0]) && results[0].IsNil() {
			return nil, nil
		}
		if errVal, ok := results[0].Interface().(error); ok {
			return nil, errVal
		}
		return results[0].Interface(), nil
	}

	return nil, nil
}

func RegisterCommandHandler[C Command, R any](bus *InMemoryCommandBus, handler CommandHandler[C, R]) {
	var cmd C
	bus.Register(cmd, handler)
}

func canBeNil(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return true
	default:
		return false
	}
}
