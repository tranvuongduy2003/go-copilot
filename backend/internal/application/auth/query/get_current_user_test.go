package authquery

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

func TestGetCurrentUserHandler_Handle(t *testing.T) {
	ctx := context.Background()

	createTestUser := func(roleIDs []uuid.UUID) *user.User {
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

	createTestRole := func(name string, permissionIDs []uuid.UUID) *role.Role {
		now := time.Now().UTC()
		testRole, _ := role.ReconstructRole(role.ReconstructRoleParams{
			ID:            uuid.New(),
			Name:          name,
			DisplayName:   name,
			Description:   "Test role",
			PermissionIDs: permissionIDs,
			IsSystem:      false,
			IsDefault:     false,
			Priority:      0,
			CreatedAt:     now,
			UpdatedAt:     now,
		})
		return testRole
	}

	createTestPermission := func(resource, action string) *permission.Permission {
		now := time.Now().UTC()
		perm, _ := permission.ReconstructPermission(permission.ReconstructPermissionParams{
			ID:          uuid.New(),
			Resource:    resource,
			Action:      action,
			Description: resource + ":" + action,
			IsSystem:    false,
			CreatedAt:   now,
			UpdatedAt:   now,
		})
		return perm
	}

	tests := []struct {
		name        string
		setupMocks  func(*testutil.MockUserRepository, *testutil.MockRoleRepository, *testutil.MockPermissionRepository) uuid.UUID
		wantErr     bool
		errContains string
		checkResult func(*testing.T, *testutil.MockUserRepository)
	}{
		{
			name: "successfully get current user with roles and permissions",
			setupMocks: func(userRepo *testutil.MockUserRepository, roleRepo *testutil.MockRoleRepository, permRepo *testutil.MockPermissionRepository) uuid.UUID {
				perm1 := createTestPermission("users", "read")
				perm2 := createTestPermission("users", "list")
				permRepo.AddPermission(perm1)
				permRepo.AddPermission(perm2)

				testRole := createTestRole("admin", []uuid.UUID{perm1.ID(), perm2.ID()})
				roleRepo.AddRole(testRole)

				testUser := createTestUser([]uuid.UUID{testRole.ID()})
				userRepo.AddUser(testUser)

				return testUser.ID()
			},
			wantErr: false,
		},
		{
			name: "successfully get current user without roles",
			setupMocks: func(userRepo *testutil.MockUserRepository, roleRepo *testutil.MockRoleRepository, permRepo *testutil.MockPermissionRepository) uuid.UUID {
				testUser := createTestUser([]uuid.UUID{})
				userRepo.AddUser(testUser)
				return testUser.ID()
			},
			wantErr: false,
		},
		{
			name: "fail when user not found",
			setupMocks: func(userRepo *testutil.MockUserRepository, roleRepo *testutil.MockRoleRepository, permRepo *testutil.MockPermissionRepository) uuid.UUID {
				return uuid.New()
			},
			wantErr:     true,
			errContains: "find user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := testutil.NewMockUserRepository()
			roleRepo := testutil.NewMockRoleRepository()
			permRepo := testutil.NewMockPermissionRepository()
			logger := testutil.NewNoopLogger()

			userID := tt.setupMocks(userRepo, roleRepo, permRepo)

			handler := NewGetCurrentUserHandler(GetCurrentUserHandlerParams{
				UserRepository:       userRepo,
				RoleRepository:       roleRepo,
				PermissionRepository: permRepo,
				Logger:               logger,
			})

			result, err := handler.Handle(ctx, GetCurrentUserQuery{UserID: userID})

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.NotEmpty(t, result.ID)
				assert.NotEmpty(t, result.Email)
			}
		})
	}
}

func TestGetCurrentUserHandler_Handle_WithMultipleRoles(t *testing.T) {
	ctx := context.Background()
	userRepo := testutil.NewMockUserRepository()
	roleRepo := testutil.NewMockRoleRepository()
	permRepo := testutil.NewMockPermissionRepository()
	logger := testutil.NewNoopLogger()

	now := time.Now().UTC()

	perm1, _ := permission.ReconstructPermission(permission.ReconstructPermissionParams{
		ID:          uuid.New(),
		Resource:    "users",
		Action:      "read",
		Description: "Read users",
		IsSystem:    false,
		CreatedAt:   now,
		UpdatedAt:   now,
	})
	perm2, _ := permission.ReconstructPermission(permission.ReconstructPermissionParams{
		ID:          uuid.New(),
		Resource:    "roles",
		Action:      "read",
		Description: "Read roles",
		IsSystem:    false,
		CreatedAt:   now,
		UpdatedAt:   now,
	})
	permRepo.AddPermission(perm1)
	permRepo.AddPermission(perm2)

	role1, _ := role.ReconstructRole(role.ReconstructRoleParams{
		ID:            uuid.New(),
		Name:          "admin",
		DisplayName:   "Admin",
		PermissionIDs: []uuid.UUID{perm1.ID()},
		IsSystem:      false,
		IsDefault:     false,
		CreatedAt:     now,
		UpdatedAt:     now,
	})
	role2, _ := role.ReconstructRole(role.ReconstructRoleParams{
		ID:            uuid.New(),
		Name:          "manager",
		DisplayName:   "Manager",
		PermissionIDs: []uuid.UUID{perm2.ID()},
		IsSystem:      false,
		IsDefault:     false,
		CreatedAt:     now,
		UpdatedAt:     now,
	})
	roleRepo.AddRole(role1)
	roleRepo.AddRole(role2)

	testUser, _ := user.ReconstructUser(user.ReconstructUserParams{
		ID:           uuid.New(),
		Email:        "test@example.com",
		PasswordHash: "$2a$10$hashedpassword",
		FullName:     "Test User",
		Status:       user.StatusActive,
		RoleIDs:      []uuid.UUID{role1.ID(), role2.ID()},
		CreatedAt:    now,
		UpdatedAt:    now,
	})
	userRepo.AddUser(testUser)

	handler := NewGetCurrentUserHandler(GetCurrentUserHandlerParams{
		UserRepository:       userRepo,
		RoleRepository:       roleRepo,
		PermissionRepository: permRepo,
		Logger:               logger,
	})

	result, err := handler.Handle(ctx, GetCurrentUserQuery{UserID: testUser.ID()})

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Len(t, result.Roles, 2)
	assert.Len(t, result.Permissions, 2)
}
