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
		} else if len(args) == 3 {
			if args[1] == "edit" {
				if _, err := strconv.Atoi(args[2]); err == nil {
					editPullRequest(args, origin)
				} else {
					fmt.Fprintln(os.Stderr, "Format: gcli pr edit <number>")
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
	case "help":
		if len(args) == 1 {
			getHelp(args, origin)
		} else {
			fmt.Fprintln(os.Stderr, "Unknown arguments for "+args[0])
			os.Exit(1)
		}
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
