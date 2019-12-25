package main

import (
	"bufio"
	"errors"
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
	for i := 1; i < 50; i++ {
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
	return "", errors.New("Reached max amount of iterations in locateDotGit()")
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
			r, _ := regexp.Compile("http[s]?://github\\.com/(.+)")
			scanner.Scan()
			origin := r.FindStringSubmatch(scanner.Text())[1]
			if strings.HasSuffix(origin, ".git") {
				origin = origin[:len(origin)-4]
			}
			return origin, err
		}
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return "", err
}
