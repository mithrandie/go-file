package file

import (
	"os"
	"time"
)

// Waiting time in seconds to wait for files to be released.
var WaitTimeout float64 = 30.0

// Interval to retry lock a file.
var RetryInterval time.Duration = 50 * time.Millisecond

// Types of Locks
type LockType int

const (
	// Shared Lock
	SHARED_LOCK LockType = iota
	// Exclusive Lock
	EXCLUSIVE_LOCK
)

// Open the file with locking. OpenToRead is the same as Open(path, os.O_RDONLY, 0400, SHARED_LOCK)
func OpenToRead(path string) (*os.File, error) {
	return Open(path, os.O_RDONLY, 0400, SHARED_LOCK)
}

// Open the file with locking. OpenToRead is the same as Open(path, os.O_WRONLY|O_TRUNC, 0600, EXCLUSIVE_LOCK)
func OpenToUpdate(path string) (*os.File, error) {
	return Open(path, os.O_WRONLY|os.O_TRUNC, 0600, EXCLUSIVE_LOCK)
}

// Open the file with locking. Create is the same as Open(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600, EXCLUSIVE_LOCK)
func Create(path string) (*os.File, error) {
	return Open(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600, EXCLUSIVE_LOCK)
}

// Open the file with locking. TryOpenToRead is the same as TryOpen(path, os.O_RDONLY, 0400, SHARED_LOCK)
func TryOpenToRead(path string) (*os.File, error) {
	return TryOpen(path, os.O_RDONLY, 0400, SHARED_LOCK)
}

// Open the file with locking. TryOpenToUpdate is the same as TryOpen(path, os.O_WRONLY|O_TRUNC, 0600, EXCLUSIVE_LOCK)
func TryOpenToUpdate(path string) (*os.File, error) {
	return TryOpen(path, os.O_WRONLY|os.O_TRUNC, 0600, EXCLUSIVE_LOCK)
}

// Open the file with file locking. If the file is already locked, waits until the file is released.
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

// Open the file with file locking. If the file is already locked, waits for up to WaitTimeout seconds.
func TryOpen(path string, flag int, perm os.FileMode, lockType LockType) (*os.File, error) {
	start := time.Now()

	fp, err := os.OpenFile(path, flag, perm)
	if err != nil {
		return nil, NewIOError(err.Error())
	}

	err = TryLock(fp, lockType)
	if err == nil {
		return fp, nil
	}

	if 0 < WaitTimeout {
		for {
			if time.Since(start).Seconds() > WaitTimeout {
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

// Places a lock on the file. If the file is already locked, waits until the file is released.
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

// Places a lock on the file. If the file is already locked, returns an error immediately.
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

// Unlocks and closes the file
func Close(fp *os.File) error {
	defer func() {
		fp.Close()
	}()

	if err := Unlock(fp); err != nil {
		return NewLockError(err.Error())
	}
	return nil
}

// Checks whether the file exists.
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
