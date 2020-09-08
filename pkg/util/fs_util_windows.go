package util

import (
	"os"
)

// DetermineTargetFileOwnership returns the user provided uid/gid combination.
// If they are set to -1, the uid/gid from the original file is used.
func DetermineTargetFileOwnership(fi os.FileInfo, uid, gid int64) (int64, int64) {
	if uid <= DoNotChangeUID {
		uid = 0
	}
	if gid <= DoNotChangeGID {
		gid = 0
	}
	return uid, gid
}

func isSame(fi1, fi2 os.FileInfo) bool {
	return fi1.Mode() == fi2.Mode() &&
		// file modification time
		fi1.ModTime() == fi2.ModTime() &&
		// file size
		fi1.Size() == fi2.Size()
}
