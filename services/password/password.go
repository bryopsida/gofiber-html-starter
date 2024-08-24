package password

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/bryopsida/gofiber-pug-starter/interfaces"
	"golang.org/x/crypto/argon2"
)

type passwordService struct {
	SaltLength uint32
	Time       uint32
	Memory     uint32
	Threads    uint8
	KeyLength  uint32
}

// NewPasswordService creates a new password service
func NewPasswordService() interfaces.IPasswordService {
	return &passwordService{
		SaltLength: 16,
		Time:       1,
		Memory:     64 * 1024,
		Threads:    4,
		KeyLength:  32,
	}
}

func (ps *passwordService) Hash(plaintext string) (string, error) {
	salt := make([]byte, ps.SaltLength)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(plaintext), salt, ps.Time, ps.Memory, ps.Threads, ps.KeyLength)
	encodedSalt := base64.RawStdEncoding.EncodeToString(salt)
	encodedHash := base64.RawStdEncoding.EncodeToString(hash)

	return fmt.Sprintf("%s:%s", encodedSalt, encodedHash), nil
}

func (ps *passwordService) Verify(plaintext, encodedHash string) (bool, error) {
	parts := strings.Split(encodedHash, ":")
	if len(parts) != 2 {
		return false, errors.New("invalid hash format")
	}

	encodedSalt := parts[0]
	encodedStoredHash := parts[1]

	salt, err := base64.RawStdEncoding.DecodeString(encodedSalt)
	if err != nil {
		return false, err
	}

	storedHash, err := base64.RawStdEncoding.DecodeString(encodedStoredHash)
	if err != nil {
		return false, err
	}

	hash := argon2.IDKey([]byte(plaintext), salt, ps.Time, ps.Memory, ps.Threads, ps.KeyLength)

	return subtle.ConstantTimeCompare(hash, storedHash) == 1, nil
}
