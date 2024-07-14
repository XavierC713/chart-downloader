package main

import (
	"fmt"
	"strconv"
	"strings"
)

func prompt(prompt string) string {
	fmt.Printf("\n%s", prompt)
	input, _ := reader.ReadString('\n')
	input = strings.Replace(input, "\n", "", -1)
	return input
}

func promptWithOptions(question string, options []string) int {
	fmt.Println(question)
	for index, option := range options {
		fmt.Printf("\n(%d) %s", index+1, option)
	}
	input, err := strconv.Atoi(prompt("Select one: "))
	if err != nil {
		fmt.Println("\nInvalid selection")
		return promptWithOptions(question, options)
	}
	// check if input is valid
	if input <= len(options) {
		return input
	}
	fmt.Println("\nInvalid selection")
	return promptWithOptions(question, options)
}
