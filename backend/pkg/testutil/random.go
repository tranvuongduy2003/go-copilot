package testutil

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	alphabet     = "abcdefghijklmnopqrstuvwxyz"
	alphanumeric = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func RandomString(length int) string {
	var sb strings.Builder
	for i := 0; i < length; i++ {
		index, _ := rand.Int(rand.Reader, big.NewInt(int64(len(alphabet))))
		sb.WriteByte(alphabet[index.Int64()])
	}
	return sb.String()
}

func RandomAlphanumeric(length int) string {
	var sb strings.Builder
	for i := 0; i < length; i++ {
		index, _ := rand.Int(rand.Reader, big.NewInt(int64(len(alphanumeric))))
		sb.WriteByte(alphanumeric[index.Int64()])
	}
	return sb.String()
}

func RandomEmail() string {
	return fmt.Sprintf("%s@%s.com", RandomString(8), RandomString(5))
}

func RandomFullName() string {
	return fmt.Sprintf("%s %s", strings.Title(RandomString(6)), strings.Title(RandomString(8)))
}

func RandomPassword() string {
	return RandomAlphanumeric(12) + "!1Aa"
}

func RandomUUID() uuid.UUID {
	return uuid.New()
}

func RandomInt(min, max int64) int64 {
	diff := max - min
	if diff <= 0 {
		return min
	}
	n, _ := rand.Int(rand.Reader, big.NewInt(diff))
	return min + n.Int64()
}

func RandomTime() time.Time {
	minTime := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	maxTime := time.Now()
	delta := maxTime.Unix() - minTime.Unix()
	randomSeconds := RandomInt(0, delta)
	return minTime.Add(time.Duration(randomSeconds) * time.Second)
}

func RandomTimeBetween(start, end time.Time) time.Time {
	delta := end.Unix() - start.Unix()
	if delta <= 0 {
		return start
	}
	randomSeconds := RandomInt(0, delta)
	return start.Add(time.Duration(randomSeconds) * time.Second)
}

func RandomBool() bool {
	n, _ := rand.Int(rand.Reader, big.NewInt(2))
	return n.Int64() == 1
}

func RandomStatus() string {
	statuses := []string{"pending", "active", "inactive", "banned"}
	index := RandomInt(0, int64(len(statuses)))
	return statuses[index]
}
