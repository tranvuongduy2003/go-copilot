package authcommand

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tranvuongduy2003/go-copilot/internal/domain/role"
	"github.com/tranvuongduy2003/go-copilot/internal/domain/user"
	"github.com/tranvuongduy2003/go-copilot/pkg/testutil"
)

func TestRegisterHandler_Handle(t *testing.T) {
	ctx := context.Background()

	createDefaultRole := func() *role.Role {
		defaultRole, _ := role.NewRole(role.NewRoleParams{
			Name:        "user",
			DisplayName: "User",
			IsDefault:   true,
		})
		return defaultRole
	}

	tests := []struct {
		name        string
		setupMocks  func(*testutil.MockUserRepository, *testutil.MockRoleRepository, *testutil.MockPermissionRepository, *testutil.MockRefreshTokenRepository, *testutil.MockTokenGenerator, *testutil.MockPasswordHasher)
		command     RegisterCommand
		wantErr     bool
		errContains string
		checkResult func(*testing.T, *testutil.MockUserRepository, *testutil.MockRefreshTokenRepository)
	}{
		{
			name: "successfully register new user",
			setupMocks: func(userRepo *testutil.MockUserRepository, roleRepo *testutil.MockRoleRepository, permissionRepo *testutil.MockPermissionRepository, tokenRepo *testutil.MockRefreshTokenRepository, tokenGen *testutil.MockTokenGenerator, passwordHasher *testutil.MockPasswordHasher) {
				defaultRole := createDefaultRole()
				roleRepo.AddRole(defaultRole)
			},
			command: RegisterCommand{
				Email:     "newuser@example.com",
				Password:  "SecurePass123!",
				FullName:  "New User",
				IPAddress: net.ParseIP("192.168.1.1"),
				UserAgent: "Mozilla/5.0",
			},
			wantErr: false,
			checkResult: func(t *testing.T, userRepo *testutil.MockUserRepository, tokenRepo *testutil.MockRefreshTokenRepository) {
				assert.Len(t, userRepo.Users, 1)
				assert.Len(t, tokenRepo.Tokens, 1)
				for _, u := range userRepo.Users {
					assert.Equal(t, "newuser@example.com", u.Email().String())
					assert.Equal(t, "New User", u.FullName().String())
					assert.Equal(t, user.StatusActive, u.Status())
				}
			},
		},
		{
			name: "fail when email already exists",
			setupMocks: func(userRepo *testutil.MockUserRepository, roleRepo *testutil.MockRoleRepository, permissionRepo *testutil.MockPermissionRepository, tokenRepo *testutil.MockRefreshTokenRepository, tokenGen *testutil.MockTokenGenerator, passwordHasher *testutil.MockPasswordHasher) {
				existingUser, _ := user.NewUser(user.NewUserParams{
					Email:        "existing@example.com",
					PasswordHash: "$2a$10$hashedpassword",
					FullName:     "Existing User",
				})
				userRepo.AddUser(existingUser)
			},
			command: RegisterCommand{
				Email:     "existing@example.com",
				Password:  "SecurePass123!",
				FullName:  "New User",
				IPAddress: net.ParseIP("192.168.1.1"),
				UserAgent: "Mozilla/5.0",
			},
			wantErr:     true,
			errContains: "already exists",
		},
		{
			name: "fail when password is too weak",
			setupMocks: func(userRepo *testutil.MockUserRepository, roleRepo *testutil.MockRoleRepository, permissionRepo *testutil.MockPermissionRepository, tokenRepo *testutil.MockRefreshTokenRepository, tokenGen *testutil.MockTokenGenerator, passwordHasher *testutil.MockPasswordHasher) {
			},
			command: RegisterCommand{
				Email:     "newuser@example.com",
				Password:  "weak",
				FullName:  "New User",
				IPAddress: net.ParseIP("192.168.1.1"),
				UserAgent: "Mozilla/5.0",
			},
			wantErr:     true,
			errContains: "password",
		},
		{
			name: "successfully register without default role",
			setupMocks: func(userRepo *testutil.MockUserRepository, roleRepo *testutil.MockRoleRepository, permissionRepo *testutil.MockPermissionRepository, tokenRepo *testutil.MockRefreshTokenRepository, tokenGen *testutil.MockTokenGenerator, passwordHasher *testutil.MockPasswordHasher) {
			},
			command: RegisterCommand{
				Email:     "newuser@example.com",
				Password:  "SecurePass123!",
				FullName:  "New User",
				IPAddress: net.ParseIP("192.168.1.1"),
				UserAgent: "Mozilla/5.0",
			},
			wantErr: false,
			checkResult: func(t *testing.T, userRepo *testutil.MockUserRepository, tokenRepo *testutil.MockRefreshTokenRepository) {
				assert.Len(t, userRepo.Users, 1)
				for _, u := range userRepo.Users {
					assert.Len(t, u.RoleIDs(), 0)
				}
			},
		},
		{
			name: "fail when email format is invalid",
			setupMocks: func(userRepo *testutil.MockUserRepository, roleRepo *testutil.MockRoleRepository, permissionRepo *testutil.MockPermissionRepository, tokenRepo *testutil.MockRefreshTokenRepository, tokenGen *testutil.MockTokenGenerator, passwordHasher *testutil.MockPasswordHasher) {
			},
			command: RegisterCommand{
				Email:     "invalid-email",
				Password:  "SecurePass123!",
				FullName:  "New User",
				IPAddress: net.ParseIP("192.168.1.1"),
				UserAgent: "Mozilla/5.0",
			},
			wantErr:     true,
			errContains: "email",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := testutil.NewMockUserRepository()
			roleRepo := testutil.NewMockRoleRepository()
			permissionRepo := testutil.NewMockPermissionRepository()
			tokenRepo := testutil.NewMockRefreshTokenRepository()
			tokenGen := testutil.NewMockTokenGenerator()
			passwordHasher := testutil.NewMockPasswordHasher()
			eventBus := testutil.NewMockEventBus()
			logger := testutil.NewNoopLogger()

			tt.setupMocks(userRepo, roleRepo, permissionRepo, tokenRepo, tokenGen, passwordHasher)

			handler := NewRegisterHandler(RegisterHandlerParams{
				UserRepository:         userRepo,
				RoleRepository:         roleRepo,
				PermissionRepository:   permissionRepo,
				RefreshTokenRepository: tokenRepo,
				TokenGenerator:         tokenGen,
				PasswordHasher:         passwordHasher,
				EventBus:               eventBus,
				RefreshTokenTTL:        24 * time.Hour,
				Logger:                 logger,
			})

			result, err := handler.Handle(ctx, tt.command)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.NotEmpty(t, result.AccessToken)
				assert.NotEmpty(t, result.RefreshToken)
				assert.NotNil(t, result.User)
				if tt.checkResult != nil {
					tt.checkResult(t, userRepo, tokenRepo)
				}
			}
		})
	}
}

