package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Response structure for the tag list endpoint
type TagsResponse struct {
	Tags []interface{} `json:"tags"`
}

// Fetch tags from Discourse API
func fetchTags(apiKey, apiUser, baseURL string) ([]interface{}, error) {
	url := fmt.Sprintf("%s/tags.json", baseURL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Api-Key", apiKey)
	req.Header.Set("Api-Username", apiUser)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch tags: %s", resp.Status)
	}

	var tagsResponse TagsResponse
	if err := json.NewDecoder(resp.Body).Decode(&tagsResponse); err != nil {
		return nil, err
	}

	return tagsResponse.Tags, nil
}

// Save tags to a JSON file
func saveTags(tags []interface{}, filename string) error {
	file, err := json.MarshalIndent(tags, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, file, 0644)
}
