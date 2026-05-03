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
			input:    "Object.new == Object.new",
			expected: false,
		},
		{
			input:    "obj = Object.new; obj == obj",
			expected: true,
		},
		{
			input:    "Object.new != Object.new",
			expected: true,
		},
		{
			input:    "obj = Object.new; obj != obj",
			expected: false,
		},
	}

	runCoreTests(t, tests)
}
