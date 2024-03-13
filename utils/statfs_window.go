//go:build windows
// +build windows

package utils

func DiskSpaceSufficient(path string, sszie uint64, PreCount1 int) bool {
	return true
}

func DiskSpaceSufficientCount(path string, sszie uint64, curPreCount1 int) uint64 {
	return 1
}
