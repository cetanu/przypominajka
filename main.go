package main

import (
	"log"
)

func main() {
	log.SetFlags(0)
	bdays, err := readBirthdays()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(bdays)
}
