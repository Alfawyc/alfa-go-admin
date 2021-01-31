package tool

import (
	"golang.org/x/crypto/bcrypt"
)

// 比对密码
func CompareHashAndPassword(original string, input string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(original), []byte(input))
	if err != nil {
		return false, err
	}
	return true, nil
}
