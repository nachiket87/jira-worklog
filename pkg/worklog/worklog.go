package worklog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	BaseURL string
	Email   string
	Token   string
}

func NewClient(email, token string) *Client {
	return &Client{
		BaseURL: "https://onestepsoftware.atlassian.net/rest/api/3",
		Email:   email,
		Token:   token,
	}
}

func ParseTimeInput(input string) (int, error) {
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

func (c *Client) AddWorklog(issueKey, comment string, timeSpentSeconds int) ([]byte, error) {
	worklog := WorklogEntry{
		Comment: Comment{
			Content: []ContentItem{
				{
					Content: []TextContent{
						{
							Text: comment,
							Type: "text",
						},
					},
					Type: "paragraph",
				},
			},
			Type:    "doc",
			Version: 1,
		},
		TimeSpentSeconds: timeSpentSeconds,
	}

	jsonData, err := json.Marshal(worklog)
	if err != nil {
		return nil, fmt.Errorf("error marshaling JSON: %w", err)
	}

	url := fmt.Sprintf("%s/issue/%s/worklog", c.BaseURL, issueKey)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(c.Email, c.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API error: %s - %s", resp.Status, string(body))
	}

	return body, nil
}
