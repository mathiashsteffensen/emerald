package compiler

import "testing"

func TestInstructionsString(t *testing.T) {
	instructions := []Instructions{
		Make(OpAdd),
		Make(OpPushConstant, 2),
		Make(OpPushConstant, 65535),
	}
	expected := `0000 OpAdd
0001 OpPushConstant 2
0004 OpPushConstant 65535
`
	concatted := Instructions{}
	for _, ins := range instructions {
		concatted = append(concatted, ins...)
	}
	if concatted.String() != expected {
		t.Errorf("instructions wrongly formatted.\nwant=%q\ngot=%q", expected, concatted.String())
	}
}
