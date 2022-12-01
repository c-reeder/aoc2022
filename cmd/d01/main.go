package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
)

type ElfCount struct {
	Name         string
	CalorieCount int
}

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

	// Initialize slice
	elfCounts := []ElfCount{
		{Name: "Elf 1"},
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if len(line) == 0 {
			// Add next Elf to slice
			elfCounts = append(elfCounts, ElfCount{
				Name: fmt.Sprintf("Elf %v", len(elfCounts)+1),
			})

		} else {
			// Add current line value to currently incrementing elf
			lineVal, err := strconv.Atoi(scanner.Text())
			if err != nil {
				log.Fatal(err)
			}
			elfCounts[len(elfCounts)-1].CalorieCount += lineVal
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// Sort elves descending by calorie count
	sort.Slice(elfCounts, func(i, j int) bool {
		return elfCounts[i].CalorieCount > elfCounts[j].CalorieCount
	})

	fmt.Println("-----------------")
	fmt.Printf("- %v has the most calories with %v\n", elfCounts[0].Name, elfCounts[0].CalorieCount)

	// Part two
	if len(elfCounts) < 3 {
		return
	}
	var topThree int
	for i := 0; i < 3; i++ {
		topThree += elfCounts[i].CalorieCount
	}

	fmt.Printf("- The top 3 elves %s, %s, and %s have %v calories\n", elfCounts[0].Name, elfCounts[1].Name, elfCounts[2].Name, topThree)
}
