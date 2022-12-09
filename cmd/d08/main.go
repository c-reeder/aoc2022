package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

var treeRgx = regexp.MustCompile(`[0-9]`)

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

	var (
		width  int
		height int
		trees  = make([][]Tree, 0)
	)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		matches := treeRgx.FindAllString(line, -1)
		treeLine := make([]Tree, len(matches))
		for i, match := range matches {
			if i == 0 {
				width = len(matches)
			} else if width != len(matches) {
				log.Fatalln("inconsistent line lengths!")
			}
			height, err := strconv.Atoi(match)
			if err != nil {
				log.Fatal("invalid tree height")
			}
			treeLine[i] = Tree{
				Height:  height,
				Visible: false,
			}
		}
		trees = append(trees, treeLine)
		height++
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	treePatch := TreePatch{
		trees:  trees,
		Height: height,
		Width:  width,
	}
	fmt.Printf("%v is the highest scenic score possible in this grid\n", treePatch.FindHighestScenicScore())
}

func (t *TreePatch) CalculateScenicScore(r, c int) uint64 {
	var (
		left  int
		right int
		above int
		below int

		height     = t.trees[r][c].Height
		localCount int
	)
	for j := c - 1; j >= 0; j-- {
		localCount++
		if t.trees[r][j].Height >= height {
			break
		}
	}
	left = localCount

	localCount = 0
	for j := c + 1; j < t.Width; j++ {
		localCount++
		if t.trees[r][j].Height >= height {
			break
		}
	}
	right = localCount

	localCount = 0
	for i := r - 1; i >= 0; i-- {
		localCount++
		if t.trees[i][c].Height >= height {
			break
		}
	}
	above = localCount

	localCount = 0
	for i := r + 1; i < t.Height; i++ {
		localCount++
		if t.trees[i][c].Height >= height {
			break
		}
	}
	below = localCount

	return uint64(left * right * above * below)
}

func (t *TreePatch) FindHighestScenicScore() uint64 {
	var max uint64
	for r := 0; r < t.Width; r++ {
		for c := 0; c < t.Height; c++ {
			score := t.CalculateScenicScore(r, c)
			if score > max {
				max = score
			}
		}
	}
	return max
}

func (t *TreePatch) CalculateExternallyVisibleTrees() int {
	var visibleCount int
	// Iterate through all rows
	for r := 0; r < t.Width; r++ {
		// Forward through row
		localMax := -1
		for c := 0; c < t.Height; c++ {
			if t.trees[r][c].Height > localMax {
				if !t.trees[r][c].Visible {
					visibleCount++
				}
				t.trees[r][c].Visible = true
				localMax = t.trees[r][c].Height
			}
		}
		// Backward through row
		localMax = -1
		for c := t.Width - 1; c >= 0; c-- {
			if t.trees[r][c].Height > localMax {
				if !t.trees[r][c].Visible {
					visibleCount++
				}
				t.trees[r][c].Visible = true
				localMax = t.trees[r][c].Height
			}
		}
	}

	// Iterate through all columns
	for c := 0; c < t.Height; c++ {
		// Top to bottom
		localMax := -1
		for r := 0; r < t.Width; r++ {
			if t.trees[r][c].Height > localMax {
				if !t.trees[r][c].Visible {
					visibleCount++
				}
				t.trees[r][c].Visible = true
				localMax = t.trees[r][c].Height
			}
		}
		// Bottom to top
		localMax = -1
		for r := t.Width - 1; r >= 0; r-- {
			if t.trees[r][c].Height > localMax {
				if !t.trees[r][c].Visible {
					visibleCount++
				}
				t.trees[r][c].Visible = true
				localMax = t.trees[r][c].Height
			}
		}
	}
	return visibleCount
}

type TreePatch struct {
	trees  [][]Tree
	Height int
	Width  int
}

type Tree struct {
	Height  int
	Visible bool
}
