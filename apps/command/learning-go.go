package main

import "fmt"

func hello(str string) {
	fmt.Println("Hello, " + str)
}

func sum(arr []int) (sum int) {
	for i := range arr {
		sum += arr[i]
	}
	return
}

func factorial(n int) int {
	if n == 1 {
		return 1
	}
	return n * factorial(n-1)
}

func main() {
	//p1
	hello("person")
	fmt.Println("Expected: Hello, person")

	//p2
	arr := []int{1, 2, 3, 4, 5}
	fmt.Println(sum(arr))
	fmt.Println("Expected: 15")

	//p3
	fmt.Println(factorial(4))
	fmt.Println("Expected: 24")
}
