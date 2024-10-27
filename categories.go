package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// Category represents a category or subcategory structure
type Category struct {
	Name          string     `json:"name"`
	Description   string     `json:"description"`
	ParentID      int        `json:"parent_category_id,omitempty"`
	Subcategories []Category `json:"subcategories,omitempty"`
}

// LoadCategories reads the JSON configuration file and returns categories
func LoadCategories(filename string) ([]Category, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var categories []Category
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &categories)
	return categories, err
}

// CreateCategory sends a POST request to Discourse to create a category
func CreateCategory(category Category, apiKey, apiUser, baseURL string, parentID int) (int, error) {
	url := fmt.Sprintf("%s/categories", baseURL)
	category.ParentID = parentID
	body, err := json.Marshal(category)
	if err != nil {
		return 0, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return 0, err
	}
	req.Header.Set("Api-Key", apiKey)
	req.Header.Set("Api-Username", apiUser)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return 0, fmt.Errorf("failed to create category: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}
	id := int(result["category"].(map[string]interface{})["id"].(float64))
	fmt.Printf("Created category '%s' with ID %d\n", category.Name, id)
	return id, nil
}

// ImportCategories recursively creates categories and subcategories
func ImportCategories(categories []Category, apiKey, apiUser, baseURL string, parentID int) error {
	for _, category := range categories {
		categoryID, err := CreateCategory(category, apiKey, apiUser, baseURL, parentID)
		if err != nil {
			log.Printf("Error creating category %s: %v\n", category.Name, err)
			continue
		}

		// Recursively create subcategories
		if len(category.Subcategories) > 0 {
			err := ImportCategories(category.Subcategories, apiKey, apiUser, baseURL, categoryID)
			if err != nil {
				log.Printf("Error creating subcategories for %s: %v\n", category.Name, err)
			}
		}
	}
	return nil
}
