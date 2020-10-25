package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
)

// Problem Q&A
type Problem struct {
	question string
	answer   string
}

// Progress track user score
type Progress struct {
	total     int
	correct   int
	incorrect int
}

func loadProblemSet(path string) []Problem {
	reader, err := os.Open(path)
	if err != nil {
		log.Fatal("No CSV found at " + path)
	}
	problems, err := csv.NewReader(reader).ReadAll()
	if err != nil {
		log.Fatal("Error with problem set at " + path)
	}
	var problemList []Problem
	for _, line := range problems {
		problem := Problem{line[0], line[1]}
		problemList = append(problemList, problem)
	}

	return problemList
}

func main() {
	problems := loadProblemSet("problems.csv")
	progress := Progress{2, 0, 0}
	reader := bufio.NewReader(os.Stdin)

	for i, p := range problems {
		if i == progress.total {
			fmt.Printf("The end! %v out of %v correct\n", progress.correct, progress.total)
			return
		}
		fmt.Printf(p.question + "=")

		response, err := reader.ReadString('\n')

		if err != nil {
			log.Fatal("Invalid response")
		}

		fmt.Println(response)
		if strings.TrimSpace(response) == p.answer {
			progress.correct++
		} else {
			progress.incorrect++
		}
		fmt.Println(progress)
	}

}
