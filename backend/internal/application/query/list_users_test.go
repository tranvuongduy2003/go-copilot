package query_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tranvuongduy2003/go-copilot/internal/application/query"
	"github.com/tranvuongduy2003/go-copilot/pkg/testutil"
)

func TestListUsersHandler_Handle(t *testing.T) {
	t.Run("returns empty list when no users exist", func(t *testing.T) {
		repository := testutil.NewMockUserRepository()
		logger := testutil.NewNoopLogger()

		handler := query.NewListUsersHandler(repository, logger)

		q := query.ListUsersQuery{
			Page:  1,
			Limit: 10,
		}

		result, err := handler.Handle(context.Background(), q)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Empty(t, result.Items)
		assert.Equal(t, int64(0), result.Total)
		assert.Equal(t, 1, result.Page)
		assert.Equal(t, 10, result.Limit)
	})

	t.Run("returns paginated users", func(t *testing.T) {
		repository := testutil.NewMockUserRepository()
		logger := testutil.NewNoopLogger()

		users := testutil.CreateUsers(5)
		for _, u := range users {
			repository.AddUser(u)
		}

		handler := query.NewListUsersHandler(repository, logger)

		q := query.ListUsersQuery{
			Page:  1,
			Limit: 3,
		}

		result, err := handler.Handle(context.Background(), q)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Len(t, result.Items, 3)
		assert.Equal(t, int64(5), result.Total)
		assert.Equal(t, 1, result.Page)
		assert.Equal(t, 3, result.Limit)
	})

	t.Run("returns second page of users", func(t *testing.T) {
		repository := testutil.NewMockUserRepository()
		logger := testutil.NewNoopLogger()

		users := testutil.CreateUsers(5)
		for _, u := range users {
			repository.AddUser(u)
		}

		handler := query.NewListUsersHandler(repository, logger)

		q := query.ListUsersQuery{
			Page:  2,
			Limit: 3,
		}

		result, err := handler.Handle(context.Background(), q)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Len(t, result.Items, 2)
		assert.Equal(t, int64(5), result.Total)
		assert.Equal(t, 2, result.Page)
	})

	t.Run("filters by status", func(t *testing.T) {
		repository := testutil.NewMockUserRepository()
		logger := testutil.NewNoopLogger()

		activeUser1 := testutil.CreateActiveUser()
		activeUser2 := testutil.CreateActiveUser()
		pendingUser := testutil.CreatePendingUser()
		repository.AddUser(activeUser1)
		repository.AddUser(activeUser2)
		repository.AddUser(pendingUser)

		handler := query.NewListUsersHandler(repository, logger)

		status := "active"
		q := query.ListUsersQuery{
			Page:   1,
			Limit:  10,
			Status: &status,
		}

		result, err := handler.Handle(context.Background(), q)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Len(t, result.Items, 2)
		assert.Equal(t, int64(2), result.Total)
		for _, u := range result.Items {
			assert.Equal(t, "active", u.Status)
		}
	})

	t.Run("ignores invalid status filter", func(t *testing.T) {
		repository := testutil.NewMockUserRepository()
		logger := testutil.NewNoopLogger()

		user1 := testutil.CreateActiveUser()
		user2 := testutil.CreatePendingUser()
		repository.AddUser(user1)
		repository.AddUser(user2)

		handler := query.NewListUsersHandler(repository, logger)

		invalidStatus := "invalid_status"
		q := query.ListUsersQuery{
			Page:   1,
			Limit:  10,
			Status: &invalidStatus,
		}

		result, err := handler.Handle(context.Background(), q)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Len(t, result.Items, 2)
	})

	t.Run("uses default pagination when page is zero", func(t *testing.T) {
		repository := testutil.NewMockUserRepository()
		logger := testutil.NewNoopLogger()

		users := testutil.CreateUsers(3)
		for _, u := range users {
			repository.AddUser(u)
		}

		handler := query.NewListUsersHandler(repository, logger)

		q := query.ListUsersQuery{
			Page:  0,
			Limit: 10,
		}

		result, err := handler.Handle(context.Background(), q)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 1, result.Page)
	})

	t.Run("returns error when repository fails", func(t *testing.T) {
		repository := testutil.NewMockUserRepository()
		repository.ListError = errors.New("database error")
		logger := testutil.NewNoopLogger()

		handler := query.NewListUsersHandler(repository, logger)

		q := query.ListUsersQuery{
			Page:  1,
			Limit: 10,
		}

		result, err := handler.Handle(context.Background(), q)

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "list users")
	})

	t.Run("returns empty page when offset exceeds total", func(t *testing.T) {
		repository := testutil.NewMockUserRepository()
		logger := testutil.NewNoopLogger()

		users := testutil.CreateUsers(3)
		for _, u := range users {
			repository.AddUser(u)
		}

		handler := query.NewListUsersHandler(repository, logger)

		q := query.ListUsersQuery{
			Page:  10,
			Limit: 10,
		}

		result, err := handler.Handle(context.Background(), q)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Empty(t, result.Items)
		assert.Equal(t, int64(3), result.Total)
	})

	t.Run("includes pagination metadata", func(t *testing.T) {
		repository := testutil.NewMockUserRepository()
		logger := testutil.NewNoopLogger()

		users := testutil.CreateUsers(25)
		for _, u := range users {
			repository.AddUser(u)
		}

		handler := query.NewListUsersHandler(repository, logger)

		q := query.ListUsersQuery{
			Page:  2,
			Limit: 10,
		}

		result, err := handler.Handle(context.Background(), q)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, int64(25), result.Total)
		assert.Equal(t, 2, result.Page)
		assert.Equal(t, 10, result.Limit)
		assert.Equal(t, 3, result.TotalPages)
		assert.True(t, result.HasNext)
		assert.True(t, result.HasPrev)
	})
}
