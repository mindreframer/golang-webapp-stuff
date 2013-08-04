package main

import (
	"regexp"
	"strconv"
	"strings"
)

var re = regexp.MustCompile("(\\d+)\\s([\\p{Latin}|\\s}]+)")

//ParseISATTypes parses a string and returns a list of ISATType objects.
func ParseISATTypes(s string) (isats []ISATType, err error) {

	//Trim all the things
	s = strings.Trim(strings.Replace(s, "  ", "", -1), " ")

	it := ISATType{}

	for _, i := range strings.Split(s, "\n") {
		i = strings.Trim(i, " ")
		if i == "" {
			continue
		}
		matches := re.FindStringSubmatch(i)
		if len(matches) > 0 {
			it.Number, _ = strconv.Atoi(matches[1])
			it.Description = matches[2]
		} else {
			it.Main = strings.Contains(i, "Aðal")
			isats = append(isats, it)
			it = ISATType{}
		}
	}
	return
}
