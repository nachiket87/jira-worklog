package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/nachiket87/jira-worklog/pkg/worklog"
)

func configure() error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter your Jira email: ")
	email, _ := reader.ReadString('\n')
	email = strings.TrimSpace(email)

	fmt.Print("Enter your Jira API token: ")
	token, _ := reader.ReadString('\n')
	token = strings.TrimSpace(token)

	// Ask for base URL (optional)
	fmt.Print("Enter your Jira base URL (press Enter for default): ")
	baseURL, _ := reader.ReadString('\n')
	baseURL = strings.TrimSpace(baseURL)

	config := &worklog.Config{
		Email:   email,
		Token:   token,
		BaseURL: baseURL,
	}

	if err := worklog.SaveConfig(config); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	fmt.Println("Configuration saved successfully!")
	configPath, _ := worklog.GetConfigPath()
	fmt.Printf("Configuration file location: %s\n", configPath)
	return nil
}

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "configure":
			if err := configure(); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			return
		case "help", "-h", "--help":
			fmt.Println("Usage:")
			fmt.Println("  worklog configure  - Configure Jira credentials")
			fmt.Println("  worklog           - Add a worklog entry")
			fmt.Println("  worklog help      - Show this help message")
			return
		}
	}

	config, err := worklog.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		fmt.Println("To configure, run: worklog configure")
		os.Exit(1)
	}

	client := worklog.NewClient(config.Email, config.Token)
	if config.BaseURL != "" {
		client.BaseURL = config.BaseURL
	}

	var issueKey string
	fmt.Print("Enter Jira issue key (e.g., MEET-1): ")
	fmt.Scanln(&issueKey)

	fmt.Print("Enter comment text: ")
	commentText, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	commentText = strings.TrimSpace(commentText)

	fmt.Print("Enter time spent (e.g., 1h 30m, 90m, 3h, 45s): ")
	timeInput, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	timeSpent, err := worklog.ParseTimeInput(timeInput)
	if err != nil {
		fmt.Printf("Error parsing time input: %v\n", err)
		return
	}

	response, err := client.AddWorklog(issueKey, commentText, timeSpent)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Worklog added successfully: %s\n", string(response))
}
