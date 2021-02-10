package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

//Number of Questions to ask
const totalQuestions = 5

type Question struct {
	question string
	answer   string
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error at start up %v", err)
		os.Exit(0)
	}
}

func readArguments() string {

	fileName := flag.String("filename", "problems.csv", "quiz question file name")

	flag.Parse()

	// returning actual value as flag package returns a pointer
	return *fileName
}

func openFile(filename string) (io.Reader, error) {
	return os.Open(filename)
}

func readCSV(f io.Reader) ([]Question, error) {
	//	i. read csv
	csvR := csv.NewReader(f)

	//  ii. prepare questions slice
	questions := make([]Question, len())
	for {
		record, err := csvR.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}

		if record[0] != "" && record[1] != "" {
			questions = append(questions, Question{record[0], record[1]})
		}
	}
	return questions, nil
}

func getInput(input chan string) {
	for {
		in := bufio.NewReader(os.Stdin)
		result, err := in.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		input <- result
	}
}

func askQuestions(questions []Question) (int, error) {
	totalScore := 0

	done := make(chan string)

	go getInput(done)

	for i := range [totalQuestions]int{} {
		ans, err := askQuestion(questions[i].question, questions[i].answer, done)
		if err != nil {
			return totalScore, nil
		}
		totalScore += ans
	}
	return totalScore, nil
}

func askQuestion(question string, answer string, done <-chan string) (int, error) {
	fmt.Printf("Question: %v\n", question)
	for {
		select {
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
	fileName := readArguments()

	// 2. Open file
	f, err := openFile(fileName)

	if err != nil {
		return err
	}

	// 3. Read csv
	questions, err := readCSV(f)

	if err != nil {
		return errors.New("Error :: problem reading CSV")
	}

	if questions == nil {
		return errors.New("Error :: No Questions")
	}
	// 4. askQuestions
	score, err := askQuestions(questions)

	if err != nil {
		return errors.New("Error :: Problem asking questions")
	}

	// 5. return results
	fmt.Printf("You got %v out of a possible %v:\n", score, totalQuestions)

	return nil
}
