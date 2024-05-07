package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

var HarnessOrg = os.Getenv("HARNESSORG")
var HarnessProject = os.Getenv("HARNESSPROJECT")
var HarnessAccountId = os.Getenv("HARNESSACCOUNTID")
var HarnessAPIKey = os.Getenv("HARNESSAPIKEY")
var GithubToken = os.Getenv("GITHUBTOKEN")
var GithubUser = os.Getenv("GITHUBUSER")
var GithubRepo = os.Getenv("GITHUBREPO")
var HarnessBaseURL = "https://app.harness.io/gateway/pm/api/v1/policies"
var GithubRepoUrl = "https://api.github.com/repos/" + GithubUser + "/" + GithubRepo

var HarnessURL = HarnessBaseURL

func main() {

	if HarnessAccountId == "" || HarnessAPIKey == "" || GithubToken == "" || GithubUser == "" || GithubRepo == "" {
		fmt.Println("Error: One or More required environment variable or environment variable values are missing.")
		return
	}

	files, err := getFilesFromGitHub(GithubRepoUrl)
	if err != nil {
		fmt.Println("Error getting files from GitHub:", err)
		return
	}

	for _, file := range files {
		if filepath.Ext(file) == ".rego" {
			// Strip the .rego extension from the file name
			fileName := file[:len(file)-len(filepath.Ext(file))]

			content, err := getFileContent(GithubRepoUrl, file)
			if err != nil {
				fmt.Println("Error getting file content:", err)
				continue
			}

			payload := map[string]interface{}{
				"identifier": fileName, // Use the stripped file name without .rego extension
				"name":       fileName, // Use the stripped file name without .rego extension
				"rego":       string(content),
			}

			if err := createPolicy(payload); err != nil {
				fmt.Println("Error creating policy in Harness:", err)
			} else {
				fmt.Println("Policy created in Harness:", file)
			}
		}
	}
}

func getFilesFromGitHub(repoURL string) ([]string, error) {
	apiURL := fmt.Sprintf("%s/contents", repoURL)
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "token "+GithubToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to fetch data: %s", resp.Status)
	}

	var items []struct {
		Path string `json:"path"`
		Type string `json:"type"`
	}
	body, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &items); err != nil {
		return nil, err
	}

	var files []string
	for _, item := range items {
		if item.Type == "file" {
			files = append(files, item.Path)
		}
	}

	return files, nil
}

func getFileContent(repoURL, filePath string) ([]byte, error) {
	apiURL := fmt.Sprintf("%s/contents/%s", repoURL, filePath)
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "token "+GithubToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to fetch file content: %s", resp.Status)
	}

	var data struct {
		Content  string `json:"content"`
		Encoding string `json:"encoding"`
	}
	responseData, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(responseData, &data); err != nil {
		return nil, err
	}

	if data.Encoding == "base64" {
		return base64.StdEncoding.DecodeString(data.Content)
	}

	return []byte(data.Content), nil
}

func createPolicy(payload map[string]interface{}) error {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", HarnessURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", HarnessAPIKey)
	query := req.URL.Query()
	query.Add("accountIdentifier", HarnessAccountId)

	if HarnessOrg != "" {
		query.Add("orgIdentifier", HarnessOrg)
	}
	if HarnessProject != "" {
		query.Add("projectIdentifier", HarnessProject)
	}

	req.URL.RawQuery = query.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		responseBody, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("failed to create policy in Harness: %s", responseBody)
	}

	return nil
}
