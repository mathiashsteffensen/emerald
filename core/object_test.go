package core_test

import "testing"

func TestObject_to_s(t *testing.T) {
	tests := []coreTestCase{
		{
			input:    "to_s",
			expected: "main",
		},
	}

	runCoreTests(t, tests)
}

func TestObject_equals(t *testing.T) {
	tests := []coreTestCase{
		{
			input:    "rt.Object.new == rt.Object.new",
			expected: false,
		},
		{
			input:    "obj = rt.Object.new; obj == obj",
			expected: true,
		},
		{
			input:    "rt.Object.new != rt.Object.new",
			expected: true,
		},
		{
			input:    "obj = rt.Object.new; obj != obj",
			expected: false,
		},
	}

	runCoreTests(t, tests)
}
