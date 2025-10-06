package jirareport

import "time"

// EpicGroup represents a group of issues under an epic
type EpicGroup struct {
	EpicKey     string
	EpicSummary string
	EpicStatus  string
	EpicURL     string
	Issues      []IssueUpdate
}

// IssueUpdate represents an issue with its updates
type IssueUpdate struct {
	Key         string
	Summary     string
	Status      string
	IssueType   string
	URL         string
	Updates     []Update
	LastUpdated time.Time
}

// Update represents a single update (comment or worklog)
type Update struct {
	Time       time.Time
	AuthorName string
	Type       string // "comment" or "worklog"
	Content    string
	TimeSpent  string // for worklogs
}

// Report contains both markdown and HTML versions of the report
type Report struct {
	Markdown string
	HTML     string
}

