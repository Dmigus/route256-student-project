package models

import "errors"

var (
	ErrNotFound           = errors.New("not found")
	ErrFailedPrecondition = errors.New("the system state is unsuitable for this operation")
)
