//go:build linux

package emerald

import "golang.org/x/sys/unix"

func applySandboxHardMemoryLimit(pid int, limitBytes int64) error {
	limit := uint64(limitBytes)
	return unix.Prlimit(pid, unix.RLIMIT_AS, &unix.Rlimit{
		Cur: limit,
		Max: limit,
	}, nil)
}

func applySandboxBestEffortMemoryLimit(limitBytes int64) error {
	return nil
}
