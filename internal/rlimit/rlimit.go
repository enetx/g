//go:build !windows

package rlimit

import "syscall"

// RlimitStack is used to adjust the maximum number of worker goroutines, taking into account the
// system's file descriptor limit.
func RlimitStack(workers int) int {
	var rLimit syscall.Rlimit

	syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)

	if uint64(workers) > rLimit.Cur {
		workers = int(float64(rLimit.Cur) * 0.7)
	}

	return workers
}
