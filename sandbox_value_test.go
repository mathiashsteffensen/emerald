package emerald

import (
	"emerald/core"
	"emerald/internal/sandboxwire"
	"emerald/object"
	"errors"
	"math"
	"testing"
)

func TestSandboxValueRoundTripsFloatBits(t *testing.T) {
	rt := core.NewRuntime()
	rt.InitSandbox()

	want := rt.NewFloat(math.Inf(1))
	encoded, err := encodeSandboxValue(want)
	if err != nil {
		t.Fatalf("encodeSandboxValue() error = %v", err)
	}
	got, err := decodeSandboxValue(encoded, maxSandboxProtocolFrameBytes)
	if err != nil {
		t.Fatalf("decodeSandboxValue() error = %v", err)
	}
	if !got.Is(object.FLOAT_VALUE) || got.Num != want.Num {
		t.Fatalf("float bits changed: got %#x, want %#x", got.Num, want.Num)
	}
}

func TestSandboxValueRejectsCompositeAlias(t *testing.T) {
	rt := core.NewRuntime()
	rt.InitSandbox()

	array := rt.NewArray([]object.EmeraldValue{rt.NewInteger(1)})
	result := rt.NewArray([]object.EmeraldValue{array, array})
	_, err := encodeSandboxValue(result)
	if !errors.Is(err, ErrSandboxUnsupportedResult) {
		t.Fatalf("expected unsupported result error, got %v", err)
	}
}

func TestDecodeSandboxValueRejectsUnknownType(t *testing.T) {
	_, err := decodeSandboxValue(&sandboxwire.Value{Type: "unknown"}, maxSandboxProtocolFrameBytes)
	if !errors.Is(err, ErrSandboxWorkerFailed) {
		t.Fatalf("expected worker failure, got %v", err)
	}
}
