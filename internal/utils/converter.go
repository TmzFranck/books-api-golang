package utils

import "strconv"

func StringToUint(s string) (uint, error) {
	num, err := strconv.ParseUint(s, 10, 64)
	return uint(num), err
}
