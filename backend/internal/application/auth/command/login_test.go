package authcommand

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tranvuongduy2003/go-copilot/internal/domain/user"
	"github.com/tranvuongduy2003/go-copilot/pkg/testutil"
)

func TestLoginHandler_Handle(t *testing.T) {
	ctx := context.Background()

	createActiveUser := func() *user.User {
		now := time.Now().UTC()
		testUser, _ := user.ReconstructUser(user.ReconstructUserParams{
			ID:           uuid.New(),
			Email:        "test@example.com",
			PasswordHash: "$2a$10$hashedpassword",
			FullName:     "Test User",
			Status:       user.StatusActive,
			CreatedAt:    now,
			UpdatedAt:    now,
		})
		return testUser
	}

	createInactiveUser := func() *user.User {
		now := time.Now().UTC()
		testUser, _ := user.ReconstructUser(user.ReconstructUserParams{
			ID:           uuid.New(),
			Email:        "inactive@example.com",
			PasswordHash: "$2a$10$hashedpassword",
			FullName:     "Inactive User",
			Status:       user.StatusInactive,
			CreatedAt:    now,
			UpdatedAt:    now,
		})
		return testUser
	}

	tests := []struct {
		name        string
		setupMocks  func(*testutil.MockUserRepository, *testutil.MockRoleRepository, *testutil.MockPermissionRepository, *testutil.MockRefreshTokenRepository, *testutil.MockTokenGenerator, *testutil.MockPasswordHasher, *testutil.MockAccountLockout)
		command     LoginCommand
		wantErr     bool
		errContains string
		checkResult func(*testing.T, *testutil.MockRefreshTokenRepository)
	}{
		{
			name: "successfully login with valid credentials",
			setupMocks: func(userRepo *testutil.MockUserRepository, roleRepo *testutil.MockRoleRepository, permissionRepo *testutil.MockPermissionRepository, tokenRepo *testutil.MockRefreshTokenRepository, tokenGen *testutil.MockTokenGenerator, passwordHasher *testutil.MockPasswordHasher, lockout *testutil.MockAccountLockout) {
				testUser := createActiveUser()
				userRepo.AddUser(testUser)
				passwordHasher.VerifyResult = true
			},
			command: LoginCommand{
				Email:     "test@example.com",
				Password:  "correctpassword",
				IPAddress: net.ParseIP("192.168.1.1"),
				UserAgent: "Mozilla/5.0",
			},
			wantErr: false,
			checkResult: func(t *testing.T, tokenRepo *testutil.MockRefreshTokenRepository) {
				assert.Len(t, tokenRepo.Tokens, 1)
			},
		},
		{
			name: "fail when user not found",
			setupMocks: func(userRepo *testutil.MockUserRepository, roleRepo *testutil.MockRoleRepository, permissionRepo *testutil.MockPermissionRepository, tokenRepo *testutil.MockRefreshTokenRepository, tokenGen *testutil.MockTokenGenerator, passwordHasher *testutil.MockPasswordHasher, lockout *testutil.MockAccountLockout) {
			},
			command: LoginCommand{
				Email:     "nonexistent@example.com",
				Password:  "anypassword",
				IPAddress: net.ParseIP("192.168.1.1"),
				UserAgent: "Mozilla/5.0",
			},
			wantErr:     true,
			errContains: "not authorized",
		},
		{
			name: "fail when password is incorrect",
			setupMocks: func(userRepo *testutil.MockUserRepository, roleRepo *testutil.MockRoleRepository, permissionRepo *testutil.MockPermissionRepository, tokenRepo *testutil.MockRefreshTokenRepository, tokenGen *testutil.MockTokenGenerator, passwordHasher *testutil.MockPasswordHasher, lockout *testutil.MockAccountLockout) {
				testUser := createActiveUser()
				userRepo.AddUser(testUser)
				passwordHasher.VerifyResult = false
			},
			command: LoginCommand{
				Email:     "test@example.com",
				Password:  "wrongpassword",
				IPAddress: net.ParseIP("192.168.1.1"),
				UserAgent: "Mozilla/5.0",
			},
			wantErr:     true,
			errContains: "not authorized",
		},
		{
			name: "fail when account is inactive",
			setupMocks: func(userRepo *testutil.MockUserRepository, roleRepo *testutil.MockRoleRepository, permissionRepo *testutil.MockPermissionRepository, tokenRepo *testutil.MockRefreshTokenRepository, tokenGen *testutil.MockTokenGenerator, passwordHasher *testutil.MockPasswordHasher, lockout *testutil.MockAccountLockout) {
				inactiveUser := createInactiveUser()
				userRepo.AddUser(inactiveUser)
				passwordHasher.VerifyResult = true
			},
			command: LoginCommand{
				Email:     "inactive@example.com",
				Password:  "correctpassword",
				IPAddress: net.ParseIP("192.168.1.1"),
				UserAgent: "Mozilla/5.0",
			},
			wantErr:     true,
			errContains: "not active",
		},
		{
			name: "fail when account is locked",
			setupMocks: func(userRepo *testutil.MockUserRepository, roleRepo *testutil.MockRoleRepository, permissionRepo *testutil.MockPermissionRepository, tokenRepo *testutil.MockRefreshTokenRepository, tokenGen *testutil.MockTokenGenerator, passwordHasher *testutil.MockPasswordHasher, lockout *testutil.MockAccountLockout) {
				testUser := createActiveUser()
				userRepo.AddUser(testUser)
				lockout.Locked = true
				lockout.RemainingTime = 10 * time.Minute
			},
			command: LoginCommand{
				Email:     "test@example.com",
				Password:  "anypassword",
				IPAddress: net.ParseIP("192.168.1.1"),
				UserAgent: "Mozilla/5.0",
			},
			wantErr:     true,
			errContains: "locked",
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
			lockout := testutil.NewMockAccountLockout()
			eventBus := testutil.NewMockEventBus()
			logger := testutil.NewNoopLogger()

			tt.setupMocks(userRepo, roleRepo, permissionRepo, tokenRepo, tokenGen, passwordHasher, lockout)

			handler := NewLoginHandler(LoginHandlerParams{
				UserRepository:         userRepo,
				RoleRepository:         roleRepo,
				PermissionRepository:   permissionRepo,
				RefreshTokenRepository: tokenRepo,
				TokenGenerator:         tokenGen,
				PasswordHasher:         passwordHasher,
				AccountLockout:         lockout,
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
					tt.checkResult(t, tokenRepo)
				}
			}
		})
	}
}

