package common

import (
	"math/rand"
	"time"
)

var encodeURL = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_")

func init() {
	rand.Seed(time.Now().Unix())
}

// RansString generates a random string of length n
func RandString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = encodeURL[rand.Intn(len(encodeURL))]
	}
	return string(b)
}
