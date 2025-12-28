package permission

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewResource(t *testing.T) {
	tests := []struct {
		name        string
		value       string
		want        string
		wantErr     bool
		errContains string
	}{
		{
			name:    "valid resource",
			value:   "users",
			want:    "users",
			wantErr: false,
		},
		{
			name:    "valid resource with underscore",
			value:   "user_profiles",
			want:    "user_profiles",
			wantErr: false,
		},
		{
			name:    "valid resource with numbers",
			value:   "users2",
			want:    "users2",
			wantErr: false,
		},
		{
			name:    "uppercase normalized to lowercase",
			value:   "USERS",
			want:    "users",
			wantErr: false,
		},
		{
			name:    "whitespace trimmed",
			value:   "  users  ",
			want:    "users",
			wantErr: false,
		},
		{
			name:        "empty resource",
			value:       "",
			wantErr:     true,
			errContains: "resource cannot be empty",
		},
		{
			name:        "whitespace only",
			value:       "   ",
			wantErr:     true,
			errContains: "resource cannot be empty",
		},
		{
			name:        "starts with number",
			value:       "123users",
			wantErr:     true,
			errContains: "must be lowercase alphanumeric",
		},
		{
			name:        "contains special characters",
			value:       "users@profiles",
			wantErr:     true,
			errContains: "must be lowercase alphanumeric",
		},
		{
			name:        "contains hyphen",
			value:       "user-profiles",
			wantErr:     true,
			errContains: "must be lowercase alphanumeric",
		},
		{
			name:        "too long",
			value:       "a" + string(make([]byte, 100)),
			wantErr:     true,
			errContains: "cannot exceed 100 characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource, err := NewResource(tt.value)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, resource.String())
			}
		})
	}
}

func TestResource_Equals(t *testing.T) {
	r1, _ := NewResource("users")
	r2, _ := NewResource("users")
	r3, _ := NewResource("posts")

	assert.True(t, r1.Equals(r2))
	assert.False(t, r1.Equals(r3))
}

func TestNewAction(t *testing.T) {
	tests := []struct {
		name        string
		value       string
		want        Action
		wantErr     bool
		errContains string
	}{
		{
			name:    "standard action - create",
			value:   "create",
			want:    ActionCreate,
			wantErr: false,
		},
		{
			name:    "standard action - read",
			value:   "read",
			want:    ActionRead,
			wantErr: false,
		},
		{
			name:    "standard action - update",
			value:   "update",
			want:    ActionUpdate,
			wantErr: false,
		},
		{
			name:    "standard action - delete",
			value:   "delete",
			want:    ActionDelete,
			wantErr: false,
		},
		{
			name:    "standard action - list",
			value:   "list",
			want:    ActionList,
			wantErr: false,
		},
		{
			name:    "standard action - manage",
			value:   "manage",
			want:    ActionManage,
			wantErr: false,
		},
		{
			name:    "standard action - admin",
			value:   "admin",
			want:    ActionAdmin,
			wantErr: false,
		},
		{
			name:    "uppercase normalized",
			value:   "CREATE",
			want:    ActionCreate,
			wantErr: false,
		},
		{
			name:    "custom valid action",
			value:   "export_csv",
			want:    Action("export_csv"),
			wantErr: false,
		},
		{
			name:        "empty action",
			value:       "",
			wantErr:     true,
			errContains: "action cannot be empty",
		},
		{
			name:        "invalid custom action format",
			value:       "export-csv",
			wantErr:     true,
			errContains: "must be lowercase alphanumeric",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			action, err := NewAction(tt.value)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, action)
			}
		})
	}
}

func TestAction_IsStandard(t *testing.T) {
	assert.True(t, ActionCreate.IsStandard())
	assert.True(t, ActionRead.IsStandard())
	assert.True(t, ActionAdmin.IsStandard())

	customAction, _ := NewAction("custom_action")
	assert.False(t, customAction.IsStandard())
}

func TestParsePermissionCode(t *testing.T) {
	tests := []struct {
		name        string
		code        string
		wantRes     string
		wantAction  string
		wantErr     bool
		errContains string
	}{
		{
			name:       "valid permission code",
			code:       "users:create",
			wantRes:    "users",
			wantAction: "create",
			wantErr:    false,
		},
		{
			name:       "valid permission code with underscore",
			code:       "user_profiles:read",
			wantRes:    "user_profiles",
			wantAction: "read",
			wantErr:    false,
		},
		{
			name:        "missing colon",
			code:        "userscreate",
			wantErr:     true,
			errContains: "must be in format",
		},
		{
			name:        "empty string",
			code:        "",
			wantErr:     true,
			errContains: "must be in format",
		},
		{
			name:        "too many colons",
			code:        "users:create:extra",
			wantErr:     true,
			errContains: "must be in format",
		},
		{
			name:        "invalid resource",
			code:        "123users:create",
			wantErr:     true,
			errContains: "invalid resource",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, err := ParsePermissionCode(tt.code)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantRes, code.Resource().String())
				assert.Equal(t, tt.wantAction, code.Action().String())
			}
		})
	}
}

func TestPermissionCode_String(t *testing.T) {
	code, _ := ParsePermissionCode("users:create")
	assert.Equal(t, "users:create", code.String())
}

func TestPermissionCode_Equals(t *testing.T) {
	code1, _ := ParsePermissionCode("users:create")
	code2, _ := ParsePermissionCode("users:create")
	code3, _ := ParsePermissionCode("users:delete")

	assert.True(t, code1.Equals(code2))
	assert.False(t, code1.Equals(code3))
}
