package core_test

import "testing"

func TestString_new_subclass(t *testing.T) {
	runCoreTests(t, []coreTestCase{
		{
			name:     "preserves subclass identity",
			input:    `SpecialString.new.class`,
			expected: "class:SpecialString",
		},
		{
			name:     "dispatches subclass methods",
			input:    `SpecialString.new.marker`,
			expected: 42,
		},
		{
			name: "retains native string storage",
			input: `
				s = SpecialString.new
				[s.size, s.upcase, s + "hello", s * 2, s == "", s.start_with?(""), s.to_s.class]
			`,
			expected: []any{0, "", "hello", "", true, true, "class:SpecialString"},
		},
		{
			name: "inherits class methods through another subclass",
			input: `
				class ChildString < SpecialString; end
				[ChildString.build.class, ChildString.build.marker]
			`,
			expected: []any{"class:ChildString", 42},
		},
		{
			name: "preserves ignored arguments without initialize",
			input: `
				s = SpecialString.new("ignored", 2, value: 3)
				[s.class, s.size]
			`,
			expected: []any{"class:SpecialString", 0},
		},
		{
			name: "forwards initializer arguments keywords and block once",
			input: `
				class InitializedString < SpecialString
					def initialize(left, right, value:)
						@values = [left, right, value, yield(left + right), size]
						99
					end
					def values; @values; end
				end
				class ChildString < InitializedString; end
				calls = 0
				s = ChildString.new(2, 3, value: 7) { |sum| calls += 1; sum * 2 }
				[s.class, s.marker, s.values, calls, s.upcase]
			`,
			expected: []any{"class:ChildString", 42, []any{2, 3, 7, 10, 0}, 1, ""},
		},
		{
			name: "initializer validates positional arity",
			input: `
				class SpecialString
					def initialize(value); end
				end
				SpecialString.new
			`,
			expected: "error:ArgumentError:wrong number of arguments (given 0, expected 1)",
		},
		{
			name: "initializer validates keywords",
			input: `
				class SpecialString
					def initialize(value:); end
				end
				SpecialString.new
			`,
			expected: "error:ArgumentError:missing keyword: :value",
		},
		{
			name: "propagates initializer exceptions",
			input: `
				class SpecialString
					def initialize; raise "initialization failed"; end
				end
				SpecialString.new
			`,
			expected: "error:RuntimeError:initialization failed",
		},
	}, `
		class SpecialString < String
			def marker; 42; end
			class << self
				def build; new; end
			end
		end
	`)
}

func TestString_new(t *testing.T) {
	runCoreTests(t, []coreTestCase{
		{
			name:     "empty string",
			input:    `[String.new.class, String.new, String.new.size]`,
			expected: []any{"class:String", "", 0},
		},
		{
			name:     "existing string remains the same object",
			input:    `s = "hello"; [String.new(s).class, String.new(s).upcase, String.new(s) != s]`,
			expected: []any{"class:String", "HELLO", false},
		},
		{
			name:     "rejects nonstrings",
			input:    `String.new(42)`,
			expected: "error:TypeError:no implicit conversion of Integer into String",
		},
		{
			name:     "rejects extra arguments",
			input:    `String.new("a", "b")`,
			expected: "error:ArgumentError:wrong number of arguments (given 2, expected 1)",
		},
		{
			name:     "rejects keywords",
			input:    `String.new(value: "a")`,
			expected: "error:ArgumentError:unknown keyword: :value",
		},
	})
}

func TestString_to_sym(t *testing.T) {
	tests := []coreTestCase{
		{
			input:    `"hello".to_sym`,
			expected: ":hello",
		},
	}

	runCoreTests(t, tests)
}

func TestString_add(t *testing.T) {
	tests := []coreTestCase{
		{
			input:    `"hello" + " " + "world"`,
			expected: "hello world",
		},
		{
			input:    `"wdaw" + 2`,
			expected: "error:TypeError:no implicit conversion of Integer into String",
		},
	}

	runCoreTests(t, tests)
}

func TestString_multiply(t *testing.T) {
	tests := []coreTestCase{
		{
			input:    `"w" * 5`,
			expected: "wwwww",
		},
		{
			input:    `"w" * ""`,
			expected: "error:TypeError:no implicit conversion of String into Integer",
		},
	}

	runCoreTests(t, tests)
}

func TestString_match(t *testing.T) {
	tests := []coreTestCase{
		{
			input:    `("a" =~ /a/).is_a?(MatchData)`,
			expected: true,
		},
		{
			input:    `"a".match(/a/).is_a?(MatchData)`,
			expected: true,
		},
	}

	runCoreTests(t, tests)
}

func TestString_upcase(t *testing.T) {
	tests := []coreTestCase{
		{
			input:    `"hello".upcase`,
			expected: "HELLO",
		},
	}

	runCoreTests(t, tests)
}

func TestString_size(t *testing.T) {
	tests := []coreTestCase{
		{
			input:    `"Hello".size`,
			expected: 5,
		},
	}

	runCoreTests(t, tests)
}

func TestString_split(t *testing.T) {
	tests := []coreTestCase{
		{
			input:    `"h e llo o".split`,
			expected: []any{"h", "e", "llo", "o"},
		},
		{
			input:    `"h e llo o".split("o")`,
			expected: []any{"h e ll", " ", ""},
		},
	}

	runCoreTests(t, tests)
}

func TestString_start_with(t *testing.T) {
	tests := []coreTestCase{
		{
			name:     "matches an exact prefix",
			input:    `"hello".start_with?("hello")`,
			expected: true,
		},
		{
			name:     "matches a shorter prefix",
			input:    `"hello".start_with?("he")`,
			expected: true,
		},
		{
			name:     "does not match text occurring later",
			input:    `"hello".start_with?("ell")`,
			expected: false,
		},
		{
			name:     "does not match a prefix longer than the receiver",
			input:    `"hi".start_with?("hello")`,
			expected: false,
		},
		{
			name:     "matches an empty prefix",
			input:    `"hello".start_with?("")`,
			expected: true,
		},
		{
			name:     "does not match a nonempty prefix on an empty receiver",
			input:    `"".start_with?("hello")`,
			expected: false,
		},
		{
			name:     "returns false without prefixes",
			input:    `"hello".start_with?`,
			expected: false,
		},
		{
			name:     "matches any prefix argument",
			input:    `"hello".start_with?("no", "he")`,
			expected: true,
		},
		{
			name:     "returns false when no prefixes match",
			input:    `"hello".start_with?("no", "x")`,
			expected: false,
		},
		{
			name:     "short circuits after a matching prefix",
			input:    `"hello".start_with?("he", 1)`,
			expected: true,
		},
		{
			name:     "requires string prefixes until a match is found",
			input:    `"hello".start_with?("no", 1)`,
			expected: "error:TypeError:no implicit conversion of Integer into String",
		},
	}

	runCoreTests(t, tests)
}
