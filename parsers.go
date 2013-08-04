package main

import (
	"errors"
	"fmt"
	"github.com/moovweb/gokogiri"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	shortForm        = "02.01.2006"
	searchTableXPath = `
        descendant-or-self::
        div[contains(concat(' ', normalize-space(@class), ' '), ' companies ')]
        /
        descendant::table
        /
        descendant::tr
        /
        descendant::td`
	tableXPath = `
        descendant-or-self::
        div[contains(concat(' ', normalize-space(@class), ' '), ' company ')]
        /
        div[contains(concat(' ', normalize-space(@class), ' '), ' boxbody ')]
        /
        table
        `
)

var (
	deregRegex         = regexp.MustCompile("(i?)Félag afskráð")
	notInBusinessRegex = regexp.MustCompile("(i?)Rekstri hætt")
	ehfRegex           = regexp.MustCompile("(i?)ehf")
)

func ParseDetails(htmlContent []byte, c *Company) (err error) {

	content := ""

	doc, err := gokogiri.ParseHtml(htmlContent)
	defer doc.Free()

	if err != nil {
		return
	}

	tables, err := doc.Search(tableXPath)
	if err != nil {
		return
	}

	if len(tables) != 4 {
		err = errors.New("Should be 4")
		return
	}

	for idx, table := range tables {

		results, _ := table.Search("tbody/tr/td")

		switch idx {

		case 0: //Main info

			for ridx, td := range results {
				content = td.Content()
				switch ridx {
				case 0:
					(*c).PostAddress, err = ParseAddress(content)
				case 1:
					(*c).LegalAddress, err = ParseAddress(content)
				case 3:
					(*c).Type = content
				}
			}

		case 1: //VATNumbers

			vnr := VATNumber{}

			for ridx, td := range results {

				content = td.Content()
				if ridx > 0 && ridx%4 == 0 {
					(*c).VATNumbers = append((*c).VATNumbers, vnr)
					vnr = VATNumber{}
				}
				switch ridx {
				case 0:
					vnr.ID, _ = strconv.Atoi(strings.Trim(content, " "))
				case 1:
					vnr.DateOpened, _ = time.Parse(shortForm, content)
				case 2:
					vnr.DateClosed, _ = time.Parse(shortForm, content)
				case 3:
					vnr.ISATTypes, _ = ParseISATTypes(content)
				}
			}
			if vnr.ID > 0 {
				(*c).VATNumbers = append((*c).VATNumbers, vnr)
			}

		}
	}

	address := c.GuessDomain()
	res, xerr := net.LookupHost(address)
	if xerr == nil {
		(*c).Domain = address
		(*c).IPS = res
	}

	return
}

func FetchDetails(c *Company) (err error) {

	content, err := ReadOrGetSSID(c.Ssid)
	if err != nil {
		return
	}
	err = ParseDetails(content, c)

	return
}

func ParseSearchResults(htmlContent []byte) (companies []Company, err error) {

	doc, err := gokogiri.ParseHtml(htmlContent)
	defer doc.Free()

	if err != nil {
		return
	}

	results, err := doc.Search(searchTableXPath)
	if err != nil {
		return
	}

	if len(results) == 0 {
		return
	}

	company := Company{Type: "Unknown"}
	content := ""
	nIdx := 0

	for idx, res := range results {
		if idx != 0 && idx%3 == 0 {
			nIdx = 0
			companies = append(companies, company)
			company = Company{}
		}
		content = res.Content()
		switch nIdx {
		case 0:
			company.Ssid = content
		case 1:
			company.Name = strings.Split(content, "\n")[0]
			if deregRegex.MatchString(content) {
				company.State = Deregistered
			} else if notInBusinessRegex.MatchString(content) {
				company.State = NotInBusiness
			}

			if ehfRegex.MatchString(content) {
				company.Type = "E1 Einkahlutafélag (ehf)"
				//company.Type = "D1 Hlutafélag, almennt (hf)"
			}

		case 2:
			company.PostAddress, _ = ParseAddress(content)
			company.LegalAddress, _ = ParseAddress(content)
		}
		nIdx++
	}

	if company.Name != "" {
		companies = append(companies, company)
	}

	for idx, c := range companies {
		if !c.ShouldGetDetails() {
			continue
		}
		FetchDetails(&c)

		companies[idx] = c
	}

	return companies, nil
}

func ScrapeStreet(street string, cc chan Company) {

	cb := func(it []Company) int {
		total := 0
		for _, company := range it {
			cc <- company
			total++
		}
		return total
	}

	content, err := ReadOrGetSearch(street)
	if err != nil {
		fmt.Println(err)
		return
	}
	c, err := ParseSearchResults(content)
	if err != nil {
		fmt.Println(err)
	}

	count := cb(c)

	if count >= 499 {

		for i := 1; i < 5; i++ {

			content, err = ReadOrGetSearch(fmt.Sprintf("%s %d", street, i))
			if err != nil {
				fmt.Println(err)
				return
			}
			c, err = ParseSearchResults(content)
			if err != nil {
				fmt.Println(err)
			}
			cb(c)
		}

	}
}
