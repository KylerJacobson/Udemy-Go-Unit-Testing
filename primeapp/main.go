package main

import "fmt"

func main() {

	n := 77

	_, msg := isPrime(n)

	fmt.Println(msg)
}

func isPrime(n int) (bool, string) {
	// 0 amd 1 are not prime numbers
	if n == 0 || n == 1 {
		return false, fmt.Sprintf("%d is not prime, by definition", n)
	}
	// negative numbers are not prime
	if n < 0 {
		return false, "Negative numbers are not prime, by definition"
	}

	// Use modulus operator repeatedly to see if we have prime number

	for i := 2; i <= n/2; i++ {
		if n%i == 0 {
			// not a prime number
			return false, fmt.Sprintf("%d is not prime, it is divisible by %d", n, i)
		}
	}
	return true, fmt.Sprintf("%d is a prime number", n)
}
