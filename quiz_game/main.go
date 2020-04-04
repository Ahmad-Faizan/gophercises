package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"
)

var problems map[string]int
var timeLimit *int

func main() {
	csvFilename := flag.String("csv", "problems.csv", "A csv file in the format of 'question,answer'")
	timeLimit = flag.Int("time", 30, "the time limit for the quiz")
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
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	inputChannel := make(chan int)
	for ques, ans := range problems {
		fmt.Printf("Question #%02d : %6s = ", index+1, ques)
		go func() {
			var userAns int
			fmt.Scanf("%d\n", &userAns)
			inputChannel <- userAns
		}()
		select {
		case <-timer.C:
			fmt.Printf("\nTime limit exceeded.")
			fmt.Printf("\nYou have answered %d questions correctly out of %d.\n", count, len(problems))
			return
		case userAns := <-inputChannel:
			if userAns == ans {
				count++
			}
		}
		index++
	}
	fmt.Printf("You have answered %d questions correctly out of %d.\n", count, len(problems))
}
