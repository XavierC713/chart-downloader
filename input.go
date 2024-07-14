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
	return strings.Trim(input, "\r")
}

func promptWithOptions(question string, options []string) int {
	fmt.Println(question)
	for index, option := range options {
		fmt.Printf("\n(%d) %s", index+1, option)
	}
	rawInput := prompt("Select one: ")
	input, err := strconv.Atoi(rawInput)
	if err != nil {
		panic(err)
	}
	// check if input is valid
	if input <= len(options) {
		return input
	}
	fmt.Println("\nInvalid selection")
	return promptWithOptions(question, options)
}

func promptConfirmation(question string) bool {
	response := strings.ToLower(prompt(fmt.Sprintf("\n%s[Y/n]: ", question)))
	if response == "" {
		return true
	}
	if response != "y" && response != "n" {
		fmt.Println("Invalid response")
		return promptConfirmation(question)
	}
	return response == "y"
}
