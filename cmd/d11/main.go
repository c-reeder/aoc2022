package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var (
	ErrImproperlyFormattedLine        = errors.New("improperly formatted line")
	ErrUnrecognizedArithmeticOperator = errors.New("unrecognized arithmetic operator")
	ErrUnrecognizedMonkeyID           = errors.New("unrecognized monkey id")
)

var paraRgx = regexp.MustCompile(
	`Monkey (\d+):
  Starting items: ([\d, ]+)
  Operation: new = (?:old|\d+) ([\*-\+\/]) (old|\d+)
  Test: divisible by (\d+)
    If true: throw to monkey (\d+)
    If false: throw to monkey (\d+)\n?`)

func main() {
	// Check args
	if len(os.Args) != 2 {
		log.Fatal("Expected 1 argument containing file name!")
	}

	monkeyList, monkeyMap := produceMonkeys(os.Args[1])

	numRounds := 10_000

	for round := 0; round < numRounds; round++ {
		for m := range monkeyList {
			// iterate through all the items the monkey has
			for _, variantList := range monkeyList[m].ItemVariantLists {
				monkeyList[m].TotalInspections++

				// Iterate through the variants for each item in the monkeys posession
				for v := range variantList {

					var arg int64
					if monkeyList[m].Operation.Arg != nil {
						arg = *monkeyList[m].Operation.Arg
					} else {
						arg = variantList[v]
					}
					switch monkeyList[m].Operation.Op {
					case Plus:
						variantList[v] = variantList[v] + arg
					case Minus:
						variantList[v] = variantList[v] - arg
					case Multiply:
						variantList[v] = variantList[v] * arg
					case Divide:
						variantList[v] = variantList[v] / arg
					default:
						log.Fatal(ErrUnrecognizedArithmeticOperator)
					}

					variantList[v] = variantList[v] % monkeyList[v].TestDivisor
				}

				// Monkey test's worry level
				var recipientMonkey *Monkey
				var ok bool
				if variantList[m] == 0 {
					recipientMonkey, ok = monkeyMap[monkeyList[m].TrueRecipient]
				} else {
					recipientMonkey, ok = monkeyMap[monkeyList[m].FalseRecipient]
				}
				if !ok {
					log.Fatal(ErrUnrecognizedMonkeyID)
				}

				// Monkey transfers item to another monkey
				recipientMonkey.ItemVariantLists = append(recipientMonkey.ItemVariantLists, variantList)
			}
			// Flush out all items since they've all been handed off to other Monkeys
			monkeyList[m].ItemVariantLists = nil
		}
	}
	sort.Slice(monkeyList, func(i, j int) bool {
		return monkeyList[i].TotalInspections > monkeyList[j].TotalInspections
	})
	fmt.Printf("Monkey Business is: %v\n", monkeyList[0].TotalInspections*monkeyList[1].TotalInspections)
}

func produceMonkeys(filename string) ([]*Monkey, map[uint]*Monkey) {
	bs, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	matches := paraRgx.FindAllSubmatch(bs, -1)
	monkeyCount := len(matches)
	monkeyList := make([]*Monkey, monkeyCount)
	monkeyMap := make(map[uint]*Monkey)
	for i, match := range matches {
		id, err := strconv.ParseUint(string(match[1]), 10, 64)
		if err != nil {
			log.Fatal(ErrImproperlyFormattedLine)
		}

		itemStrs := strings.Split(string(match[2]), ", ")
		variantLists := make([][]int64, len(itemStrs))
		for i := range itemStrs {
			item, err := strconv.ParseInt(itemStrs[i], 10, 64)
			if err != nil {
				log.Fatal(ErrImproperlyFormattedLine)
			}
			variantLists[i] = make([]int64, monkeyCount)
			for j := range variantLists[i] {
				variantLists[i][j] = item
			}
		}

		var arg *int64
		argStr := string(match[4])
		if argStr != "old" {
			argInt, err := strconv.ParseInt(argStr, 10, 64)
			if err != nil {
				log.Fatal(ErrImproperlyFormattedLine)
			}
			arg = &argInt
		}
		var divisor int64
		divisor, err = strconv.ParseInt(string(match[5]), 10, 64)
		trueRecipient, err := strconv.ParseUint(string(match[6]), 10, 64)
		if err != nil {
			log.Fatal(ErrImproperlyFormattedLine)
		}
		falseRecipient, err := strconv.ParseUint(string(match[7]), 10, 64)
		if err != nil {
			log.Fatal(ErrImproperlyFormattedLine)
		}
		m := &Monkey{
			ID:               uint(id),
			ItemVariantLists: variantLists,
			Operation:        MonkeyOperation{Op: Operator(match[3]), Arg: arg},
			TestDivisor:      divisor,
			TrueRecipient:    uint(trueRecipient),
			FalseRecipient:   uint(falseRecipient),
			TotalInspections: 0,
		}
		monkeyList[i] = m
		monkeyMap[m.ID] = m
	}
	return monkeyList, monkeyMap
}

type Operator string

const (
	Plus     Operator = "+"
	Minus    Operator = "-"
	Divide   Operator = "/"
	Multiply Operator = "*"
)

type MonkeyOperation struct {
	Op  Operator
	Arg *int64
}

type Monkey struct {
	ID               uint
	ItemVariantLists [][]int64
	Operation        MonkeyOperation
	TestDivisor      int64
	TrueRecipient    uint
	FalseRecipient   uint
	TotalInspections uint
}
