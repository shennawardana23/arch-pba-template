package random

import (
	"crypto/rand"
	"math/big"
)

const (
	Numeric       = "0123456789"
	Alpha         = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	AlphaLower    = "abcdefghijklmnopqrstuvwxyz"
	AlphaUpper    = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Alphanum      = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	AlphanumLower = "0123456789abcdefghijklmnopqrstuvwxyz"
	AlphanumUpper = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func GenerateRandomString(letters string, n int) (string, error) {
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		ret[i] = letters[num.Int64()]
	}

	return string(ret), nil
}
