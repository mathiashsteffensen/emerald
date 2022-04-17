package compiler

import "testing"

func TestInstructionsString(t *testing.T) {
	instructions := []Instructions{
		Make(OpAdd),
		Make(OpGetLocal, 1),
		Make(OpPushConstant, 2),
		Make(OpPushConstant, 65535),
		Make(OpCloseBlock, 65535, 255),
	}
	expected := `0000 OpAdd
0001 OpGetLocal 1
0003 OpPushConstant 2
0006 OpPushConstant 65535
0009 OpCloseBlock 65535 255
`
	concatted := Instructions{}
	for _, ins := range instructions {
		concatted = append(concatted, ins...)
	}
	if concatted.String() != expected {
		t.Errorf("instructions wrongly formatted.\nwant=%q\ngot=%q", expected, concatted.String())
	}
}
