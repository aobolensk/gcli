package main

import (
	"fmt"
	"os"
)

func getInfo(args []string, origin string) {
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
}
