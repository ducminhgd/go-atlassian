package responsetypes

type Transition struct {
	Expand        string                 `json:"expand,omitempty"`
	Fields        map[string]interface{} `json:"fields,omitempty"`
	HasScreen     bool                   `json:"hasScreen,omitempty"`
	ID            string                 `json:"id,omitempty"`
	IsAvailable   bool                   `json:"isAvailable,omitempty"`
	IsConditional bool                   `json:"isConditional,omitempty"`
	IsGlobal      bool                   `json:"isGlobal,omitempty"`
	IsInitial     bool                   `json:"isInitial,omitempty"`
	Looped        bool                   `json:"looped,omitempty"`
	Name          string                 `json:"name,omitempty"`
	To            StatusDetails          `json:"to,omitempty"`
}
