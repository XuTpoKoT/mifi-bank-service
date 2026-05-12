package security

import (
	"os"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
)

var (
	publicKeyRing  *crypto.KeyRing
	privateKeyRing *crypto.KeyRing
)

func InitPGP() error {
	pubPath := os.Getenv("PGP_PUBLIC_KEY_PATH")
	privPath := os.Getenv("PGP_PRIVATE_KEY_PATH")
	passphrase := os.Getenv("PGP_PASSPHRASE")

	pubArmored, err := os.ReadFile(pubPath)
	if err != nil {
		return err
	}

	privArmored, err := os.ReadFile(privPath)
	if err != nil {
		return err
	}

	pubKeyObj, err := crypto.NewKeyFromArmored(
		string(pubArmored),
	)
	if err != nil {
		return err
	}

	privKeyObj, err := crypto.NewKeyFromArmored(
		string(privArmored),
	)
	if err != nil {
		return err
	}

	// unlock private key
	unlockedPrivKey, err := privKeyObj.Unlock(
		[]byte(passphrase),
	)
	if err != nil {
		return err
	}

	publicKeyRing, err = crypto.NewKeyRing(
		pubKeyObj,
	)
	if err != nil {
		return err
	}

	privateKeyRing, err = crypto.NewKeyRing(
		unlockedPrivKey,
	)
	if err != nil {
		return err
	}

	return nil
}

func Encrypt(value string) (string, error) {
	msg := crypto.NewPlainMessage(
		[]byte(value),
	)

	enc, err := publicKeyRing.Encrypt(
		msg,
		nil,
	)
	if err != nil {
		return "", err
	}

	armored, err := enc.GetArmored()
	if err != nil {
		return "", err
	}

	return armored, nil
}

func Decrypt(value string) (string, error) {
	msg, err := crypto.NewPGPMessageFromArmored(
		value,
	)
	if err != nil {
		return "", err
	}

	dec, err := privateKeyRing.Decrypt(
		msg,
		nil,
		0,
	)
	if err != nil {
		return "", err
	}

	return string(dec.GetBinary()), nil
}
