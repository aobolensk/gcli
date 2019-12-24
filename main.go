package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func process(args []string) {
	err := locateDotGit()
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
	case "issue":
		if len(args) == 1 {
			// Get list of opened issues
			fmt.Println("List of opened issues for " + origin + ":")
			var result []map[string]interface{}
			for page := 1; ; page++ {
				resp, err := queryList(
					"GET",
					"https://api.github.com/repos/"+origin+"/issues?per_page=100&page="+
						strconv.Itoa(page))
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}
				if len(resp) == 0 {
					break
				}
				result = append(result, resp...)
			}
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			for _, issue := range result {
				link := issue["html_url"].(string)
				if strings.Contains(link, "issues") {
					fmt.Println(link)
				}
			}
		} else if len(args) == 2 {
			if _, err := strconv.Atoi(args[1]); err == nil {
				// Get issue information by number
				result, err := queryObject(
					"GET",
					"https://api.github.com/repos/"+origin+"/issues/"+args[1])
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}
				link := result["html_url"].(string)
				splittedLink := strings.Split(link, "/")
				if splittedLink[len(splittedLink)-1] == args[1] {
					fmt.Println(result["title"])
					fmt.Println(result["body"])
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
		fmt.Println("List of opened pull requests for " + origin + ":")
		var result []map[string]interface{}
		for page := 1; ; page++ {
			resp, err := queryList(
				"GET",
				"https://api.github.com/repos/"+origin+"/issues?per_page=100&page="+
					strconv.Itoa(page))
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			if len(resp) == 0 {
				break
			}
			result = append(result, resp...)
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		for _, issue := range result {
			link := issue["html_url"].(string)
			if strings.Contains(link, "pull") {
				fmt.Println(link)
			}
		}
	case "help":
		fmt.Println(
			"Usage:\n" +
				"\tgcli <command> [arguments]\n\n" +
				"The commands are:\n" +
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
