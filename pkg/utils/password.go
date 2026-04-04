package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(pwd string) (string, error) {
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return "", nil
	}
	return string(hashedPwd), nil
}

func CheckHashPassWord(hashedPwd, pwd string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(pwd))
}
