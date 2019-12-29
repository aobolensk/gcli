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
		resp, err := queryObject(
			"GET",
			"https://api.github.com/repos/"+origin)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Printf("Link: %s -> %d ⭐; %d ⑂; %d ❗; %d 👁️\n",
			resp["html_url"].(string),
			uint(resp["stargazers_count"].(float64)),
			uint(resp["forks"].(float64)),
			uint(resp["open_issues"].(float64)),
			uint(resp["subscribers_count"].(float64)))
		fmt.Println("Owner: " + resp["owner"].(map[string]interface{})["html_url"].(string))
		fmt.Println("Last update: " + resp["updated_at"].(string))
	case "commit":
		if len(args) == 1 {
			// Get list of opened commits
			fmt.Println("List of commits for " + origin + ":")
			var result []map[string]interface{}
			for page := 1; ; page++ {
				resp, err := queryList(
					"GET",
					"https://api.github.com/repos/"+origin+"/commits?per_page=100&page="+
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
			for _, commit := range result {
				fmt.Printf("%s %-*.*s %s\n",
					commit["sha"].(string)[0:7],
					20, 20,
					commit["commit"].(map[string]interface{})["author"].(map[string]interface{})["name"].(string),
					strings.Split(commit["commit"].(map[string]interface{})["message"].(string), "\n")[0])
			}
		} else {
			fmt.Fprintln(os.Stderr, "Unknown arguments for "+args[0])
			os.Exit(1)
		}
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
					fmt.Printf("%-*.*s > %s\n", 50, 50, issue["title"].(string), link)
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
				if result["pull_request"] != nil {
					fmt.Fprintln(os.Stderr, args[1]+" is a pull request")
					os.Exit(1)
				}
				link := result["html_url"].(string)
				splittedLink := strings.Split(link, "/")
				state := result["state"].(string)
				if state == "open" {
					state = "\033[32m" + state + "\033[0m"
				} else if state == "closed" {
					state = "\033[31m" + state + "\033[0m"
				}
				if splittedLink[len(splittedLink)-1] == args[1] {
					fmt.Println("\033[1m" + result["title"].(string) +
						" [" + state + "]\033[0m")
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
	case "pr":
		if len(args) == 1 {
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
					fmt.Printf("%-*.*s > %s\n", 50, 50, issue["title"].(string), link)
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
				if splittedLink[len(splittedLink)-1] == args[1] {
					fmt.Println("\033[1m" + result["title"].(string) +
						" [" + state + "]\033[0m")
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
