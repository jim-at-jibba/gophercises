package main

import (
	"fmt"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error at start up %v", err)
		os.Exit(0)
	}
}

func run() error {
	fmt.Printf("Yay - cooking on gas")

	return nil
}
