package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Reads an array of binary instructions and iterates over each instruction.
type Decoder struct {
	data []byte
	pos  int
}

func (d *Decoder) HasNext() bool {
	return (d.pos + 1) < len(d.data)
}

func (d *Decoder) Next() ([]byte, error) {
	if !d.HasNext() {
		return nil, io.EOF
	}

	// Creates an slice that contains the bytes of the instruction
	inst := d.data[d.pos : d.pos+2] // end of slice range is exlusive

	d.pos += 2
	return inst, nil
}

type Opcode byte

const (
	// MOV destination, sourcce
	// Register/memory to/from register
	Mov Opcode = 0b0010_0010
)

func byteToOpcodeString(value byte) string {
	switch value {
	case byte(Mov):
		return "mov"
	}

	return ""
}

// RegisterW1 when W = 1, (Instruction operates on word data)
type RegisterW1 byte

const (
	AX RegisterW1 = 0b000
	CX RegisterW1 = 0b001
	DX RegisterW1 = 0b010
	BX RegisterW1 = 0b011
	SP RegisterW1 = 0b100
	BP RegisterW1 = 0b101
	SI RegisterW1 = 0b110
	DI RegisterW1 = 0b111
)

func byteToRegisterW1String(value byte) string {
	switch value {
	case byte(AX):
		return "ax"
	case byte(CX):
		return "cx"
	case byte(DX):
		return "dx"
	case byte(BX):
		return "bx"
	case byte(SP):
		return "sp"
	case byte(BP):
		return "bp"
	case byte(SI):
		return "si"
	case byte(DI):
		return "di"
	}

	return ""
}

// RegisterW0 when W = 0, (Instruction operates on byte data instead of word data)
type RegisterW0 byte

const (
	AL RegisterW0 = 0b000
	CL RegisterW0 = 0b001
	DL RegisterW0 = 0b010
	BL RegisterW0 = 0b011
	AH RegisterW0 = 0b100
	CH RegisterW0 = 0b101
	DH RegisterW0 = 0b110
	BH RegisterW0 = 0b111
)

func byteToRegisterW0String(value byte) string {
	switch value {
	case byte(AL):
		return "al"
	case byte(CL):
		return "cl"
	case byte(DL):
		return "dl"
	case byte(BL):
		return "bl"
	case byte(AH):
		return "ah"
	case byte(CH):
		return "ch"
	case byte(DH):
		return "dh"
	case byte(BH):
		return "bh"
	}

	return ""
}

func decodeMovIntruction(instruction []byte) string {
	fmt.Printf("%b\n", instruction)

	var builder strings.Builder // Zero value is ready to use

	// 1. Decode the opcode:
	// Op code is in the first 6 bits of first byte
	op := instruction[0] >> 2

	// The bit 8 of first byte determine the w field:
	// when 0, instruction operates on byte data
	// when 1, instructions operate on word data
	w := instruction[0] & 0b1
	// 2. Decode the source registry (when bit 7 of first byte is 0, reg is the source)
	// Destination is in in second byte.
	regField := (instruction[1] >> 3) & 0b0000_0111

	// 3. Decode the destination registry.
	rmField := instruction[1] & 0b0000_0111

	builder.WriteString(fmt.Sprint(byteToOpcodeString(op), " "))
	if w == 1 {
		builder.WriteString(
			fmt.Sprintf(
				"%v, %v",
				byteToRegisterW1String(rmField),
				byteToRegisterW1String(regField),
			),
		)
	} else {
		builder.WriteString(
			fmt.Sprintf(
				"%v, %v",
				byteToRegisterW0String(rmField),
				byteToRegisterW0String(regField),
			),
		)
	}

	return builder.String()
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("The filename arg is required.")
	}

	fileName := os.Args[1]

	// Get the absolute path to the executable
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	path := filepath.Join(wd, "decoder", fileName)
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("; %v\n", fileName))
	// Tells assembler we intent to run assembly for old 8086 architecture.
	builder.WriteString("bits 16\n")

	decoder := Decoder{
		data: data,
	}
	for {
		if !decoder.HasNext() {
			break
		}

		instr, err := decoder.Next()
		if err != nil {
			log.Panic("Error decoding instruciton", err)
		}
		intrString := decodeMovIntruction(instr)
		builder.WriteString(intrString + "\n")
	}

	// Shows assembly code:
	asm := builder.String()
	fmt.Println(asm)

	// Saves the final assembly into disk.
	newAsmFile := filepath.Join(wd, "decoder", "result.asm")
	err = os.WriteFile(newAsmFile, []byte(asm), 0644)
	if err != nil {
		log.Fatal(err)
	}
}
