package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func locateDotGit() error {
	path, err := os.Getwd()
	if err != nil {
		return err
	}
	_, err = os.Stat(filepath.Join(path, ".git"))
	return err
}

func query(repository string) ([]map[string]interface{}, error) {
	client := http.Client{}
	var result []map[string]interface{}
	for page := 1; ; page++ {
		req, err := http.NewRequest(
			"GET",
			"https://api.github.com/repos/"+repository+"/issues?per_page=100&page="+
				strconv.Itoa(page),
			nil)
		if err != nil {
			return nil, err
		}
		req.Header.Add("Authorization", "token "+os.Getenv("GITHUB_TOKEN"))
		req.Header.Add("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		var current []map[string]interface{}
		defer resp.Body.Close()
		err = json.NewDecoder(resp.Body).Decode(&current)
		if err != nil {
			return nil, err
		}
		if len(current) == 0 {
			break
		}
		result = append(result, current...)
	}
	return result, nil
}

func main() {
	err := locateDotGit()
	if err != nil {
		fmt.Println("Could not find .git folder")
		os.Exit(1)
	}
	if os.Getenv("GITHUB_TOKEN") == "" {
		fmt.Fprintln(os.Stderr, "Please, provide GITHUB_TOKEN as environment variable")
		os.Exit(1)
	}
	repository := flag.String("repo", "gooddoog/gcli", "Repository: author/repo")
	flag.Parse()
	fmt.Println("Getting list of opened issues for " + *repository + ":")
	resp, err := query(*repository)
	for _, issue := range resp {
		fmt.Println(issue["html_url"])
	}
	if err != nil {
		fmt.Println(err)
	}
}
