package core_test

import (
	"math"
	"testing"
)

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
			name:     "adding a Float",
			input:    "2+1.5",
			expected: 3.5,
		},
		{
			name:     "subtracting a Float",
			input:    "2-1.5",
			expected: 0.5,
		},
		{
			name:     "multiplying by a Float",
			input:    "2*1.5",
			expected: 3.0,
		},
		{
			name:     "dividing by a Float",
			input:    "3/1.5",
			expected: 2.0,
		},
		{
			name:     "dividing by Float zero",
			input:    "3/0.0",
			expected: math.Inf(1),
		},
		{
			name:     "shifting left",
			input:    "6<<2",
			expected: 24,
		},
		{
			name:     "shifting left by a negative count",
			input:    "8<<-1",
			expected: 4,
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
		{
			input:    "6 + false",
			expected: "error:TypeError:false can't be coerced into Integer",
		},
		{
			input:    `6 + "1"`,
			expected: "error:TypeError:String can't be coerced into Integer",
		},
		{
			input:    "6.+()",
			expected: "error:ArgumentError:wrong number of arguments (given 0, expected 1)",
		},
		{
			input:    "6 << nil",
			expected: "error:TypeError:no implicit conversion of NilClass into Integer",
		},
		{
			input:    "6.<<()",
			expected: "error:ArgumentError:wrong number of arguments (given 0, expected 1)",
		},
		{
			input:    "6 === 6",
			expected: true,
		},
		{
			input:    "6 === 6.0",
			expected: true,
		},
		{
			input:    "6 === Integer",
			expected: false,
		},
	}

	runCoreTests(t, tests)
}

func TestInteger_spaceship(t *testing.T) {
	tests := []coreTestCase{
		{
			input:    "1 <=> 2",
			expected: -1,
		},
		{
			input:    "3 <=> 2",
			expected: 1,
		},
		{
			input:    "2 <=> 2",
			expected: 0,
		},
		{
			input:    "-2 <=> 2",
			expected: -1,
		},
		{
			input:    "70 <=> 100.0",
			expected: -1,
		},
		{
			input:    "70 <=> 70.0",
			expected: 0,
		},
		{
			input:    "100 <=> 70.0",
			expected: 1,
		},
		{
			input:    `-2 <=> "boop"`,
			expected: nil,
		},
		{
			name:     "with NaN",
			input:    "1 <=> (0.0 / 0.0)",
			expected: nil,
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

func TestInteger_to_f(t *testing.T) {
	tests := []coreTestCase{
		{
			input:    "1231.to_f",
			expected: 1231.0,
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
