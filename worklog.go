package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Comment struct {
	Content []ContentItem `json:"content"`
	Type    string        `json:"type"`
	Version int           `json:"version"`
}

type ContentItem struct {
	Content []TextContent `json:"content"`
	Type    string        `json:"type"`
}

type TextContent struct {
	Text string `json:"text"`
	Type string `json:"type"`
}

type WorklogEntry struct {
	Comment          Comment `json:"comment"`
	TimeSpentSeconds int     `json:"timeSpentSeconds"`
}

func parseTimeInput(input string) (int, error) {
	input = strings.ToLower(strings.TrimSpace(input))
	if !strings.Contains(input, "h") && !strings.Contains(input, "m") && !strings.Contains(input, "s") {
		input = input + "s"
	}

	duration, err := time.ParseDuration(input)
	if err != nil {
		return 0, err
	}

	return int(duration.Seconds()), nil
}

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

	var issueKey string

	fmt.Print("Enter Jira issue key (e.g., MEET-1): ")
	fmt.Scanln(&issueKey)

	fmt.Print("Enter comment text: ")
	commentText, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	commentText = strings.TrimSpace(commentText)

	fmt.Print("Enter time spent (e.g., 1h 30m, 90m, 3h, 45s): ")
	timeInput, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	timeSpent, err := parseTimeInput(timeInput)
	if err != nil {
		fmt.Printf("Error parsing time input: %v\n", err)
		return
	}

	worklog := WorklogEntry{
		Comment: Comment{
			Content: []ContentItem{
				{
					Content: []TextContent{
						{
							Text: commentText,
							Type: "text",
						},
					},
					Type: "paragraph",
				},
			},
			Type:    "doc",
			Version: 1,
		},
		TimeSpentSeconds: timeSpent,
	}

	jsonData, err := json.Marshal(worklog)
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		return
	}

	url := fmt.Sprintf("https://onestepsoftware.atlassian.net/rest/api/3/issue/%s/worklog", issueKey)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(email, apiToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		return
	}

	fmt.Printf("Response Status: %s\n", resp.Status)
	fmt.Printf("Response Body: %s\n", string(body))
}
