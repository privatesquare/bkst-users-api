package utils

import (
	"fmt"
	"math/rand"
	"time"
	"unicode"
)

// GetRandomPassword generates a random string of upper + lower case alphabets and digits
// which is 23 bits long and returns the string
func GetRandomPassword() string {
	rand.Seed(time.Now().UnixNano())
	digits := "0123456789"
	all := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz" + digits
	length := 23
	buf := make([]byte, length)
	buf[0] = digits[rand.Intn(len(digits))]
	for i := 1; i < length; i++ {
		buf[i] = all[rand.Intn(len(all))]
	}
	rand.Shuffle(len(buf), func(i, j int) {
		buf[i], buf[j] = buf[j], buf[i]
	})
	return string(buf)
}

func VerifyPassword(password string) error {
	var (
		numOfLetters                  = 0
		number, upper, lower, special bool
	)
	for _, c := range password {
		switch {
		case unicode.IsNumber(c):
			number = true
			numOfLetters++
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			special = true
			numOfLetters++
		case unicode.IsUpper(c):
			upper = true
			numOfLetters++
		case unicode.IsLower(c) || c == ' ':
			lower = true
			numOfLetters++
		}
	}
	fmt.Println(numOfLetters, number, upper, lower, special)
	if numOfLetters > 8 && number && upper && lower && special {
		return nil
	} else {
		return InvalidPasswordError
	}
}
