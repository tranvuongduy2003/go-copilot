package query_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tranvuongduy2003/go-copilot/internal/application/query"
	"github.com/tranvuongduy2003/go-copilot/internal/domain/shared"
	"github.com/tranvuongduy2003/go-copilot/pkg/testutil"
)

func TestGetUserHandler_Handle(t *testing.T) {
	t.Run("returns user successfully", func(t *testing.T) {
		repository := testutil.NewMockUserRepository()
		logger := testutil.NewNoopLogger()

		existingUser := testutil.CreateTestUser()
		repository.AddUser(existingUser)

		handler := query.NewGetUserHandler(repository, logger)

		q := query.GetUserQuery{
			UserID: existingUser.ID(),
		}

		result, err := handler.Handle(context.Background(), q)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, existingUser.ID(), result.ID)
		assert.Equal(t, existingUser.Email().String(), result.Email)
		assert.Equal(t, existingUser.FullName().String(), result.FullName)
	})

	t.Run("returns not found error when user does not exist", func(t *testing.T) {
		repository := testutil.NewMockUserRepository()
		logger := testutil.NewNoopLogger()

		handler := query.NewGetUserHandler(repository, logger)

		q := query.GetUserQuery{
			UserID: uuid.New(),
		}

		result, err := handler.Handle(context.Background(), q)

		require.Error(t, err)
		assert.Nil(t, result)
		assert.True(t, shared.IsNotFoundError(err))
	})

	t.Run("returns error when repository fails", func(t *testing.T) {
		repository := testutil.NewMockUserRepository()
		repository.FindError = errors.New("database error")
		logger := testutil.NewNoopLogger()

		handler := query.NewGetUserHandler(repository, logger)

		q := query.GetUserQuery{
			UserID: uuid.New(),
		}

		result, err := handler.Handle(context.Background(), q)

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "find user")
	})
}

func TestGetUserByEmailHandler_Handle(t *testing.T) {
	t.Run("returns user by email successfully", func(t *testing.T) {
		repository := testutil.NewMockUserRepository()
		logger := testutil.NewNoopLogger()

		existingUser := testutil.NewUserBuilder().
			WithEmail("test@example.com").
			MustBuild()
		repository.AddUser(existingUser)

		handler := query.NewGetUserByEmailHandler(repository, logger)

		q := query.GetUserByEmailQuery{
			Email: "test@example.com",
		}

		result, err := handler.Handle(context.Background(), q)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "test@example.com", result.Email)
	})

	t.Run("returns not found error when email does not exist", func(t *testing.T) {
		repository := testutil.NewMockUserRepository()
		logger := testutil.NewNoopLogger()

		handler := query.NewGetUserByEmailHandler(repository, logger)

		q := query.GetUserByEmailQuery{
			Email: "nonexistent@example.com",
		}

		result, err := handler.Handle(context.Background(), q)

		require.Error(t, err)
		assert.Nil(t, result)
		assert.True(t, shared.IsNotFoundError(err))
	})

	t.Run("returns error when repository fails", func(t *testing.T) {
		repository := testutil.NewMockUserRepository()
		repository.FindError = errors.New("database error")
		logger := testutil.NewNoopLogger()

		handler := query.NewGetUserByEmailHandler(repository, logger)

		q := query.GetUserByEmailQuery{
			Email: "test@example.com",
		}

		result, err := handler.Handle(context.Background(), q)

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "find user by email")
	})
}

func TestCheckEmailExistsHandler_Handle(t *testing.T) {
	t.Run("returns true when email exists", func(t *testing.T) {
		repository := testutil.NewMockUserRepository()
		logger := testutil.NewNoopLogger()

		existingUser := testutil.NewUserBuilder().
			WithEmail("existing@example.com").
			MustBuild()
		repository.AddUser(existingUser)

		handler := query.NewCheckEmailExistsHandler(repository, logger)

		q := query.CheckEmailExistsQuery{
			Email: "existing@example.com",
		}

		exists, err := handler.Handle(context.Background(), q)

		require.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("returns false when email does not exist", func(t *testing.T) {
		repository := testutil.NewMockUserRepository()
		logger := testutil.NewNoopLogger()

		handler := query.NewCheckEmailExistsHandler(repository, logger)

		q := query.CheckEmailExistsQuery{
			Email: "nonexistent@example.com",
		}

		exists, err := handler.Handle(context.Background(), q)

		require.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("returns error when repository fails", func(t *testing.T) {
		repository := testutil.NewMockUserRepository()
		repository.FindError = errors.New("database error")
		logger := testutil.NewNoopLogger()

		handler := query.NewCheckEmailExistsHandler(repository, logger)

		q := query.CheckEmailExistsQuery{
			Email: "test@example.com",
		}

		_, err := handler.Handle(context.Background(), q)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "check email exists")
	})
}
