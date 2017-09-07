package file

import (
	"os"
	"runtime"
	"testing"
)

func TestOpen(t *testing.T) {
	var err error

	notexistpath := GetTestFilePath("notexist.txt")
	_, err = OpenToRead(notexistpath)
	if err == nil {
		t.Fatal("no error, want IOError")
	}
	if _, ok := err.(*IOError); !ok {
		t.Fatal("error is not a IOError")
	}

	_, err = OpenToUpdate(notexistpath)
	if err == nil {
		t.Fatal("no error, want IOError")
	}
	if _, ok := err.(*IOError); !ok {
		t.Fatal("error is not a IOError")
	}

	_, err = OpenToReadWithTimeout(notexistpath)
	if err == nil {
		t.Fatal("no error, want IOError")
	}
	if _, ok := err.(*IOError); !ok {
		t.Fatal("error is not a IOError")
	}

	_, err = OpenToUpdateWithTimeout(notexistpath)
	if err == nil {
		t.Fatal("no error, want IOError")
	}
	if _, ok := err.(*IOError); !ok {
		t.Fatal("error is not a IOError")
	}

	switch runtime.GOOS {
	case "darwin", "dragonfly", "freebsd", "linux", "netbsd", "openbsd", "windows":
		WaitTimeout = 0.1

		shpath := GetTestFilePath("lock_sh.txt")
		expath := GetTestFilePath("lock_ex.txt")

		shfp, _ := os.OpenFile(shpath, os.O_CREATE, 0600)
		shfp.Close()
		exfp, _ := os.OpenFile(expath, os.O_CREATE, 0600)
		exfp.Close()

		shfp1, err := OpenToRead(shpath)
		defer Close(shfp1)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		exfp1, err := OpenToUpdate(expath)
		defer Close(exfp1)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		shfp2, err := OpenToReadWithTimeout(shpath)
		defer Close(shfp2)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		exfp2, err := OpenToUpdateWithTimeout(expath)
		defer Close(exfp2)
		if err == nil {
			t.Fatal("no error, want error for duplicate exclusive lock")
		}
		if _, ok := err.(*TimeoutError); !ok {
			t.Fatal("error is not a TimeoutError")
		}

		err = Lock(nil, SHARED_LOCK)
		if err == nil {
			t.Fatal("no error, want error for invalid file descriptor")
		}
		if _, ok := err.(*LockError); !ok {
			t.Fatal("error is not a LockError")
		}
		err = TryLock(nil, SHARED_LOCK)
		if err == nil {
			t.Fatal("no error, want error for invalid file descriptor")
		}
		if _, ok := err.(*LockError); !ok {
			t.Fatal("error is not a LockError")
		}
	}
}

func TestClose(t *testing.T) {
	path := GetTestFilePath("closetest.txt")
	fp, _ := Create(path)
	err := fp.Close()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	err = Close(fp)
	if err == nil {
		t.Fatal("no error, want error for invalid file descriptor")
	}
	if _, ok := err.(*LockError); !ok {
		t.Fatal("error is not a LockError")
	}
}
