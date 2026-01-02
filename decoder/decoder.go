package decoder

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

// Details of the mov operation are after page 160 of 8086 user's manual

// Reads an array of binary instructions and iterates over each instruction.
type Decoder struct {
	Data []byte
	pos  int
}

func (d *Decoder) HasNext() bool {
	return (d.pos + 1) < len(d.Data)
}

func (d *Decoder) Next() (Opcode, []byte, error) {
	if !d.HasNext() {
		return 0, nil, io.EOF
	}

	// We pass in the next two bytes, to try to analyze the opcode.
	// Creates an slice that contains the bytes of the instruction
	inst := d.Data[d.pos : d.pos+2] // end of slice range is exlusive
	opcode, bytesToRead, error := d.analyzeOpCode(inst)

	if error != nil {
		return 0, nil, error
	}

	fullInstruction := d.Data[d.pos : d.pos+bytesToRead]
	d.pos += bytesToRead
	return opcode, fullInstruction, nil
}

type Opcode byte

const (
	// MOV destination, sourcce
	// Register/memory to/from register
	MovRegisterMemoryToFromRegister Opcode = 0b0010_0010

	// Immediate to register.
	MovImmediateToRegister Opcode = 0b1011
)

// Returns the opcode name and the lenght of bytes to read for this opcode.
func (d *Decoder) analyzeOpCode(instruction []byte) (Opcode, int, error) {
	firstByte := instruction[0]

	// op code is usually encoded in the first 6 bits of the first byte.
	if firstByte>>2 == byte(MovRegisterMemoryToFromRegister) {
		// Register mode/Memory mode with displacement length
		modField := instruction[1] >> 6
		var bytesToRead int = 0

		switch modField {
		case 0b00:
			// Memory mode, no displacement follows.
			// Except when R/M field = 110, then, 16-bit displacement follwos.
			bytesToRead = 2
		case 0b01:
			// Memory mode, 8 bit displacement follows
			bytesToRead = 3 // An additional byte
		case 0b10:
			// Memory mode, 16 bit displacement follows
			bytesToRead = 4 // Two additional bytes.
		case 0b11:
			// Register mode (no displacement)
			bytesToRead = 2
		}

		return MovRegisterMemoryToFromRegister, bytesToRead, nil
	} else if firstByte>>4 == byte(MovImmediateToRegister) {
		var bytesToRead int = 2
		var isWord bool = (firstByte>>3)&0b00001 == 1
		if isWord {
			bytesToRead = 3
		}

		return MovImmediateToRegister, bytesToRead, nil
	}

	return 0, 0, errors.New("cannot identify instruction")
}

// Register/memory to/from register
// [1 0 0 0 1 0 d w]
// [mod(2 bits) reg(3 bits) rm(3 bits)]
// [Displacement Low (8 bits)]
// [Displacement Hight (8 bits)]
func decodeMovRegisterMemoryToFromRegister(instruction []byte) string {
	var builder strings.Builder // Zero value is ready to use

	// The bit 8 of first byte determine the w field:
	// when 0, instruction operates on byte data
	// when 1, instructions operate on word data
	w := instruction[0]&0b1 == 1
	// 2. Decode the source registry (when bit 7 of first byte is 0, reg is the source)
	// Destination is in in second byte.
	regField := (instruction[1] >> 3) & 0b0000_0111

	// 3. Decode the destination registry.
	rmField := instruction[1] & 0b0000_0111

	builder.WriteString("mov ")
	builder.WriteString(
		fmt.Sprintf(
			"%v, %v",
			byteToRegisterString(w, rmField),
			byteToRegisterString(w, regField),
		),
	)

	return builder.String()
}

// Immediate to register.
// [1 0 1 1 w reg(3 bits)]
// [data(8 bits)]
// [data(8 bits - if w = 1)]
func decodeMovImmediateToRegister(instruction []byte) string {
	var builder strings.Builder // Zero value is ready to use

	// The bit 8 of first byte determine the isWord field:
	// when 0, instruction operates on byte data
	// when 1, instructions operate on word data
	isWord := (instruction[0]>>3)&0b1 == 1
	// 2. Decode the source registry (when bit 7 of first byte is 0, reg is the source)
	// Destination is in in second byte.
	regField := instruction[0] & 0b0000_0111

	var data uint16 = uint16(instruction[1])
	// 16-bit immediate-to-register
	if isWord {
		data = data | (uint16(instruction[2]) << 8)
	}

	builder.WriteString("mov ")
	builder.WriteString(
		fmt.Sprintf(
			"%v, %v",
			byteToRegisterString(isWord, regField),
			data,
		),
	)

	return builder.String()
}

func (d *Decoder) AsmString(opcode Opcode, instruction []byte) string {
	// op code is usually encoded in the first 6 bits of the first byte.
	switch opcode {
	case MovRegisterMemoryToFromRegister:
		return decodeMovRegisterMemoryToFromRegister(instruction)
	case MovImmediateToRegister:
		return decodeMovImmediateToRegister(instruction)
	}

	return ""
}
