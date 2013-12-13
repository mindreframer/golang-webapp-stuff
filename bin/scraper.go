package main

import (
	"encoding/csv"
	"github.com/StefanKjartansson/fyrirtaekjaskra"
	iconv "github.com/djimenez/iconv-go"
	"github.com/howbazaar/loggo"
	"io"
	"log"
	"os"
	"time"
)

func ImportStreets(filename string) (s []string, err error) {

	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	idx := 0
	x, _ := iconv.NewReader(file, "iso-8859-1", "utf-8")
	reader := csv.NewReader(x)
	reader.Comma = ';'
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if idx != 0 {
			s = append(s, record[3])
		}
		idx++
	}
	return
}

func main() {

	streets, err := ImportStreets("./gotuskra.csv")
	if err != nil {
		log.Fatal(err)
	}

	loggo.GetLogger("fyrirtaekjaskra").SetLogLevel(loggo.DEBUG)

	scraper := fyrirtaekjaskra.NewScraper()
	scraper.ScrapeList(streets[0:120])

	cnt := 0

L:
	for {
		select {
		case _ = <-scraper.CompanyChan:
			//log.Printf("%+v\n", c)
			cnt++
		case err := <-scraper.ErrChan:
			log.Fatal(err)
		case <-time.After(10 * time.Second):
			break L
		}
	}

	log.Printf("Scraped %d companies", cnt)

}
