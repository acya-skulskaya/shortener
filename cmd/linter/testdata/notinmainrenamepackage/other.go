package main

import (
	renamedLog "log"
	renamedOS "os"
)

func other() {
	renamedLog.Fatal("test") // want "Using log.Fatal function outside of main function of main package is discouraged"
	renamedOS.Exit(1)        // want "Using os.Exit function outside of main function of main package is discouraged"
}
