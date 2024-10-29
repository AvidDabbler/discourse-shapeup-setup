package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// Category represents a Discourse category, including subcategories and settings
type Category struct {
	ID               int         `json:"id"`
	Name             string      `json:"name"`
	Description      string      `json:"description"`
	Slug             string      `json:"slug"`
	ParentCategoryID *int        `json:"parent_category_id"`
	Position         int         `json:"position"`
	Subcategories    []Category  `json:"subcategory_list,omitempty"`
	Settings         interface{} `json:"settings,omitempty"` // Category-specific settings
}

// CategoriesData structure to match the JSON format
type CategoriesData struct {
	Categories []Category `json:"categories"`
}

// Fetch base categories list from Discourse
func fetchCategories(apiKey, apiUser, baseURL string) ([]Category, error) {
	url := fmt.Sprintf("%s/categories.json", baseURL)
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
		return nil, fmt.Errorf("failed to fetch categories: %s", resp.Status)
	}

	var rootResponse struct {
		CategoryList struct {
			Categories []Category `json:"categories"`
		} `json:"category_list"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&rootResponse); err != nil {
		return nil, err
	}

	return rootResponse.CategoryList.Categories, nil
}

// Fetch detailed category info, including settings, for each category
func fetchCategoryDetails(category Category, apiKey, apiUser, baseURL string) (Category, error) {
	url := fmt.Sprintf("%s/c/%s/%d.json", baseURL, category.Slug, category.ID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return category, err
	}
	req.Header.Set("Api-Key", apiKey)
	req.Header.Set("Api-Username", apiUser)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return category, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return category, fmt.Errorf("failed to fetch category details for %s: %s", category.Name, resp.Status)
	}

	var detailedCategory struct {
		Category struct {
			Subcategories []Category  `json:"subcategory_list"`
			Settings      interface{} `json:"settings"`
		} `json:"category"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&detailedCategory); err != nil {
		return category, err
	}

	// Update the main category with detailed info
	category.Subcategories = detailedCategory.Category.Subcategories
	category.Settings = detailedCategory.Category.Settings

	return category, nil
}

// Fetch all categories with detailed information, including settings
func fetchAllCategoriesWithDetails(apiKey, apiUser, baseURL string) ([]Category, error) {
	categories, err := fetchCategories(apiKey, apiUser, baseURL)
	if err != nil {
		return nil, err
	}

	var detailedCategories []Category
	for _, category := range categories {
		detailedCategory, err := fetchCategoryDetails(category, apiKey, apiUser, baseURL)
		if err != nil {
			log.Printf("Error fetching details for category %s: %v\n", category.Name, err)
			continue
		}
		detailedCategories = append(detailedCategories, detailedCategory)

		// Optional: Add a delay to avoid rate limits
		time.Sleep(500 * time.Millisecond)
	}

	return detailedCategories, nil
} // Save categories to a JSON file
func saveCategories(categories []Category, filename string) error {
	file, err := json.MarshalIndent(categories, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, file, 0644)
}
