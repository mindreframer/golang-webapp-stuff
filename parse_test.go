package fyrirtaekjaskra

import (
	"fmt"
	"io"
	"os"
	"testing"
	"time"
)

func read(filename string) io.Reader {

	fi, err := os.Open(filename)
	if err != nil {
		panic("unable to read a file")
	}
	return fi
}

func TestParseDetails(t *testing.T) {

	ssid := "5902697199"

	c := Company{
		Ssid: ssid,
	}

	scraper := NewScraper()
	err := scraper.ParseDetails(read(fmt.Sprintf("./test/fskra-%s.html", ssid)), &c)

	if err != nil {
		t.Error("Parsing has error: %s", err.Error())
		return
	}

	if c.VATNumbers[0].ID != 10487 {
		t.Errorf("vatnumber is weird: %v\n", c.VATNumbers)
	}

}

func TestXpathSearchTable(t *testing.T) {

	scraper := NewScraper()
	go scraper.ParseSearchResults(read("./test/fskra-leit.html"))

	c := <-scraper.CompanyChan
	t.Logf("%+v", c)
	if c.Ssid != "5407051000" ||
		c.Name != "A Einn ehf" {
		t.Errorf("Parsing error, %+v", c)
		return
	}

	select {
	case x := <-scraper.CompanyChan:
		t.Logf("%+v", x)
	case <-time.After(2 * time.Second):
		break
	}
}
