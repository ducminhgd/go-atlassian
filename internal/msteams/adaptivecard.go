package msteams

// AdaptiveCardElement represents a generic element in an AdaptiveCard
type AdaptiveCardElement struct {
	Type      string                `json:"type"`
	Text      string                `json:"text,omitempty"`
	Size      string                `json:"size,omitempty"`
	Weight    string                `json:"weight,omitempty"`
	Color     string                `json:"color,omitempty"`
	Wrap      bool                  `json:"wrap,omitempty"`
	Spacing   string                `json:"spacing,omitempty"`
	Separator bool                  `json:"separator,omitempty"`
	Items     []AdaptiveCardElement `json:"items,omitempty"`
	Columns   []AdaptiveCardElement `json:"columns,omitempty"`
	Width     any                   `json:"width,omitempty"` // can be string or int
	Actions   []AdaptiveCardAction  `json:"actions,omitempty"`
	URL       string                `json:"url,omitempty"`
	Title     string                `json:"title,omitempty"`
	Style     string                `json:"style,omitempty"`
	Inlines   []AdaptiveCardInline  `json:"inlines,omitempty"`
}

// AdaptiveCardAction represents an action in an AdaptiveCard
type AdaptiveCardAction struct {
	Type  string `json:"type"`
	URL   string `json:"url,omitempty"`
	Title string `json:"title,omitempty"`
}

// AdaptiveCardInline represents inline content in RichTextBlock
type AdaptiveCardInline struct {
	Type   string `json:"type"`
	Text   string `json:"text,omitempty"`
	URL    string `json:"url,omitempty"`
	Weight string `json:"weight,omitempty"`
	Color  string `json:"color,omitempty"`
}

// AdaptiveCard represents the full AdaptiveCard structure
type AdaptiveCard struct {
	Type    string                `json:"type"`
	Schema  string                `json:"$schema"`
	Version string                `json:"version"`
	Body    []AdaptiveCardElement `json:"body"`
	MSTeams *MSTeamsProperties    `json:"msteams,omitempty"`
}

// MSTeamsProperties contains Teams-specific properties
type MSTeamsProperties struct {
	Width string `json:"width"`
}

// NewAdaptiveCard creates a new AdaptiveCard with default settings for Teams
func NewAdaptiveCard() AdaptiveCard {
	return AdaptiveCard{
		Type:    "AdaptiveCard",
		Schema:  "http://adaptivecards.io/schemas/adaptive-card.json",
		Version: "1.5",
		Body:    []AdaptiveCardElement{},
		MSTeams: &MSTeamsProperties{
			Width: "Full",
		},
	}
}

// AddTextBlock adds a text block to the AdaptiveCard
func (ac *AdaptiveCard) AddTextBlock(text, size, weight string, wrap bool) {
	element := AdaptiveCardElement{
		Type:   "TextBlock",
		Text:   text,
		Wrap:   wrap,
	}
	
	if size != "" {
		element.Size = size
	}
	
	if weight != "" {
		element.Weight = weight
	}
	
	ac.Body = append(ac.Body, element)
}

// AddRichTextBlock adds a rich text block with inline elements
func (ac *AdaptiveCard) AddRichTextBlock(inlines []AdaptiveCardInline, spacing string, separator bool) {
	element := AdaptiveCardElement{
		Type:      "RichTextBlock",
		Inlines:   inlines,
		Spacing:   spacing,
		Separator: separator,
	}
	
	ac.Body = append(ac.Body, element)
}

// AddContainer adds a container with child items
func (ac *AdaptiveCard) AddContainer(items []AdaptiveCardElement, spacing, style string) {
	element := AdaptiveCardElement{
		Type:    "Container",
		Items:   items,
		Spacing: spacing,
		Style:   style,
	}
	
	ac.Body = append(ac.Body, element)
}

// CreateTextRun creates a text run for use in RichTextBlock
func CreateTextRun(text, weight, color string) AdaptiveCardInline {
	inline := AdaptiveCardInline{
		Type: "TextRun",
		Text: text,
	}
	
	if weight != "" {
		inline.Weight = weight
	}
	
	if color != "" {
		inline.Color = color
	}
	
	return inline
}

// CreateOpenUrlAction creates an action to open a URL
func CreateOpenUrlAction(title, url string) AdaptiveCardAction {
	return AdaptiveCardAction{
		Type:  "Action.OpenUrl",
		Title: title,
		URL:   url,
	}
}
