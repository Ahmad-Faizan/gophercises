package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strconv"
)

var problems map[string]int

func main() {
	csvFilename := flag.String("csv", "problems.csv", "A csv file in the format of 'question,answer'")
	flag.Parse()
	loadCSVfile(csvFilename)
}

func loadCSVfile(csvFilename *string) {
	file, err := os.Open(*csvFilename)
	if err != nil {
		fmt.Printf("The program cannot find the CSV file specified: %s", *csvFilename)
		os.Exit(1)
	}
	defer file.Close()

	problems, err := parseCSVfile(file)
	if err != nil {
		fmt.Printf("Unable to parse the CSV file")
		os.Exit(1)
	}
	startGame(problems)
}

func parseCSVfile(file *os.File) (map[string]int, error) {
	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return nil, err
	}
	problems := make(map[string]int, len(records[0]))
	for _, record := range records {
		if problems[record[0]], err = strconv.Atoi(record[1]); err != nil {
			return problems, err
		}
	}
	return problems, nil
}

func startGame(problems map[string]int) {
	var index, count int
	for ques, ans := range problems {
		var userAns int
		fmt.Printf("Question #%02d : %6s = ", index+1, ques)
		fmt.Scanf("%d\n", &userAns)
		if userAns == ans {
			count++
		}
		index++
	}
	fmt.Printf("You have answered %d questions correctly out of %d.\n", count, len(problems))
}
