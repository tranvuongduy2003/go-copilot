package security

import (
	"golang.org/x/crypto/bcrypt"
)

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hashedPassword, password string) error
}

type bcryptPasswordHasher struct {
	cost int
}

func NewPasswordHasher(cost int) PasswordHasher {
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		cost = bcrypt.DefaultCost
	}
	return &bcryptPasswordHasher{cost: cost}
}

func NewDefaultPasswordHasher() PasswordHasher {
	return NewPasswordHasher(bcrypt.DefaultCost)
}

func (hasher *bcryptPasswordHasher) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), hasher.cost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (hasher *bcryptPasswordHasher) Compare(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
