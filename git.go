package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func locateDotGit() (string, error) {
	path, err := os.Getwd()
	if err != nil {
		return "", err
	}
	const MAXITERATIONS = 32
	for i := 0; i < MAXITERATIONS; i++ {
		if _, err := os.Stat(path); err != nil {
			break
		}
		_, err := os.Stat(filepath.Join(path, ".git"))
		if err == nil {
			return filepath.Join(path, ".git"), err
		}
		newPath := filepath.Join(path, "..")
		if newPath == path {
			return "", errors.New("Reached root directory and could not find .git folder")
		}
		path = newPath
	}
	return "", fmt.Errorf(
		"Reached max amount of iterations (%d) in locateDotGit function", MAXITERATIONS)
}

func extractOrigin() (string, error) {
	path, err := locateDotGit()
	if err != nil {
		return "", err
	}
	file, err := os.Open(filepath.Join(path, "config"))
	if err != nil {
		return "", err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		remoteOriginMatch, _ := regexp.MatchString("\\[remote \"origin\"\\]", scanner.Text())
		if remoteOriginMatch {
			scanner.Scan()
			r, _ := regexp.Compile("(git@github\\.com:(.+))|(http[s]?://github\\.com/(.+))")
			match := r.FindStringSubmatch(scanner.Text())
			origin := ""
			// Extracted from "git@github.com:"
			if len(match) > 2 && len(match[2]) > 0 {
				origin = match[2]
			}
			// Extracted from "https://github.com:"
			if len(match) > 4 && len(match[4]) > 0 {
				origin = match[4]
			}
			// Truncate extra .git extension
			if strings.HasSuffix(origin, ".git") {
				origin = origin[:len(origin)-4]
			}
			if len(origin) == 0 {
				return "", fmt.Errorf("Origin is empty")
			}
			return origin, err
		}
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return "", err
}
