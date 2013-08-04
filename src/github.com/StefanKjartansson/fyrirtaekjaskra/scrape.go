package main

import (
	"fmt"
	"labix.org/v2/mgo"
	"log"
	"strings"
)

var (
	BufferSize = 32
	cc         = make(chan Company, BufferSize)
)

func main() {

	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	c := session.DB("fyrirtaekjaskra").C("companies")

	streets, err := ImportStreets("./gotuskra.csv")
	if err != nil {
		fmt.Println(err)
	}

	go func() {
		for _, s := range streets {
			if len(strings.Split(s, " ")) == 1 {
				log.Println(s)
				ScrapeStreet(s, cc)
			}
		}
	}()

	go func() {
		for {
			select {
			case ev := <-cc:
				err = c.Insert(ev)
				if err != nil {
					panic(err)
				}
			}
		}
	}()

	select {}
}
