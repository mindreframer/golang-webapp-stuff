package main

import (
	"fmt"
	"testing"
)

func TestParseAddress(t *testing.T) {

	expected := Address{
		Street:      "Sætúni",
		HouseNumber: 10,
		Postcode:    105,
		Place:       "Reykjavík",
	}

	a, err := ParseAddress("Sætúni 10, 105 Reykjavík")

	if err != nil {
		t.Error("Parsing has error:", err)
		return
	}

	if a != expected {
		t.Errorf("ParseAddress: %v, expected: %v.", a, expected)
	}

	expected = Address{
		Street:      "Litla-Fjarðarhorn",
		HouseNumber: 0,
		Postcode:    510,
		Place:       "Hólmavík",
	}

	a, err = ParseAddress("Litla-Fjarðarhorn  510 Hólmavík")

	if a != expected {
		t.Errorf("ParseAddress: %v, expected: %v.", a, expected)
	}

	testCases := []string{
		"Fornustekkum II  781 Höfn í Hornafirði",
		"Dunhaga 5 Tæknigarði  107 Reykjavík",
		"Skútuvogi 1 b  104 Reykjavík",
	}

	for _, s := range testCases {
		a, err = ParseAddress(s)
		fmt.Printf("%v\n", a)
		if err != nil {
			t.Errorf("Error: %v parsing %s (%v).\n", err, s, a)
			return
		}
	}

}
