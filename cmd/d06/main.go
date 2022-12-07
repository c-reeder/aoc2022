package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"unicode"
)

var ErrEncounteredBadRune = errors.New("encountered bad rune")
var ErrBadfile = errors.New("bad file")

var markerSize int = 4

func main() {
	// Parse flags
	{
		partBFlag := flag.Bool("b", false, "To switch to part b")
		flag.Parse()
		if partBFlag != nil && *partBFlag {
			markerSize = 14
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

	i := 0                             // Index of rune in transmission
	buffer := make([]rune, markerSize) // buffer of whatever desired marker size
	start := 0                         // The index within the buffer at which the current run of distint runes starts
	length := 0                        // The length of the current run of distinct runes stored in the buffer

	reader := bufio.NewReader(file)
	for {
		// Read next rune
		r, _, err := reader.ReadRune()
		// Validate
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(ErrBadfile)
		}
		if r == unicode.ReplacementChar {
			log.Fatal(ErrEncounteredBadRune)
		}

		// start main logic
		// Loop over currently captured distinct runes to ensure
		// the new addition doesn't match one of them
		// If it does then increate the start to just past the match
		// and decrease the length accordingly
		for x := 0; x < length; x++ {
			if buffer[(x+start)%markerSize] == r {
				start = (start + x + 1) % markerSize
				length = length - (x + 1)
				break
			}
		}

		// Add the new rune to the buffer and increment size
		buffer[(start+length)%markerSize] = r
		length++

		// Stop when we have a full marker
		if length == markerSize {
			break
		}
		i++
		// end main logic
	}
	for x := 0; x < markerSize; x++ {
		fmt.Printf("%s", string(buffer[(start+x)%markerSize]))
	}
	fmt.Printf("\nbuffer: length: %v\n", i+1)
}
