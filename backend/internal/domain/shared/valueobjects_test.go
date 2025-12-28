package shared

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name        string
		password    string
		expectError bool
		errContains string
	}{
		{
			name:        "valid password with all requirements",
			password:    "Password123!",
			expectError: false,
		},
		{
			name:        "valid complex password",
			password:    "MySecure@Pass99",
			expectError: false,
		},
		{
			name:        "too short password",
			password:    "Pass1!",
			expectError: true,
			errContains: "at least 8 characters",
		},
		{
			name:        "exactly 8 characters valid password",
			password:    "Pass12!a",
			expectError: false,
		},
		{
			name:        "missing uppercase letter",
			password:    "password123!",
			expectError: true,
			errContains: "uppercase letter",
		},
		{
			name:        "missing lowercase letter",
			password:    "PASSWORD123!",
			expectError: true,
			errContains: "lowercase letter",
		},
		{
			name:        "missing number",
			password:    "Password!!",
			expectError: true,
			errContains: "number",
		},
		{
			name:        "missing special character",
			password:    "Password123",
			expectError: true,
			errContains: "special character",
		},
		{
			name:        "missing multiple requirements",
			password:    "password",
			expectError: true,
			errContains: "uppercase letter",
		},
		{
			name:        "too long password",
			password:    "Password123!" + string(make([]byte, 120)),
			expectError: true,
			errContains: "128 characters",
		},
		{
			name:        "empty password",
			password:    "",
			expectError: true,
			errContains: "at least 8 characters",
		},
		{
			name:        "only spaces",
			password:    "        ",
			expectError: true,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			err := ValidatePassword(testCase.password)

			if testCase.expectError {
				assert.Error(t, err)
				if testCase.errContains != "" {
					assert.Contains(t, err.Error(), testCase.errContains)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidatePassword_SpecialCharacters(t *testing.T) {
	specialChars := []string{"!", "@", "#", "$", "%", "^", "&", "*", "(", ")", "_", "+", "-", "=", "[", "]", "{", "}", "|", ";", "'", ":", "\"", ",", ".", "/", "<", ">", "?"}

	for _, char := range specialChars {
		t.Run("special char "+char, func(t *testing.T) {
			password := "Password1" + char
			err := ValidatePassword(password)
			assert.NoError(t, err, "password with special character '%s' should be valid", char)
		})
	}
}

func TestNewEmail(t *testing.T) {
	tests := []struct {
		name        string
		email       string
		expectError bool
	}{
		{
			name:        "valid email",
			email:       "test@example.com",
			expectError: false,
		},
		{
			name:        "valid email with subdomain",
			email:       "user@mail.example.com",
			expectError: false,
		},
		{
			name:        "valid email with plus",
			email:       "user+tag@example.com",
			expectError: false,
		},
		{
			name:        "empty email",
			email:       "",
			expectError: true,
		},
		{
			name:        "invalid format - no at sign",
			email:       "testexample.com",
			expectError: true,
		},
		{
			name:        "invalid format - no domain",
			email:       "test@",
			expectError: true,
		},
		{
			name:        "invalid format - no local part",
			email:       "@example.com",
			expectError: true,
		},
		{
			name:        "email with spaces",
			email:       "  test@example.com  ",
			expectError: false,
		},
		{
			name:        "uppercase email normalized to lowercase",
			email:       "TEST@EXAMPLE.COM",
			expectError: false,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			email, err := NewEmail(testCase.email)

			if testCase.expectError {
				assert.Error(t, err)
				assert.True(t, email.IsZero())
			} else {
				assert.NoError(t, err)
				assert.False(t, email.IsZero())
			}
		})
	}
}

func TestEmail_Normalization(t *testing.T) {
	email, err := NewEmail("  TEST@EXAMPLE.COM  ")
	assert.NoError(t, err)
	assert.Equal(t, "test@example.com", email.String())
}

func TestEmail_Domain(t *testing.T) {
	email, _ := NewEmail("user@example.com")
	assert.Equal(t, "example.com", email.Domain())
}

func TestEmail_LocalPart(t *testing.T) {
	email, _ := NewEmail("user@example.com")
	assert.Equal(t, "user", email.LocalPart())
}

func TestEmail_Equals(t *testing.T) {
	email1, _ := NewEmail("test@example.com")
	email2, _ := NewEmail("test@example.com")
	email3, _ := NewEmail("other@example.com")

	assert.True(t, email1.Equals(email2))
	assert.False(t, email1.Equals(email3))
}

func TestNewFullName(t *testing.T) {
	tests := []struct {
		name        string
		fullName    string
		expectError bool
	}{
		{
			name:        "valid full name",
			fullName:    "John Doe",
			expectError: false,
		},
		{
			name:        "minimum length name",
			fullName:    "Jo",
			expectError: false,
		},
		{
			name:        "empty name",
			fullName:    "",
			expectError: true,
		},
		{
			name:        "too short name",
			fullName:    "J",
			expectError: true,
		},
		{
			name:        "whitespace only",
			fullName:    "   ",
			expectError: true,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			name, err := NewFullName(testCase.fullName)

			if testCase.expectError {
				assert.Error(t, err)
				assert.True(t, name.IsZero())
			} else {
				assert.NoError(t, err)
				assert.False(t, name.IsZero())
			}
		})
	}
}

func TestNewPasswordHash(t *testing.T) {
	tests := []struct {
		name        string
		hashValue   string
		expectError bool
	}{
		{
			name:        "valid hash",
			hashValue:   "$2a$10$somehashedvalue",
			expectError: false,
		},
		{
			name:        "empty hash",
			hashValue:   "",
			expectError: true,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			hash, err := NewPasswordHash(testCase.hashValue)

			if testCase.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.hashValue, hash.Value())
			}
		})
	}
}

