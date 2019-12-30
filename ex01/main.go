package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

var (
	csvFileName string
	timeLimit   int
)

type question struct {
	question string
	answer   string
}

func parseQuiz(lines [][]string) []question {
	ret := make([]question, len(lines))
	for i, line := range lines {
		ret[i] = question{
			question: strings.TrimSpace(line[0]),
			answer:   strings.TrimSpace(line[1]),
		}
	}
	return ret
}

func readQuizFile(fileName string) [][]string {
	file, err := os.Open(csvFileName)
	defer file.Close()

	if err != nil {
		exit(fmt.Sprintf("Failed to open file %s\n", csvFileName))
	}

	r := csv.NewReader(file)
	lines, err := r.ReadAll()

	if err != nil {
		exit("An error ocurred while trying to parse the csv file.")
	}

	return lines
}

func quiz(questions []question, timer *time.Timer) (score int) {
	for i, q := range questions {
		fmt.Printf("Problem #%d: %s = \n", i+1, q.question)

		answerCh := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer
		}()

		select {
		case <-timer.C:
			fmt.Println("\nTimeout!")
			return score
		case answer := <-answerCh:
			if answer == q.answer {
				score++
			}
		}
	}
	return score
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

func main() {
	flag.StringVar(&csvFileName, "csv", "problems.csv", "csv file in the format of question,answer.")
	flag.IntVar(&timeLimit, "limit", 30, "time limit for the quiz in seconds.")
	flag.Parse()

	lines := readQuizFile(csvFileName)
	questions := parseQuiz(lines)
	timer := time.NewTimer(time.Duration(timeLimit) * time.Second)
	score := quiz(questions, timer)
	fmt.Printf("You scored %d out of %d.\n", score, len(questions))
}
