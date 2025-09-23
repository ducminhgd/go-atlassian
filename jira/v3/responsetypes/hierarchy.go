package responsetypes

// Hierarchy represents the hierarchy of issue types in a project
type Hierarchy struct {
	BaseLevelID int              `json:"baseLevelId,omitempty"`
	Levels      []HierarchyLevel `json:"levels,omitempty"`
}

// HierarchyLevel represents a single level in the issue type hierarchy
type HierarchyLevel struct {
	AboveLevelID         int    `json:"aboveLevelId,omitempty"`
	BelowLevelID         int    `json:"belowLevelId,omitempty"`
	ExternalUUID         string `json:"externalUuid,omitempty"`
	HierarchyLevelNumber int    `json:"hierarchyLevelNumber,omitempty"`
	ID                   int    `json:"id,omitempty"`
	IssueTypeIDs         []int  `json:"issueTypeIds,omitempty"`
	Level                int    `json:"level,omitempty"`
	Name                 string `json:"name,omitempty"`
}
