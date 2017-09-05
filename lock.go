// +build !darwin,!dragonfly,!freebsd,!linux,!netbsd,!openbsd,!windows

package file

import (
	"os"
)

func LockSH(fp *os.File) error {
	return nil
}

func LockEX(fp *os.File) error {
	return nil
}

func TryLockSH(fp *os.File) error {
	return nil
}

func TryLockEX(fp *os.File) error {
	return nil
}

func Unlock(fp *os.File) error {
	return nil
}
