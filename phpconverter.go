package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func main() {
	// Open the input file
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	outputFile, err := os.Create("output.txt")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer outputFile.Close()

	writer := bufio.NewWriter(outputFile)
	defer writer.Flush()

	scanner := bufio.NewScanner(file)
	re := regexp.MustCompile(`\("([^,]+)", (\d+), "(\d*)n?(.*)"\);`)

	for scanner.Scan() {
		line := scanner.Text()
		match := re.FindStringSubmatch(line)
		if match != nil {
			name := match[1]
			level := match[2]
			nodes := match[3]
			if nodes == "" {
				nodes = "0"
			}
			comment := strings.TrimSpace(match[4])

			if comment == "" {
				fmt.Fprintf(writer, "%s\t%s\t%s\n", name, level, nodes)
			} else {
				fmt.Fprintf(writer, "%s\t%s\t%s\t%s\n", name, level, nodes, comment)
			}
		} else {
			fmt.Println("No match for line:", line)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}

	fmt.Println("Data successfully written to output.txt")
}
