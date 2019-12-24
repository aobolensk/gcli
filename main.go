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
			return origin, err
		}
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return "", err
}

func query(httpMethod string, query string) ([]map[string]interface{}, error) {
	client := http.Client{}
	req, err := http.NewRequest(
		httpMethod,
		query,
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
	var result []map[string]interface{}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func process(args []string) {
	err := locateDotGit()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Could not find .git folder", err)
		os.Exit(1)
	}
	origin, err := extractOrigin()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Could not extract origin remote", err)
		os.Exit(1)
	}
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Use: gcli help")
		os.Exit(1)
	}
	switch args[0] {
	case "issues":
		fmt.Println("List of opened issues for " + origin + ":")
		var result []map[string]interface{}
		for page := 1; ; page++ {
			resp, err := query(
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
			fmt.Println(issue["html_url"])
		}
	case "help":
		fmt.Println(
			"Usage:\n" +
				"\tgcli <command> [arguments]\n\n" +
				"The commands are:\n" +
				"\tissues\t\tget list of issues\n" +
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