func TestLoginHandler_Handle_ResetsLockoutOnSuccess(t *testing.T) {
	ctx := context.Background()
	userRepo := testutil.NewMockUserRepository()
	roleRepo := testutil.NewMockRoleRepository()
	permissionRepo := testutil.NewMockPermissionRepository()
	tokenRepo := testutil.NewMockRefreshTokenRepository()
	tokenGen := testutil.NewMockTokenGenerator()
	passwordHasher := testutil.NewMockPasswordHasher()
	lockout := testutil.NewMockAccountLockout()
	eventBus := testutil.NewMockEventBus()
	logger := testutil.NewNoopLogger()

	now := time.Now().UTC()
	testUser, _ := user.ReconstructUser(user.ReconstructUserParams{
		ID:           uuid.New(),
		Email:        "test@example.com",
		PasswordHash: "$2a$10$hashedpassword",
		FullName:     "Test User",
		Status:       user.StatusActive,
		CreatedAt:    now,
		UpdatedAt:    now,
	})
	userRepo.AddUser(testUser)
	passwordHasher.VerifyResult = true
	lockout.AttemptCount = 3

	handler := NewLoginHandler(LoginHandlerParams{
		UserRepository:         userRepo,
		RoleRepository:         roleRepo,
		PermissionRepository:   permissionRepo,
		RefreshTokenRepository: tokenRepo,
		TokenGenerator:         tokenGen,
		PasswordHasher:         passwordHasher,
		AccountLockout:         lockout,
		EventBus:               eventBus,
		RefreshTokenTTL:        24 * time.Hour,
		Logger:                 logger,
	})

	_, err := handler.Handle(ctx, LoginCommand{
		Email:     "test@example.com",
		Password:  "correctpassword",
		IPAddress: net.ParseIP("192.168.1.1"),
		UserAgent: "Mozilla/5.0",
	})

	require.NoError(t, err)
	assert.Equal(t, 0, lockout.AttemptCount)
}

