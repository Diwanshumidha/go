package main

import (
	"errors"
	"fmt"
	"math/rand"
)

type Difficulty struct {
	Name string
	Chances int
}

var difficultyChoiceMap = map[int]Difficulty{
	1: {"Easy", 10},
	2: {"Medium", 5},
	3: {"Hard", 3},
	4: {"Insane", 1},
}

func main() {
	display()
	chances := selectDifficulty()
	playGame(chances, false)
}


func playGame(chances int, isDebug bool) {

	randomNumber := rand.Intn(100) + 1

	if isDebug {
		fmt.Printf("[DEBUG] The number is %d\n", randomNumber)
	}

	tries := 0

	for {
		if tries >= chances {
			fmt.Printf("You lost! The number was %d\n", randomNumber)
			return
		}
		tries++

		var guess int
		fmt.Printf("Enter your guess: ")
		fmt.Scan(&guess)

		if guess == randomNumber {
			fmt.Printf("Congratulations! You guessed the correct number in %d tries!\n", tries)
			return
		}

		diff := randomNumber - guess


		switch {
			case diff < 0:
				fmt.Printf("Incorrect! The number is less than %d.\n", guess)
			case diff > 0:
				fmt.Printf("Incorrect! The number is greater than %d.\n", guess)
			default:
				fmt.Printf("Invalid input. Please try again.\n")
				continue;
		}
	}
}

func display() {
	fmt.Println("Welcome to the Number Guessing Game!")
	fmt.Println("I'm thinking of a number between 1 and 100.")
	fmt.Println("Can you guess what it is?")
	fmt.Println("")
}


func selectDifficulty() int {
	maxDifficulty := len(difficultyChoiceMap)
	fmt.Println("Please select the difficulty level:")
	for i, val := range difficultyChoiceMap {
		fmt.Printf("%d. %s (%d Chances)\n", i, val.Name, val.Chances)
	}
	fmt.Println("")

	for {
		var difficulty int
		fmt.Printf("Enter your choice: ")
		fmt.Scan(&difficulty)
		fmt.Println("")

		if(difficulty >= 1 && difficulty <= maxDifficulty) {
			difficulty, err := getDifficulty(difficulty)
			if err != nil {
				fmt.Println(err)
			}

			fmt.Printf("Great! You have selected the %s difficulty level.\n", difficulty.Name)
			fmt.Println("Let's start the game!")
			fmt.Println("")

			return difficulty.Chances
		}

		fmt.Printf("Invalid input. Please try again. Must be between 1 and %d.\n", maxDifficulty)
	}
}

func getDifficulty(difficulty int) (Difficulty, error) {
	val, ok := difficultyChoiceMap[difficulty]
	if !ok {
		return Difficulty{}, errors.New("Invalid difficulty level")
	}

	return val, nil
}
