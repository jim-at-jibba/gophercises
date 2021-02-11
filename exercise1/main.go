package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

//Number of Questions to ask
const totalQuestions = 5

type question struct {
	question string
	answer   string
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error at start up %v", err)
		os.Exit(0)
	}
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

func readArguments() (string, int) {

	fileName := flag.String("filename", "problems.csv", "quiz question file name")
	timeLimit := flag.Int("timeLimit", 30, "Time allowed for quiz")

	flag.Parse()

	// returning actual value as flag package returns a pointer
	return *fileName, *timeLimit
}

func openFile(filename string) (io.Reader, error) {
	return os.Open(filename)
}

func readCSV(f io.Reader) ([]question, error) {
	//	i. read csv
	csvR := csv.NewReader(f)

	// we can confidently use readAll because it is unlikely there
	// are gonna be millions of lines
	lines, err := csvR.ReadAll()

	if err != nil {
		exit("error reading csv")
	}
	//  ii. prepare questions slice
	questions := make([]question, len(lines))
	for i, line := range lines {
		questions[i] = question{
			question: line[0],
			answer:   line[1],
		}
	}
	return questions, nil
}

func getInput(input chan string) {
	for {
		in := bufio.NewReader(os.Stdin)
		result, err := in.ReadString('\n')
		if err != nil {
			exit(err.Error())
		}

		input <- result
	}
}

func askQuestions(questions []question, timeLimit int) (int, error) {
	totalScore := 0
	timer := time.NewTimer(time.Duration(timeLimit) * time.Second)

	done := make(chan string)

	go getInput(done)

	for i := range [totalQuestions]int{} {
		ans, err := askQuestion(questions[i].question, questions[i].answer, done, timer.C)
		if err != nil && ans == -1 {
			return totalScore, nil
		}
		totalScore += ans
	}
	return totalScore, nil
}

func askQuestion(question string, answer string, done <-chan string, timer <-chan time.Time) (int, error) {
	fmt.Printf("Question: %v\n", question)
	for {
		select {
		case <-timer:
			return -1, fmt.Errorf("Time out")
		case ans := <-done:
			score := 0
			if strings.Compare(strings.Trim(strings.ToLower(ans), "\n"), answer) == 0 {
				score = 1
			} else {
				return 0, nil
			}

			return score, nil
		}
	}
}

func run() error {
	fmt.Printf("Yay - cooking on gas\n")

	// 1. read args
	fileName, timeLimit := readArguments()

	// 2. Open file
	f, err := openFile(fileName)

	if err != nil {
		exit(err.Error())
	}

	// 3. Read csv
	questions, err := readCSV(f)

	if err != nil {
		exit("Error :: problem reading CSV")
	}

	if questions == nil {
		exit("Error :: No Questions")
	}
	// 4. askQuestions
	score, err := askQuestions(questions, timeLimit)

	if err != nil {
		exit("Error :: Problem asking questions")
	}

	// 5. return results
	fmt.Printf("You got %v out of a possible %v:\n", score, totalQuestions)

	return nil
}
