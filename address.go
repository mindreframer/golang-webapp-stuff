package main

import (
	"errors"
	"strconv"
	"strings"
)

// ParseAddress parses a string and returns an address
func ParseAddress(s string) (a Address, err error) {

	parts := []string{}
	if strings.Contains(s, ",") {
		parts = strings.Split(s, ",")
	} else {
		parts = strings.Split(s, "  ")
	}

	if len(parts) != 2 {
		err = errors.New("Parts split expected to have length of 2.")
		return
	}

	for idx, p := range parts {
		x := strings.Split(strings.Trim(p, " "), " ")
		switch idx {
		case 0:
			for sidx, i := range x {
				if i == "" {
					continue
				}
				switch sidx {
				case 0:
					a.Street = x[0]
				case 1:
					a.HouseNumber, err = strconv.Atoi(x[1])
				}
			}
		case 1:
			for sidx, i := range x {
				if i == "" {
					continue
				}
				switch sidx {
				case 0:
					a.Postcode, err = strconv.Atoi(x[0])
				case 1:
					a.Place = x[1]
				}
			}
		}
	}
	return
}
