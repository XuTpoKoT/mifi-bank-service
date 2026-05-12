package security

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"os"
)

func ComputeHMAC(value string) string {
	key := os.Getenv("CARD_HMAC_SECRET")

	mac := hmac.New(
		sha256.New,
		[]byte(key),
	)

	mac.Write([]byte(value))

	return hex.EncodeToString(
		mac.Sum(nil),
	)
}
