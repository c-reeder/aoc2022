package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

var (
	ErrImproperlyFormattedLine = errors.New("improperly formatted line")
	ErrUnrecognizedInstruction = errors.New("unrecognized instruction")
)

type Instruction string

const (
	NoOp Instruction = "noop"
	AddX Instruction = "addx"
)

var partB bool

func main() {
	// Parse flags
	{
		partBFlag := flag.Bool("b", false, "To switch to part b")
		flag.Parse()
		if partBFlag != nil {
			partB = *partBFlag
		}

	}

	// Check args
	if len(flag.Args()) != 1 {
		log.Fatal("Expected 1 argument containing file name!")
	}

	// Open file
	file, err := os.Open(flag.Args()[0])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	answerSum := 0

	device := NewDevice(file)
	cycleCount := 0
	for {
		x := device.ReadRegisterX()
		next := device.Tick()

		// Part A
		if !partB {
			if (cycleCount-19)%40 == 0 {
				fmt.Printf("Cycle #%v * X:%v = %v\n", (cycleCount + 1), x, (cycleCount+1)*x)
				answerSum += (cycleCount + 1) * x
			}
		}

		if partB {
			if cycleCount%40 == 0 {
				fmt.Println()
			}
			if x-1 == (cycleCount%40) ||
				x+1 == (cycleCount%40) ||
				x == (cycleCount%40) {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}

		if !next {
			break
		}
		cycleCount++
	}

	// Part A cont.
	if !partB {
		fmt.Printf("Sum is %v\n", answerSum)
	}
}

func NewDevice(file *os.File) *Device {
	scanner := bufio.NewScanner(file)
	return &Device{
		x:       1,
		scanner: scanner,
	}

}

type Device struct {
	x                   int
	scanner             *bufio.Scanner
	instrCyclesRemining int
	instruction         Instruction
	instrArg            int
}

// Tick carries out the next clock cycle tick
// Returns whether it was the last cycle
func (d *Device) Tick() bool {
	if d.instrCyclesRemining == 0 {
		if !d.readNextInstruction() {
			return false
		}
	}
	if d.instrCyclesRemining == 1 {
		d.carryOutCurrInstr()
	}
	d.instrCyclesRemining--
	return true
}

// readNextInstruction loads up the device state with the next instruction
// this includes the instruction, it's argument(if any), and how many cycles
// it will need to complete
func (d *Device) readNextInstruction() bool {
	//before := d.instrCyclesRemining
	if !d.scanner.Scan() {
		return false
	}
	line := d.scanner.Text()
	err := d.scanner.Err()
	if err != nil {
		log.Fatal(err)
	}
	tokens := strings.Split(line, " ")
	if len(tokens) < 1 {
		log.Fatalln(ErrImproperlyFormattedLine)
	}
	d.instruction = Instruction(tokens[0])
	switch d.instruction {
	case NoOp:
		if len(tokens) != 1 {
			log.Fatalln(ErrImproperlyFormattedLine)
		}
		d.instrCyclesRemining = 1
	case AddX:
		if len(tokens) != 2 {
			log.Fatalln(ErrImproperlyFormattedLine)
		}
		d.instrArg, err = strconv.Atoi(tokens[1])
		if err != nil {
			log.Fatalln(ErrImproperlyFormattedLine)
		}
		d.instrCyclesRemining = 2
	default:
		log.Fatalln(ErrUnrecognizedInstruction)
	}
	return true
}

// carryOutCurrInstr carries out the instruction on the last cycle
// of it's execution
func (d *Device) carryOutCurrInstr() {
	switch d.instruction {
	case AddX:
		d.x += d.instrArg
	default:
	}
}

// ReadRegisterX Reads the current contents of register X
func (d *Device) ReadRegisterX() int {
	return d.x
}
