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
)

type OverlapType int

const (
	None OverlapType = iota
	Partial
	Containing
)

type SectionAssignment struct {
	Start uint64
	End   uint64
}

var lineRegex = regexp.MustCompile(`^([0-9]+)-([0-9]+),([0-9]+)-([0-9]+)$`)
var ErrOversizedLine = errors.New("line exceeds maximum length")
var ErrImproperlyFormattedLine = errors.New("improperly formatted line")
var ErrInvalidRanges = errors.New("end of assignment range was before start")

const MaxByesPerLine = 3 * 1024 // 3kB max line length

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

	// Total count of overlapping assignments
	var total int

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Validate max line length
		if len(line) > MaxByesPerLine {
			log.Fatal(ErrOversizedLine)
		}

		// Extract assignments from line
		assigns, err := getRangesFromLine(line)
		if err != nil || len(assigns) != 2 {
			log.Fatal(err)
		}

		overlapType := determineOverlap(assigns[0], assigns[1])

		if !partB {
			// Part A
			if overlapType == Containing {
				total++
			}
		} else {
			// Part B
			if overlapType == Containing || overlapType == Partial {
				total++
			}
		}

	}
	fmt.Println("-----------------")
	fmt.Printf("Total is: %v\n", total)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

}

// getRangesFromLine validates a line and extracts the two section assignments from it
func getRangesFromLine(line string) (assigns []SectionAssignment, err error) {
	matches := lineRegex.FindStringSubmatch(line)
	if len(matches) != 5 {
		return nil, ErrImproperlyFormattedLine
	}

	assigns = make([]SectionAssignment, 2)

	assigns[0].Start, err = strconv.ParseUint(matches[1], 10, 64)
	if err != nil {
		return nil, ErrImproperlyFormattedLine
	}
	assigns[0].End, err = strconv.ParseUint(matches[2], 10, 64)
	if err != nil {
		return nil, ErrImproperlyFormattedLine
	}
	assigns[1].Start, err = strconv.ParseUint(matches[3], 10, 64)
	if err != nil {
		return nil, ErrImproperlyFormattedLine
	}
	assigns[1].End, err = strconv.ParseUint(matches[4], 10, 64)
	if err != nil {
		return nil, ErrImproperlyFormattedLine
	}

	if assigns[0].End < assigns[0].Start {
		return nil, ErrInvalidRanges
	}

	if assigns[1].End < assigns[1].Start {
		return nil, ErrInvalidRanges
	}

	return assigns, nil
}

// determineOverlap determines whether two assignments overlap partially, one contains
// the other, or not at all
func determineOverlap(a, b SectionAssignment) OverlapType {
	if a.Start < b.Start {
		if a.End < b.Start {
			return None
		}
		if a.End < b.End {
			return Partial
		}
		return Containing
	}

	if b.Start < a.Start {
		if b.End < a.Start {
			return None
		}
		if b.End < a.End {
			return Partial
		}
		return Containing
	}

	return Containing
}
