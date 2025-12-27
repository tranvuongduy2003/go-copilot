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
}

type InMemoryCommandBus struct {
	handlers map[reflect.Type]interface{}
	mutex    sync.RWMutex
	logger   logger.Logger
}

func NewInMemoryCommandBus(log logger.Logger) *InMemoryCommandBus {
	return &InMemoryCommandBus{
		handlers: make(map[reflect.Type]interface{}),
		logger:   log,
	}
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
		if !results[0].IsNil() {
			result = results[0].Interface()
		}

		var err error
		if !results[1].IsNil() {
			err = results[1].Interface().(error)
		}

		return result, err
	}

	if len(results) == 1 {
		if results[0].IsNil() {
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
