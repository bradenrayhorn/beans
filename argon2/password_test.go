package argon2_test

import (
	"testing"

	"github.com/bradenrayhorn/beans/argon2"
	"github.com/stretchr/testify/require"
)

func TestHashCompare(t *testing.T) {
	hash, err := argon2.GenerateHash("password")
	require.Nil(t, err)

	t.Run("equal when passwords equal", func(t *testing.T) {
		equal, err := argon2.CompareHashAndPassword(hash, "password")
		require.Nil(t, err)
		require.True(t, equal)
	})

	t.Run("not equal when passwords not equal", func(t *testing.T) {
		equal, err := argon2.CompareHashAndPassword(hash, "password2")
		require.Nil(t, err)
		require.False(t, equal)
	})

	t.Run("not equal with blank password", func(t *testing.T) {
		equal, err := argon2.CompareHashAndPassword(hash, "")
		require.Nil(t, err)
		require.False(t, equal)
	})

	t.Run("cannot use invalid hash", func(t *testing.T) {
		_, err := argon2.CompareHashAndPassword("$b$v=19$m=1,t=1,p=1$YQ$YQ", "password")
		require.Errorf(t, err, "invalid hash")
	})

	t.Run("cannot use invalid hash version", func(t *testing.T) {
		_, err := argon2.CompareHashAndPassword("$argon2id$v=20$m=1,t=1,p=1$YQ$YQ", "password")
		require.Errorf(t, err, "invalid hash")
	})

	t.Run("cannot use more than max memory", func(t *testing.T) {
		_, err := argon2.CompareHashAndPassword("$argon2id$v=19$m=65537,t=1,p=1$YQ$YQ", "password")
		require.Errorf(t, err, "invalid hash")
	})

	t.Run("cannot use more than max iterations", func(t *testing.T) {
		_, err := argon2.CompareHashAndPassword("$argon2id$v=19$m=1,t=5,p=1$YQ$YQ", "password")
		require.Errorf(t, err, "invalid hash")
	})

	t.Run("cannot use more than max threads", func(t *testing.T) {
		_, err := argon2.CompareHashAndPassword("$argon2id$v=19$m=1,t=1,p=3$YQ$YQ", "password")
		require.Errorf(t, err, "invalid hash")
	})
}
