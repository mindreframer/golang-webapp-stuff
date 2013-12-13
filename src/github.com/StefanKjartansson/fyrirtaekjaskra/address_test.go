package fyrirtaekjaskra

import (
	"testing"
)

func TestParseAddress(t *testing.T) {

	testCases := make(map[string]Address)
	testCases["Litla-Fjarðarhorn  510 Hólmavík"] = Address{
		Street:   "Litla-Fjarðarhorn",
		Postcode: 510,
		Place:    "Hólmavík",
	}
	testCases["Sætúni 10, 105 Reykjavík"] = Address{
		Street:      "Sætúni",
		HouseNumber: "10",
		Postcode:    105,
		Place:       "Reykjavík",
	}
	testCases["Fornustekkum II  781 Höfn í Hornafirði"] = Address{
		Street:   "Fornustekkum II",
		Postcode: 781,
		Place:    "Höfn í Hornafirði",
	}
	testCases["Dunhaga 5 Tæknigarði  107 Reykjavík"] = Address{
		Street:      "Dunhaga",
		HouseNumber: "5",
		Postcode:    107,
		Place:       "Reykjavík",
	}
	testCases["Skútuvogi 1 b  104 Reykjavík"] = Address{
		Street:      "Skútuvogi",
		HouseNumber: "1b",
		Postcode:    104,
		Place:       "Reykjavík",
	}
	testCases["Domus Medica  Egilsgötu 3  101 Reykjavík"] = Address{
		Street:      "Egilsgötu",
		HouseNumber: "3",
		Postcode:    101,
		Place:       "Reykjavík",
	}
	testCases["Domus Medica, Egilsgötu 3  101 Reykjavík"] = Address{
		Street:      "Egilsgötu",
		HouseNumber: "3",
		Postcode:    101,
		Place:       "Reykjavík",
	}
	testCases["Fluggörðum 30d  101 Reykjavík"] = Address{
		Street:      "Fluggörðum",
		HouseNumber: "30d",
		Postcode:    101,
		Place:       "Reykjavík",
	}
	testCases["Hafnarstræti 91-95, 600 Akureyri"] = Address{
		Street:      "Hafnarstræti",
		HouseNumber: "91-95",
		Postcode:    600,
		Place:       "Akureyri",
	}
	testCases["Hafnarstræti 20 4.hæð, 101 Reykjavík"] = Address{
		Street:      "Hafnarstræti",
		HouseNumber: "20",
		Postcode:    101,
		Place:       "Reykjavík",
	}
	testCases["Austurstræti 17 (6.h), 101 Reykjavík"] = Address{
		Street:      "Austurstræti",
		HouseNumber: "17",
		Postcode:    101,
		Place:       "Reykjavík",
	}
	testCases["Klapparstíg 25.27, 105 Reykjavík"] = Address{
		Street:      "Klapparstíg",
		HouseNumber: "25-27",
		Postcode:    105,
		Place:       "Reykjavík",
	}
	testCases["Hringbraut Landsp., 101 Reykjavík"] = Address{
		Street:   "Hringbraut Landsp.",
		Postcode: 101,
		Place:    "Reykjavík",
	}
	testCases["Laufásvegi  12, 101 Reykjavík"] = Address{
		Street:      "Laufásvegi",
		HouseNumber: "12",
		Postcode:    101,
		Place:       "Reykjavík",
	}
	testCases["Lindargötu Fjármálar., 150 Reykjavík"] = Address{
		Street:   "Lindargötu Fjármálar.",
		Postcode: 150,
		Place:    "Reykjavík",
	}
	testCases["Kirkjustræti Austurv., 101 Reykjavík"] = Address{
		Street:   "Kirkjustræti Austurv.",
		Postcode: 101,
		Place:    "Reykjavík",
	}

	for s, exp := range testCases {
		a, err := ParseAddress(s)
		if a != exp {
			t.Errorf("ParseAddress: %v, expected: %v.", a, exp)
		}
		if err != nil {
			t.Errorf("Error: %v parsing %s (%v).\n", err, s, a)
			return
		}
	}

}
