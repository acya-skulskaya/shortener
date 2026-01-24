package haspanic2

func main() {
	run()
}

func run() {
	panic("test for panic") // want "Using panic function is discouraged"
}
