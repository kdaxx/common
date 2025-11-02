package errs

import (
	"errors"
	"fmt"
	"go.uber.org/multierr"
)

func Causef(cause error, format string, message ...any) error {
	msg := fmt.Errorf(format, message...)
	if cause == nil {
		return msg
	}
	return fmt.Errorf("%s: %w", msg, cause)
}

func Cause(cause error, message ...any) error {
	msg := fmt.Sprint(message...)
	if cause == nil {
		return errors.New(msg)
	}
	return fmt.Errorf("%s: %w", msg, cause)
}

func New(text string) error {
	return errors.New(text)
}

func Newf(format string, msg ...any) error {
	return fmt.Errorf(format, msg...)
}

// Combine skips over nil arguments
func Combine(errs ...error) error {
	return multierr.Combine(errs...)
}
