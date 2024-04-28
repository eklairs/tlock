package utils

import (
	"strconv"

	"github.com/pquerna/otp"
)

// Converts the given string to an integer, ignoring any errors
func ToInt(str string) int {
	res, _ := strconv.Atoi(str)

	return res
}

func ToHashFunction(value string) otp.Algorithm {
	switch value {
	case "SHA1":
		return otp.AlgorithmSHA1
	case "SHA512":
		return otp.AlgorithmSHA512
	case "MD5":
		return otp.AlgorithmMD5
	default:
		return otp.AlgorithmSHA256
	}
}

func Or(left string, right int) int {
	if left == "" {
		return right
	}

	return ToInt(left)
}
