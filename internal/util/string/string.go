package helper

import (
	"encoding/base64"
	"strconv"

	"github.com/mochammadshenna/arch-pba-template/internal/util/json"
)

func EscapeJsonString(i string) string {
	b, err := json.Marshal(i)
	if err != nil {
		return ""
	}
	if len(b) < 2 {
		return ""
	}
	// Trim the beginning and trailing " character
	return string(b[1 : len(b)-1])
}

func HumanizeAmount(i int64) string {
	s := strconv.FormatInt(i, 10)
	if len(s) < 4 {
		return s
	}
	var result string
	for i := len(s) - 1; i >= 0; i-- {
		result = string(s[i]) + result
		if (len(s)-i)%3 == 0 && i != 0 {
			result = "." + result
		}
	}
	return result
}

func PadLeft(s string, length int, pad string) string {
	if len(s) >= length {
		return s
	}
	return PadLeft(pad+s, length, pad)
}

func EncodePointer(id int64) string {
	return base64.StdEncoding.EncodeToString([]byte(strconv.FormatInt(id, 10)))
}

func DecodePointer(s string) (int64, error) {
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(string(b), 10, 64)
}
