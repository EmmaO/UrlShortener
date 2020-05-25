package stringextensions

import "math/rand"

//RandomAlphaNumeric returns a string which contains alphanumeric characters(a-z, A-Z, 0-9)
func RandomAlphaNumeric(length int) string {
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

	returnValue := ""
	for i := 0; i < length; i++ {
		returnValue += string(chars[rand.Intn(len(chars))])
	}

	return returnValue
}
