package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var ErrOversizedLine = errors.New("line exceeds maximum length")
var ErrImproperlyFormattedLine = errors.New("improperly formatted line")

const MaxByesPerLine = 3 * 1024 // 3kB max line length
var commandRgx = regexp.MustCompile(`move (\d+) from (\d+) to (\d+)`)

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

	// slice to be used as a stack of the lines of the
	// text file containing the diagram
	dgrmLines := make([]string, 0)

	// Fill the dgrmLines stack
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Validate max line length
		if len(line) > MaxByesPerLine {
			log.Fatal(ErrOversizedLine)
		}

		if len(line) == 0 {
			break
		}

		dgrmLines = append(dgrmLines, line)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// Calculate how many stacks/piles of crates their are
	numStacks := len(strings.Fields(dgrmLines[len(dgrmLines)-1]))

	// 2D slice containing them
	stacks := make([][]rune, numStacks)

	// How long each line in the diagram should be (in runes)
	lineLength := 4*numStacks - 1

	// unwind the stack of diagram lines to fill the 2D slice with runes
	for i := len(dgrmLines) - 2; i >= 0; i-- {
		runes := []rune(dgrmLines[i])
		if len(runes) != lineLength {
			log.Fatal(ErrImproperlyFormattedLine)
		}
		for j := 0; j < numStacks; j++ {
			r := runes[j*4+1]
			if r != ' ' {
				stacks[j] = append(stacks[j], r)
			}
		}
	}

	// Iterate through the command lines to mutate the 2D slice data
	for scanner.Scan() {
		line := scanner.Text()

		// Validate max line length
		if len(line) > MaxByesPerLine {
			log.Fatal(ErrOversizedLine)
		}
		matches := commandRgx.FindStringSubmatch(line)
		if len(matches) != 4 {
			log.Fatal(ErrImproperlyFormattedLine)
		}
		howMany, err := strconv.Atoi(matches[1])
		if err != nil {
			log.Fatal(ErrImproperlyFormattedLine)
		}
		whence, err := strconv.Atoi(matches[2])
		if err != nil {
			log.Fatal(ErrImproperlyFormattedLine)
		}
		whither, err := strconv.Atoi(matches[3])
		if err != nil {
			log.Fatal(ErrImproperlyFormattedLine)
		}
		if !partB {
			moveCratesIndiv(stacks, howMany, whence-1, whither-1)
		} else {
			moveCratesInBulk(stacks, howMany, whence-1, whither-1)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// Iterate through finalized 2D slice data to print out the last item in each sub-slice
	for i := 0; i < numStacks; i++ {
		fmt.Printf("%v", string(stacks[i][len(stacks[i])-1]))
	}
	fmt.Println("\n-----------------")

}

func moveCratesIndiv(stacks [][]rune, howMany, whence, whither int) {
	for i := 0; i < howMany; i++ {
		lastIdx := len(stacks[whence]) - 1
		stacks[whither] = append(stacks[whither], stacks[whence][lastIdx])
		stacks[whence] = stacks[whence][:lastIdx]
	}
}
func moveCratesInBulk(stacks [][]rune, howMany, whence, whither int) {
	whenceLen := len(stacks[whence])
	stacks[whither] = append(stacks[whither], stacks[whence][whenceLen-howMany:whenceLen]...)
	stacks[whence] = stacks[whence][:whenceLen-howMany]
}
