# Discourse Shapeup bootstrap

## About

This repo is a Golang setup script to bootstrap a new discourse server for startups. This script is meant to setup a server so that I will have categories that will handle, customer service, feature development and requests as well as handle the cycles. The script bootstraps all of the categories, tags, groups, with its corresponding json file.

## Categories

The Categories in [categories.json](./categories.json) are setup with subcategories and can be customized to your needs.

It is important to note that if the category is already setup in the server (like General by default), it will fail with a 422. This is normal.

## Tags and Group Tags

While developing this project there was an issue with creating the tags themselves and I had to eventually resort to importing the tags via csv. This csv file is also included in this repo.

## Pinned posts

There was some attempt to updating the pinned posts after the creation, but this was a failed effort. I will just be updating them manually with the info in the [pinned_posts.json](./pinned_posts.json) file.

## Helpful tips

When it came to looking up what I referred to as "cards" I came to find that this is what is referred to as "boxes" in the discourse world.

List of components I use

- Clickable Topic
- Showcase Categories
- Search Banner
- Discourse Loading Slider
