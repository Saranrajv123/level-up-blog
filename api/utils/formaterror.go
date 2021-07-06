package utils

import (
	"errors"
	"strings"
)

func FormatError(err string) error {
	if strings.Contains(err, "full_name") {
		return errors.New("full name already Taken")
	}

	if strings.Contains(err, "email") {
		return errors.New("email already Taken")
	}
	if strings.Contains(err, "title") {
		return errors.New("title already Taken")
	}
	if strings.Contains(err, "hashPassword") {
		return errors.New("incorrect password")
	}

	return errors.New("incorrect Details")
}
