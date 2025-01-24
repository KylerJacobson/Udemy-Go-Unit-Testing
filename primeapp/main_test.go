package main

import "testing"

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
