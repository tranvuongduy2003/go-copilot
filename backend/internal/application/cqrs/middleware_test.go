package cqrs

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tranvuongduy2003/go-copilot/pkg/testutil"
	"github.com/tranvuongduy2003/go-copilot/pkg/validator"
)

type testCommand struct {
	Value string `validate:"required"`
}

type testQuery struct {
	ID string `validate:"required,uuid"`
}

func TestLoggingCommandMiddleware(t *testing.T) {
	log := testutil.NewNoopLogger()
	middleware := LoggingCommandMiddleware(log)

	tests := []struct {
		name       string
		dispatcher CommandDispatcher
		wantErr    bool
	}{
		{
			name: "logs successful command execution",
			dispatcher: func(ctx context.Context, cmd Command) (interface{}, error) {
				return "success", nil
			},
			wantErr: false,
		},
		{
			name: "logs failed command execution",
			dispatcher: func(ctx context.Context, cmd Command) (interface{}, error) {
				return nil, errors.New("command failed")
			},
			wantErr: true,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			wrappedDispatcher := middleware(testCase.dispatcher)
			result, err := wrappedDispatcher(context.Background(), testCommand{Value: "test"})

			if testCase.wantErr {
				require.Error(t, err)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.Equal(t, "success", result)
			}
		})
	}
}

func TestLoggingQueryMiddleware(t *testing.T) {
	log := testutil.NewNoopLogger()
	middleware := LoggingQueryMiddleware(log)

	tests := []struct {
		name       string
		dispatcher QueryDispatcher
		wantErr    bool
	}{
		{
			name: "logs successful query execution",
			dispatcher: func(ctx context.Context, query Query) (interface{}, error) {
				return "result", nil
			},
			wantErr: false,
		},
		{
			name: "logs failed query execution",
			dispatcher: func(ctx context.Context, query Query) (interface{}, error) {
				return nil, errors.New("query failed")
			},
			wantErr: true,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			wrappedDispatcher := middleware(testCase.dispatcher)
			result, err := wrappedDispatcher(context.Background(), testQuery{ID: "550e8400-e29b-41d4-a716-446655440000"})

			if testCase.wantErr {
				require.Error(t, err)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.Equal(t, "result", result)
			}
		})
	}
}

func TestValidationCommandMiddleware(t *testing.T) {
	validate := validator.New()
	middleware := ValidationCommandMiddleware(validate)

	tests := []struct {
		name        string
		command     Command
		wantErr     bool
		errContains string
	}{
		{
			name:    "passes valid command",
			command: testCommand{Value: "test"},
			wantErr: false,
		},
		{
			name:        "rejects invalid command",
			command:     testCommand{Value: ""},
			wantErr:     true,
			errContains: "validation failed",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			dispatcher := func(ctx context.Context, cmd Command) (interface{}, error) {
				return "executed", nil
			}

			wrappedDispatcher := middleware(dispatcher)
			result, err := wrappedDispatcher(context.Background(), testCase.command)

			if testCase.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), testCase.errContains)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.Equal(t, "executed", result)
			}
		})
	}
}

func TestValidationQueryMiddleware(t *testing.T) {
	validate := validator.New()
	middleware := ValidationQueryMiddleware(validate)

	tests := []struct {
		name        string
		query       Query
		wantErr     bool
		errContains string
	}{
		{
			name:    "passes valid query",
			query:   testQuery{ID: "550e8400-e29b-41d4-a716-446655440000"},
			wantErr: false,
		},
		{
			name:        "rejects invalid query - empty",
			query:       testQuery{ID: ""},
			wantErr:     true,
			errContains: "validation failed",
		},
		{
			name:        "rejects invalid query - not uuid",
			query:       testQuery{ID: "not-a-uuid"},
			wantErr:     true,
			errContains: "validation failed",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			dispatcher := func(ctx context.Context, query Query) (interface{}, error) {
				return "executed", nil
			}

			wrappedDispatcher := middleware(dispatcher)
			result, err := wrappedDispatcher(context.Background(), testCase.query)

			if testCase.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), testCase.errContains)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.Equal(t, "executed", result)
			}
		})
	}
}

type mockTransactionManager struct {
	beginError    error
	commitError   error
	rollbackError error
	committed     bool
	rolledBack    bool
}

func (manager *mockTransactionManager) Begin(ctx context.Context) (context.Context, error) {
	if manager.beginError != nil {
		return nil, manager.beginError
	}
	return context.WithValue(ctx, "transaction", true), nil
}

func (manager *mockTransactionManager) Commit(ctx context.Context) error {
	manager.committed = true
	return manager.commitError
}

func (manager *mockTransactionManager) Rollback(ctx context.Context) error {
	manager.rolledBack = true
	return manager.rollbackError
}

