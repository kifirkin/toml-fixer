package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

var (
	inputFile  *string
	outputFile *string
	prefix     *string
	action     *string
)

func init() {
	inputFile = flag.String("i", "", "Input file path.")
	outputFile = flag.String("o", "", "Output file path.")
	prefix = flag.String("p", "", "Prefix to search for target line.")
	action = flag.String("action", "", "Predefined action to be perfomed. Available: `price:fix`, `quote:escape`")
}

func main() {
	flag.Parse()

	if *action == "" {
		fmt.Println("Please specify wich action needed to be perfomed with flag `-action`")
		flag.PrintDefaults()
	}

	if *inputFile == "" {
		fmt.Println("Please set input file with flag `-i`")
		flag.PrintDefaults()
	}

	if *outputFile == "" {
		fmt.Println("Please set output file with flag `-o`")
		flag.PrintDefaults()
	}

	if *prefix == "" {
		fmt.Println("Please set string prefix with flag `-p`")
		flag.PrintDefaults()
	}

	input, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")

	type callbackFn func(line string, i int) (result string)

	var cb callbackFn

	switch *action {
	case "price:fix":
		cb = func(line string, i int) string {
			parts := strings.Split(line, ".")

			if len(parts) == 1 {
				return strings.TrimSpace(line) + ".0"
			}
			return line
		}
	case "quote:escape":
		cb = func(line string, i int) string {
			return strings.Replace(line, "\"", "\\\"", 0)
		}
	}

	for i, line := range lines {
		if strings.HasPrefix(line, *prefix) {
			lines[i] = cb(line, i)
		}
	}

	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(*outputFile, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}
}
