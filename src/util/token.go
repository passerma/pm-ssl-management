package util

import (
	goutils "github.com/typa01/go-utils"
)

var token = goutils.GUID()

func GenerateToken() string {
	token = goutils.GUID()
	return token
}

func ValidateToken(v string) bool {
	if v == "" {
		return false
	}
	return v == token
}

func GetToken() string {
	return token
}
