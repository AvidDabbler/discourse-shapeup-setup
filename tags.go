package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// TagGroup structure based on Discourse API
type TagGroup struct {
	Name            string   `json:"name"`
	Tags            []string `json:"tags"`
	VisibilityLevel int      `json:"visibility_level"`
}

// Config for JSON import
type Config struct {
	TagGroups []TagGroup `json:"tag_groups"`
}

// Load configuration from tags_and_groups.json
func loadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(data, &config)
	return &config, err
}

// CreateTopic creates a topic with the specified tag to initialize it in Discourse
func createTopic(tag string, apiKey, apiUser, baseURL string) error {
	url := fmt.Sprintf("%s/posts", baseURL)
	body, _ := json.Marshal(map[string]interface{}{
		"title":    fmt.Sprintf("Initializing tag: %s", tag),
		"raw":      fmt.Sprintf("Temporary post to initialize tag: %s", tag),
		"category": 1, // Set this to the category ID where you want to create the topic
		"tags":     []string{tag},
	})

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Api-Key", apiKey)
	req.Header.Set("Api-Username", apiUser)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to create topic for tag '%s': %s", tag, resp.Status)
	}
	fmt.Printf("Initialized tag '%s' with a temporary topic\n", tag)
	return nil
}

// CreateTagGroup creates a tag group with existing tags
func createTagGroup(tagGroup TagGroup, apiKey, apiUser, baseURL string) error {
	url := fmt.Sprintf("%s/tag_groups", baseURL)
	body, err := json.Marshal(tagGroup)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Api-Key", apiKey)
	req.Header.Set("Api-Username", apiUser)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to create tag group '%s': %s", tagGroup.Name, resp.Status)
	}
	fmt.Printf("Tag group '%s' created successfully\n", tagGroup.Name)
	return nil
}
