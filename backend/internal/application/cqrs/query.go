package cqrs

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	"github.com/tranvuongduy2003/go-copilot/pkg/logger"
)

type Query interface{}

type QueryHandler[Q Query, R any] interface {
	Handle(ctx context.Context, query Q) (R, error)
}

type QueryBus interface {
	Dispatch(ctx context.Context, query Query) (interface{}, error)
	Register(queryType Query, handler interface{})
}

type InMemoryQueryBus struct {
	handlers map[reflect.Type]interface{}
	mutex    sync.RWMutex
	logger   logger.Logger
}

func NewInMemoryQueryBus(log logger.Logger) *InMemoryQueryBus {
	return &InMemoryQueryBus{
		handlers: make(map[reflect.Type]interface{}),
		logger:   log,
	}
}

func (bus *InMemoryQueryBus) Register(queryType Query, handler interface{}) {
	bus.mutex.Lock()
	defer bus.mutex.Unlock()

	qryType := reflect.TypeOf(queryType)
	bus.handlers[qryType] = handler

	bus.logger.Debug("registered query handler",
		logger.String("query_type", qryType.String()),
	)
}

func (bus *InMemoryQueryBus) Dispatch(ctx context.Context, query Query) (interface{}, error) {
	bus.mutex.RLock()
	qryType := reflect.TypeOf(query)
	handler, exists := bus.handlers[qryType]
	bus.mutex.RUnlock()

	if !exists {
		return nil, fmt.Errorf("no handler registered for query type: %s", qryType.String())
	}

	bus.logger.Debug("dispatching query",
		logger.String("query_type", qryType.String()),
	)

	handlerValue := reflect.ValueOf(handler)
	handleMethod := handlerValue.MethodByName("Handle")

	if !handleMethod.IsValid() {
		return nil, fmt.Errorf("handler does not have Handle method for query type: %s", qryType.String())
	}

	args := []reflect.Value{
		reflect.ValueOf(ctx),
		reflect.ValueOf(query),
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

func RegisterQueryHandler[Q Query, R any](bus *InMemoryQueryBus, handler QueryHandler[Q, R]) {
	var qry Q
	bus.Register(qry, handler)
}
