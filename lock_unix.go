// +build darwin dragonfly freebsd linux netbsd openbsd

package file

import (
	"os"
	"syscall"
)

// Places a shared lock on the file. If the file is already locked, waits until the file is released.
func LockSH(fp *os.File) error {
	return syscall.Flock(int(fp.Fd()), syscall.LOCK_SH)
}

// Places an exclusive lock on the file. If the file is already locked, waits until the file is released.
func LockEX(fp *os.File) error {
	return syscall.Flock(int(fp.Fd()), syscall.LOCK_EX)
}

// Places a shared lock on the file. If the file is already locked, returns an error immediately.
func TryLockSH(fp *os.File) error {
	return syscall.Flock(int(fp.Fd()), syscall.LOCK_SH|syscall.LOCK_NB)
}

// Places an exclusive lock on the file. If the file is already locked, returns an error immediately.
func TryLockEX(fp *os.File) error {
	return syscall.Flock(int(fp.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
}

// Unlock the file.
func Unlock(fp *os.File) error {
	return syscall.Flock(int(fp.Fd()), syscall.LOCK_UN)
}
