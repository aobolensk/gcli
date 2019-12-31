package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
)

func queryObject(httpMethod string, query string, body map[string]interface{}) (map[string]interface{}, error) {
	client := http.Client{}
	bodyBuffer := new(bytes.Buffer)
	if body != nil {
		err := json.NewEncoder(bodyBuffer).Encode(body)
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest(
		httpMethod,
		query,
		bodyBuffer)
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

func queryList(httpMethod string, query string, body map[string]interface{}) ([]map[string]interface{}, error) {
	client := http.Client{}
	bodyBuffer := new(bytes.Buffer)
	if body != nil {
		err := json.NewEncoder(bodyBuffer).Encode(body)
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest(
		httpMethod,
		query,
		bodyBuffer)
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
