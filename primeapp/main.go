package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func main() {
	// print welcome message
	intro()
	// create a channel to indicate when a user wants to quit
	doneChan := make(chan bool)

	// create a goroutine to read in a user input and run progrma
	go readUserInput(os.Stdin, doneChan)
	// block until doneChan get a value
	<-doneChan
	// close the channel
	close(doneChan)
	// say goodbye
	fmt.Println("Goodbye.")
}

func readUserInput(in io.Reader, doneChan chan bool) {
	scanner := bufio.NewScanner(in)
	for {
		res, done := checkNumbers(scanner)
		if done {
			doneChan <- true
			return
		}
		fmt.Println(res)
		prompt()
	}
} 

func checkNumbers(scanner *bufio.Scanner) (string, bool) {
	// read user input
	scanner.Scan()

	// Check if a user wants to quit

	if strings.EqualFold(scanner.Text(), "q") {
		return "", true
	}
	// try to convert what the user typed into an int

	numToCheck, err := strconv.Atoi(scanner.Text())
	if err != nil {
		return "Please enter a whole number!", false
	}
	_, msg := isPrime(numToCheck)
	return msg, false
}

func intro() {
	fmt.Println("Is it Prime?")
	fmt.Println("------------")
	fmt.Println("Enter a whole number and we'll tell you if it's prime or not. Enter q to quit.")
	prompt()
}

func prompt() {
	fmt.Print("-> ")
}

func isPrime(n int) (bool, string) {
	// 0 and 1 are not prime

	if n == 0 || n == 1 {
		return false, fmt.Sprintf("%d is not prime by definition", n)
	}
	if n < 0 {
		return false, "Negative numbers are not prime"
	}

	for i := 2; i < n; i++ {
		if n%i == 0 {
			return false, fmt.Sprintf("%d is not prime, it's divisible by %d", n, i)
		}
	}
	return true, fmt.Sprintf("%d is prime", n)
}
