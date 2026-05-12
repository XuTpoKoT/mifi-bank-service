package security

import (
	"encoding/base64"
)

func Encrypt(value string) string {
	return base64.StdEncoding.EncodeToString(
		[]byte(value),
	)
}

func Decrypt(value string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
