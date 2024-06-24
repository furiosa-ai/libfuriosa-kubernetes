package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

var sourceFile = flag.String("source", "", "Path to the source file")
var outputFile = flag.String("output", "", "Path to the output file")

func main() {
	flag.Parse()

	if *sourceFile == "" || *outputFile == "" {
		fmt.Println("source and output flags are required")
		flag.Usage()
		return
	}

	// read source file
	inFile, err := os.Open(*sourceFile)
	if err != nil {
		fmt.Printf("Error opening source file: %v\n", err)
		return
	}
	defer inFile.Close()

	// create outputFile to flush the filtered output
	outFile, err := os.Create(*outputFile)
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		return
	}
	defer outFile.Close()

	// scanner to read infile
	scanner := bufio.NewScanner(inFile)
	// writer to write to outputFile
	writer := bufio.NewWriter(outFile)
	for scanner.Scan() {
		line := scanner.Text()
		// skip lines with ref* and allocs*
		if strings.Contains(line, "ref") || strings.Contains(line, "allocs") {
			continue
		}
		if _, err := writer.WriteString(line + "\n"); err != nil {
			fmt.Printf("Error writing to output file: %v\n", err)
			return
		}
	}

	// Ensure all buffered data is written to the output file
	writer.Flush()
}
