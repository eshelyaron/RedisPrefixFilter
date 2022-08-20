package utils

import "math/rand"

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func GetRandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = letters[rand.Int63()%int64(len(letters))]
	}
	return string(b)
}

func GetRandomStrings(minLength, maxLenght, numberOfStrings int) []string {
	res := make([]string, numberOfStrings)
	diffLength := maxLenght - minLength
	for i := range res {
		res[i] = GetRandomString(rand.Intn(diffLength) + minLength)
	}
	return res
}
