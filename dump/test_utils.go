package dump

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func GenerateRandomNumber() (string, error) {
	max := big.NewInt(9999999999)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}

	numberString := fmt.Sprintf("%010d", n)
	return numberString, nil
}
