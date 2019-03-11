package file

import (
	"context"
	"os"
	"runtime"
	"testing"
	"time"
)

func TestOpen(t *testing.T) {
	var err error
	ctx, _ := context.WithTimeout(context.Background(), 100*time.Millisecond)
	retryDelay := 50 * time.Millisecond

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

	_, err = OpenToReadContext(ctx, retryDelay, notexistpath)
	if err == nil {
		t.Fatal("no error, want IOError")
	}
	if _, ok := err.(*IOError); !ok {
		t.Fatal("error is not a IOError")
	}

	_, err = OpenToReadContext(ctx, retryDelay, notexistpath)
	if err == nil {
		t.Fatal("no error, want IOError")
	}
	if _, ok := err.(*IOError); !ok {
		t.Fatal("error is not a IOError")
	}

	switch runtime.GOOS {
	case "darwin", "dragonfly", "freebsd", "linux", "netbsd", "openbsd", "windows":
		shpath := GetTestFilePath("lock_sh.txt")
		expath := GetTestFilePath("lock_ex.txt")

		shfp, _ := os.OpenFile(shpath, os.O_CREATE, 0600)
		_ = shfp.Close()
		exfp, _ := os.OpenFile(expath, os.O_CREATE, 0600)
		_ = exfp.Close()

		shfp1, err := OpenToRead(shpath)
		defer func() { _ = Close(shfp1) }()
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		exfp1, err := OpenToUpdate(expath)
		defer func() { _ = Close(exfp1) }()
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		shfp2, err := OpenToReadContext(ctx, retryDelay, shpath)
		defer func() { _ = Close(shfp2) }()
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		exfp2, err := OpenToUpdateContext(ctx, retryDelay, expath)
		defer func() { _ = Close(exfp2) }()
		if err == nil {
			t.Fatal("no error, want error for duplicate exclusive lock")
		}
		if _, ok := err.(*TimeoutError); !ok {
			t.Fatal("error is not a TimeoutError")
		}

		exfp3, err := TryOpenToRead(expath)
		defer func() { _ = Close(exfp3) }()
		if err == nil {
			t.Fatal("no error, want error for duplicate exclusive lock")
		}
		if _, ok := err.(*LockError); !ok {
			t.Fatal("error is not a LockError")
		}

		exfp4, err := TryOpenToUpdate(expath)
		defer func() { _ = Close(exfp4) }()
		if err == nil {
			t.Fatal("no error, want error for duplicate exclusive lock")
		}
		if _, ok := err.(*LockError); !ok {
			t.Fatal("error is not a LockError")
		}

		err = RLock(nil)
		if err == nil {
			t.Fatal("no error, want error for invalid file descriptor")
		}
		if _, ok := err.(*LockError); !ok {
			t.Fatal("error is not a LockError")
		}
	case "solaris":
		// maybe write later
	}
}

func TestClose(t *testing.T) {
	path := GetTestFilePath("closetest.txt")
	fp, _ := Create(path)
	err := fp.Close()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	switch runtime.GOOS {
	case "darwin", "dragonfly", "freebsd", "linux", "netbsd", "openbsd", "solaris", "windows":
		err = Close(fp)
		if err == nil {
			t.Fatal("no error, want error for invalid file descriptor")
		}
		if _, ok := err.(*LockError); !ok {
			t.Fatal("error is not a LockError")
		}
	}
}
