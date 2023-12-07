package argon2

import (
	"encoding/base64"
	"errors"
	"math"
	"strconv"
	"strings"

	"golang.org/x/crypto/argon2"
)

type hashed struct {
	version    int
	memory     uint32
	iterations uint32
	threads    uint8
	salt       []byte
	key        []byte
	keyLength  uint32
}

func (h *hashed) decodeVersion(r string) error {
	if len(r) < 3 || r[0:2] != "v=" {
		return errInvalidHash
	}

	version, err := strconv.ParseInt(r[2:], 10, 32)
	if err != nil {
		return err
	}

	h.version = int(version)
	return nil
}

func (h *hashed) decodeConfig(r string) error {
	parts := strings.Split(r, ",")
	if len(parts) != 3 {
		return errInvalidHash
	}

	// memory
	if len(parts[0]) < 3 || parts[0][0:2] != "m=" {
		return errInvalidHash
	}

	memory, err := strconv.ParseUint(parts[0][2:], 10, 32)
	if err != nil {
		return err
	}
	h.memory = uint32(memory)

	// iterations
	if len(parts[1]) < 3 || parts[1][0:2] != "t=" {
		return errInvalidHash
	}

	iterations, err := strconv.ParseUint(parts[1][2:], 10, 32)
	if err != nil {
		return err
	}
	h.iterations = uint32(iterations)

	// threads
	if len(parts[2]) < 3 || parts[2][0:2] != "p=" {
		return errInvalidHash
	}

	threads, err := strconv.ParseUint(parts[2][2:], 10, 8)
	if err != nil {
		return err
	}
	h.threads = uint8(threads)

	return nil
}

func (h *hashed) decodeSalt(r string) error {
	salt, err := base64.RawStdEncoding.DecodeString(r)
	if err != nil {
		return err
	}

	h.salt = salt
	return nil
}

func (h *hashed) decodeKey(r string) error {
	key, err := base64.RawStdEncoding.DecodeString(r)
	if err != nil {
		return err
	}
	if len(key) > math.MaxUint32 {
		return errors.New("key too long")
	}

	h.key = key
	h.keyLength = uint32(len(key))
	return nil
}

func (h *hashed) validate() error {
	if h.iterations > argonConfig.iterations ||
		h.memory > argonConfig.memory ||
		h.threads > argonConfig.threads ||
		h.version != argon2.Version {
		return errInvalidHash
	}

	return nil
}
