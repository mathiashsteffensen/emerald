package core_test

import "testing"

func TestIO_sysopen(t *testing.T) {
	tests := []coreTestCase{
		{
			input: `
				file_descriptor = IO.sysopen "fixtures/require_test.rb"
				io = IO.new file_descriptor
				b = io.getbyte
				io.close
				b
			`,
			expected: 100,
		},
	}

	runCoreTests(t, tests)
}

func TestIO_open(t *testing.T) {
	tests := []coreTestCase{
		{
			name: "when called without a block",
			input: `
				file_descriptor = IO.sysopen "fixtures/require_test.rb"
				io = IO.open file_descriptor
				b = io.getbyte
				io.close
				b
			`,
			expected: 100,
		},
		{
			name: "when called with a block",
			input: `
				file_descriptor = IO.sysopen "fixtures/require_test.rb"
				IO.open(file_descriptor) do |io|
					io.getbyte
				end
			`,
			expected: 100,
		},
	}

	runCoreTests(t, tests)
}
