package main

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
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
