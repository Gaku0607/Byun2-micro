package util

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func GenerateCode() string {

	code := []byte("1234567890abcdefghijklmnopqrstuvwxyz")
	size := int32(len(code))

	randomcode := make([]byte, 8)

	for index := range randomcode {
		randomcode[index] = code[rand.Int31n(size)]
	}

	return string(randomcode)
}
