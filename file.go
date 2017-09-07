package file

import (
	"os"
	"time"
)

var Timeout float64 = 10 //Seconds
var RetryInterval = 50 * time.Millisecond

type LockType int8

const (
	SHARED_LOCK LockType = iota
	EXCLUSIVE_LOCK
)

func OpenForRead(path string) (*os.File, error) {
	return Open(path, os.O_RDONLY, 0400, SHARED_LOCK)
}

func OpenForUpdate(path string) (*os.File, error) {
	return Open(path, os.O_WRONLY|os.O_TRUNC, 0600, EXCLUSIVE_LOCK)
}

func OpenForCreate(path string) (*os.File, error) {
	return Open(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600, EXCLUSIVE_LOCK)
}

func OpenNBForRead(path string) (*os.File, error) {
	return OpenNB(path, os.O_RDONLY, 0400, SHARED_LOCK)
}

func OpenNBForUpdate(path string) (*os.File, error) {
	return OpenNB(path, os.O_WRONLY|os.O_TRUNC, 0600, EXCLUSIVE_LOCK)
}

func Open(path string, flag int, perm os.FileMode, lockType LockType) (*os.File, error) {
	fp, err := os.OpenFile(path, flag, perm)
	if err != nil {
		return nil, NewIOError(err.Error())
	}

	err = Lock(fp, lockType)
	if err != nil {
		fp.Close()
		return nil, err
	}
	return fp, nil
}

func OpenNB(path string, flag int, perm os.FileMode, lockType LockType) (*os.File, error) {
	start := time.Now()

	fp, err := os.OpenFile(path, flag, perm)
	if err != nil {
		return nil, NewIOError(err.Error())
	}

	err = TryLock(fp, lockType)
	if err == nil {
		return fp, nil
	}

	if 0 < Timeout {
		for {
			if time.Since(start).Seconds() > Timeout {
				err = NewTimeoutError(path)
				break
			}
			time.Sleep(RetryInterval)

			if err = TryLock(fp, lockType); err == nil {
				break
			}
		}
	}
	if err != nil {
		fp.Close()
		return nil, err
	}

	return fp, nil
}

func Lock(fp *os.File, lockType LockType) error {
	var err error
	switch lockType {
	case EXCLUSIVE_LOCK:
		err = LockEX(fp)
	default:
		err = LockSH(fp)
	}

	if err != nil {
		return NewLockError(err.Error())
	}
	return nil
}

func TryLock(fp *os.File, lockType LockType) error {
	var err error
	switch lockType {
	case EXCLUSIVE_LOCK:
		err = TryLockEX(fp)
	default:
		err = TryLockSH(fp)
	}

	if err != nil {
		return NewLockError(err.Error())
	}
	return nil
}

func Close(fp *os.File) error {
	defer func() {
		fp.Close()
	}()

	if err := Unlock(fp); err != nil {
		return NewLockError(err.Error())
	}
	return nil
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