func TestPagination(t *testing.T) {
	tests := []struct {
		name          string
		page          int
		limit         int
		expectedPage  int
		expectedLimit int
	}{
		{
			name:          "valid pagination",
			page:          1,
			limit:         20,
			expectedPage:  1,
			expectedLimit: 20,
		},
		{
			name:          "negative page defaults to 1",
			page:          -1,
			limit:         20,
			expectedPage:  1,
			expectedLimit: 20,
		},
		{
			name:          "zero page defaults to 1",
			page:          0,
			limit:         20,
			expectedPage:  1,
			expectedLimit: 20,
		},
		{
			name:          "negative limit defaults to 20",
			page:          1,
			limit:         -1,
			expectedPage:  1,
			expectedLimit: 20,
		},
		{
			name:          "limit exceeding max capped to max",
			page:          1,
			limit:         500,
			expectedPage:  1,
			expectedLimit: 100,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			pagination := NewPagination(testCase.page, testCase.limit)

			assert.Equal(t, testCase.expectedPage, pagination.Page())
			assert.Equal(t, testCase.expectedLimit, pagination.Limit())
		})
	}
}

func TestPagination_Offset(t *testing.T) {
	pagination := NewPagination(3, 20)
	assert.Equal(t, 40, pagination.Offset())
}

func TestPagination_TotalPages(t *testing.T) {
	pagination := NewPagination(1, 20)

	assert.Equal(t, 0, pagination.TotalPages(0))
	assert.Equal(t, 1, pagination.TotalPages(20))
	assert.Equal(t, 2, pagination.TotalPages(21))
	assert.Equal(t, 5, pagination.TotalPages(100))
}

func TestPagination_HasNext(t *testing.T) {
	pagination := NewPagination(1, 20)

	assert.True(t, pagination.HasNext(50))
	assert.False(t, pagination.HasNext(20))
	assert.False(t, pagination.HasNext(0))
}

func TestPagination_HasPrev(t *testing.T) {
	assert.False(t, NewPagination(1, 20).HasPrev())
	assert.True(t, NewPagination(2, 20).HasPrev())
}

func TestSortOrder(t *testing.T) {
	assert.True(t, SortOrderAsc.IsValid())
	assert.True(t, SortOrderDesc.IsValid())
	assert.False(t, SortOrder("invalid").IsValid())

	assert.Equal(t, "ASC", SortOrderAsc.SQL())
	assert.Equal(t, "DESC", SortOrderDesc.SQL())

	assert.Equal(t, SortOrderAsc, NewSortOrder("asc"))
	assert.Equal(t, SortOrderAsc, NewSortOrder("ASC"))
	assert.Equal(t, SortOrderDesc, NewSortOrder("desc"))
	assert.Equal(t, SortOrderAsc, NewSortOrder("invalid"))
}

func TestNewPaginatedResult(t *testing.T) {
	items := []string{"a", "b", "c"}
	pagination := NewPagination(1, 20)

	result := NewPaginatedResult(items, 50, pagination)

	assert.Equal(t, items, result.Items)
	assert.Equal(t, int64(50), result.Total)
	assert.Equal(t, 1, result.Page)
	assert.Equal(t, 20, result.Limit)
	assert.Equal(t, 3, result.TotalPages)
	assert.True(t, result.HasNext)
	assert.False(t, result.HasPrev)
}

func TestNewPaginatedResult_NilItems(t *testing.T) {
	pagination := NewPagination(1, 20)
	result := NewPaginatedResult[string](nil, 0, pagination)

	assert.NotNil(t, result.Items)
	assert.Empty(t, result.Items)
}
