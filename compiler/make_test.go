package compiler

import (
	"testing"
)

func TestMake(t *testing.T) {
	tests := []struct {
		name     string
		operator Opcode
		operands []int
		expected Instructions
	}{
		{
			"OpPushConstant",
			OpPushConstant,
			[]int{65534},
			[]byte{byte(OpPushConstant), 255, 254},
		},
		{
			"OpAdd",
			OpAdd,
			[]int{},
			[]byte{byte(OpAdd)},
		},
		{
			"OpGetLocal",
			OpGetLocal,
			[]int{255},
			[]byte{byte(OpGetLocal), 255},
		},
		{
			"OpCloseBlock",
			OpCloseBlock,
			[]int{65534, 255},
			[]byte{byte(OpCloseBlock), 255, 254, 255},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			instruction := Make(tt.operator, tt.operands...)

			if len(instruction) != len(tt.expected) {
				t.Fatalf("instruction has wrong length. want=%d, got=%d", len(tt.expected), len(instruction))
			}

			for i, b := range tt.expected {
				if instruction[i] != tt.expected[i] {
					t.Errorf("wrong byte at pos %d. want=%d, got=%d", i, b, instruction[i])
				}
			}
		})
	}
}
