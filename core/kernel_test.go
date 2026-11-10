package core_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestKernel_raise(t *testing.T) {
	tests := []coreTestCase{
		{
			name:     "raising with just a string",
			input:    `raise "You done fucked up this time"`,
			expected: "error:RuntimeError:You done fucked up this time",
		},
		{
			name:     "raising with a specific error class",
			input:    `raise TypeError, "I wanted an Integer yo"`,
			expected: "error:TypeError:I wanted an Integer yo",
		},
		{
			name: "raising with just a custom error class",
			input: `
				class MyError < StandardError; end
				def rescued_class
					raise MyError
				rescue MyError => e
					e.class
				end
				rescued_class.to_s
			`,
			expected: "MyError",
		},
		{
			name: "raising with just a built-in error class",
			input: `
				def rescued_class
					raise StandardError
				rescue StandardError => e
					e.class
				end
				rescued_class.to_s
			`,
			expected: "StandardError",
		},
		{
			name:     "raising with just a non-exception class",
			input:    `raise Object`,
			expected: "error:TypeError:exception class/object expected",
		},
		{
			name:     "raising with a non-exception class and a message",
			input:    `raise Object, "x"`,
			expected: "error:TypeError:exception class/object expected",
		},
		{
			name: "raising a non-exception class does not call its constructor",
			input: `
				class NeedsArgument
					def initialize(a); end
				end
				raise NeedsArgument
			`,
			expected: "error:TypeError:exception class/object expected",
		},
		{
			name:     "raising with just a non-class",
			input:    `raise 3`,
			expected: "error:TypeError:exception class/object expected",
		},
		{
			name:     "raising with just a module",
			input:    `raise Kernel`,
			expected: "error:TypeError:exception class/object expected",
		},
		{
			name:     "raising inside a rescue body",
			input:    `def f; raise StandardError, "a"; rescue => e; raise TypeError, "b"; end; f`,
			expected: "error:TypeError:b",
		},
		{
			name:     "raising with a non-class and a message",
			input:    `raise 3, "x"`,
			expected: "error:TypeError:exception class/object expected",
		},
		{
			name: "raising a class whose initialize needs an argument",
			input: `
				class Needy < StandardError
					def initialize(a); end
				end
				raise Needy
			`,
			expected: "error:ArgumentError:wrong number of arguments (given 0, expected 1)",
		},
		{
			name: "raising a class whose initialize raises",
			input: `
				class Exploding < StandardError
					def initialize(message)
						raise TypeError, "from initialize"
					end
				end
				raise Exploding, "x"
			`,
			expected: "error:TypeError:from initialize",
		},
		{
			name: "a non-exception result from new does not panic",
			input: `
				class Impostor < StandardError
					class << self
						def new
							3
						end
					end
				end
				raise Impostor
			`,
			expected: "error:TypeError:exception object expected",
		},
	}
	runCoreTests(t, tests)
}

func TestKernel_require_relative(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	rootDir := filepath.Join(dir, "../")

	tests := []coreTestCase{
		{
			input: `
				require_relative("fixtures/require_test")
				sleep 0.01
				require_relative("../spec/fixtures/require_test")
			`,
			expected: true,
		},
		{
			name: "requiring same file twice",
			input: `
				require_relative("fixtures/require_test")
				require_relative("fixtures/require_test")
			`,
			expected: false,
		},
		{
			name:     "when file doesn't exist",
			input:    `require_relative("../lib/main")`,
			expected: fmt.Sprintf("error:LoadError:cannot load such file -- %s/lib/main", rootDir),
		},
		{
			name: "resolving namespaced constant name in required file",
			input: `
				require_relative("fixtures/namespaced_class")
				
				module A
					module C
						B.name
					end
				end
			`,
			expected: "A::B",
		},
		{
			name: "return in a block at the top level of the required file",
			input: `
				require_relative("fixtures/return_in_block")
				[$return_in_block_loaded, :main_goes_on]
			`,
			expected: []any{true, ":main_goes_on"},
		},
	}

	runCoreTests(t, tests)
}

func TestKernel_class(t *testing.T) {
	tests := []coreTestCase{
		{
			input:    "Object.new.class",
			expected: "class:Object",
		},
		{
			input:    "Object.class",
			expected: "class:Class",
		},
	}

	runCoreTests(t, tests)
}

func TestKernel_kind_of(t *testing.T) {
	tests := []coreTestCase{
		{
			name:     "when self is instance of class",
			input:    "Object.new.kind_of?(Object)",
			expected: true,
		},
		{
			name:     "when self is instance of subclass",
			input:    "String.new.kind_of?(Object)",
			expected: true,
		},
		{
			name:     "singleton is instance of Class",
			input:    "String.kind_of?(Class)",
			expected: true,
		},
		{
			name: "when passed a not included module",
			input: `
			module MyMod; end
			class MyClass; end
			MyClass.new.kind_of?(MyMod)`,
			expected: false,
		},
		{
			name: "when passed an included module",
			input: `
			module MyMod; end
			class MyClass
				include(MyMod)
			end
			MyClass.new.kind_of?(MyMod)`,
			expected: true,
		},
		{
			name:     "when passed wrong type of arg",
			input:    "String.kind_of?(23)",
			expected: "error:TypeError:class or module required",
		},
	}

	runCoreTests(t, tests)
}

func TestKernel_puts(t *testing.T) {
	tests := []coreTestCase{
		{
			name:     "with a single string argument",
			input:    `puts("Hello World!")`,
			expected: nil,
		},
		{
			name:     "with multiple string arguments",
			input:    `puts("Hello", "World!")`,
			expected: nil,
		},
		{
			name:     "with nil argument",
			input:    `puts(nil)`,
			expected: nil,
		},
	}

	runCoreTests(t, tests)
}

func TestKernel_include(t *testing.T) {
	tests := []coreTestCase{
		{
			name: "in main object",
			input: `module MyMod
				def hello; "Hello"; end
			end
			include(MyMod)
			hello`,
			expected: "Hello",
		},
		{
			name: "in custom class",
			input: `module MyMod
				def hello; "Hello"; end
			end
		
			class MyClass
				include(MyMod)
			end
		
			MyClass.new.hello`,
			expected: "Hello",
		},
		{
			name: "in custom static class",
			input: `module MyMod
				def hello; "Hello"; end
			end
			
			class MyClass
				class << self
					include(MyMod)
				end
			end

			MyClass.hello`,
			expected: "Hello",
		},
		{
			name: "wrong argument type",
			input: `
				class C
					include "boop"
				end
			`,
			expected: "error:TypeError:wrong argument type String (expected Module)",
		},
	}

	runCoreTests(t, tests)
}
