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
	primeTest := []struct {
		name     string
		testNum  int
		expected bool
		msg      string
	}{
		{"prime", 7, true, "7 is prime"},
		{"not prime", 8, false, "8 is not prime, it's divisible by 2"},
		{"zero", 0, false, "0 is not prime by definition"},
		{"negative", -1, false, "Negative numbers are not prime"},
	}

	for _, e := range primeTest {
		result, msg := isPrime(e.testNum)
		if e.expected && !result {
			t.Errorf("%s: expected true but got false", e.name)
		}
		if !e.expected && result {
			t.Errorf("%s: expected false but got true", e.name)
		}
		if msg != e.msg {
			t.Errorf("%s: expected %s but got %s", e.name, e.msg, msg)
		}
	}
}

func Test_prompt(t *testing.T){
	// save a copy of os.Stdout
	oldOut := os.Stdout
	// create a read and write pipe
	r, w, _ := os.Pipe()

	// set os.Stdout to our write pipe
	os.Stdout = w
	prompt()

	// close our writer
	_ = w.Close()

	// reset os. Stdout to what is was before
	os.Stdout = oldOut

	// read the output of our prompt() func from our read pipe
	out, _ := io.ReadAll(r)

	if string(out) != "-> " {
		t.Errorf("incorrect prompt: expected -> but go %s", string(out))
	}
}

func Test_intro(t *testing.T){
	// save a copy of os.Stdout
	oldOut := os.Stdout
	// create a read and write pipe
	r, w, _ := os.Pipe()

	// set os.Stdout to our write pipe
	os.Stdout = w
	intro()

	// close our writer
	_ = w.Close()

	// reset os. Stdout to what is was before
	os.Stdout = oldOut

	// read the output of our prompt() func from our read pipe
	out, _ := io.ReadAll(r)

	if !strings.Contains(string(out), "Enter a whole number"){
		t.Errorf("Intro text is not correct, go %s", string(out))
	}
}

func Test_checkNumbers(t *testing.T) {

	checkNumbersTest := []struct {
		name string
		input string
		response string
	}{
		{"prime", "7", "7 is prime"},
		{"not prime", "4", "4 is not prime, it's divisible by 2"},
		{"zero", "0", "0 is not prime by definition"},
		{"one", "1", "1 is not prime by definition"},
		{"negative", "-11", "Negative numbers are not prime"},
		{"decimal", "1.5", "Please enter a whole number!"},
		{"empty", "", "Please enter a whole number!"},
		{"typed", "three", "Please enter a whole number!"},
		{"quit", "q", ""},
		
	}
	for _, e := range checkNumbersTest {
		input := strings.NewReader(e.input)
		reader := bufio.NewScanner(input)
		res, _ := checkNumbers(reader)
		
		if res != e.response{
			t.Errorf("expected %s but got %s", e.response, res)
		}
	}
}

func Test_readUserInput(t *testing.T) {
	// to test this func we need a channel, and an instance fo an io.Reader

	doneChan := make(chan bool)

	// create a reference to a bytes.Buffer
	var stdin bytes.Buffer

	stdin.Write([]byte("1\nq\n"))
	go readUserInput(&stdin,doneChan)
	<-doneChan
	close(doneChan)
}