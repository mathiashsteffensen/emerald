//go:build darwin

package emerald

import (
	"golang.org/x/sys/unix"
	"syscall"
)

func applySandboxHardMemoryLimit(pid int, limitBytes int64) error {
	return nil
}

func applySandboxBestEffortMemoryLimit(limitBytes int64) error {
	limit := uint64(limitBytes)
	return syscall.Setrlimit(unix.RLIMIT_RSS, &syscall.Rlimit{
		Cur: limit,
		Max: limit,
	})
}
