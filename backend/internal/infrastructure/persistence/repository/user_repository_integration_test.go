//go:build integration

package repository

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tranvuongduy2003/go-copilot/internal/domain/user"
	"github.com/tranvuongduy2003/go-copilot/pkg/testutil"
)

func TestUserRepository_Integration(t *testing.T) {
	suite, cleanup := testutil.SetupIntegrationTest(t)
	defer cleanup()

	repository := NewUserRepository(suite.DatabasePool)

	t.Run("Create", func(t *testing.T) {
		suite.CleanAllTables(t)

		testUser, err := user.NewUser(user.NewUserParams{
			Email:        "integration@test.com",
			PasswordHash: "$2a$10$hashedpassword",
			FullName:     "Integration Test User",
		})
		require.NoError(t, err)

		err = repository.Create(context.Background(), testUser)
		require.NoError(t, err)

		foundUser, err := repository.FindByID(context.Background(), testUser.ID())
		require.NoError(t, err)
		assert.Equal(t, testUser.ID(), foundUser.ID())
		assert.Equal(t, testUser.Email().String(), foundUser.Email().String())
		assert.Equal(t, testUser.FullName().String(), foundUser.FullName().String())
	})

	t.Run("Create duplicate email fails", func(t *testing.T) {
		suite.CleanAllTables(t)

		testUser1, err := user.NewUser(user.NewUserParams{
			Email:        "duplicate@test.com",
			PasswordHash: "$2a$10$hashedpassword",
			FullName:     "User One",
		})
		require.NoError(t, err)

		err = repository.Create(context.Background(), testUser1)
		require.NoError(t, err)

		testUser2, err := user.NewUser(user.NewUserParams{
			Email:        "duplicate@test.com",
			PasswordHash: "$2a$10$hashedpassword",
			FullName:     "User Two",
		})
		require.NoError(t, err)

		err = repository.Create(context.Background(), testUser2)
		assert.ErrorIs(t, err, user.ErrEmailAlreadyExists)
	})

	t.Run("FindByID returns not found for non-existent user", func(t *testing.T) {
		suite.CleanAllTables(t)

		_, err := repository.FindByID(context.Background(), uuid.New())
		assert.ErrorIs(t, err, user.ErrUserNotFound)
	})

	t.Run("FindByEmail case insensitive", func(t *testing.T) {
		suite.CleanAllTables(t)

		testUser, err := user.NewUser(user.NewUserParams{
			Email:        "CaSeTest@example.com",
			PasswordHash: "$2a$10$hashedpassword",
			FullName:     "Case Test User",
		})
		require.NoError(t, err)

		err = repository.Create(context.Background(), testUser)
		require.NoError(t, err)

		foundUser, err := repository.FindByEmail(context.Background(), "casetest@example.com")
		require.NoError(t, err)
		assert.Equal(t, testUser.ID(), foundUser.ID())
	})

	t.Run("Update modifies user", func(t *testing.T) {
		suite.CleanAllTables(t)

		testUser, err := user.NewUser(user.NewUserParams{
			Email:        "update@test.com",
			PasswordHash: "$2a$10$hashedpassword",
			FullName:     "Original Name",
		})
		require.NoError(t, err)

		err = repository.Create(context.Background(), testUser)
		require.NoError(t, err)

		err = testUser.UpdateProfile("Updated Name")
		require.NoError(t, err)

		err = repository.Update(context.Background(), testUser)
		require.NoError(t, err)

		foundUser, err := repository.FindByID(context.Background(), testUser.ID())
		require.NoError(t, err)
		assert.Equal(t, "Updated Name", foundUser.FullName().String())
	})

	t.Run("Delete soft deletes user", func(t *testing.T) {
		suite.CleanAllTables(t)

		testUser, err := user.NewUser(user.NewUserParams{
			Email:        "delete@test.com",
			PasswordHash: "$2a$10$hashedpassword",
			FullName:     "Delete Test User",
		})
		require.NoError(t, err)

		err = repository.Create(context.Background(), testUser)
		require.NoError(t, err)

		err = repository.Delete(context.Background(), testUser.ID())
		require.NoError(t, err)

		_, err = repository.FindByID(context.Background(), testUser.ID())
		assert.ErrorIs(t, err, user.ErrUserNotFound)
	})

	t.Run("ExistsByEmail returns correct result", func(t *testing.T) {
		suite.CleanAllTables(t)

		exists, err := repository.ExistsByEmail(context.Background(), "nonexistent@test.com")
		require.NoError(t, err)
		assert.False(t, exists)

		testUser, err := user.NewUser(user.NewUserParams{
			Email:        "exists@test.com",
			PasswordHash: "$2a$10$hashedpassword",
			FullName:     "Exists Test User",
		})
		require.NoError(t, err)

		err = repository.Create(context.Background(), testUser)
		require.NoError(t, err)

		exists, err = repository.ExistsByEmail(context.Background(), "exists@test.com")
		require.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("List with pagination", func(t *testing.T) {
		suite.CleanAllTables(t)

		for i := 0; i < 5; i++ {
			testUser, err := user.NewUser(user.NewUserParams{
				Email:        "list" + string(rune('0'+i)) + "@test.com",
				PasswordHash: "$2a$10$hashedpassword",
				FullName:     "List User " + string(rune('0'+i)),
			})
			require.NoError(t, err)
			err = repository.Create(context.Background(), testUser)
			require.NoError(t, err)
		}

		users, total, err := repository.List(context.Background(), user.ListFilter{}, user.ListPagination{
			Page:  1,
			Limit: 2,
		})
		require.NoError(t, err)
		assert.Equal(t, int64(5), total)
		assert.Len(t, users, 2)
	})

	t.Run("List with status filter", func(t *testing.T) {
		suite.CleanAllTables(t)

		activeUser, err := user.NewUser(user.NewUserParams{
			Email:        "active@test.com",
			PasswordHash: "$2a$10$hashedpassword",
			FullName:     "Active User",
		})
		require.NoError(t, err)
		err = activeUser.Activate()
		require.NoError(t, err)
		err = repository.Create(context.Background(), activeUser)
		require.NoError(t, err)

		pendingUser, err := user.NewUser(user.NewUserParams{
			Email:        "pending@test.com",
			PasswordHash: "$2a$10$hashedpassword",
			FullName:     "Pending User",
		})
		require.NoError(t, err)
		err = repository.Create(context.Background(), pendingUser)
		require.NoError(t, err)

		activeStatus := user.StatusActive
		users, total, err := repository.List(context.Background(), user.ListFilter{
			Status: &activeStatus,
		}, user.ListPagination{
			Page:  1,
			Limit: 10,
		})
		require.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Len(t, users, 1)
		assert.Equal(t, "active@test.com", users[0].Email().String())
	})
}
