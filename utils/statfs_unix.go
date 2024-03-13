//go:build !windows
// +build !windows

package utils

import (
	"log"
	"syscall"
)

func DiskSpaceSufficient(path string, sszie uint64, PreCount1 int) bool {
	var stat syscall.Statfs_t
	if err := syscall.Statfs(path, &stat); err != nil {
		log.Println("syscall.Statfs %v", err)
		return false
	}
	if uint64(stat.Bavail)*uint64(stat.Bsize) < sszie*7*uint64(PreCount1) {
		log.Println("Insufficient hard disk space！", uint64(stat.Bavail)*uint64(stat.Bsize), sszie*7*uint64(PreCount1))
		return false
	}
	return true
}

func DiskSpaceSufficientCount(path string, sszie uint64, curPreCount1 int) uint64 {
	var stat syscall.Statfs_t
	if err := syscall.Statfs(path, &stat); err != nil {
		log.Println("syscall.Statfs %v", err)
		return 0
	}

	count := (uint64(stat.Bavail)*uint64(stat.Bsize) - sszie*14*uint64(curPreCount1)) / (sszie * 14)
	if count < 0 {
		count = 0
	}
	log.Println("Free hard disk space：", uint64(stat.Bavail)*uint64(stat.Bsize), "Reserve space：", sszie*14*uint64(curPreCount1), "Support sector：", count)
	return count
}
