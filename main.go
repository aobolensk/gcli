package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
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

func extractOrigin() (string, error) {
	path, err := os.Getwd()
	if err != nil {
		return "", err
	}
	file, err := os.Open(filepath.Join(path, ".git", "config"))
	if err != nil {
		return "", err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		remoteOriginMatch, _ := regexp.MatchString("\\[remote \"origin\"\\]", scanner.Text())
		if remoteOriginMatch {
			r, _ := regexp.Compile("http[s]?://github\\.com/(.+)\\.")
			scanner.Scan()
			origin := r.FindStringSubmatch(scanner.Text())[1]
			fmt.Println(origin)
			return origin, err
		}
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return "", err
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
		fmt.Println("Could not find .git folder", err)
		os.Exit(1)
	}
	origin, err := extractOrigin()
	fmt.Println(origin)
	if err != nil {
		fmt.Println("Could not extract origin remote", err)
		os.Exit(1)
	}
	if os.Getenv("GITHUB_TOKEN") == "" {
		fmt.Fprintln(os.Stderr, "Please, provide GITHUB_TOKEN as environment variable")
		os.Exit(1)
	}
	fmt.Println("Getting list of opened issues for " + origin + ":")
	resp, err := query(origin)
	for _, issue := range resp {
		fmt.Println(issue["html_url"])
	}
	if err != nil {
		fmt.Println(err)
	}
}
