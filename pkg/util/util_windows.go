package util

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/minio/highwayhash"
	"io"
	"os"
	"strconv"
	"sync"
)

// Hasher returns a hash function, used in snapshotting to determine if a file has changed
func Hasher() func(string) (string, error) {
	pool := sync.Pool{
		New: func() interface{} {
			b := make([]byte, highwayhash.Size*10*1024)
			return &b
		},
	}
	key := make([]byte, highwayhash.Size)
	hasher := func(p string) (string, error) {
		h, _ := highwayhash.New(key)
		fi, err := os.Lstat(p)
		if err != nil {
			return "", err
		}
		h.Write([]byte(fi.Mode().String()))
		h.Write([]byte(fi.ModTime().String()))

		if fi.Mode().IsRegular() {
			f, err := os.Open(p)
			if err != nil {
				return "", err
			}
			defer f.Close()
			buf := pool.Get().(*[]byte)
			defer pool.Put(buf)
			if _, err := io.CopyBuffer(h, f, *buf); err != nil {
				return "", err
			}
		}

		return hex.EncodeToString(h.Sum(nil)), nil
	}
	return hasher
}

// CacheHasher takes into account everything the regular hasher does except for mtime
func CacheHasher() func(string) (string, error) {
	hasher := func(p string) (string, error) {
		h := md5.New()
		fi, err := os.Lstat(p)
		if err != nil {
			return "", err
		}
		h.Write([]byte(fi.Mode().String()))

		if fi.Mode().IsRegular() {
			f, err := os.Open(p)
			if err != nil {
				return "", err
			}
			defer f.Close()
			if _, err := io.Copy(h, f); err != nil {
				return "", err
			}
		}

		return hex.EncodeToString(h.Sum(nil)), nil
	}
	return hasher
}


// RedoHasher returns a hash function, which looks at mtime, size, filemode, owner uid and gid
// Note that the mtime can lag, so it's possible that a file will have changed but the mtime may look the same.
func RedoHasher() func(string) (string, error) {
	hasher := func(p string) (string, error) {
		h := md5.New()
		fi, err := os.Lstat(p)
		if err != nil {
			return "", err
		}
		h.Write([]byte(fi.Mode().String()))
		h.Write([]byte(fi.ModTime().String()))
		h.Write([]byte(strconv.FormatInt(fi.Size(), 16)))

		return hex.EncodeToString(h.Sum(nil)), nil
	}
	return hasher
}

