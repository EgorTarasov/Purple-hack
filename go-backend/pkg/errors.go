package pkg

import (
	"strings"
)

const (
	DuplicateKeyPrefix = "ERROR: duplicate key"
	PageNotFoundPrefix = "Cannot"
)

func CheckErrDuplicateKey(err error) bool {
	return strings.HasPrefix(err.Error(), DuplicateKeyPrefix)
}

func CheckErrPageNotFound(err error) bool {
	return strings.HasPrefix(err.Error(), PageNotFoundPrefix)
}
