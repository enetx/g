//go:build !windows

package rlimit

import "syscall"

// RlimitStack adjusts the maximum number of worker goroutines, taking into account the
// system's open-file-descriptor limit (RLIMIT_NOFILE). If the limit cannot be read,
// the requested worker count is returned unchanged.
func RlimitStack(workers int) int {
	var rLimit syscall.Rlimit

	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		return workers
	}

	if uint64(workers) > rLimit.Cur {
		workers = int(float64(rLimit.Cur) * 0.7)
	}

	return workers
}
