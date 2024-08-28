package main

import "fmt"

func Panic() {
	panic("panic")
}

func main() {
	defer func() {
		// if r := recover(); r != nil {
		// 	fmt.Println("Recovered in f", r)
		// }
	}()
	Panic()
	fmt.Println("Hello World") // never reached
}
