package responsetypes

type ProjectLandingPageInfo struct {
	Attributes    map[string]interface{} `json:"attributes,omitempty"`
	BoardID       int                    `json:"boardId,omitempty"`
	BoardName     string                 `json:"boardName,omitempty"`
	ProjectKey    string                 `json:"projectKey,omitempty"`
	ProjectType   string                 `json:"projectType,omitempty"`
	QueueCategory string                 `json:"queueCategory,omitempty"`
	QueueID       int                    `json:"queueId,omitempty"`
	QueueName     string                 `json:"queueName,omitempty"`
	SimpleBoard   bool                   `json:"simpleBoard,omitempty"`
	Simplified    bool                   `json:"simplified,omitempty"`
	Url           string                 `json:"url,omitempty"`
}
