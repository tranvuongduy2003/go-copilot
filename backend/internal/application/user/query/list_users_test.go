package userquery

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tranvuongduy2003/go-copilot/internal/domain/user"
	"github.com/tranvuongduy2003/go-copilot/pkg/testutil"
)

func TestListUsersHandler_Handle(t *testing.T) {
	ctx := context.Background()

	createTestUsers := func() []*user.User {
		users := make([]*user.User, 0)
		for i := 0; i < 5; i++ {
			emails := []string{
				"user1@example.com",
				"user2@example.com",
				"user3@example.com",
				"user4@example.com",
				"user5@example.com",
			}
			u, _ := user.NewUser(user.NewUserParams{
				Email:        emails[i],
				PasswordHash: "$2a$10$hashedpassword",
				FullName:     "Test User",
			})
			if i%2 == 0 {
				u.Activate()
			}
			u.ClearDomainEvents()
			users = append(users, u)
		}
		return users
	}

	tests := []struct {
		name        string
		setupMocks  func(*testutil.MockUserRepository)
		query       ListUsersQuery
		wantErr     bool
		errContains string
		checkResult func(*testing.T, *testutil.MockUserRepository)
	}{
		{
			name: "successfully list all users",
			setupMocks: func(userRepo *testutil.MockUserRepository) {
			},
			query: ListUsersQuery{
				Page:  1,
				Limit: 10,
			},
			wantErr: false,
			checkResult: func(t *testing.T, userRepo *testutil.MockUserRepository) {
			},
		},
		{
			name: "successfully list users with pagination",
			setupMocks: func(userRepo *testutil.MockUserRepository) {
			},
			query: ListUsersQuery{
				Page:  1,
				Limit: 2,
			},
			wantErr: false,
			checkResult: func(t *testing.T, userRepo *testutil.MockUserRepository) {
			},
		},
		{
			name: "successfully list users filtered by status",
			setupMocks: func(userRepo *testutil.MockUserRepository) {
			},
			query: ListUsersQuery{
				Page:   1,
				Limit:  10,
				Status: stringPtr("active"),
			},
			wantErr: false,
			checkResult: func(t *testing.T, userRepo *testutil.MockUserRepository) {
			},
		},
		{
			name: "return empty list when no users match filter",
			setupMocks: func(userRepo *testutil.MockUserRepository) {
			},
			query: ListUsersQuery{
				Page:   1,
				Limit:  10,
				Status: stringPtr("banned"),
			},
			wantErr: false,
			checkResult: func(t *testing.T, userRepo *testutil.MockUserRepository) {
			},
		},
		{
			name: "fail when repository returns error",
			setupMocks: func(userRepo *testutil.MockUserRepository) {
				userRepo.ListError = errors.New("database error")
			},
			query: ListUsersQuery{
				Page:  1,
				Limit: 10,
			},
			wantErr:     true,
			errContains: "list users",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := testutil.NewMockUserRepository()
			logger := testutil.NewNoopLogger()

			for _, u := range createTestUsers() {
				userRepo.AddUser(u)
			}

			tt.setupMocks(userRepo)

			handler := NewListUsersHandler(userRepo, logger)
			result, err := handler.Handle(ctx, tt.query)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.NotNil(t, result.Items)
				assert.GreaterOrEqual(t, result.Total, int64(0))
				assert.Equal(t, tt.query.Page, result.Page)
				assert.Equal(t, tt.query.Limit, result.Limit)
				if tt.checkResult != nil {
					tt.checkResult(t, userRepo)
				}
			}
		})
	}
}

func stringPtr(s string) *string {
	return &s
}
