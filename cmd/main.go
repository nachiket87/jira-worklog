package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/nachiket87/jira-worklog/pkg/worklog"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		fmt.Println("Please create a .env file with JIRA_EMAIL and JIRA_API_TOKEN")
		fmt.Println("Example .env file contents:")
		fmt.Println("JIRA_EMAIL=your.email@example.com")
		fmt.Println("JIRA_API_TOKEN=your-api-token")
		os.Exit(1)
	}

	email := os.Getenv("JIRA_EMAIL")
	apiToken := os.Getenv("JIRA_API_TOKEN")

	if email == "" || apiToken == "" {
		fmt.Println("Error: JIRA_EMAIL and JIRA_API_TOKEN must be set in your .env file")
		os.Exit(1)
	}

	client := worklog.NewClient(email, apiToken)

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
