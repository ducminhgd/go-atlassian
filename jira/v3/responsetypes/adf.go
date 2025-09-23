package responsetypes

// AtlassianDocumentFormat represents the Atlassian Document Format (ADF) structure
type AtlassianDocumentFormat struct {
	Type    string         `json:"type"`
	Version int            `json:"version"`
	Content []DocumentNode `json:"content"`
}

// DocumentNode represents a node in the Atlassian Document Format
type DocumentNode struct {
	Type    string        `json:"type"`
	Content []NodeContent `json:"content,omitempty"`
	Marks   []Mark        `json:"marks,omitempty"`
	Attrs   *NodeAttrs    `json:"attrs,omitempty"`
	Text    string        `json:"text,omitempty"`
}

// NodeContent represents the content of a node in the document
type NodeContent struct {
	Type    string        `json:"type"`
	Content []NodeContent `json:"content,omitempty"`
	Attrs   *NodeAttrs    `json:"attrs,omitempty"`
	Text    string        `json:"text,omitempty"`
	Marks   []Mark        `json:"marks,omitempty"`
}

// Mark represents formatting marks applied to text
type Mark struct {
	Type  string     `json:"type"`
	Attrs *MarkAttrs `json:"attrs,omitempty"`
}

// MarkAttrs represents attributes that can be applied to marks
type MarkAttrs struct {
	Href        string `json:"href,omitempty"`
	Title       string `json:"title,omitempty"`
	Level       int    `json:"level,omitempty"`
	Language    string `json:"language,omitempty"`
	ID          string `json:"id,omitempty"`
	Collection  string `json:"collection,omitempty"`
	Color       string `json:"color,omitempty"`
	Type        string `json:"type,omitempty"`
	Align       string `json:"align,omitempty"`
	LocalID     string `json:"localId,omitempty"`
	AccessLevel string `json:"accessLevel,omitempty"`
}

// NodeAttrs represents attributes that can be applied to nodes
type NodeAttrs struct {
	ID            string                 `json:"id,omitempty"`
	Text          string                 `json:"text,omitempty"`
	Level         int                    `json:"level,omitempty"`
	Language      string                 `json:"language,omitempty"`
	LocalID       string                 `json:"localId,omitempty"`
	Layout        string                 `json:"layout,omitempty"`
	Width         int                    `json:"width,omitempty"`
	PanelType     string                 `json:"panelType,omitempty"`
	ExtensionType string                 `json:"extensionType,omitempty"`
	ExtensionKey  string                 `json:"extensionKey,omitempty"`
	Parameters    map[string]interface{} `json:"parameters,omitempty"`
	AccessLevel   string                 `json:"accessLevel,omitempty"`
	Collection    string                 `json:"collection,omitempty"`
	Timestamp     string                 `json:"timestamp,omitempty"`
	ShortName     string                 `json:"shortName,omitempty"`
	State         string                 `json:"state,omitempty"`
	Title         string                 `json:"title,omitempty"`
	URL           string                 `json:"url,omitempty"`
}

// Common ADF node types constants
const (
	NodeTypeDoc         = "doc"
	NodeTypeParagraph   = "paragraph"
	NodeTypeText        = "text"
	NodeTypeHeading     = "heading"
	NodeTypeBlockquote  = "blockquote"
	NodeTypeBulletList  = "bulletList"
	NodeTypeOrderedList = "orderedList"
	NodeTypeListItem    = "listItem"
	NodeTypeCodeBlock   = "codeBlock"
	NodeTypeMediaSingle = "mediaSingle"
	NodeTypeMedia       = "media"
	NodeTypeHardBreak   = "hardBreak"
	NodeTypeMention     = "mention"
	NodeTypeEmoji       = "emoji"
	NodeTypeStatus      = "status"
	NodeTypeTable       = "table"
	NodeTypeTableRow    = "tableRow"
	NodeTypeTableCell   = "tableCell"
	NodeTypeTableHeader = "tableHeader"
)

// Common mark types constants
const (
	MarkTypeStrong      = "strong"
	MarkTypeEm          = "em"
	MarkTypeStrike      = "strike"
	MarkTypeUnderline   = "underline"
	MarkTypeCode        = "code"
	MarkTypeLink        = "link"
	MarkTypeTextColor   = "textColor"
	MarkTypeAlignment   = "alignment"
	MarkTypeIndentation = "indentation"
)
