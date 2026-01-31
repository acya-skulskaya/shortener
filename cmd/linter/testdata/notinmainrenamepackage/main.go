package main

import (
	"log"
	"os"
)

func main() {
	other()

	log.Fatal("test")
	os.Exit(1)
}
