package random

import (
	"math/rand"
	"time"
)

func GenetareRandomCode() string {
	rand.Seed(time.Now().UnixNano())
	const charset = "01234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789"
	b := make([]byte, 6)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}

	return string(b)
}
