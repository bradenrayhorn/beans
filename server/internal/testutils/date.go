package testutils

import (
	"math/rand"
	"time"
)

func RandomTime() time.Time {
	offset := int64((time.Hour * 24 * 365 * 5).Seconds())

	return time.Unix(rand.Int63n((time.Now().Unix()-offset)+offset), 0)
}
