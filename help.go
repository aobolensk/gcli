package main

import "fmt"

func getHelp(_ []string, _ string) {
	fmt.Println("Usage")
	fmt.Println("\tgcli <command> [arguments]")
	fmt.Println("\nThe commands are:")
	const SIZE = 24
	fmt.Printf("\t%-*.*s %s\n", SIZE, SIZE,
		"commit", "get list of commits in master branch")
	fmt.Printf("\t%-*.*s %s\n", SIZE, SIZE,
		"commit <SHA>", "get info about particular commit")
	fmt.Printf("\t%-*.*s %s\n", SIZE, SIZE,
		"info", "get info about this repo")
	fmt.Printf("\t%-*.*s %s\n", SIZE, SIZE,
		"issue", "get list of issues")
	fmt.Printf("\t%-*.*s %s\n", SIZE, SIZE,
		"issue <number>", "get info about particular issue")
	fmt.Printf("\t%-*.*s %s\n", SIZE, SIZE,
		"issue create", "create a new issue")
	fmt.Printf("\t%-*.*s %s\n", SIZE, SIZE,
		"issue edit <number>", "edit the issue")
	fmt.Printf("\t%-*.*s %s\n", SIZE, SIZE,
		"issue close <number>", "close the issue")
	fmt.Printf("\t%-*.*s %s\n", SIZE, SIZE,
		"issue reopen <number>", "reopen the issue")
	fmt.Printf("\t%-*.*s %s\n", SIZE, SIZE,
		"pr", "get list of pull requests")
	fmt.Printf("\t%-*.*s %s\n", SIZE, SIZE,
		"pr <number>", "get info about particular pull request")
	fmt.Printf("\t%-*.*s %s\n", SIZE, SIZE,
		"help", "get this help message")
}
