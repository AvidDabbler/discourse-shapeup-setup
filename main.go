package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Retrieve environment variables
	apiKey := os.Getenv("DISCOURSE_API_KEY")
	apiUser := os.Getenv("DISCOURSE_API_USER")
	baseURL := os.Getenv("DISCOURSE_BASE_URL")

	// Check for missing environment variables
	if apiKey == "" || apiUser == "" || baseURL == "" {
		log.Fatalf("Missing one or more required environment variables: DISCOURSE_API_KEY, DISCOURSE_API_USER, DISCOURSE_BASE_URL")
	}

	// Load tag groups from JSON file
	config, err := loadConfig("tags_and_groups.json")
	if err != nil {
		log.Fatalf("Error loading tags_and_groups.json: %v", err)
	}

	// Initialize each tag by creating a temporary topic for it
	for _, group := range config.TagGroups {
		for _, tag := range group.Tags {
			err := createTopic(tag, apiKey, apiUser, baseURL)
			if err != nil {
				log.Printf("Error initializing tag %s: %v\n", tag, err)
			}
		}
	}

	// Create tag groups with initialized tags
	for _, tagGroup := range config.TagGroups {
		err := createTagGroup(tagGroup, apiKey, apiUser, baseURL)
		if err != nil {
			log.Printf("Error creating tag group %s: %v\n", tagGroup.Name, err)
		}
	}

	// Load categories from JSON file
	categories, err := LoadCategories("categories.json")
	if err != nil {
		log.Fatalf("Error loading categories.json: %v", err)
	}

	// Import categories and subcategories
	err = ImportCategories(categories, apiKey, apiUser, baseURL, 0)
	if err != nil {
		log.Fatalf("Error importing categories: %v", err)
	}

	// TODO: fix the pinned messages updates
	//
	// // Load pinned messages from JSON file
	// pinnedMessages, err := loadPinnedMessages("pinned_messages.json")
	// if err != nil {
	// 	log.Fatalf("Error loading pinned_messages.json: %v", err)
	// }
	//
	// var exportedPosts []ExportedPost
	// for _, message := range pinnedMessages {
	// 	// Retrieve or set category ID for each category (for demo, assume 1 for all)
	// 	categoryID := 1 // Replace this with the actual category ID lookup
	// 	exportedPost, err := createOrUpdatePinnedPost(message.Category, message.PinnedMessage.Title, message.PinnedMessage.Content, categoryID, apiKey, apiUser, baseURL)
	// 	if err != nil {
	// 		log.Printf("Error creating or updating pinned post in %s category: %v\n", message.Category, err)
	// 		continue
	// 	}
	// 	exportedPosts = append(exportedPosts, exportedPost)
	// }
	//
	// // Export the pinned posts to a JSON file
	// if err := exportPinnedPosts(exportedPosts, "exported_pinned_posts.json"); err != nil {
	// 	log.Fatalf("Error exporting pinned posts: %v", err)
	// }
	//
	// fmt.Println("Pinned posts export completed.")
}
