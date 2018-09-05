package utils

import (
	"bytes"
	"math/rand"
)

func BytesToString(data []byte) string {
	var n = bytes.IndexByte(data, 0)
	return string(data[:n])
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func GenerateHash(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
