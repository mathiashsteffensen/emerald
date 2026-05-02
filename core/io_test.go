package core_test

import "testing"

func TestIO_sysopen(t *testing.T) {
	tests := []coreTestCase{
		{
			input: `
				file_descriptor = rt.IO.sysopen "fixtures/require_test.rb"
				io = rt.IO.new file_descriptor
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
				file_descriptor = rt.IO.sysopen "fixtures/require_test.rb"
				io = rt.IO.open file_descriptor
				b = io.getbyte
				io.close
				b
			`,
			expected: 100,
		},
		{
			name: "when called with a block",
			input: `
				file_descriptor = rt.IO.sysopen "fixtures/require_test.rb"
				rt.IO.open(file_descriptor) do |io|
					io.getbyte
				end
			`,
			expected: 100,
		},
	}

	runCoreTests(t, tests)
}

func TestIO_read(t *testing.T) {
	tests := []coreTestCase{
		{
			name:     "reads the whole file",
			input:    `rt.IO.read("core/fixtures/require_test.rb").size`,
			file:     "../io.rb",
			expected: 242,
		},
	}

	runCoreTests(t, tests)
}
