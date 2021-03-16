package vm

import "fmt"

const (
	Load  = 0x01
	Store = 0x02
	Add   = 0x03
	Sub   = 0x04
	Halt  = 0xff
)

// Stretch goals
const (
	Addi = 0x05
	Subi = 0x06
	Jump = 0x07
	Beqz = 0x08
)

// Given a 256 byte array of "memory", run the stored program
// to completion, modifying the data in place to reflect the result
//
// The memory format is:
//
// 00 01 02 03 04 05 06 07 08 09 0a 0b 0c 0d 0e 0f ... ff
// __ __ __ __ __ __ __ __ __ __ __ __ __ __ __ __ ... __
// ^==DATA===============^ ^==INSTRUCTIONS==============^
//
func compute(memory []byte) {
	registers := []byte{0x08, 0x00, 0x00}
	pc := registers[0]
	for {
		op := memory[pc]
		switch op {
		// begin "normal" ops
		// byte immediately after pc is a register
		// byte after register is a memory address
		case Load:
			registers[memory[pc+1]] = memory[memory[pc+2]]
		case Store:
			addr := memory[pc+2]
			if addr > 7 {
				panic(fmt.Sprintf("invalid memory access %x", addr))
			}
			memory[addr] = registers[memory[pc+1]]
		case Add:
			registers[memory[pc+1]] = registers[memory[pc+1]] + registers[memory[pc+2]]
		case Sub:
			registers[memory[pc+1]] = registers[memory[pc+1]] - registers[memory[pc+2]]
			// begin "immediate" ops
			// byte immediately after pc is a register
			// byte after register is a value to be added/subtracted
		case Addi:
			registers[memory[pc+1]] = registers[memory[pc+1]] + memory[pc+2]
		case Subi:
			registers[memory[pc+1]] = registers[memory[pc+1]] - memory[pc+2]
			// begin "flow control" ops
		case Jump:
			// byte immediately after pc is a memory address to jump to
			// skip the normal pc increment
			pc = memory[pc+1]
			continue
		case Beqz:
			// byte immediately after pc is a register to check
			// byte after register is the memory address to jump to if reg value is 0
			if registers[memory[pc+1]] == 0 {
				imm := memory[pc+2]
				pc = pc + imm
			}
		case Halt:
			return
		default:
			panic(fmt.Sprintf("unexpected op code: %d. exiting as this program may not do what you want", op))
		}
		// most opcodes + operands are 3 bytes long, except in cases where we explicitly manipulate the pc or halt execution.
		pc = pc + 3
	}
}
