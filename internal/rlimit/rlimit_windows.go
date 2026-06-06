package rlimit

// RlimitStack adjusts the maximum number of worker goroutines based on the system's
// open-file-descriptor limit. Windows has no RLIMIT_NOFILE, so the count is unchanged.
func RlimitStack(workers int) int { return workers }
