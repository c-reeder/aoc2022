package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
)

var ErrOversizedLine = errors.New("line exceeds maximum length")
var ErrInvalidItem = errors.New("invalid item detected")
var ErrNoBadgeFound = errors.New("no badge was found for group")

const MaxByesPerLine = 3 * 1024 // 3kB max line length
const SectionsPerLine = 2

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

	// Sum of all the priorities repeated between sections in a line (for part A)
	// or sum of all badge priorities (for part B)
	var sum int

	var groupLines [3]string
	var lineNum uint64

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Validate max line length
		if len(line) > MaxByesPerLine {
			log.Fatal(ErrOversizedLine)
		}

		if !partB {
			// Part A
			// Add line score to running total
			repeats, err := getPrioritiesRepeatedBetweenSections(line, 2)
			if err != nil {
				log.Fatal(err)
			}

			// Add priorities of repeats to running sum
			for _, priority := range repeats {
				sum += priority
			}

		} else {
			// Part B

			mod := lineNum % 3
			groupLines[mod] = line
			lineNum++

			if mod == 2 {
				groupBadgePriority, err := findBadgePriorityForGroup(groupLines)
				if err != nil {
					log.Fatal(err)
				}
				sum += groupBadgePriority
			}
		}

	}
	fmt.Println("-----------------")
	fmt.Printf("Sum is: %v\n", sum)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

}

// findBadgePriorityForGroup determines the item in common amongst
// a group of three and returns the priority for it
func findBadgePriorityForGroup(groupLines [3]string) (int, error) {
	// maps priority to the # of lines in the group which contain at least one
	groupMap := make(map[int]int)
	for i := range groupLines {
		// map of all the items we've already seen in this line
		lineMap := make(map[int]bool)
		for _, r := range []rune(groupLines[i]) {
			p, err := runeToPriority(r)
			if err != nil {
				return 0, err
			}
			if !lineMap[p] {
				lineMap[p] = true
				groupMap[p]++
				// If we've seen this item in more than 2 lines
				if groupMap[p] > 2 {
					return p, nil
				}
			}
		}
	}
	return 0, ErrNoBadgeFound
}

// getPrioritiesRepeatedBetweenHalves breaks the string in X sections
// and returns a slice of integers representing the priorities of
// any runes that appear in more than one section
func getPrioritiesRepeatedBetweenSections(line string, sections int) (priorities []int, err error) {
	// priorsToSecs maps priorities to section indices
	priorsToSecs := make(map[int]int)

	// repeatedMap is where we mark that an item has been
	// repeated. This is to prevent duplicates in the
	// priorities list returned
	repeatedMap := make(map[int]bool)

	// currSec is the index of the current section
	var currSec int

	lineRunes := []rune(line)

	// secLen is the length of a single section
	// (rounded up in case the line doesn't evenly divide)
	secLen := int(math.Ceil(float64(len(lineRunes)) / float64(sections)))

	// Loop over runes in the entire line
	for i, r := range lineRunes {
		// Signal that we've started the next section
		if i%secLen == 0 {
			currSec++
		}
		// Validate rune and convert to priority
		p, err := runeToPriority(r)
		if err != nil {
			return nil, err
		}
		// Check if this rune was already in another section
		// If so, then add it to the slice to return
		if s, ok := priorsToSecs[p]; ok &&
			s != currSec && !repeatedMap[p] {
			priorities = append(priorities, p)
			repeatedMap[p] = true
		}
		// Mark that this rune is in the current section
		priorsToSecs[p] = currSec
	}
	return priorities, nil
}

// runeToPriority validates a rune and converts it to a priority
func runeToPriority(r rune) (int, error) {
	if r > '@' && r < '[' {
		return int(r - 'A' + 27), nil
	}
	if r > '`' && r < '{' {
		return int(r - 'a' + 1), nil
	}
	return 0, ErrInvalidItem
}
