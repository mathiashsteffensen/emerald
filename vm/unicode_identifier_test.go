package vm

import "testing"

func TestUnicodeIdentifiers(t *testing.T) {
	tests := []vmTestCase{
		{
			name:     "locals",
			input:    "café = 2; π = 3; café + π",
			expected: 5,
		},
		{
			name: "method with predicate suffix, keyword argument and block",
			input: `
				def naïve?(ключ:)
					yield(ключ + 1)
				end

				naïve?(ключ: 7) { |зн| зн * 2 }
			`,
			expected: 16,
		},
		{
			name: "instance variable and setter in a class with a non-ASCII constant name",
			input: `
				class Ärger
					def 名前=(value)
						@名前 = value
					end

					def 名前
						@名前
					end
				end

				ärger = Ärger.new
				ärger.名前 = 4
				ärger.名前
			`,
			expected: 4,
		},
		{
			name:     "global",
			input:    "$π = 3; $π",
			expected: 3,
		},
		{
			name:     "symbol as hash key",
			input:    "{ключ: 1, :日本 => 2}[:日本]",
			expected: 2,
		},
		{
			name:     "quoted and bare symbols are the same symbol",
			input:    `:"日本" == :日本`,
			expected: true,
		},
		{
			name:     "symbol to_s",
			input:    ":café?.to_s",
			expected: "café?",
		},
		{
			name: "scoped constant",
			input: `
				module Ünit
					Ω = 9
				end

				Ünit::Ω
			`,
			expected: 9,
		},
	}

	runVmTests(t, tests)
}