func TestRegisterHandler_Handle_AssignsDefaultRole(t *testing.T) {
	ctx := context.Background()
	userRepo := testutil.NewMockUserRepository()
	roleRepo := testutil.NewMockRoleRepository()
	permissionRepo := testutil.NewMockPermissionRepository()
	tokenRepo := testutil.NewMockRefreshTokenRepository()
	tokenGen := testutil.NewMockTokenGenerator()
	passwordHasher := testutil.NewMockPasswordHasher()
	eventBus := testutil.NewMockEventBus()
	logger := testutil.NewNoopLogger()

	defaultRole, _ := role.NewRole(role.NewRoleParams{
		Name:        "user",
		DisplayName: "User",
		IsDefault:   true,
	})
	roleRepo.AddRole(defaultRole)

	handler := NewRegisterHandler(RegisterHandlerParams{
		UserRepository:         userRepo,
		RoleRepository:         roleRepo,
		PermissionRepository:   permissionRepo,
		RefreshTokenRepository: tokenRepo,
		TokenGenerator:         tokenGen,
		PasswordHasher:         passwordHasher,
		EventBus:               eventBus,
		RefreshTokenTTL:        24 * time.Hour,
		Logger:                 logger,
	})

	result, err := handler.Handle(ctx, RegisterCommand{
		Email:     "newuser@example.com",
		Password:  "SecurePass123!",
		FullName:  "New User",
		IPAddress: net.ParseIP("192.168.1.1"),
		UserAgent: "Mozilla/5.0",
	})

	require.NoError(t, err)
	require.NotNil(t, result)

	for _, u := range userRepo.Users {
		assert.Len(t, u.RoleIDs(), 1)
		assert.True(t, u.HasRole(defaultRole.ID()))
	}
}

func TestRegisterHandler_Handle_PublishesEvents(t *testing.T) {
	ctx := context.Background()
	userRepo := testutil.NewMockUserRepository()
	roleRepo := testutil.NewMockRoleRepository()
	permissionRepo := testutil.NewMockPermissionRepository()
	tokenRepo := testutil.NewMockRefreshTokenRepository()
	tokenGen := testutil.NewMockTokenGenerator()
	passwordHasher := testutil.NewMockPasswordHasher()
	eventBus := testutil.NewMockEventBus()
	logger := testutil.NewNoopLogger()

	handler := NewRegisterHandler(RegisterHandlerParams{
		UserRepository:         userRepo,
		RoleRepository:         roleRepo,
		PermissionRepository:   permissionRepo,
		RefreshTokenRepository: tokenRepo,
		TokenGenerator:         tokenGen,
		PasswordHasher:         passwordHasher,
		EventBus:               eventBus,
		RefreshTokenTTL:        24 * time.Hour,
		Logger:                 logger,
	})

	_, err := handler.Handle(ctx, RegisterCommand{
		Email:     "newuser@example.com",
		Password:  "SecurePass123!",
		FullName:  "New User",
		IPAddress: net.ParseIP("192.168.1.1"),
		UserAgent: "Mozilla/5.0",
	})

	require.NoError(t, err)
	assert.NotEmpty(t, eventBus.PublishedEvents)
}
