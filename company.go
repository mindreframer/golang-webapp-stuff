package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type CompanyState int

const (
	Active CompanyState = iota
	Deregistered
	NotInBusiness
)

var (
	dismiss = []string{"húsfélag", "sjóður", "félag", "flokkur",
		"samband", "starfsmannafélg", "skóli", "kirkj", "samtök"}
	suffixes = []string{"ehf", "hf", "sf", "slf", "ohf"}
)

type Address struct {
	Street      string `json:"street"`
	HouseNumber int    `json:"number"`
	Postcode    int    `json:"postcode"`
	Place       string `json:"place"`
}

type ISATType struct {
	Number      int    `json:"number"`
	Description string `json:"description"`
	Main        bool   `json:"is_main"`
}

type VATNumber struct {
	ID         int        `json:"id"`
	DateOpened time.Time  `json:"date_opened"`
	DateClosed time.Time  `json:"date_closed,omitempty"`
	ISATTypes  []ISATType `json:"isat_types"`
}

type Company struct {
	Ssid         string       `json:"ssid"`
	Name         string       `json:"name"`
	Domain       string       `json:"domain,omitempty"`
	IPS          []string     `json:"ip_addresses,omitempty"`
	PostAddress  Address      `json:"post_address,omitempty"`
	LegalAddress Address      `json:"legal_address,omitempty"`
	Type         string       `json:"company_type"`
	VATNumbers   []VATNumber  `json:"vat_numbers,omitempty"`
	State        CompanyState `json:"company_state"`
}

// ShouldGetDetails determines whether a company is interesting
// enough to fetch it's details page
func (c Company) ShouldGetDetails() bool {
	if c.State != Active {
		return false
	}

	// Filter out individuals' ssids
	x, _ := strconv.Atoi(string(c.Ssid[0]))
	if x < 4 {
		return false
	}

	// We don't care about anything that matches the dismissed strings
	n := strings.ToLower(c.Name)
	for _, d := range dismiss {
		if strings.Contains(n, d) {
			return false
		}
	}
	return true
}

// GuessDomain returns a string containing a guess of which domain
// the company has based on it's name
func (c Company) GuessDomain() string {

	n := Asciify(c.Name)

    // Trim the suffixes
	for _, s := range suffixes {
		n = strings.TrimSuffix(n, s)
		n = strings.TrimSuffix(n, s + ".")
	}

    // Trim spaces and illegal characters.
	n = strings.TrimSpace(n)
	for _, s := range []string{" ", ",", "."} {
		n = strings.Replace(n, s, "", -1)
	}

	return fmt.Sprintf("%s.is", n)
}
