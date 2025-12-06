package main

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// func VerifyPassword(storedHash string, password string) error {
// 	err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password))
// 	if err != nil {
// 		// password is wrong (or hash is invalid)
// 		// treat as invalid credentials
// 	} else {
// 		// password is correct
// 	}
// }
