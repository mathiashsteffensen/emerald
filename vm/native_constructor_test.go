package vm

import (
	"fmt"
	"testing"
)

func TestInheritedNativeConstructor(t *testing.T) {
	for _, parent := range []string{"Exception", "StandardError", "ArgumentError", "TypeError", "NameError", "NoMethodError", "RuntimeError", "LoadError", "String", "Regexp", "Range", "Time", "IO", "TCPServer", "Class"} {
		t.Run(parent, func(t *testing.T) {
			runVmTests(t, []vmTestCase{
				{name: "no argument initialize", input: `class Custom < ` + parent + `
					def initialize; @value = 42; end
					def value; @value; end
				end
				item = Custom.new; [item.class == Custom, item.value]`, expected: []any{true, 42}},
				{name: "forward arguments keywords and block", input: `class Custom < ` + parent + `
					def initialize(value, key:); @value = [value, key, yield]; end
					def value; @value; end
				end
				class Child < Custom; end
				item = Child.new(42, key: 7) { 9 }; [item.class == Child, item.value]`, expected: []any{true, []any{42, 7, 9}}},
				{name: "inherited user new", input: fmt.Sprintf(`class Custom < %s
					class << self; def new; 42; end; end
				end
				class Child < Custom; end
				Child.new`, parent), expected: 42},
			})
		})
	}
}

func TestInheritedNativeConstructorStorage(t *testing.T) {
	runVmTests(t, []vmTestCase{
		{name: "exception", input: `class Custom < StandardError; def marker; 42; end; end
			item = Custom.new("boom"); [item.class == Custom, item.marker, item.to_s]`, expected: []any{true, 42, ""}},
		{name: "regexp", input: `class Custom < Regexp; def initialize; @value = 42; end; end
			item = Custom.new; [item.class == Custom, item.inspect, item === "abc"]`, expected: []any{true, "//", true}},
		{name: "string", input: `class Custom < String; end; [Custom.new.class == Custom, Custom.new.size]`, expected: []any{true, 0}},
		{name: "range", input: `class Custom < Range; end
			item = Custom.new; values = []; item.each { |value| values << value }; [item.class == Custom, values]`, expected: []any{true, []any{nil}}},
		{name: "time", input: `class Custom < Time; end; item = Custom.new; [item.class == Custom, item - item]`, expected: []any{true, float64(0)}},
		{name: "class", input: `class Custom < Class; end; item = Custom.new; [item.class == Custom, item.name, item.new.class == item]`, expected: []any{true, "", true}},
	})
}
