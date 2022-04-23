package main

import (
	"fmt"
	"os"

	_ "embed"
)

//go:embed data.txt
var sampleData []byte

func main() {
	if len(os.Args) < 1 {
		printHelpAndExit()
	}

	args := os.Args[1:]

	url := ""

	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "-u" {
			if i == len(args)-1 {
				fmt.Fprintf(os.Stderr, "Argument %s requires a value\n", arg)
				printHelpAndExit()
			}
			url = args[i+1]
			i++
		} else {
			fmt.Fprintf(os.Stderr, "Unknown argument %s\n", arg)
			printHelpAndExit()
		}
	}

	performTestSuite(url)
}

func printHelpAndExit() {
	fmt.Printf("Usage: %s -u <absolute URL>\n", os.Args[0])
	os.Exit(1)
}
