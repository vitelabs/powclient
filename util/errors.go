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

	ErrClientPowFailed = PowError{
		ErrMsg: errors.New("client pow failed"),
		Code:   -10011,
	}

	ErrLengthNotEnough = PowError{
		ErrMsg: errors.New("the length is not enough"),
		Code:   -20010,
	}

	ErrThresholdParsingFailed = PowError{
		ErrMsg: errors.New("parse threshold failed"),
		Code:   -20011,
	}

	ErrHashParsingFailed = PowError{
		ErrMsg: errors.New("parse datahash failed"),
		Code:   -20012,
	}
)
