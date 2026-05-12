package security

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"
)

func GeneratePAN() (string, error) {
	// тестовый BIN
	base := "427638"

	// добиваем до 15 цифр
	for len(base) < 15 {
		n, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		base += n.String()
	}

	checkDigit := luhnCheckDigit(base)

	return base + fmt.Sprintf("%d", checkDigit), nil
}

func GenerateExpiry() string {
	t := time.Now().AddDate(3, 0, 0)
	return t.Format("01/06")
}

func GenerateCVV() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1000))
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%03d", n.Int64()), nil
}

func luhnCheckDigit(number string) int {
	sum := 0
	double := true

	for i := len(number) - 1; i >= 0; i-- {
		d := int(number[i] - '0')

		if double {
			d *= 2
			if d > 9 {
				d -= 9
			}
		}

		sum += d
		double = !double
	}

	return (10 - sum%10) % 10
}
