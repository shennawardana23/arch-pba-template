package password

import (
	"context"

	"github.com/mochammadshenna/arch-pba-template/internal/util/helper"
	"golang.org/x/crypto/bcrypt"
)

func CreateHashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	helper.PanicError(err)
	return string(bytes)
}

func CheckHashPassword(ctx context.Context, password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err
}
