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
		} else {
			fmt.Fprintln(os.Stderr, "Unknown arguments for "+args[0])
			os.Exit(1)
		}
	case "issue":
		if len(args) == 1 {
			getOpenIssues(args, origin)
		} else if len(args) == 2 {
			if _, err := strconv.Atoi(args[1]); err == nil {
				getIssueByNumber(args, origin)
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
		fmt.Println(
			"Usage:\n" +
				"\tgcli <command> [arguments]\n\n" +
				"The commands are:\n" +
				"\tinfo\t\tget info about this repo\n" +
				"\tissue\t\tget list of issues\n" +
				"\tpr\t\tget list of pull requests\n" +
				"\thelp\t\tget this help message\n")
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
