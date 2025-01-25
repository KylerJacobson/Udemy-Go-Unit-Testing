package main

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func Test_isPrime(t *testing.T) {
	primeTests := []struct {
		name     string
		testNum  int
		expected bool
		msg      string
	}{
		{
			"Prime", 7, true, "7 is a prime number",
		},
		{
			"Negative", -1, false, "Negative numbers are not prime, by definition",
		},
		{
			"Not Prime", 77, false, "77 is not prime, it is divisible by 7",
		},
		{
			"Zero", 0, false, "0 is not prime, by definition",
		},
	}

	for _, e := range primeTests {
		result, msg := isPrime(e.testNum)
		if e.expected != result {
			t.Errorf("%s expected %t but got %t for %d", e.name, e.expected, result, e.testNum)
		}

		if msg != e.msg {
			t.Errorf("%s expected message '%s' but got '%s'", e.name, e.msg, msg)
		}
	}
}

func Test_prompt(t *testing.T) {
	// Save a copy of os.Stdout so we can restore it later
	oldOut := os.Stdout

	// create a read and write pipe
	r, w, _ := os.Pipe()

	// Set os.Stdout to our write pipe
	os.Stdout = w

	prompt()

	// Close the write pipe so we can read from it
	_ = w.Close()
	os.Stdout = oldOut

	// Read the output of our prompt function from our read pipe
	out, _ := io.ReadAll(r)

	// perform the test
	if string(out) != "-> " {
		t.Errorf("Expected prompt to be '-> ' but got '%s'", string(out))
	}
}

func Test_intro(t *testing.T) {
	oldOut := os.Stdout

	r, w, _ := os.Pipe()

	os.Stdout = w

	intro()

	_ = w.Close()
	os.Stdout = oldOut

	out, _ := io.ReadAll(r)

	if !strings.Contains(string(out), "Enter a whole number") {
		t.Errorf("intro text is not correct; got '%s'", string(out))
	}
}

func Test_checkNumbers(t *testing.T) {

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			"Prime", "7", "7 is a prime number",
		},
		{
			"Not a Number", "Seven", "Please enter a whole number",
		},
		{
			"Quit", "q", "",
		},
		{
			name:     "Zero",
			input:    "0",
			expected: "0 is not prime, by definition",
		},
		{
			name:     "One",
			input:    "1",
			expected: "1 is not prime, by definition",
		},
		{
			name:     "Three",
			input:    "3",
			expected: "3 is a prime number",
		},
		{
			name:     "Negative",
			input:    "-1",
			expected: "Negative numbers are not prime, by definition",
		},
	}

	for _, e := range tests {
		input := strings.NewReader(e.input)
		reader := bufio.NewScanner(input)
		res, _ := checkNumbers(reader)
		if res != e.expected {
			t.Errorf("%s expected '%s' but got '%s'", e.name, e.expected, res)
		}
	}

	input := strings.NewReader("7")
	reader := bufio.NewScanner(input)
	res, _ := checkNumbers(reader)
	if !strings.EqualFold(res, "7 is a prime number") {
		t.Errorf("incorrect value returned; got '%s'", res)
	}
}

func Test_readUserInput(t *testing.T) {
	doneChan := make(chan bool)
	var stdin bytes.Buffer
	stdin.Write([]byte("1\nq\n"))
	go readUserInput(&stdin, doneChan)
	<-doneChan
	close(doneChan)
}
