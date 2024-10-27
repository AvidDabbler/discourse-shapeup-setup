package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

// PinnedMessage represents a pinned post structure in each category
type PinnedMessage struct {
	Category      string `json:"category"`
	PinnedMessage struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	} `json:"pinned_message"`
}

// ExportedPost represents the structure for exporting pinned post information
type ExportedPost struct {
	Category string `json:"category"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	URL      string `json:"url"`
}

// LoadPinnedMessages loads the pinned messages configuration from JSON
func loadPinnedMessages(filename string) ([]PinnedMessage, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var pinnedMessages struct {
		Categories []PinnedMessage `json:"categories"`
	}
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &pinnedMessages)
	return pinnedMessages.Categories, err
}

// CreateOrUpdatePinnedPost creates or updates a pinned post in Discourse
func createOrUpdatePinnedPost(categoryName, title, content string, categoryID int, apiKey, apiUser, baseURL string) (ExportedPost, error) {
	// Check if a pinned post exists by listing topics in the category
	checkURL := fmt.Sprintf("%s/c/%d.json", baseURL, categoryID)
	req, err := http.NewRequest("GET", checkURL, nil)
	if err != nil {
		log.Fatalf("Error creating GET request: %v", err)
	}

	// Add headers for authentication
	req.Header.Set("Api-Key", apiKey)
	req.Header.Set("Api-Username", apiUser)

	// Execute the GET request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ExportedPost{}, err
	}
	defer resp.Body.Close()

	// Save response body to a file named "response_body.json" for debugging
	file, err := os.Create("response_body.json")
	if err != nil {
		log.Fatalf("Error creating file: %v", err)
	}
	defer file.Close()
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		log.Fatalf("Error saving response body to file: %v", err)
	}

	// Resetting response body for further processing
	resp.Body.Close()
	resp, err = client.Do(req)
	if err != nil {
		return ExportedPost{}, err
	}
	defer resp.Body.Close()

	// Decode the response JSON
	var existingPosts map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&existingPosts); err != nil {
		return ExportedPost{}, err
	}

	// Search for an existing pinned post by title
	for _, topic := range existingPosts["topic_list"].(map[string]interface{})["topics"].([]interface{}) {
		topicData := topic.(map[string]interface{})
		if strings.EqualFold(topicData["title"].(string), title) {
			// Update the post if it already exists
			topicID := int(topicData["id"].(float64))
			updateURL := fmt.Sprintf("%s/t/%d.json", baseURL, topicID)
			body, _ := json.Marshal(map[string]interface{}{
				"title":  title,
				"raw":    content,
				"pinned": true,
			})
			req, _ := http.NewRequest("PUT", updateURL, bytes.NewBuffer(body))
			req.Header.Set("Api-Key", apiKey)
			req.Header.Set("Api-Username", apiUser)
			req.Header.Set("Content-Type", "application/json")

			resp, err := client.Do(req)
			if err != nil {
				return ExportedPost{}, err
			}
			defer resp.Body.Close()
			fmt.Printf("Updated pinned post in '%s' category.\n", categoryName)
			return ExportedPost{Category: categoryName, Title: title, Content: content, URL: fmt.Sprintf("%s/t/%d", baseURL, topicID)}, nil
		}
	}

	// Create a new pinned post if none exists
	createURL := fmt.Sprintf("%s/posts.json", baseURL)
	body, _ := json.Marshal(map[string]interface{}{
		"title":    title,
		"raw":      content,
		"category": categoryID,
	})
	req, err = http.NewRequest("POST", createURL, bytes.NewBuffer(body))
	if err != nil {
		return ExportedPost{}, err
	}
	req.Header.Set("Api-Key", apiKey)
	req.Header.Set("Api-Username", apiUser)
	req.Header.Set("Content-Type", "application/json")

	resp, err = client.Do(req)
	if err != nil {
		return ExportedPost{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return ExportedPost{}, fmt.Errorf("failed to create pinned post: %s", resp.Status)
	}

	var newPost map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&newPost); err != nil {
		return ExportedPost{}, err
	}
	topicID := int(newPost["topic_id"].(float64))

	fmt.Printf("Created new pinned post in '%s' category.\n", categoryName)
	return ExportedPost{Category: categoryName, Title: title, Content: content, URL: fmt.Sprintf("%s/t/%d", baseURL, topicID)}, nil
}

// ExportPinnedPosts exports all pinned posts details to JSON
func exportPinnedPosts(posts []ExportedPost, filename string) error {
	file, err := json.MarshalIndent(posts, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, file, 0644)
}
