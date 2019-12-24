package main

import (
	"encoding/json"
	"net/http"
	"os"
)

func queryObject(httpMethod string, query string) (map[string]interface{}, error) {
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
	var result map[string]interface{}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func queryList(httpMethod string, query string) ([]map[string]interface{}, error) {
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
