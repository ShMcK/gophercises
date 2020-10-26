package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

// Problem Q&A
type Problem struct {
	question string
	answer   string
}

// Progress track user score
type Progress struct {
	total   int
	correct int
}

func loadProblemSet(path string) []Problem {
	reader, err := os.Open(path)
	if err != nil {
		log.Fatal("No CSV found at " + path)
	}
	lines, err := csv.NewReader(reader).ReadAll() // assumes small vsc
	if err != nil {
		log.Fatal("Error with problem set at " + path)
	}
	var problems = make([]Problem, len(lines))
	for i, line := range lines {
		problems[i] = Problem{
			question: line[0],
			answer:   line[1],
		}
	}

	return problems
}

func validateResponse(response string, answer string, progress *Progress) {
	if strings.TrimSpace(response) == strings.TrimSpace(answer) {
		fmt.Println("Correct!")
		progress.correct++
	} else {
		fmt.Println("Wrong!")
	}
}

func captureResponse(reader *bufio.Reader) string {
	// read response until "enter" is pressed
	response, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal("Invalid response")
	}
	return response
}

func main() {
	// parse flags
	csvFilename := flag.String("csv", "problems.csv", "a csv file with 'questions,answers' format")
	timeout := flag.Int("timeout", 30, "time limit for the quiz in seconds")
	flag.Parse()
	// load csv
	problems := loadProblemSet(*csvFilename)

	progress := Progress{
		total:   3,
		correct: 0,
	}

	reader := bufio.NewReader(os.Stdin)

	timer := time.NewTimer(time.Duration(*timeout) * time.Second).C

	for i, p := range problems {
		// question
		fmt.Printf("Problem #%d: %s = \n", i+1, p.question)
		responseChannel := make(chan string)

		go func() {
			responseChannel <- captureResponse(reader)
		}()

		select {
		case <-timer:
			fmt.Println("Time over!")
			return
		case response := <-responseChannel:
			// check answer
			validateResponse(response, p.answer, &progress)
		}

		// check is complete
		if i == progress.total-1 {
			fmt.Printf("The end! %v out of %v correct\n", progress.correct, progress.total)
			return
		}
	}

}
