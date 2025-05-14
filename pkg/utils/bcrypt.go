package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	hp, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hp), nil
}

func CheckPassword(userPassword string, storedPassword string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(userPassword)); err != nil {
		return err
	}
	return nil
}
