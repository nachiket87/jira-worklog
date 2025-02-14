package worklog

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
