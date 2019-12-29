package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func getOpenCommits(args []string, origin string) {
	fmt.Println("List of commits for " + origin + ":")
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
		for _, commit := range resp {
			fmt.Printf("%s %-*.*s %s\n",
				commit["sha"].(string)[0:7],
				20, 20,
				commit["commit"].(map[string]interface{})["author"].(map[string]interface{})["name"].(string),
				strings.Split(commit["commit"].(map[string]interface{})["message"].(string), "\n")[0])
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}
