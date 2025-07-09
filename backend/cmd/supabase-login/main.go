package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/supabase-community/gotrue-go"
)

const API_URL = ""
const ANNO_KEY = ""

func main() {
	email := ""
	password := ""
	token, err := getUserToken(email, password)
	if err != nil {
		fmt.Println("Error getting token:", err)
		os.Exit(1)
	}
	fmt.Println("Bearer", token)
}

// getUserToken retrieves a user's Bearer token using email and password authentication
func getUserToken(email, password string) (string, error) {
	// Extract project reference from URL
	projectRef := extractProjectRef(API_URL)
	if projectRef == "" {
		return "", fmt.Errorf("could not extract project reference from URL")
	}

	// Initialize GoTrue client with project reference and anon key
	authClient := gotrue.New(projectRef, ANNO_KEY)

	// Sign in with email and password
	resp, err := authClient.SignInWithEmailPassword(email, password)
	if err != nil {
		return "", fmt.Errorf("sign in failed: %w", err)
	}

	// Return access token (Bearer token)
	return resp.AccessToken, nil
}

// extractProjectRef extracts the project reference from a Supabase URL
func extractProjectRef(url string) string {
	// Remove protocol prefix (https://)
	url = strings.TrimPrefix(url, "https://")

	// Extract the subdomain (project reference)
	parts := strings.Split(url, ".")
	if len(parts) > 0 {
		return parts[0]
	}

	return ""
}
