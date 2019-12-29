package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
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
			fmt.Println("List of opened pull requests for " + origin + ":")
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
				for _, issue := range resp {
					link := issue["html_url"].(string)
					if strings.Contains(link, "pull") {
						fmt.Printf("%-*.*s > %s\n", 50, 50, issue["title"].(string), link)
					}
				}
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}
			}
		} else if len(args) == 2 {
			if _, err := strconv.Atoi(args[1]); err == nil {
				// Get issue information by number
				result, err := queryObject(
					"GET",
					"https://api.github.com/repos/"+origin+"/pulls/"+args[1])
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}
				if result["message"] == "Not Found" {
					fmt.Fprintln(os.Stderr, args[1]+" is an issue")
					os.Exit(1)
				}
				link := result["html_url"].(string)
				splittedLink := strings.Split(link, "/")
				state := result["state"].(string)
				if state == "open" {
					state = "\033[32m" + state + "\033[0m"
				} else if state == "closed" {
					state = "\033[31m" + state + "\033[0m"
					if result["merged_at"] != nil {
						state = "\033[35mmerged\033[0m"
					}
				}
				labels := [](string){}
				for _, label := range result["labels"].([]interface{}) {
					labels = append(labels, label.(map[string]interface{})["name"].(string))
				}
				if splittedLink[len(splittedLink)-1] == args[1] {
					fmt.Println("\033[1m" + result["title"].(string) +
						" [" + state + "]\033[0m")
					if len(labels) > 0 {
						fmt.Println("Labels: " + strings.Join(labels[:], ", "))
					}
					fmt.Println(result["body"].(string))
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
