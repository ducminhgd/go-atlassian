package responsetypes

type Group struct {
	GroupID string `json:"groupId,omitempty"`
	Name    string `json:"name,omitempty"`
	Self    string `json:"self,omitempty"`
}

type SimpleListWrapperGroup struct {
	Items []Group `json:"items,omitempty"`
	Size  int     `json:"size,omitempty"`
}
