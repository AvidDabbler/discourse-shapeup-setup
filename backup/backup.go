package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Load environment variables
func loadEnv() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func main() {
	loadEnv()

	apiKey := os.Getenv("DISCOURSE_API_KEY")
	apiUser := os.Getenv("DISCOURSE_API_USER")
	baseURL := os.Getenv("DISCOURSE_BASE_URL")

	categories, err := fetchAllCategoriesWithDetails(apiKey, apiUser, baseURL)
	if err != nil {
		log.Fatalf("Error fetching categories with details: %v", err)
	}

	err = saveCategories(categories, "categories.json")
	if err != nil {
		log.Fatalf("Error saving categories to file: %v", err)
	}

	fmt.Println("Categories backup saved to categories_backup.json")

	tags, err := fetchTags(apiKey, apiUser, baseURL)
	if err != nil {
		log.Fatalf("Error fetching tags: %v", err)
	}

	err = saveTags(tags, "tags.json")
	if err != nil {
		log.Fatalf("Error saving tags to file: %v", err)
	}

	fmt.Println("Tags backup saved to tags_backup.json")
}
