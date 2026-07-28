//go:build linux

package emerald

import (
	"os/exec"
	"testing"

	"golang.org/x/sys/unix"
)

func TestApplySandboxHardMemoryLimit(t *testing.T) {
	cmd := exec.Command("/bin/sleep", "5")
	if err := cmd.Start(); err != nil {
		t.Fatalf("start child: %v", err)
	}
	defer func() {
		_ = cmd.Process.Kill()
		_ = cmd.Wait()
	}()

	const limit = 256 << 20
	if err := applySandboxHardMemoryLimit(cmd.Process.Pid, limit); err != nil {
		t.Fatalf("applySandboxHardMemoryLimit() error = %v", err)
	}

	var got unix.Rlimit
	if err := unix.Prlimit(cmd.Process.Pid, unix.RLIMIT_AS, nil, &got); err != nil {
		t.Fatalf("read child rlimit: %v", err)
	}
	if got.Cur != limit || got.Max != limit {
		t.Fatalf("RLIMIT_AS = {%d, %d}, want {%d, %d}", got.Cur, got.Max, uint64(limit), uint64(limit))
	}
}
