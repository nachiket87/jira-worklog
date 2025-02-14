# Jira Worklog CLI

A command-line tool for logging work in Jira. This tool allows you to quickly add work logs to Jira issues with support for human-readable time formats like "1h 30m".

## Installation

```bash
go install github.com/nachiket87/jira-worklog/cmd/worklog
```

## Setup

1. First, get your Jira API token:
   - Go to https://id.atlassian.com/manage-profile/security/api-tokens
   - Click "Create API token"
   - Give it a meaningful label (e.g., "Worklog CLI")
   - Copy the generated token

2. Configure the tool:
```bash
worklog configure
```

This will prompt you for:
- Your Jira email (e.g., your-name@onestepsoftware.com)
- Your API token (from step 1)
- Base URL (optional - press Enter for default)

The configuration will be saved in:
- Mac/Linux: `~/.config/jira-worklog/config.json`
- Windows: `%USERPROFILE%\.config\jira-worklog\config.json`

## Usage

Simply run:
```bash
worklog
```

The tool will prompt you for:
- Jira issue key (e.g., MEET-1)
- Comment text
- Time spent

Supported time formats:
- "1h 30m" (1 hour 30 minutes)
- "45m" (45 minutes)
- "3h" (3 hours)
- "90m" (90 minutes)
- "1.5h" (1 hour 30 minutes)
- "30s" (30 seconds)

## Commands

- `worklog` - Add a worklog entry
- `worklog configure` - Set up or update your Jira credentials
- `worklog help` - Show help message

## Development Setup

If you want to contribute or modify the tool:

1. Clone the repository:
```bash
git clone https://github.com/nachiket87/jira-worklog.git
cd jira-worklog
```

2. Install dependencies:
```bash
go mod tidy
```

3. Build:
```bash
go build ./cmd/worklog
```

## Security Note

Your API token is stored securely in your user configuration directory with read/write permissions only for your user account. Never share your API token or commit it to version control.
