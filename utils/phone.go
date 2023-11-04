package utils

import "regexp"

func IsValidPhoneNumber(phoneNumber string) bool {
	re := regexp.MustCompile(`^[0-9]{10}$`)
	return re.MatchString(phoneNumber)
}
