package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func getOpenPullRequests(args []string, origin string) {
	fmt.Println("List of opened pull requests for " + origin + ":")
	for page := 1; ; page++ {
		resp, err := queryList(
			"GET",
			"https://api.github.com/repos/"+origin+"/issues?per_page=100&page="+
				strconv.Itoa(page),
			nil)
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
}

func getPullRequestByNumber(args []string, origin string) {
	result, err := queryObject(
		"GET",
		"https://api.github.com/repos/"+origin+"/pulls/"+args[1],
		nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if result["message"] != nil {
		fmt.Fprintln(os.Stderr, result["message"])
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
}
