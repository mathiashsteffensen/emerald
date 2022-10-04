package core_test

import "testing"

func TestFalseClass_to_s(t *testing.T) {
	tests := []coreTestCase{
		{
			input:    "false.to_s",
			expected: "false",
		},
	}

	runCoreTests(t, tests)
}
