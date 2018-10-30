package util

import (
	"github.com/pkg/errors"
)

type PowError struct {
	ErrMsg error
	Code   int
}

var (
	ErrServerPostFailed = PowError{
		ErrMsg: errors.New("request to server failed"),
		Code:   -10010,
	}

	ErrLengthNotEnough = PowError{
		ErrMsg: errors.New("the length is not enough"),
		Code:   -20010,
	}

	ErrLengthOfHashIncorrect = PowError{
		ErrMsg: errors.New("the length doesn't match hash_size"),
		Code:   -20011,
	}
)
