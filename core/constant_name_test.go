package core_test

import (
	"emerald/core"
	"testing"
)

func TestIsConstantName(t *testing.T) {
	tests := []struct {
		name     string
		expected bool
	}{
		{"Foo", true},
		{"Ärger", true},
		{"ǅx", true},
		{"Ω", true},
		{"Ⅰ", true},
		{"Ⓐx", true},
		{"foo", false},
		{"_Foo", false},
		{"éclair", false},
		{"π", false},
		{"日本", false},
		{"", false},
	}

	for _, tt := range tests {
		if actual := core.IsConstantName(tt.name); actual != tt.expected {
			t.Errorf("IsConstantName(%q): expected %t, got %t", tt.name, tt.expected, actual)
		}
	}
}
