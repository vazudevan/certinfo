package main

import (
	"fmt"
	"os"
)

func main() {
	if err := newApp().Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s ERROR: %s\n", os.Args[0], err)
		os.Exit(1)
	}
}