func TestTransactionCommandMiddleware(t *testing.T) {
	tests := []struct {
		name            string
		setupManager    func() *mockTransactionManager
		dispatcher      CommandDispatcher
		wantErr         bool
		errContains     string
		expectCommit    bool
		expectRollback  bool
	}{
		{
			name: "commits on success",
			setupManager: func() *mockTransactionManager {
				return &mockTransactionManager{}
			},
			dispatcher: func(ctx context.Context, cmd Command) (interface{}, error) {
				assert.Equal(t, true, ctx.Value("transaction"))
				return "success", nil
			},
			wantErr:        false,
			expectCommit:   true,
			expectRollback: false,
		},
		{
			name: "rolls back on command error",
			setupManager: func() *mockTransactionManager {
				return &mockTransactionManager{}
			},
			dispatcher: func(ctx context.Context, cmd Command) (interface{}, error) {
				return nil, errors.New("command failed")
			},
			wantErr:        true,
			errContains:    "command failed",
			expectCommit:   false,
			expectRollback: true,
		},
		{
			name: "returns error when begin fails",
			setupManager: func() *mockTransactionManager {
				return &mockTransactionManager{
					beginError: errors.New("begin failed"),
				}
			},
			dispatcher: func(ctx context.Context, cmd Command) (interface{}, error) {
				return "success", nil
			},
			wantErr:        true,
			errContains:    "begin transaction",
			expectCommit:   false,
			expectRollback: false,
		},
		{
			name: "returns error when commit fails",
			setupManager: func() *mockTransactionManager {
				return &mockTransactionManager{
					commitError: errors.New("commit failed"),
				}
			},
			dispatcher: func(ctx context.Context, cmd Command) (interface{}, error) {
				return "success", nil
			},
			wantErr:        true,
			errContains:    "commit transaction",
			expectCommit:   true,
			expectRollback: false,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			manager := testCase.setupManager()
			middleware := TransactionCommandMiddleware(manager)

			wrappedDispatcher := middleware(testCase.dispatcher)
			result, err := wrappedDispatcher(context.Background(), testCommand{Value: "test"})

			if testCase.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), testCase.errContains)
			} else {
				require.NoError(t, err)
				assert.Equal(t, "success", result)
			}

			assert.Equal(t, testCase.expectCommit, manager.committed)
			assert.Equal(t, testCase.expectRollback, manager.rolledBack)
		})
	}
}

type mockMetricsRecorder struct {
	commandRecordings []commandRecording
	queryRecordings   []queryRecording
}

type commandRecording struct {
	commandType string
	duration    time.Duration
	success     bool
}

type queryRecording struct {
	queryType string
	duration  time.Duration
	success   bool
}

func (recorder *mockMetricsRecorder) RecordCommandDuration(commandType string, duration time.Duration, success bool) {
	recorder.commandRecordings = append(recorder.commandRecordings, commandRecording{
		commandType: commandType,
		duration:    duration,
		success:     success,
	})
}

func (recorder *mockMetricsRecorder) RecordQueryDuration(queryType string, duration time.Duration, success bool) {
	recorder.queryRecordings = append(recorder.queryRecordings, queryRecording{
		queryType: queryType,
		duration:  duration,
		success:   success,
	})
}

func TestMetricsCommandMiddleware(t *testing.T) {
	recorder := &mockMetricsRecorder{}
	middleware := MetricsCommandMiddleware(recorder)

	tests := []struct {
		name        string
		dispatcher  CommandDispatcher
		wantSuccess bool
	}{
		{
			name: "records successful command",
			dispatcher: func(ctx context.Context, cmd Command) (interface{}, error) {
				return "success", nil
			},
			wantSuccess: true,
		},
		{
			name: "records failed command",
			dispatcher: func(ctx context.Context, cmd Command) (interface{}, error) {
				return nil, errors.New("failed")
			},
			wantSuccess: false,
		},
	}

	for index, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			wrappedDispatcher := middleware(testCase.dispatcher)
			_, _ = wrappedDispatcher(context.Background(), testCommand{Value: "test"})

			require.Len(t, recorder.commandRecordings, index+1)
			assert.Contains(t, recorder.commandRecordings[index].commandType, "testCommand")
			assert.Equal(t, testCase.wantSuccess, recorder.commandRecordings[index].success)
			assert.GreaterOrEqual(t, recorder.commandRecordings[index].duration, time.Duration(0))
		})
	}
}

func TestMetricsQueryMiddleware(t *testing.T) {
	recorder := &mockMetricsRecorder{}
	middleware := MetricsQueryMiddleware(recorder)

	tests := []struct {
		name        string
		dispatcher  QueryDispatcher
		wantSuccess bool
	}{
		{
			name: "records successful query",
			dispatcher: func(ctx context.Context, query Query) (interface{}, error) {
				return "result", nil
			},
			wantSuccess: true,
		},
		{
			name: "records failed query",
			dispatcher: func(ctx context.Context, query Query) (interface{}, error) {
				return nil, errors.New("failed")
			},
			wantSuccess: false,
		},
	}

	for index, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			wrappedDispatcher := middleware(testCase.dispatcher)
			_, _ = wrappedDispatcher(context.Background(), testQuery{ID: "550e8400-e29b-41d4-a716-446655440000"})

			require.Len(t, recorder.queryRecordings, index+1)
			assert.Contains(t, recorder.queryRecordings[index].queryType, "testQuery")
			assert.Equal(t, testCase.wantSuccess, recorder.queryRecordings[index].success)
			assert.GreaterOrEqual(t, recorder.queryRecordings[index].duration, time.Duration(0))
		})
	}
}

