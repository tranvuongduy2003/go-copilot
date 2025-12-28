package userquery

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tranvuongduy2003/go-copilot/internal/domain/role"
	"github.com/tranvuongduy2003/go-copilot/internal/domain/user"
	"github.com/tranvuongduy2003/go-copilot/pkg/testutil"
)

func TestGetUserRolesHandler_Handle(t *testing.T) {
	ctx := context.Background()

	editorRole, _ := role.NewRole(role.NewRoleParams{
		Name:        "editor",
		DisplayName: "Editor",
	})

	viewerRole, _ := role.NewRole(role.NewRoleParams{
		Name:        "viewer",
		DisplayName: "Viewer",
	})

	createTestUserWithRoles := func(roleIDs []uuid.UUID) *user.User {
		now := time.Now().UTC()
		testUser, _ := user.ReconstructUser(user.ReconstructUserParams{
			ID:           uuid.New(),
			Email:        "test@example.com",
			PasswordHash: "$2a$10$hashedpassword",
			FullName:     "Test User",
			Status:       user.StatusActive,
			RoleIDs:      roleIDs,
			CreatedAt:    now,
			UpdatedAt:    now,
		})
		return testUser
	}

	createTestUserWithoutRoles := func() *user.User {
		testUser, _ := user.NewUser(user.NewUserParams{
			Email:        "test@example.com",
			PasswordHash: "$2a$10$hashedpassword",
			FullName:     "Test User",
		})
		return testUser
	}

	tests := []struct {
		name        string
		setupMocks  func(*testutil.MockUserRepository, *testutil.MockRoleRepository) *user.User
		query       func(*user.User) GetUserRolesQuery
		wantErr     bool
		errContains string
		checkResult func(*testing.T, int)
	}{
		{
			name: "successfully get user roles",
			setupMocks: func(userRepo *testutil.MockUserRepository, roleRepo *testutil.MockRoleRepository) *user.User {
				roleRepo.AddRole(editorRole)
				roleRepo.AddRole(viewerRole)
				testUser := createTestUserWithRoles([]uuid.UUID{editorRole.ID(), viewerRole.ID()})
				userRepo.AddUser(testUser)
				return testUser
			},
			query: func(u *user.User) GetUserRolesQuery {
				return GetUserRolesQuery{UserID: u.ID()}
			},
			wantErr: false,
			checkResult: func(t *testing.T, count int) {
				assert.Equal(t, 2, count)
			},
		},
		{
			name: "return empty list when user has no roles",
			setupMocks: func(userRepo *testutil.MockUserRepository, roleRepo *testutil.MockRoleRepository) *user.User {
				testUser := createTestUserWithoutRoles()
				userRepo.AddUser(testUser)
				return testUser
			},
			query: func(u *user.User) GetUserRolesQuery {
				return GetUserRolesQuery{UserID: u.ID()}
			},
			wantErr: false,
			checkResult: func(t *testing.T, count int) {
				assert.Equal(t, 0, count)
			},
		},
		{
			name: "fail when user not found",
			setupMocks: func(userRepo *testutil.MockUserRepository, roleRepo *testutil.MockRoleRepository) *user.User {
				testUser := createTestUserWithoutRoles()
				userRepo.AddUser(testUser)
				return testUser
			},
			query: func(u *user.User) GetUserRolesQuery {
				return GetUserRolesQuery{UserID: uuid.New()}
			},
			wantErr:     true,
			errContains: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := testutil.NewMockUserRepository()
			roleRepo := testutil.NewMockRoleRepository()
			logger := testutil.NewNoopLogger()

			testUser := tt.setupMocks(userRepo, roleRepo)

			handler := NewGetUserRolesHandler(userRepo, roleRepo, logger)
			query := tt.query(testUser)

			result, err := handler.Handle(ctx, query)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				if tt.checkResult != nil {
					tt.checkResult(t, len(result))
				}
			}
		})
	}
}
