package utils

import (
	"math/rand"
	"strings"
)

func BytesToString(data []byte) string {
	return string(data)
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func GenerateHash(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func RemoveUnnecessarySymbols(data string) string {
	var removeSymbols = []string{"\n", "\t", "\r"}
	for _, symbol := range removeSymbols {
		data = strings.Replace(data, symbol, "", -1)
	}
	return data
}
