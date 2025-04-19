package utils

import "golang.org/x/crypto/bcrypt"

func HashaPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func CheckPasswordHash(passwordDB, passwordInput string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(passwordDB), []byte(passwordInput))
	if err != nil {
		return false, err
	}

	return true, nil
}
