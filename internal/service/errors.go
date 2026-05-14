package service

import "errors"

var ErrAccessDenied = errors.New(
	"access denied",
)
