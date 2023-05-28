package helper

import (
	"crypto/rand"
	"math/big"
	"strings"
)

func IsInArray(arr []int, val int) bool {
	for _, item := range arr {
		if item == val {
			return true
		}
	}
	return false
}

func IsPegawai(role string) bool {
	roles := []string{"guru", "staff-ict", "guru-bk", "tata-usaha"}
	return strings.Contains(strings.Join(roles, ","), role)
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func RandomString(length int) string {
	var result string
	for i := 0; i < length; i++ {
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			panic(err)
		}
		result += string(charset[randomIndex.Int64()])
	}
	return result
}