func TestChainCommandMiddleware(t *testing.T) {
	executionOrder := []string{}

	middleware1 := func(next CommandDispatcher) CommandDispatcher {
		return func(ctx context.Context, cmd Command) (interface{}, error) {
			executionOrder = append(executionOrder, "middleware1-before")
			result, err := next(ctx, cmd)
			executionOrder = append(executionOrder, "middleware1-after")
			return result, err
		}
	}

	middleware2 := func(next CommandDispatcher) CommandDispatcher {
		return func(ctx context.Context, cmd Command) (interface{}, error) {
			executionOrder = append(executionOrder, "middleware2-before")
			result, err := next(ctx, cmd)
			executionOrder = append(executionOrder, "middleware2-after")
			return result, err
		}
	}

	coreDispatcher := func(ctx context.Context, cmd Command) (interface{}, error) {
		executionOrder = append(executionOrder, "core")
		return "result", nil
	}

	chained := ChainCommandMiddleware(middleware1, middleware2)
	finalDispatcher := chained(coreDispatcher)

	result, err := finalDispatcher(context.Background(), testCommand{Value: "test"})

	require.NoError(t, err)
	assert.Equal(t, "result", result)

	expectedOrder := []string{
		"middleware1-before",
		"middleware2-before",
		"core",
		"middleware2-after",
		"middleware1-after",
	}
	assert.Equal(t, expectedOrder, executionOrder)
}

func TestChainQueryMiddleware(t *testing.T) {
	executionOrder := []string{}

	middleware1 := func(next QueryDispatcher) QueryDispatcher {
		return func(ctx context.Context, query Query) (interface{}, error) {
			executionOrder = append(executionOrder, "middleware1-before")
			result, err := next(ctx, query)
			executionOrder = append(executionOrder, "middleware1-after")
			return result, err
		}
	}

	middleware2 := func(next QueryDispatcher) QueryDispatcher {
		return func(ctx context.Context, query Query) (interface{}, error) {
			executionOrder = append(executionOrder, "middleware2-before")
			result, err := next(ctx, query)
			executionOrder = append(executionOrder, "middleware2-after")
			return result, err
		}
	}

	coreDispatcher := func(ctx context.Context, query Query) (interface{}, error) {
		executionOrder = append(executionOrder, "core")
		return "result", nil
	}

	chained := ChainQueryMiddleware(middleware1, middleware2)
	finalDispatcher := chained(coreDispatcher)

	result, err := finalDispatcher(context.Background(), testQuery{ID: "550e8400-e29b-41d4-a716-446655440000"})

	require.NoError(t, err)
	assert.Equal(t, "result", result)

	expectedOrder := []string{
		"middleware1-before",
		"middleware2-before",
		"core",
		"middleware2-after",
		"middleware1-after",
	}
	assert.Equal(t, expectedOrder, executionOrder)
}

func TestCommandBusWithMiddleware(t *testing.T) {
	log := testutil.NewNoopLogger()
	bus := NewInMemoryCommandBus(log)

	middlewareExecuted := false
	testMiddleware := func(next CommandDispatcher) CommandDispatcher {
		return func(ctx context.Context, cmd Command) (interface{}, error) {
			middlewareExecuted = true
			return next(ctx, cmd)
		}
	}

	bus.Use(testMiddleware)

	handler := &testCommandHandler{}
	RegisterCommandHandler[testCommand, string](bus, handler)

	result, err := bus.Dispatch(context.Background(), testCommand{Value: "test"})

	require.NoError(t, err)
	assert.Equal(t, "handled: test", result)
	assert.True(t, middlewareExecuted)
}

func TestQueryBusWithMiddleware(t *testing.T) {
	log := testutil.NewNoopLogger()
	bus := NewInMemoryQueryBus(log)

	middlewareExecuted := false
	testMiddleware := func(next QueryDispatcher) QueryDispatcher {
		return func(ctx context.Context, query Query) (interface{}, error) {
			middlewareExecuted = true
			return next(ctx, query)
		}
	}

	bus.Use(testMiddleware)

	handler := &testQueryHandler{}
	RegisterQueryHandler[testQuery, string](bus, handler)

	result, err := bus.Dispatch(context.Background(), testQuery{ID: "550e8400-e29b-41d4-a716-446655440000"})

	require.NoError(t, err)
	assert.Equal(t, "found: 550e8400-e29b-41d4-a716-446655440000", result)
	assert.True(t, middlewareExecuted)
}

type testCommandHandler struct{}

func (handler *testCommandHandler) Handle(ctx context.Context, command testCommand) (string, error) {
	return "handled: " + command.Value, nil
}

type testQueryHandler struct{}

func (handler *testQueryHandler) Handle(ctx context.Context, query testQuery) (string, error) {
	return "found: " + query.ID, nil
}
