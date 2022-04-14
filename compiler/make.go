package compiler

import "encoding/binary"

// Make - Generates VM instructions from an OpCode & operands
func Make(operator Opcode, operands ...int) Instructions {
	// Look up OpCode definition to figure out how many operands it takes
	def, ok := definitions[operator]
	if !ok {
		return []byte{}
	}

	// Add number of operands to instructionLen, so we know how big of a byte slice to allocate
	// Starts at 1 to make room for the operator
	instructionLen := 1
	for _, w := range def.OperandWidths {
		instructionLen += w
	}

	// Allocate the byte slice and assign the operator to the first position
	instruction := make([]byte, instructionLen)
	instruction[0] = byte(operator)
	offset := 1

	for i, o := range operands {
		width := def.OperandWidths[i]
		switch width {
		case 2: // Operators with a width of 2 are uint16
			binary.BigEndian.PutUint16(instruction[offset:], uint16(o))
		case 1: // Width of 1 are simple bytes
			instruction[offset] = byte(o)
		}
		offset += width
	}

	return instruction
}
