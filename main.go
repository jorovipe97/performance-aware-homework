package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	deco "github.com/jorovipe97/performance-aware-homework/decoder"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("The filename arg is required.")
	}

	fileName := os.Args[1]

	// Get the working directory
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	path := filepath.Join(wd, "listings", fileName)
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("; %v\n", fileName))
	// Tells assembler we intent to run assembly for old 8086 architecture.
	builder.WriteString("bits 16\n")

	decoder := deco.Decoder{
		Data: data,
	}
	for {
		if !decoder.HasNext() {
			break
		}

		opcode, instr, err := decoder.Next()
		if err != nil {
			log.Print(err)
			break
		}
		instrAsmString := decoder.AsmString(opcode, instr)
		builder.WriteString(instrAsmString + "\n")
	}

	// Shows assembly code:
	asm := builder.String()
	fmt.Println(asm)

	// Saves the final assembly into disk.
	newAsmFile := filepath.Join(wd, "result.asm")
	err = os.WriteFile(newAsmFile, []byte(asm), 0644)
	if err != nil {
		log.Fatal(err)
	}
}
