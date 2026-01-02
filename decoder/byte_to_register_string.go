package decoder

// Computes the register name that should use.
//
// isWordOperation: If true instruction operate on word data, else operates on
// byte data. A word are 2 bytes in 8086.
//
// register: Can be either reg or r/m fields.
func byteToRegisterString(isWordOperation bool, register byte) string {
	if isWordOperation {
		return byteToRegisterW1String(register)
	} else {
		return byteToRegisterW0String(register)
	}
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
