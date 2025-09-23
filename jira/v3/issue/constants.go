package issue

const (
	// Issue API Group endpoints
	ISSUE_GET_ENDPOINT    = "/rest/api/3/issue/%s"
	ISSUE_CREATE_ENDPOINT = "/rest/api/3/issue"
	ISSUE_UPDATE_ENDPOINT = "/rest/api/3/issue/%s"
	ISSUE_DELETE_ENDPOINT = "/rest/api/3/issue/%s"

	// Issue Search API Group endpoints
	ISSUE_SEARCH_ENDPOINT     = "/rest/api/3/search"
	ISSUE_SEARCH_JQL_ENDPOINT = "/rest/api/3/search/jql"

	// Issue transitions
	ISSUE_TRANSITIONS_ENDPOINT = "/rest/api/3/issue/%s/transitions"

	// Issue comments
	ISSUE_COMMENTS_ENDPOINT       = "/rest/api/3/issue/%s/comment"
	ISSUE_COMMENT_DETAIL_ENDPOINT = "/rest/api/3/issue/%s/comment/%s"

	// Issue attachments
	ISSUE_ATTACHMENTS_ENDPOINT = "/rest/api/3/issue/%s/attachments"

	// Issue watchers
	ISSUE_WATCHERS_ENDPOINT = "/rest/api/3/issue/%s/watchers"

	// Issue votes
	ISSUE_VOTES_ENDPOINT = "/rest/api/3/issue/%s/votes"

	// Issue worklog
	ISSUE_WORKLOG_ENDPOINT        = "/rest/api/3/issue/%s/worklog"
	ISSUE_WORKLOG_DETAIL_ENDPOINT = "/rest/api/3/issue/%s/worklog/%s"

	// Issue links
	ISSUE_LINKS_ENDPOINT = "/rest/api/3/issueLink"

	// Issue types
	ISSUE_TYPES_ENDPOINT = "/rest/api/3/issuetype"

	// Issue priorities
	ISSUE_PRIORITIES_ENDPOINT = "/rest/api/3/priority"

	// Issue statuses
	ISSUE_STATUSES_ENDPOINT = "/rest/api/3/status"

	// Issue fields
	ISSUE_FIELDS_ENDPOINT = "/rest/api/3/field"
)
