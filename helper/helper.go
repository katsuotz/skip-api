package helper

import "strings"

func IsInArray(arr []int, val int) bool {
	for _, item := range arr {
		if item == val {
			return true
		}
	}
	return false
}

func IsGuru(role string) bool {
	roles := []string{"guru", "staff-ict", "guru-bk", "tata-usaha"}
	return strings.Contains(strings.Join(roles, ","), role)
}
