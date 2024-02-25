package argon2

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"testing"

	"golang.org/x/crypto/argon2"
)

var argonConfig = struct {
	iterations uint32
	memory     uint32
	threads    uint8
	keyLength  uint32
	saltLength int
}{
	iterations: 4,
	memory:     64 * 1024,
	threads:    2,
	keyLength:  64,
	saltLength: 32,
}

var errInvalidHash = errors.New("invalid argon2id hash")

func GenerateHash(password string) (string, error) {
	salt := make([]byte, argonConfig.saltLength)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	// Use lightweight and unsecure config while running a test
	if testing.Testing() {
		argonConfig.iterations = 1
		argonConfig.memory = 512
		argonConfig.keyLength = 8
		argonConfig.saltLength = 8
	}

	key := argon2.IDKey(
		[]byte(password),
		salt,
		argonConfig.iterations,
		argonConfig.memory,
		argonConfig.threads,
		argonConfig.keyLength,
	)

	hash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		argonConfig.memory,
		argonConfig.iterations,
		argonConfig.threads,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(key),
	)

	return hash, nil
}

func CompareHashAndPassword(hash string, password string) (bool, error) {
	hashed, err := decodeHash(hash)
	if err != nil {
		return false, err
	}

	if err = hashed.validate(); err != nil {
		return false, err
	}

	key := argon2.IDKey(
		[]byte(password),
		hashed.salt,
		hashed.iterations,
		hashed.memory,
		hashed.threads,
		hashed.keyLength,
	)

	return bytes.Equal(key, hashed.key), nil
}

func decodeHash(hash string) (*hashed, error) {
	parts := strings.Split(hash, "$")
	if len(parts) != 6 || parts[0] != "" || parts[1] != "argon2id" {
		return nil, errInvalidHash
	}

	hashed := &hashed{}
	if err := hashed.decodeVersion(parts[2]); err != nil {
		return nil, errInvalidHash
	}

	if err := hashed.decodeConfig(parts[3]); err != nil {
		return nil, errInvalidHash
	}

	if err := hashed.decodeSalt(parts[4]); err != nil {
		return nil, errInvalidHash
	}

	if err := hashed.decodeKey(parts[5]); err != nil {
		return nil, errInvalidHash
	}

	return hashed, nil
}
