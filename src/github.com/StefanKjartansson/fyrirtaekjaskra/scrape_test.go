package main

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func read(filename string) []byte {

	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		panic("unable to read a file")
	}
	return contents
}

func TestParseDetails(t *testing.T) {

	ssid := "5902697199"

	c := Company{
		Ssid: ssid,
	}

	err := ParseDetails(read(fmt.Sprintf("./test/fskra-%s.html", ssid)), &c)

	if err != nil {
		t.Error("Parsing has error:", err)
		return
	}

	if c.VATNumbers[0].ID != 10487 {
		t.Errorf("vatnumber is weird: %v\n", c.VATNumbers)
	}

}

func TestXpathSearchTable(t *testing.T) {

	companies, err := ParseSearchResults(read("./test/fskra-leit.html"))
	if err != nil {
		t.Error("Parsing has error:", err)
		return
	}

	if companies[0].Ssid != "5407051000" ||
		companies[0].Name != "A Einn ehf" {

		t.Error("Parsing error", companies[0])
		return
	}

}
