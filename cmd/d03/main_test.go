package main

import "testing"

func TestPassingRuneToPriority(t *testing.T) {
	const sections = 2
	goodRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	for i, r := range goodRunes {
		p, err := runeToPriority(r)
		if err != nil {
			t.Errorf("Validation should have passed for rune: %v", r)
		}
		if p != i+1 {
			t.Errorf("Priority for %v should have been %v but was %v\n", r, i+1, p)
		}
	}
}

func TestFailingRuneToPrioritytion(t *testing.T) {
	const sections = 2
	evilRunes := []rune(" !\"#$%&'()*+,-./0123456789:;<=>?@[\\]^_`{|}~")

	for _, r := range evilRunes {
		if _, err := runeToPriority(r); err == nil {
			t.Errorf("Validation should not have passed for rune: %v", r)
		}
	}
}

func TestGetRepeatedPrioritiesBetweenSections(t *testing.T) {
	const sections = 2
	lines := []string{
		"vJrwpWtwJgWrhcsFMMfFFhFp",
		"jqHRNqRjqzjGDLGLrsFMfFZSrLrFZsSL",
		"PmmdzqPrVvPwwTWBwg",
		"wMqvLMZHhHMvwLHjbvcjnnSBnvTQFn",
		"ttgJtRGJQctTZtZT",
		"CrZsJsPPZsGzwwsLwLmpwMDw",
	}
	expectedRepeats := [][]int{
		{16},
		{38},
		{42},
		{22},
		{20},
		{19},
	}

	for i := range lines {
		repeats, err := getPrioritiesRepeatedBetweenSections(lines[i], sections)
		if err != nil {
			t.Log(err)
			t.Errorf("Error encountered on line: %v", lines[i])
		}
		if !checkSlicesEqual(expectedRepeats[i], repeats) {
			t.Errorf("Expected %v but received %v for line: %v\n", expectedRepeats[i], repeats, lines[i])
		}

	}
}

// checkSlicesEqual does a deep comparison on two int slices
func checkSlicesEqual(a []int, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestFindBadgePriorityForGroup(t *testing.T) {

	groups := [][3]string{
		{
			"vJrwpWtwJgWrhcsFMMfFFhFp",
			"jqHRNqRjqzjGDLGLrsFMfFZSrLrFZsSL",
			"PmmdzqPrVvPwwTWBwg",
		},
		{
			"wMqvLMZHhHMvwLHjbvcjnnSBnvTQFn",
			"ttgJtRGJQctTZtZT",
			"CrZsJsPPZsGzwwsLwLmpwMDw",
		},
	}

	badgePriorities := []int{
		18,
		52,
	}

	for i := range groups {
		p, err := findBadgePriorityForGroup(groups[i])
		if err != nil {
			t.Log(err)
			t.Errorf("Received error for group: %v\n", groups[i])
		}
		if p != badgePriorities[i] {
			t.Errorf("Expected %v for group %v but got %v\n", badgePriorities[i], groups[i], p)
		}
	}

}

func TestFailingFindBadgePriorityForGroup(t *testing.T) {
	evilGroups := [][3]string{
		{
			"vJwpWtwJgWhcsFMMfFFhFp",
			"jqHRNqRjqzjGDLGLrsFMfFZSrLrFZsSL",
			"PmmdzqPrVvPwwTWBwg",
		},
		{
			"wMqvLMZHhHMvwLHjbvcjnnSBnvTQFn",
			"ttgJtRGJQctTtT",
			"CrZsJsPPZsGzwwsLwLmpwMDw",
		},
		{
			"wMqvLMZHhHMvwLHjbvcjnnSBnvTQFn",
			"ttgJtRGJQctTZtZT",
			"CrsJsPPsGzwwsLwLmpwMDw",
		},
	}

	for i := range evilGroups {
		if _, err := findBadgePriorityForGroup(evilGroups[i]); err == nil {
			t.Errorf("Validation should not have passed for line: %v", evilGroups[i])
		}
	}
}
