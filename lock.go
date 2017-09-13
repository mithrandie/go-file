// +build !darwin,!dragonfly,!freebsd,!linux,!netbsd,!openbsd,!solaris,!windows

package file

import (
	"os"
)

// Do nothing
func LockSH(fp *os.File) error {
	return nil
}

// Do nothing
func LockEX(fp *os.File) error {
	return nil
}

// Do nothing
func TryLockSH(fp *os.File) error {
	return nil
}

// Do nothing
func TryLockEX(fp *os.File) error {
	return nil
}

// Do nothing
func Unlock(fp *os.File) error {
	return nil
}
