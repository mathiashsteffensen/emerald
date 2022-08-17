package object

import (
	"testing"
	"unsafe"
)

func TestEmeraldValueSize(t *testing.T) {
	var value EmeraldValue

	if unsafe.Sizeof(value) != 16 {
		t.Errorf("blank EmeraldValue takes up %d bytes", unsafe.Sizeof(value))
	}
}