func TestLoginHandler_Handle_RecordsFailedAttempts(t *testing.T) {
	ctx := context.Background()
	userRepo := testutil.NewMockUserRepository()
	roleRepo := testutil.NewMockRoleRepository()
	permissionRepo := testutil.NewMockPermissionRepository()
	tokenRepo := testutil.NewMockRefreshTokenRepository()
	tokenGen := testutil.NewMockTokenGenerator()
	passwordHasher := testutil.NewMockPasswordHasher()
	lockout := testutil.NewMockAccountLockout()
	eventBus := testutil.NewMockEventBus()
	logger := testutil.NewNoopLogger()

	now := time.Now().UTC()
	testUser, _ := user.ReconstructUser(user.ReconstructUserParams{
		ID:           uuid.New(),
		Email:        "test@example.com",
		PasswordHash: "$2a$10$hashedpassword",
		FullName:     "Test User",
		Status:       user.StatusActive,
		CreatedAt:    now,
		UpdatedAt:    now,
	})
	userRepo.AddUser(testUser)
	passwordHasher.VerifyResult = false

	handler := NewLoginHandler(LoginHandlerParams{
		UserRepository:         userRepo,
		RoleRepository:         roleRepo,
		PermissionRepository:   permissionRepo,
		RefreshTokenRepository: tokenRepo,
		TokenGenerator:         tokenGen,
		PasswordHasher:         passwordHasher,
		AccountLockout:         lockout,
		EventBus:               eventBus,
		RefreshTokenTTL:        24 * time.Hour,
		Logger:                 logger,
	})

	_, _ = handler.Handle(ctx, LoginCommand{
		Email:     "test@example.com",
		Password:  "wrongpassword",
		IPAddress: net.ParseIP("192.168.1.1"),
		UserAgent: "Mozilla/5.0",
	})

	assert.Equal(t, 1, lockout.AttemptCount)
}

func TestLoginHandler_Handle_PublishesLoginEvent(t *testing.T) {
	ctx := context.Background()
	userRepo := testutil.NewMockUserRepository()
	roleRepo := testutil.NewMockRoleRepository()
	permissionRepo := testutil.NewMockPermissionRepository()
	tokenRepo := testutil.NewMockRefreshTokenRepository()
	tokenGen := testutil.NewMockTokenGenerator()
	passwordHasher := testutil.NewMockPasswordHasher()
	lockout := testutil.NewMockAccountLockout()
	eventBus := testutil.NewMockEventBus()
	logger := testutil.NewNoopLogger()

	now := time.Now().UTC()
	testUser, _ := user.ReconstructUser(user.ReconstructUserParams{
		ID:           uuid.New(),
		Email:        "test@example.com",
		PasswordHash: "$2a$10$hashedpassword",
		FullName:     "Test User",
		Status:       user.StatusActive,
		CreatedAt:    now,
		UpdatedAt:    now,
	})
	userRepo.AddUser(testUser)
	passwordHasher.VerifyResult = true

	handler := NewLoginHandler(LoginHandlerParams{
		UserRepository:         userRepo,
		RoleRepository:         roleRepo,
		PermissionRepository:   permissionRepo,
		RefreshTokenRepository: tokenRepo,
		TokenGenerator:         tokenGen,
		PasswordHasher:         passwordHasher,
		AccountLockout:         lockout,
		EventBus:               eventBus,
		RefreshTokenTTL:        24 * time.Hour,
		Logger:                 logger,
	})

	_, err := handler.Handle(ctx, LoginCommand{
		Email:     "test@example.com",
		Password:  "correctpassword",
		IPAddress: net.ParseIP("192.168.1.1"),
		UserAgent: "Mozilla/5.0",
	})

	require.NoError(t, err)
	assert.NotEmpty(t, eventBus.PublishedEvents)
}
