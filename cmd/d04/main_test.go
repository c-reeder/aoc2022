package main

import "testing"

func TestPassingLines(t *testing.T) {
	goodLines := []string{
		"2-4,6-8",
		"2-3,4-5",
		"5-7,7-9",
		"2-8,3-7",
		"6-6,4-6",
		"2-6,4-8",
	}
	goodAssigns := [][]SectionAssignment{
		{SectionAssignment{2, 4}, SectionAssignment{6, 8}},
		{SectionAssignment{2, 3}, SectionAssignment{4, 5}},
		{SectionAssignment{5, 7}, SectionAssignment{7, 9}},
		{SectionAssignment{2, 8}, SectionAssignment{3, 7}},
		{SectionAssignment{6, 6}, SectionAssignment{4, 6}},
		{SectionAssignment{2, 6}, SectionAssignment{4, 8}},
	}

	for i, line := range goodLines {
		assigns, err := getRangesFromLine(line)
		if err != nil {
			t.Errorf("Validation should have passed for line: %v", line)
		}
		if !checkSlicesEqual(goodAssigns[i], assigns) {
			t.Errorf("Assignments for %v should have been %v but were %v\n", line, goodAssigns[i], assigns)
		}
	}
}

func TestFailingLines(t *testing.T) {
	evilLines := []string{
		"24,6-8",
		"2-3,45",
		"5-77-9",
		"2-,3-7",
		"6-6,-6",
		"2-6,4-",
		"",
		"12",
		",",
		"9-2, 10-12",
		"5-34,58-2",
	}

	for _, line := range evilLines {
		if _, err := getRangesFromLine(line); err == nil {
			t.Errorf("Validation should not have passed for line: %v", line)
		}
	}
}

func TestFullyOverlapping(t *testing.T) {
	assigns := [][2]SectionAssignment{
		{{1, 8}, {4, 6}},
		{{56789, 78904}, {60345, 71045}},
		{{12, 28}, {20, 28}},
		{{56, 68}, {56, 59}},
		{{5002, 8002}, {3008, 9104}},
	}
	for _, a := range assigns {
		if determineOverlap(a[0], a[1]) != Containing {
			t.Errorf("%v and %v should overlap, but received false", a[0], a[1])
		}
	}
}

func TestNotFullyOverlapping(t *testing.T) {
	assigns := [][2]SectionAssignment{
		{{21, 38}, {38, 40}},
		{{101, 301}, {201, 401}},
		{{255, 400}, {100, 280}},
		{{255, 400}, {100, 255}},
		{{555, 666}, {777, 888}},
	}
	for _, a := range assigns {
		if determineOverlap(a[0], a[1]) == Containing {
			t.Errorf("%v and %v do not overlap, but received true", a[0], a[1])
		}
	}
}

// checkSlicesEqual does a deep comparison on two SectionAssignment slices
func checkSlicesEqual(a []SectionAssignment, b []SectionAssignment) bool {
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
