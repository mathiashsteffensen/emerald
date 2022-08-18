package core_test

import "testing"

func TestIntegerOperators(t *testing.T) {
	tests := []coreTestCase{
		{
			name:     "adding",
			input:    "1+2",
			expected: 3,
		},
		{
			name:     "adding with .",
			input:    "1.+(2)",
			expected: 3,
		},
		{
			name:     "subtracting",
			input:    "1-2",
			expected: -1,
		},
		{
			name:     "subtracting with .",
			input:    "1.-(2)",
			expected: -1,
		},
		{
			name:     "dividing",
			input:    "6/2",
			expected: 3,
		},
		{
			name:     "dividing with .",
			input:    "6./(2)",
			expected: 3,
		},
		{
			name:     "dividing with remainder",
			input:    "5/2",
			expected: 2.5,
		},
		{
			name:     "Multiplying",
			input:    "6*2",
			expected: 12,
		},
		{
			name:     "multiplying with .",
			input:    "6.*(2)",
			expected: 12,
		},
		{
			name:     "equals",
			input:    "6==2",
			expected: false,
		},
		{
			name:     "equals with .",
			input:    "6.==(2)",
			expected: false,
		},
		{
			name:     "not equals",
			input:    "6!=2",
			expected: true,
		},
		{
			name:     "not equals with .",
			input:    "6.!=(2)",
			expected: true,
		},
		{
			name:     "less than",
			input:    "6<2",
			expected: false,
		},
		{
			name:     "less than with .",
			input:    "6.<(2)",
			expected: false,
		},
		{
			name:     "greater than",
			input:    "6>2",
			expected: true,
		},
		{
			name:     "greater than with .",
			input:    "6.>(2)",
			expected: true,
		},
		{
			name:     "less than or equals",
			input:    "6<=2",
			expected: false,
		},
		{
			name:     "less than or equals with .",
			input:    "6.<=(2)",
			expected: false,
		},
		{
			name:     "greater than or equals",
			input:    "6>=2",
			expected: true,
		},
		{
			name:     "greater than or equals with .",
			input:    "6.>=(2)",
			expected: true,
		},
	}

	runCoreTests(t, tests)
}

func TestInteger_to_s(t *testing.T) {
	tests := []coreTestCase{
		{
			input:    "1231.to_s",
			expected: "1231",
		},
	}

	runCoreTests(t, tests)
}

func TestInteger_times(t *testing.T) {
	tests := []coreTestCase{
		{
			input: `
			count = 0
			1000.times { count = count + 1 }
			count
			`,
			expected: 1000,
		},
		{
			input: `
			count = 0
			5.times { |i| count = count + i }
			count
			`,
			expected: 10,
		},
	}

	runCoreTests(t, tests)
}
