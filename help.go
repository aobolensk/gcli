package main

import "fmt"

func helpPrinter(command string, message string) {
	const SIZE = 35
	fmt.Printf("\t%-*.*s %s\n", SIZE, SIZE, command, message)
}

func getHelp(args []string, origin string) {
	fmt.Println("Usage")
	fmt.Println("\tgcli <command> [arguments]")
	fmt.Println("\nThe commands are:")
	helpPrinter("commit", "get list of commits in master branch")
	helpPrinter("commit <SHA>", "get info about particular commit")
	helpPrinter("info", "get info about this repo")
	helpPrinter("issue", "get list of issues")
	helpPrinter("issue <number>", "get info about particular issue")
	helpPrinter("issue create", "create a new issue")
	helpPrinter("issue edit <number>", "edit the issue")
	helpPrinter("issue close <number>", "close the issue")
	helpPrinter("issue reopen <number>", "reopen the issue")
	helpPrinter("issue assign <number> <assignee>", "assign user to the issue")
	helpPrinter("issue unassign <number> <assignee>", "unassign user from the issue")
	helpPrinter("pr", "get list of pull requests")
	helpPrinter("pr <number>", "get info about particular pull request")
	helpPrinter("pr edit <number>", "edit the pull request")
	helpPrinter("help", "get this help message")
}
