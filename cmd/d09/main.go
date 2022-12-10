package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
)

var (
	ErrImproperlyFormattedLine = errors.New("improperly formatted line")
)

var cmdRgx = regexp.MustCompile(`^([LRUD]) ([0-9]+)$`)

func main() {
	// Check args
	if len(os.Args[1:]) != 1 {
		log.Fatal("Expected 1 argument containing file name!")
	}

	// Open file
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	ropeGrid := NewRopeGrid(10)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		matches := cmdRgx.FindStringSubmatch(line)
		if len(matches) != 3 {
			log.Fatalln(ErrImproperlyFormattedLine)
		}
		cmd := RopeCommand(matches[1])
		arg, err := strconv.Atoi(matches[2])
		if err != nil {
			log.Fatalln(ErrImproperlyFormattedLine)
		}
		ropeGrid.ProcessCommand(cmd, arg)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("The tail has visited %v positions!\n", ropeGrid.TailVisitCount())
}

type RopeCommand string

const (
	Up    RopeCommand = "U"
	Down  RopeCommand = "D"
	Left  RopeCommand = "L"
	Right RopeCommand = "R"
)

type KnotPos struct {
	X int64
	Y int64
}

type RopeGrid struct {
	TailPositions map[KnotPos]bool
	knots         []KnotPos
}

func NewRopeGrid(ropeSize int) RopeGrid {
	return RopeGrid{
		TailPositions: make(map[KnotPos]bool),
		knots:         make([]KnotPos, ropeSize),
	}
}

func (g *RopeGrid) TailVisitCount() int {
	var count int
	for range g.TailPositions {
		count++
	}
	return count
}

func (g *RopeGrid) ProcessCommand(cmd RopeCommand, cmdArg int) {
	for i := 0; i < cmdArg; i++ {
		// Move the Head
		switch cmd {
		case Up:
			g.knots[0].Y++
		case Down:
			g.knots[0].Y--
		case Left:
			g.knots[0].X--
		case Right:
			g.knots[0].X++
		}

		// Iterate through subsequent knots
		for j := 1; j < len(g.knots); j++ {
			dx := g.knots[j-1].X - g.knots[j].X
			dy := g.knots[j-1].Y - g.knots[j].Y

			var (
				rX int64
				rY int64
			)

			if dx != 0 {
				rX = dx / int64(math.Abs(float64(dx)))
			}
			if dy != 0 {
				rY = dy / int64(math.Abs(float64(dy)))
			}

			if int64(math.Sqrt(math.Pow(float64(dx), 2)+math.Pow(float64(dy), 2))) > 1 {
				g.knots[j].X += rX
				g.knots[j].Y += rY
			}
		}

		g.TailPositions[g.knots[len(g.knots)-1]] = true
	}
}
