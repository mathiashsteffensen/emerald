package object

import (
	"testing"
)

func TestContext_ValidateMethodVisibility(t *testing.T) {
	self := &Instance{BaseEmeraldValue: &BaseEmeraldValue{}}
	other := &Instance{BaseEmeraldValue: &BaseEmeraldValue{}}
	ctx := &Context{Self: NewHeapObject(self)}

	if !ctx.ValidateMethodVisibility(NewHeapObject(self), PUBLIC, true) {
		t.Error("public method should be visible to self")
	}

	if !ctx.ValidateMethodVisibility(NewHeapObject(other), PUBLIC, true) {
		t.Error("public method should be visible to other")
	}

	if !ctx.ValidateMethodVisibility(NewHeapObject(self), PRIVATE, true) {
		t.Error("private method should be visible to self")
	}

	if ctx.ValidateMethodVisibility(NewHeapObject(other), PRIVATE, true) {
		t.Error("private method should not be visible to other")
	}

	if ctx.ValidateMethodVisibility(NewHeapObject(self), PROTECTED, true) {
		t.Error("protected method visibility not yet implemented, should be false")
	}
}

func TestContext_SetDefaultMethodVisibility(t *testing.T) {
	ctx := &Context{}
	ctx.SetDefaultMethodVisibility(PRIVATE)
	if ctx.DefaultMethodVisibility != PRIVATE {
		t.Error("failed to set default method visibility")
	}
}
