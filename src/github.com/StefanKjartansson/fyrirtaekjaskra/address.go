package fyrirtaekjaskra

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	reNumber       = regexp.MustCompile(`\d+-?\.?`)
	rePostcode     = regexp.MustCompile(`\d{3}\s+`)
	reRemainder    = regexp.MustCompile(`^[a-zA-Z]{1}$`)
	reStrictNumber = regexp.MustCompile(`^\d+$`)
	excemptionList = []string{
		"Domus",
		"Medica",
	}
)

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// ParseAddress parses a string and returns an address
func ParseAddress(s string) (a Address, err error) {

	pcl := rePostcode.FindStringSubmatchIndex(s)
	if pcl == nil {
		err = fmt.Errorf("No postcode found for: '%s'", s)
		return
	}

	if len(pcl) < 2 {
		err = fmt.Errorf("Postcode location error: '%s', %+v", s, pcl)
		return
	}

	pstart := pcl[0]
	pend := pcl[1]

	a.Postcode, err = strconv.Atoi(strings.TrimSpace(s[pstart:pend]))

	if err != nil {
		return
	}

	// Municipality follows the postcode
	a.Place = strings.TrimSpace(s[pend:])

	// Isolate the address part
	addressPart := strings.Trim(s[:pstart], ", ")

	// Find the house number
	anl := reNumber.FindStringSubmatchIndex(addressPart)

	// No house number, set the street and return
	if len(anl) == 0 {
		a.Street = addressPart
		return
	}

	a.HouseNumber = addressPart[anl[0]:anl[1]]

	// The address part trailing the number is larger than the capturing regex, 
	// this indicates that there's either a house character in the housenumber 
	// or a range of building numbers
	if len(addressPart) > anl[1] {

		remainder := strings.TrimSpace(addressPart[anl[1]:])

		// We only care about trailing house characters and building ranges
		if reRemainder.MatchString(remainder) || reStrictNumber.MatchString(remainder) {
			a.HouseNumber += remainder
		}

		// Some building ranges are delimited by a dot, replace with a dash
		a.HouseNumber = strings.Replace(a.HouseNumber, ".", "-", -1)
	}

	// Street name part, usually there is just a single name but in some cases
	// this part is a place name (not unusual to encounter farm names here).
	for _, s := range strings.Split(strings.TrimSpace(addressPart[:anl[0]]), " ") {

		s = strings.Trim(s, ", ")

		// Ignore empty strings and excempt strings
		// TODO: Expand excemption list to return f.i. the address of mall instead of
		// it's print name.
		if s == "" || stringInSlice(s, excemptionList) {
			continue
		}

		// Add space if there are more than one parts
		if a.Street != "" {
			a.Street += " "
		}
		a.Street += s
	}

	// Trim trailing spaces
	a.Street = strings.TrimSpace(a.Street)

	return
}
