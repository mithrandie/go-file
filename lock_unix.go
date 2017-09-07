// +build darwin dragonfly freebsd linux netbsd openbsd

package file

import (
	"os"
	"syscall"
)

func LockSH(fp *os.File) error {
	return syscall.Flock(int(fp.Fd()), syscall.LOCK_SH)
}

func LockEX(fp *os.File) error {
	return syscall.Flock(int(fp.Fd()), syscall.LOCK_EX)
}

func TryLockSH(fp *os.File) error {
	return syscall.Flock(int(fp.Fd()), syscall.LOCK_SH|syscall.LOCK_NB)
}

func TryLockEX(fp *os.File) error {
	return syscall.Flock(int(fp.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
}

func Unlock(fp *os.File) error {
	return syscall.Flock(int(fp.Fd()), syscall.LOCK_UN)
}
