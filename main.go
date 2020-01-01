package main

import (
	"fmt"
	"os"
	"strconv"
)

func process(args []string) {
	_, err := locateDotGit()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Could not find .git folder")
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	origin, err := extractOrigin()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Could not extract origin remote")
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Use: gcli help")
		os.Exit(1)
	}
	switch args[0] {
	case "info":
		getInfo(args, origin)
	case "commit":
		if len(args) == 1 {
			getOpenCommits(args, origin)
		} else if len(args) == 2 {
			getCommitBySHA(args, origin)
		} else {
			fmt.Fprintln(os.Stderr, "Unknown arguments for "+args[0])
			os.Exit(1)
		}
	case "issue":
		if len(args) == 1 {
			getOpenIssues(args, origin)
		} else if len(args) == 2 {
			if args[1] == "create" {
				createIssue(args, origin)
			} else if _, err := strconv.Atoi(args[1]); err == nil {
				getIssueByNumber(args, origin)
			} else {
				fmt.Fprintln(os.Stderr, "Unknown arguments for "+args[0])
				os.Exit(1)
			}
		} else if len(args) == 3 {
			if args[1] == "edit" {
				if _, err := strconv.Atoi(args[2]); err == nil {
					editIssue(args, origin)
				} else {
					fmt.Fprintln(os.Stderr, "Format: gcli issue edit <number>")
					os.Exit(1)
				}
			} else if args[1] == "close" {
				if _, err := strconv.Atoi(args[2]); err == nil {
					closeIssue(args, origin)
				} else {
					fmt.Fprintln(os.Stderr, "Format: gcli issue close <number>")
					os.Exit(1)
				}
			} else if args[1] == "reopen" {
				if _, err := strconv.Atoi(args[2]); err == nil {
					reopenIssue(args, origin)
				} else {
					fmt.Fprintln(os.Stderr, "Format: gcli issue reopen <number>")
					os.Exit(1)
				}
			} else {
				fmt.Fprintln(os.Stderr, "Unknown arguments for "+args[0])
				os.Exit(1)
			}
		} else {
			fmt.Fprintln(os.Stderr, "Unknown arguments for "+args[0])
			os.Exit(1)
		}
	case "pr":
		if len(args) == 1 {
			getOpenPullRequests(args, origin)
		} else if len(args) == 2 {
			if _, err := strconv.Atoi(args[1]); err == nil {
				getPullRequestByNumber(args, origin)
			} else {
				fmt.Fprintln(os.Stderr, "Unknown arguments for "+args[0])
				os.Exit(1)
			}
		} else {
			fmt.Fprintln(os.Stderr, "Unknown arguments for "+args[0])
			os.Exit(1)
		}
	case "help":
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
	default:
		fmt.Fprintln(os.Stderr, "Unknown command. Use: gcli help")
		os.Exit(1)
	}
}

func main() {
	if os.Getenv("GITHUB_TOKEN") == "" {
		fmt.Fprintln(os.Stderr, "Please, provide GITHUB_TOKEN as environment variable")
		os.Exit(1)
	}
	process(os.Args[1:])
}
