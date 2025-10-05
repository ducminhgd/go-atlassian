package responsetypes

// IssueComment represents a single comment on a Jira issue
type IssueComment struct {
	Self         string                  `json:"self"`
	ID           string                  `json:"id"`
	Author       User                    `json:"author"`
	Body         AtlassianDocumentFormat `json:"body"`
	UpdateAuthor User                    `json:"updateAuthor"`
	Created      string                  `json:"created"`
	Updated      string                  `json:"updated"`
	JsdPublic    bool                    `json:"jsdPublic"`
	ParentID     string                  `json:"parentId"`
}

// PagedComment represents a paged list of comments with pagination information
type PagedComment struct {
	StartAt    int            `json:"startAt"`
	MaxResults int            `json:"maxResults"`
	Total      int            `json:"total"`
	Comments   []IssueComment `json:"comments"`
}
