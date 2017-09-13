// +build solaris

package file

import (
	"os"
	"syscall"
	"time"
)

// Places a shared lock on the file. If the file is already locked, waits until the file is released.
func LockSH(fp *os.File) error {
	for {
		err := TryLockSH(fp)
		if err == nil {
			break
		}
		if err != syscall.EAGAIN {
			return err
		}
		time.Sleep(RetryInterval)
	}

	return nil
}

// Places an exclusive lock on the file. If the file is already locked, waits until the file is released.
func LockEX(fp *os.File) error {
	for {
		err := TryLockEX(fp)
		if err == nil {
			break
		}
		if err != syscall.EAGAIN {
			return err
		}
		time.Sleep(RetryInterval)
	}

	return nil
}

// Places a shared lock on the file. If the file is already locked, returns an error immediately.
func TryLockSH(fp *os.File) error {
	var lock syscall.Flock_t
	lock.Start = 0
	lock.Len = 0
	lock.Whence = 0
	lock.Pid = 0
	lock.Type = syscall.F_RDLCK
	return syscall.FcntlFlock(uintptr(fp.Fd()), syscall.F_SETLK, &lock)
}

// Places an exclusive lock on the file. If the file is already locked, returns an error immediately.
func TryLockEX(fp *os.File) error {
	var lock syscall.Flock_t
	lock.Start = 0
	lock.Len = 0
	lock.Whence = 0
	lock.Pid = 0
	lock.Type = syscall.F_WRLCK
	return syscall.FcntlFlock(uintptr(fp.Fd()), syscall.F_SETLK, &lock)
}

// Unlock the file.
func Unlock(fp *os.File) error {
	var lock syscall.Flock_t
	lock.Start = 0
	lock.Len = 0
	lock.Whence = 0
	lock.Type = syscall.F_UNLCK
	return syscall.FcntlFlock(uintptr(fp.Fd()), syscall.F_SETLK, &lock)
}
