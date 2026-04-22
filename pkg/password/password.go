package password

import "golang.org/x/crypto/bcrypt"

func Hash(value string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(value), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func Compare(hashed, value string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(value))
}
