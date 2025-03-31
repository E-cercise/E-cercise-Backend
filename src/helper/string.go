package helper

import (
	"crypto/rand"
	"math/big"
	"strings"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func RandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			panic(err) // Handle error properly in production code
		}
		b[i] = charset[n.Int64()]
	}
	return string(b)
}

func AbbreviateEquipmentName(equipmentName, optionName string) string {
	words := strings.Fields(equipmentName)

	if len(words) <= 5 {
		return equipmentName + " " + optionName
	}

	abbreviated := strings.Join(words[:5], " ") + "..."
	return abbreviated + " " + optionName
}
