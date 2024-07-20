package auth

import "golang.org/x/crypto/bcrypt"

const DefaultCost = 12

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), DefaultCost)

	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func ComparePasswordToHash(originalHash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(originalHash), []byte(password))
	return err != nil
}
