package nomain

import (
	"log"
	"os"
)

func main() {
	run()
	other()

	log.Fatal("test") // want "Using log.Fatal function outside of main function of main package is discouraged"
	os.Exit(1)        // want "Using os.Exit function outside of main function of main package is discouraged"
}

func run() {
	log.Fatal("test") // want "Using log.Fatal function outside of main function of main package is discouraged"
	os.Exit(1)        // want "Using os.Exit function outside of main function of main package is discouraged"
}
