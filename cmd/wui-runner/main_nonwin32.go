//go:build !windows

package main

import "log"

func main() {
	log.Println("Does not work on non-windows build")
}
