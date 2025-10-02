package utils

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

// Argon2 parameters - these are recommended values
var (
	// Time parameter (iterations)
	argon2Time uint32 = 1
	// Memory parameter (in KB)
	argon2Memory uint32 = 64 * 1024
	// Parallelism parameter
	argon2Threads uint8 = 4
	// Salt length
	saltLength uint32 = 16
	// Hash length
	hashLength uint32 = 32
)

// HashPassword creates an Argon2id hash of the password
func HashPassword(password string) (string, error) {
	// Generate a random salt
	salt := make([]byte, saltLength)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	// Generate the hash
	hash := argon2.IDKey([]byte(password), salt, argon2Time, argon2Memory, argon2Threads, hashLength)

	// Encode salt and hash to base64
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// Return format: $argon2id$v=19$m=65536,t=1,p=4$salt$hash
	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, argon2Memory, argon2Time, argon2Threads, b64Salt, b64Hash)

	return encodedHash, nil
}

// VerifyPassword verifies a password against an Argon2id hash
func VerifyPassword(password, encodedHash string) (bool, error) {
	// Parse the encoded hash
	salt, hash, params, err := decodeHash(encodedHash)
	if err != nil {
		return false, err
	}

	// Generate hash from the provided password using the same parameters
	testHash := argon2.IDKey([]byte(password), salt, params.time, params.memory, params.threads, uint32(len(hash)))

	// Compare hashes using constant time comparison
	return subtle.ConstantTimeCompare(hash, testHash) == 1, nil
}

// argon2Params holds the parameters for Argon2
type argon2Params struct {
	memory  uint32
	time    uint32
	threads uint8
}

// decodeHash decodes an Argon2id hash string
func decodeHash(encodedHash string) (salt, hash []byte, params *argon2Params, err error) {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return nil, nil, nil, errors.New("invalid hash format")
	}

	if parts[1] != "argon2id" {
		return nil, nil, nil, errors.New("unsupported hash type")
	}

	version := -1
	_, err = fmt.Sscanf(parts[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, errors.New("incompatible argon2 version")
	}

	params = &argon2Params{}
	_, err = fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &params.memory, &params.time, &params.threads)
	if err != nil {
		return nil, nil, nil, err
	}

	salt, err = base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return nil, nil, nil, err
	}

	hash, err = base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return nil, nil, nil, err
	}

	return salt, hash, params, nil
}

// GenerateSecureToken generates a cryptographically secure random token
func GenerateSecureToken(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}
