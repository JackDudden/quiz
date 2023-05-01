package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

type Question struct {
	prompt string
	answer string
}

func main() {
	pathPtr := flag.String("file", "problems.csv", "Path to problems csv file")
	timePtr := flag.Int("time", 30, "Time to complete the quiz")
	shufflePtr := flag.Bool("shuffle", false, "Whether to shuffle the quiz")
	flag.Parse()

	quiz, err := loadQuiz(pathPtr, shufflePtr)
	if err != nil {
		panic(err)
	}
	total := len(quiz)
	score := 0

	timer := time.NewTimer(time.Duration(*timePtr) * time.Second)
	for _, question := range quiz {
		fmt.Println("Question: " + question.prompt)
		answerCh := make(chan string)
		go func() {
			var answer string
			fmt.Scanln(&answer)
			answerCh <- strings.TrimSpace(answer)
		}()

		select {
		case <-timer.C:
			fmt.Println("Total Score: ", score, "/", total)
			return
		case answer := <-answerCh:
			if answer == question.answer {
				score++
			}
		}
	}
	fmt.Println("Total Score: ", score, "/", total)
}

func loadQuiz(fpath *string, shuffle *bool) ([]Question, error) {
	fileContent, err := os.Open(*fpath)
	if err != nil {
		return nil, err
	}
	defer fileContent.Close()

	lines, err := csv.NewReader(fileContent).ReadAll()
	if err != nil {
		return nil, err
	}

	var store []Question
	for _, line := range lines {
		store = append(store, Question{
			line[0],
			strings.TrimSpace(line[1]),
		})
	}

	if *shuffle {
		rand.Shuffle(len(store), func(i, j int) {
			store[i], store[j] = store[j], store[i]
		})
	}

	return store, nil
}
