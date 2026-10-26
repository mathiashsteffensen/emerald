package vm

import (
	"emerald/compiler"
	"emerald/core"
	"emerald/object"
	"reflect"
	"testing"
)

func TestBacktrace(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name: "nested methods",
			input: `def inner
  raise "boom"
end

def outer
  inner
end

outer`,
			expected: []string{"test:2:in 'inner'", "test:6:in 'outer'", "test:9:in '<main>'"},
		},
		{
			name: "blocks are labelled by their method",
			input: `def outer
  [1].each do |x|
    [2].each do |y|
      raise "boom"
    end
  end
end

outer`,
			expected: []string{
				"test:4:in 'block (2 levels) in outer'",
				"test:3:in 'block in outer'",
				"test:2:in 'outer'",
				"test:9:in '<main>'",
			},
		},
		{
			name: "errors raised by native methods point at the caller",
			input: `def call_missing
  nil.missing
end

call_missing`,
			expected: []string{"test:2:in 'call_missing'", "test:5:in '<main>'"},
		},
		{
			name: "arity errors point at the call",
			input: `def one(arg); end

def caller
  one
end

caller`,
			expected: []string{"test:4:in 'caller'", "test:7:in '<main>'"},
		},
		{
			name: "re-raising keeps the original backtrace",
			input: `def inner
  raise "boom"
end

def outer
  inner
rescue StandardError => error
  raise error
end

outer`,
			expected: []string{"test:2:in 'inner'", "test:6:in 'outer'", "test:11:in '<main>'"},
		},
		{
			name: "raising a fresh error from a rescue starts a new backtrace",
			input: `def inner
  raise "boom"
end

def outer
  inner
rescue StandardError
  raise "again"
end

outer`,
			expected: []string{"test:8:in 'outer'", "test:11:in '<main>'"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rt := core.NewRuntime()
			rt.Init()

			vm := New("test", compiler.Compile("test", tt.input, rt), rt)
			vm.Run()

			exception := rt.Heap.GetGlobalVariableString("$!")
			if exception.IsNil() {
				t.Fatal("expected an uncaught exception")
			}

			backtrace := exception.Heap.(object.EmeraldError).Backtrace()
			if !reflect.DeepEqual(backtrace, tt.expected) {
				t.Errorf("backtrace\nwant=%q\ngot =%q", tt.expected, backtrace)
			}
		})
	}
}

func TestExceptionBacktraceMethod(t *testing.T) {
	runVmTests(t, []vmTestCase{
		{name: "nil before raise", input: `RuntimeError.new("boom").backtrace`, expected: nil},
		{
			name: "rescued error exposes its frames",
			input: `def inner
  raise "boom"
end

def outer
  inner
rescue StandardError => error
  error.backtrace
end

outer`,
			expected: []any{"test:2:in 'inner'", "test:6:in 'outer'", "test:11:in '<main>'"},
		},
		{name: "message", input: "def f; raise \"boom\"\nrescue => e\ne.message; end; f", expected: "boom"},
		{name: "raise with an exception instance", input: `def f; raise TypeError.new("typed"); rescue TypeError => e; e.message; end; f`, expected: "typed"},
	})
}
