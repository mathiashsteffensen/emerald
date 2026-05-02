package core_test

import "testing"

func TestFile_is_absolute_path(t *testing.T) {
	tests := []coreTestCase{
		{
			input:    `rt.File.absolute_path? "./hello.rb"`,
			expected: false,
		},
		{
			input:    `rt.File.absolute_path? "hello.rb"`,
			expected: false,
		},
		{
			input:    `rt.File.absolute_path? "/hello.rb"`,
			expected: true,
		},
	}

	runCoreTests(t, tests)
}
