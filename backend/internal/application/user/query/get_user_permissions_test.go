package userquery

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tranvuongduy2003/go-copilot/internal/domain/permission"
	"github.com/tranvuongduy2003/go-copilot/internal/domain/role"
	"github.com/tranvuongduy2003/go-copilot/internal/domain/user"
	"github.com/tranvuongduy2003/go-copilot/pkg/testutil"
)

func TestGetUserPermissionsHandler_Handle(t *testing.T) {
	ctx := context.Background()

	createPerm, _ := permission.NewPermission(permission.NewPermissionParams{
		Resource:    "users",
		Action:      "create",
		Description: "Create users",
	})

	readPerm, _ := permission.NewPermission(permission.NewPermissionParams{
		Resource:    "users",
		Action:      "read",
		Description: "Read users",
	})

	editorRole, _ := role.NewRole(role.NewRoleParams{
		Name:          "editor",
		DisplayName:   "Editor",
		PermissionIDs: []uuid.UUID{createPerm.ID(), readPerm.ID()},
	})

	viewerRole, _ := role.NewRole(role.NewRoleParams{
		Name:          "viewer",
		DisplayName:   "Viewer",
		PermissionIDs: []uuid.UUID{readPerm.ID()},
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
		setupMocks  func(*testutil.MockUserRepository, *testutil.MockRoleRepository, *testutil.MockPermissionRepository) *user.User
		query       func(*user.User) GetUserPermissionsQuery
		wantErr     bool
		errContains string
		checkResult func(*testing.T, int)
	}{
		{
			name: "successfully get user permissions from single role",
			setupMocks: func(userRepo *testutil.MockUserRepository, roleRepo *testutil.MockRoleRepository, permRepo *testutil.MockPermissionRepository) *user.User {
				permRepo.AddPermission(createPerm)
				permRepo.AddPermission(readPerm)
				roleRepo.AddRole(editorRole)
				testUser := createTestUserWithRoles([]uuid.UUID{editorRole.ID()})
				userRepo.AddUser(testUser)
				return testUser
			},
			query: func(u *user.User) GetUserPermissionsQuery {
				return GetUserPermissionsQuery{UserID: u.ID()}
			},
			wantErr: false,
			checkResult: func(t *testing.T, count int) {
				assert.Equal(t, 2, count)
			},
		},
		{
			name: "aggregate permissions from multiple roles without duplicates",
			setupMocks: func(userRepo *testutil.MockUserRepository, roleRepo *testutil.MockRoleRepository, permRepo *testutil.MockPermissionRepository) *user.User {
				permRepo.AddPermission(createPerm)
				permRepo.AddPermission(readPerm)
				roleRepo.AddRole(editorRole)
				roleRepo.AddRole(viewerRole)
				testUser := createTestUserWithRoles([]uuid.UUID{editorRole.ID(), viewerRole.ID()})
				userRepo.AddUser(testUser)
				return testUser
			},
			query: func(u *user.User) GetUserPermissionsQuery {
				return GetUserPermissionsQuery{UserID: u.ID()}
			},
			wantErr: false,
			checkResult: func(t *testing.T, count int) {
				assert.Equal(t, 2, count)
			},
		},
		{
			name: "return empty list when user has no roles",
			setupMocks: func(userRepo *testutil.MockUserRepository, roleRepo *testutil.MockRoleRepository, permRepo *testutil.MockPermissionRepository) *user.User {
				testUser := createTestUserWithoutRoles()
				userRepo.AddUser(testUser)
				return testUser
			},
			query: func(u *user.User) GetUserPermissionsQuery {
				return GetUserPermissionsQuery{UserID: u.ID()}
			},
			wantErr: false,
			checkResult: func(t *testing.T, count int) {
				assert.Equal(t, 0, count)
			},
		},
		{
			name: "fail when user not found",
			setupMocks: func(userRepo *testutil.MockUserRepository, roleRepo *testutil.MockRoleRepository, permRepo *testutil.MockPermissionRepository) *user.User {
				testUser := createTestUserWithoutRoles()
				userRepo.AddUser(testUser)
				return testUser
			},
			query: func(u *user.User) GetUserPermissionsQuery {
				return GetUserPermissionsQuery{UserID: uuid.New()}
			},
			wantErr:     true,
			errContains: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := testutil.NewMockUserRepository()
			roleRepo := testutil.NewMockRoleRepository()
			permRepo := testutil.NewMockPermissionRepository()
			logger := testutil.NewNoopLogger()

			testUser := tt.setupMocks(userRepo, roleRepo, permRepo)

			handler := NewGetUserPermissionsHandler(userRepo, roleRepo, permRepo, logger)
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
