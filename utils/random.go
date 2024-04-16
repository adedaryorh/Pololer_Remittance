package utils

import (
	"math/rand"
	"time"
)

var alphabets string = "abcdefghijklmnopqrstuvwxyz"

func RandomString(r int) string {
	//r specify length of strings
	bits := []rune{}
	k := len(alphabets)

	for i := 0; i < r; i++ {
		index := rand.Intn(k)
		bits = append(bits, rune(alphabets[index]))
	}
	return string(bits)
}

func RandomEmail() string {
	return RandomString(8) + "@gmail.com"
}

func RandomInt(min, max int) int {
	rand.Seed(time.Now().UnixNano()) // Seed the random number generator
	return rand.Intn(max-min+1) + min
}
