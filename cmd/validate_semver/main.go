package main

import (
	"log"

	"golang.org/x/mod/semver"
)

var version = ""

func main() {
	if !semver.IsValid(version) {
		log.Fatalln("semver invalid")
	}
}
