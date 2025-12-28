package security

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPasswordHasher_Hash(t *testing.T) {
	hasher := NewDefaultPasswordHasher()

	password := "TestPassword123!"
	hash, err := hasher.Hash(password)

	require.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.NotEqual(t, password, hash)
}

func TestPasswordHasher_HashProducesDifferentHashes(t *testing.T) {
	hasher := NewDefaultPasswordHasher()

	password := "TestPassword123!"
	hash1, _ := hasher.Hash(password)
	hash2, _ := hasher.Hash(password)

	assert.NotEqual(t, hash1, hash2)
}

func TestPasswordHasher_Compare(t *testing.T) {
	hasher := NewDefaultPasswordHasher()

	password := "TestPassword123!"
	hash, _ := hasher.Hash(password)

	err := hasher.Compare(hash, password)
	assert.NoError(t, err)
}

func TestPasswordHasher_CompareWrongPassword(t *testing.T) {
	hasher := NewDefaultPasswordHasher()

	password := "TestPassword123!"
	wrongPassword := "WrongPassword123!"
	hash, _ := hasher.Hash(password)

	err := hasher.Compare(hash, wrongPassword)
	assert.Error(t, err)
}

func TestPasswordHasher_Verify(t *testing.T) {
	hasher := NewDefaultPasswordHasher()

	password := "TestPassword123!"
	hash, _ := hasher.Hash(password)

	tests := []struct {
		name           string
		hashedPassword string
		plainPassword  string
		expectedValid  bool
		expectError    bool
	}{
		{
			name:           "correct password",
			hashedPassword: hash,
			plainPassword:  password,
			expectedValid:  true,
			expectError:    false,
		},
		{
			name:           "wrong password",
			hashedPassword: hash,
			plainPassword:  "WrongPassword123!",
			expectedValid:  false,
			expectError:    false,
		},
		{
			name:           "empty password",
			hashedPassword: hash,
			plainPassword:  "",
			expectedValid:  false,
			expectError:    false,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			valid, err := hasher.Verify(testCase.hashedPassword, testCase.plainPassword)

			if testCase.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.expectedValid, valid)
			}
		})
	}
}

func TestPasswordHasher_CustomCost(t *testing.T) {
	tests := []struct {
		name         string
		cost         int
		expectedCost int
	}{
		{
			name:         "valid cost",
			cost:         12,
			expectedCost: 12,
		},
		{
			name:         "cost below minimum defaults to default",
			cost:         bcrypt.MinCost - 1,
			expectedCost: bcrypt.DefaultCost,
		},
		{
			name:         "cost above maximum defaults to default",
			cost:         bcrypt.MaxCost + 1,
			expectedCost: bcrypt.DefaultCost,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			hasher := NewPasswordHasher(testCase.cost)
			password := "TestPassword123!"

			hash, err := hasher.Hash(password)
			require.NoError(t, err)

			cost, err := bcrypt.Cost([]byte(hash))
			require.NoError(t, err)

			assert.Equal(t, testCase.expectedCost, cost)
		})
	}
}

func TestDefaultPasswordHasher_UsesBcryptDefaultCost(t *testing.T) {
	hasher := NewDefaultPasswordHasher()
	password := "TestPassword123!"

	hash, err := hasher.Hash(password)
	require.NoError(t, err)

	cost, err := bcrypt.Cost([]byte(hash))
	require.NoError(t, err)

	assert.Equal(t, bcrypt.DefaultCost, cost)
}

func TestPasswordHasher_InvalidHash(t *testing.T) {
	hasher := NewDefaultPasswordHasher()

	valid, err := hasher.Verify("invalid-hash", "password")

	assert.Error(t, err)
	assert.False(t, valid)
}

func TestPasswordHasher_TimingAttackResistance(t *testing.T) {
	hasher := NewDefaultPasswordHasher()

	password := "TestPassword123!"
	hash, _ := hasher.Hash(password)

	hasher.Compare(hash, password)
	hasher.Compare(hash, "completely-wrong-password-that-is-very-different")
}
