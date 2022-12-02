package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

var ErrImproperlyFormattedLine = errors.New("Improperly formatted line!")

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

	// Score for all rounds
	var totalScore int

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Validate, capitalize, and split line
		if len(line) != 3 {
			log.Fatal(ErrImproperlyFormattedLine)
		}
		line = strings.ToUpper(line)
		tokens := strings.Split(line, " ")
		if len(tokens) != 2 {
			log.Fatal(ErrImproperlyFormattedLine)
		}

		// Add line score to running total
		lineScore, err := determineLineScore(tokens[0], tokens[1])
		if err != nil {
			log.Fatal(err)
		}
		totalScore += lineScore

	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("-----------------")
	fmt.Printf("Total score is: %v\n", totalScore)

}

// determineLineScore calculates the score for a single line
// taking in strings for each of the two values on that line
func determineLineScore(first, sec string) (int, error) {
	firstVal, err := runeOffsetFromBase([]rune(first)[0], '@')
	if err != nil {
		return 0, err
	}
	secVal, err := runeOffsetFromBase([]rune(sec)[0], 'W')
	if err != nil {
		return 0, err
	}

	var outcomePoints int
	var selectionPoints int

	if !partB {
		// Part A
		// firstVal is opponent's selection
		// secVal is our selection
		selectionPoints = secVal
		switch secVal - firstVal {
		case 0:
			outcomePoints = 3
		case 2, -1:
			outcomePoints = 0
		default:
			outcomePoints = 6
		}

	} else {
		// Part B
		// firstVal is opponent's selection
		// secVal is the outcome
		outcomePoints = 3 * (secVal - 1)
		switch secVal {
		case 1:
			selectionPoints = ((firstVal + 1) % 3) + 1
		case 2:
			selectionPoints = firstVal
		default:
			selectionPoints = (firstVal % 3) + 1
		}
	}

	//fmt.Printf("%v %v : %v %v : %v %v\n",
	//	first, sec, firstVal, secVal, outcomePoints, outcomePoints+selectionPoints)

	return outcomePoints + selectionPoints, nil
}

// runeOffsetFromBase determines the offset from a given base
// rune to the value provided and ensures that the values
// are restricted to 1, 2, and 3
func runeOffsetFromBase(val rune, base rune) (int, error) {
	diff := val - base
	if diff < 0 || diff > 3 {
		return 0, ErrImproperlyFormattedLine
	}
	return int(diff), nil
}
