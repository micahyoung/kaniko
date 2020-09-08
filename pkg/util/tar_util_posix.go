package util

import (
	"os"
	"syscall"
)

func getSyscallStatT(i os.FileInfo) *syscall.Stat_t {
	if sys := i.Sys(); sys != nil {
		if stat, ok := sys.(*syscall.Stat_t); ok {
			return stat
		}
	}
	return nil
}

