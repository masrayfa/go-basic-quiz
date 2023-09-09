package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	// 1. input the name of the file
	fName := flag.String("f", "quiz.csv", "path of csv file")

	// 2. set duration of the timer
	timer := flag.Int("limit", 30, "the time limit for the quiz in seconds")
	flag.Parse()

	// 3. pull the problems from the file
	problems, err := problemPuller(*fName)
	if err != nil {
		exit(fmt.Sprintf("Failed to parse the provided CSV file. %s\n", err.Error()))
	}
	// 4. create a variable to count the correct answers
	correctAnswers := 0

	// 5. using the duration of the timer, initiate the timer
	timerObj := time.NewTimer(time.Duration(*timer) * time.Second)
	ansChan := make(chan string)

	// 6. create a for loop to iterate through the problems
	// 6.1 create a label for the loop
	problemLoop:
	for i, problem := range problems {
		// 7. create a variable to store the answer
		var answer string
		// 8. print the problem question to the terminal
		fmt.Printf("Problem #%d: %s = ", i+1, problem.question)

		// 9. create a goroutine to read the answer
		go func() {
			fmt.Scanf("%s\n", &answer)
			// 10. send the answer to the channel
			ansChan <- answer
		}()

		// 11. create a select statement to check if the timer has expired
		select {
		case <- timerObj.C:
			fmt.Printf("\nYou scored %d out of %d.\n", correctAnswers, len(problems))
			break problemLoop
		// 12. if the timer has not expired, check the answer
		case iAns := <- ansChan:
			// 13. if the answer is correct, increment the correct answer variable
			if iAns == problem.answer {
				correctAnswers++
			}
			if i == len(problems) - 1 {
				// 15. close the channel
				close(ansChan)
			}
		}
	}
	// 14. print the number of correct answers
	fmt.Printf("You scored %d out of %d.\n", correctAnswers, len(problems))
	fmt.Println("Thank you for playing!")
}

// create a struct to store the problems
type problem struct {
	question string
	answer   string
}

func problemPuller(filename string) ([]problem, error) {
	if fObj, err := os.Open(filename); err == nil {
		csvReader := csv.NewReader(fObj)
		if cLines, err := csvReader.ReadAll(); err == nil {
			defer fObj.Close()
			return parseProblem(cLines), nil
		} else {
			return nil, fmt.Errorf("Failed to parse the provided CSV file. %s\n", err.Error())
		}
	} else {
		return nil, fmt.Errorf("Failed to open the provided CSV file. %s\n", err.Error())
	}
}

func parseProblem(lines [][]string) []problem {
	// 1. create a variable to store the problems
	problemLines := make([]problem, len(lines))
	// 2. iterate through the lines
	for line := range lines {
		problemLines[line] = problem{
			// 0 is the question, 1 is the answer
			question: lines[line][0],
			answer:   lines[line][1],
		}
	}
	return problemLines
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}