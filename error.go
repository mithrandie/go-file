package file

import "fmt"

type IOError struct {
	message string
}

func NewIOError(message string) error {
	return &IOError{
		message: message,
	}
}

func (e IOError) Error() string {
	return e.message
}

type LockError struct {
	message string
}

func NewLockError(message string) error {
	return &LockError{
		message: message,
	}
}

func (e LockError) Error() string {
	return e.message
}

type TimeoutError struct {
	message string
}

func NewTimeoutError(filepath string) error {
	return &TimeoutError{
		message: fmt.Sprintf("lock file %s: timeout", filepath),
	}
}

func (e TimeoutError) Error() string {
	return e.message
}
