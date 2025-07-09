package notion

import "github.com/bytedance/sonic"

type UpdatePageRequest struct {
	// Map where keys are property names or IDs and values are Property Value Objects
	// structured according to the Notion API documentation for updating properties.
	// Example: {"My Select Property": {"select": {"id": "option_id"}}}
	// Use interface{} for flexibility as update structures differ per property type.
	Properties map[string]interface{} `json:"properties,omitempty"`
	Archived   *bool                  `json:"archived,omitempty"`
	Icon       *Icon                  `json:"icon,omitempty"`  // Reuse Icon type (Emoji, External, File)
	Cover      *Cover                 `json:"cover,omitempty"` // Reuse Cover type (External, File)
}

// CreatePageRequest defines the body for POST /v1/pages
type CreatePageRequest struct {
	Parent Parent `json:"parent"` // Required: { "page_id": "..." } or { "database_id": "..." }

	// Page properties. Structure depends on whether parent is page or database.
	// For database parent, keys MUST match database schema properties.
	// Must include the database's title property.
	// For page parent, only 'title' is valid: {"title": {"title": [{"type": "text", "text": {"content": "..."}}]}}
	Properties map[string]interface{} `json:"properties"` // Required

	// Optional: Array of block objects to create as children.
	// Use map[string]interface{} as block structure for creation is simpler than the full response Block struct.
	Children *[]map[string]interface{} `json:"children,omitempty"`
	Icon     *Icon                     `json:"icon,omitempty"`  // Reuse Icon type
	Cover    *Cover                    `json:"cover,omitempty"` // Reuse Cover type
}

// SearchSort defines the sorting object for search requests.
type SearchSort struct {
	Direction string `json:"direction"` // "ascending" or "descending"
	Timestamp string `json:"timestamp"` // "last_edited_time"
}

// SearchFilter defines the filtering object for search requests.
type SearchFilter struct {
	Value    string `json:"value"`    // "page" or "database"
	Property string `json:"property"` // "object"
}

// SearchRequest defines the body for POST /v1/search
type SearchRequest struct {
	Query       *string       `json:"query,omitempty"`
	Sort        *SearchSort   `json:"sort,omitempty"`
	Filter      *SearchFilter `json:"filter,omitempty"`
	StartCursor *string       `json:"start_cursor,omitempty"`
	PageSize    *int          `json:"page_size,omitempty"` // Notion API uses integer
}

// PropertySchemaCreateObject defines the structure needed when *creating* a database property.
// Note: This is subtly different from the response PropertySchemaObject.
type PropertySchemaCreateObject struct {
	// Top-level key IS the type ("title", "rich_text", "select", etc.)
	// The value is the configuration object for that type.
	// Use map[string]interface{} for maximum flexibility when constructing the request.
	// Example for a Select property:
	// "My Select": { "select": { "options": [ {"name": "A", "color": "red"} ] } }
	// Example for Title:
	// "Name": { "title": {} }
	// Type assertion/construction happens before marshaling this.
	TypeConfig map[string]interface{} `json:"-"` // This field is conceptual, the map is built dynamically
}

// CreateDatabaseRequest defines the body for POST /v1/databases
type CreateDatabaseRequest struct {
	Parent Parent `json:"parent"` // Required: { "page_id": "..." } or { "workspace": true }

	Title       []RichText  `json:"title"` // Required: Title rich text array
	Icon        *Icon       `json:"icon,omitempty"`
	Cover       *Cover      `json:"cover,omitempty"`
	IsInline    *bool       `json:"is_inline,omitempty"`
	Description *[]RichText `json:"description,omitempty"` // Added description based on markdown
	// Properties defines the schema. Key is the desired property name.
	// Value is the property schema definition (e.g., {"select": { "options": [...] }} )
	Properties map[string]interface{} `json:"properties"` // Required. Use interface{} for flexibility.
}

// QueryDatabaseRequest defines the body for POST /v1/databases/{database_id}/query
type QueryDatabaseRequest struct {
	// Notion filter object structure. See Notion API docs.
	Filter map[string]interface{} `json:"filter,omitempty"`
	// Array of Notion sort objects. See Notion API docs.
	Sorts       []map[string]interface{} `json:"sorts,omitempty"`
	StartCursor *string                  `json:"start_cursor,omitempty"`
	PageSize    *int                     `json:"page_size,omitempty"` // Notion API uses integer
}

// AppendBlockChildrenRequest defines the body for PATCH /v1/blocks/{block_id}/children
type AppendBlockChildrenRequest struct {
	// Array of block objects to append. Use map for flexibility in creation.
	Children []map[string]interface{} `json:"children"`        // Required
	After    *string                  `json:"after,omitempty"` // Added based on API docs
}

// UpdateBlockRequest defines the body for PATCH /v1/blocks/{block_id}
type UpdateBlockRequest struct {
	// A map containing *only one* key, which is the block type string (e.g., "paragraph").
	// The value is an object containing the fields to update for that block type.
	// Example: {"paragraph": {"rich_text": [...], "color": "default"}}
	BlockUpdateContent map[string]interface{} `json:"-"` // Use map, key is block type, value is update object

	// Optional field to archive/unarchive the block.
	Archived *bool `json:"archived,omitempty"`
}

// CreateCommentRequest defines the body for POST /v1/comments
type CreateCommentRequest struct {
	// Provide EITHER Parent OR DiscussionID
	Parent       *Parent    `json:"parent,omitempty"` // e.g., { "page_id": "..." }
	DiscussionID *string    `json:"discussion_id,omitempty"`
	RichText     []RichText `json:"rich_text"` // Required
}

// UpdateDatabaseRequest defines the body for PATCH /v1/databases/{database_id}
type UpdateDatabaseRequest struct {
	Title       *[]RichText `json:"title,omitempty"`
	Description *[]RichText `json:"description,omitempty"`
	// Used to update the database property schema. Key is property name or ID.
	// Value contains the updated schema definition (e.g., updating select options).
	Properties map[string]interface{} `json:"properties,omitempty"`
	Icon       *Icon                  `json:"icon,omitempty"`
	Cover      *Cover                 `json:"cover,omitempty"`
	Archived   *bool                  `json:"archived,omitempty"` // Added archived based on markdown
}

// Helper function to marshal UpdateBlockRequest correctly
func (r *UpdateBlockRequest) MarshalJSON() ([]byte, error) {
	// Create a temporary map
	tempMap := make(map[string]interface{})

	// Copy the type-specific update content
	for key, value := range r.BlockUpdateContent {
		tempMap[key] = value
	}

	// Add the archived field if it's not nil
	if r.Archived != nil {
		tempMap["archived"] = r.Archived
	}

	// Marshal the temporary map
	return sonic.Marshal(tempMap)
}

// PartialUser represents a user object containing only the ID.
type PartialUser struct {
	Object string `json:"object"` // Should be "user"
	ID     string `json:"id"`     // User UUID
}

// User represents a Notion user (person or bot).
type User struct {
	Object        string     `json:"object"` // "user"
	ID            string     `json:"id"`
	Name          *string    `json:"name,omitempty"`       // Nullable
	AvatarURL     *string    `json:"avatar_url,omitempty"` // Nullable
	Type          string     `json:"type"`                 // "person" or "bot"
	Person        *Person    `json:"person,omitempty"`
	Bot           *Bot       `json:"bot,omitempty"`
	Owner         *UserOwner `json:"owner,omitempty"`          // Only for bot owner info
	WorkspaceName *string    `json:"workspace_name,omitempty"` // Only for bot
}

// Person details for a person user type.
type Person struct {
	Email string `json:"email"` // Note: Requires specific capabilities
}

// Bot details for a bot user type.
type Bot struct {
	Owner         *UserOwner `json:"owner"` // Owner details
	WorkspaceName *string    `json:"workspace_name"`
}

// UserOwner describes the owner of a bot.
type UserOwner struct {
	Type      string       `json:"type"` // "user" or "workspace"
	User      *PartialUser `json:"user,omitempty"`
	Workspace *bool        `json:"workspace,omitempty"` // Likely always true when type is "workspace"
}

// RichTextObject is the base interface for different rich text types.
// We'll use a struct containing all possible types for simplicity, similar to Block.
type RichText struct {
	Type        string       `json:"type"` // "text", "mention", "equation"
	Text        *TextContent `json:"text,omitempty"`
	Mention     *Mention     `json:"mention,omitempty"`
	Equation    *Equation    `json:"equation,omitempty"`
	Annotations *Annotations `json:"annotations,omitempty"`
	PlainText   string       `json:"plain_text,omitempty"`
	Href        *string      `json:"href,omitempty"` // Nullable
}

// TextContent holds the content for a 'text' rich text object.
type TextContent struct {
	Content string `json:"content"`
	Link    *Link  `json:"link,omitempty"` // Nullable
}

// Link object within text content.
type Link struct {
	Type string  `json:"type"`          // Deprecated, use Href on RichText instead
	URL  *string `json:"url,omitempty"` // Use Href on RichText instead
}

// Mention object within rich text.
type Mention struct {
	Type        string           `json:"type"` // "user", "page", "database", "date", "link_preview", "template_mention"
	User        *PartialUser     `json:"user,omitempty"`
	Page        *PageMention     `json:"page,omitempty"`
	Database    *DatabaseMention `json:"database,omitempty"`
	Date        *DateValue       `json:"date,omitempty"` // Re-use DateValue from properties
	LinkPreview *LinkPreview     `json:"link_preview,omitempty"`
	// TemplateMention fields omitted for brevity, add if needed
}

type PageMention struct {
	ID string `json:"id"` // Page UUID
}

type DatabaseMention struct {
	ID string `json:"id"` // Database UUID
}

type LinkPreview struct {
	URL string `json:"url"`
}

// Equation object within rich text.
type Equation struct {
	Expression string `json:"expression"`
}

// Annotations provide styling information for rich text.
type Annotations struct {
	Bold          bool   `json:"bold"`
	Italic        bool   `json:"italic"`
	Strikethrough bool   `json:"strikethrough"`
	Underline     bool   `json:"underline"`
	Code          bool   `json:"code"`
	Color         string `json:"color"` // e.g., "default", "gray", "red_background"
}

// FileObject represents a file hosted by Notion.
type FileObject struct {
	Type string `json:"type"` // "file"
	File *File  `json:"file"`
	Name string `json:"name,omitempty"` // Used in File property values
}

// File details for a Notion-hosted file.
type File struct {
	URL        string `json:"url"`
	ExpiryTime string `json:"expiry_time"` // ISO 8601
}

// ExternalFileObject represents an externally hosted file.
type ExternalFileObject struct {
	Type     string     `json:"type"` // "external"
	External *External  `json:"external"`
	Caption  []RichText `json:"caption,omitempty"` // Used in Image blocks
}

// External file details.
type External struct {
	URL string `json:"url"`
}

// EmojiObject represents an emoji character used as an icon.
type EmojiObject struct {
	Type  string `json:"type"` // "emoji"
	Emoji string `json:"emoji"`
}

// Icon represents a page or database icon. Can be Emoji, External, or File.
type Icon struct {
	Type     string    `json:"type"` // "emoji", "external", or "file"
	Emoji    *string   `json:"emoji,omitempty"`
	External *External `json:"external,omitempty"`
	File     *File     `json:"file,omitempty"`
}

// Cover represents a page or database cover image. Can be External or File.
type Cover struct {
	Type     string    `json:"type"` // "external" or "file"
	External *External `json:"external,omitempty"`
	File     *File     `json:"file,omitempty"`
}

// Parent object indicating the location of a page, database, or block.
type Parent struct {
	Type       string `json:"type"` // "database_id", "page_id", "workspace", "block_id"
	DatabaseID string `json:"database_id,omitempty"`
	PageID     string `json:"page_id,omitempty"`
	Workspace  *bool  `json:"workspace,omitempty"`
	BlockID    string `json:"block_id,omitempty"` // Added for blocks
}

// SelectOption represents an option in a Select or MultiSelect property.
type SelectOption struct {
	ID    *string `json:"id,omitempty"` // Optional ID assigned by Notion
	Name  string  `json:"name"`
	Color string  `json:"color"` // e.g., "blue", "red"
}

// DateValue represents a date property value, which might include time and timezone.
type DateValue struct {
	Start    string  `json:"start"`               // ISO 8601 Date or DateTime
	End      *string `json:"end,omitempty"`       // ISO 8601 Date or DateTime (Nullable)
	TimeZone *string `json:"time_zone,omitempty"` // Nullable
}

// FormulaValue holds the result of a formula calculation.
type FormulaValue struct {
	Type    string     `json:"type"` // "string", "number", "boolean", "date"
	String  *string    `json:"string,omitempty"`
	Number  *float64   `json:"number,omitempty"`
	Boolean *bool      `json:"boolean,omitempty"`
	Date    *DateValue `json:"date,omitempty"`
}

// RelationValue represents a relation property value (an array of page references).
type RelationValue struct {
	ID string `json:"id"` // The ID of the related page
}

// RollupValue holds the result of a rollup calculation.
type RollupValue struct {
	Type        string         `json:"type"` // "number", "date", "array"
	Number      *float64       `json:"number,omitempty"`
	Date        *DateValue     `json:"date,omitempty"`
	Array       *[]interface{} `json:"array,omitempty"` // Array of property values (structure depends on function)
	Function    string         `json:"function"`
	Incomplete  *string        `json:"incomplete,omitempty"`  // Present if results are partial
	Unsupported *string        `json:"unsupported,omitempty"` // Present if type is unsupported
}

// PropertyValueObject holds the value for a specific property type on a Page.
// We use pointers to structs for each type to handle the polymorphism.
type PropertyValueObject struct {
	ID   string `json:"id"`   // Property ID
	Type string `json:"type"` // "title", "rich_text", "number", etc.

	Title          *[]RichText                `json:"title,omitempty"`
	RichText       *[]RichText                `json:"rich_text,omitempty"`
	Number         *float64                   `json:"number,omitempty"`
	Select         *SelectOption              `json:"select,omitempty"` // Single select option
	MultiSelect    *[]SelectOption            `json:"multi_select,omitempty"`
	Date           *DateValue                 `json:"date,omitempty"`
	People         *[]PartialUser             `json:"people,omitempty"`
	Files          *[]FileObject              `json:"files,omitempty"` // Can contain File or External files
	Checkbox       *bool                      `json:"checkbox,omitempty"`
	URL            *string                    `json:"url,omitempty"`
	Email          *string                    `json:"email,omitempty"`
	PhoneNumber    *string                    `json:"phone_number,omitempty"`
	Formula        *FormulaValue              `json:"formula,omitempty"`
	Relation       *[]RelationValue           `json:"relation,omitempty"`
	Rollup         *RollupValue               `json:"rollup,omitempty"`
	CreatedTime    *string                    `json:"created_time,omitempty"` // ISO 8601
	CreatedBy      *PartialUser               `json:"created_by,omitempty"`
	LastEditedTime *string                    `json:"last_edited_time,omitempty"` // ISO 8601
	LastEditedBy   *PartialUser               `json:"last_edited_by,omitempty"`
	Status         *SelectOption              `json:"status,omitempty"`       // Status is similar to Select
	UniqueID       *UniqueIDPropertyValue     `json:"unique_id,omitempty"`    // Added UniqueID
	Verification   *VerificationPropertyValue `json:"verification,omitempty"` // Added Verification
}

// UniqueIDPropertyValue specifically for the 'unique_id' property type.
type UniqueIDPropertyValue struct {
	Prefix *string `json:"prefix"` // Nullable
	Number int     `json:"number"`
}

// VerificationPropertyValue specifically for the 'verification' property type.
type VerificationPropertyValue struct {
	State      *string    `json:"state,omitempty"`
	VerifiedBy *User      `json:"verified_by,omitempty"` // Can be null
	Date       *DateValue `json:"date,omitempty"`        // Can be null
}

// Page represents a Notion page object.
type Page struct {
	Object         string                         `json:"object"` // "page"
	ID             string                         `json:"id"`
	CreatedTime    string                         `json:"created_time"`
	LastEditedTime string                         `json:"last_edited_time"`
	CreatedBy      PartialUser                    `json:"created_by"`
	LastEditedBy   PartialUser                    `json:"last_edited_by"`
	Cover          *Cover                         `json:"cover"` // Nullable
	Icon           *Icon                          `json:"icon"`  // Nullable
	Parent         Parent                         `json:"parent"`
	Archived       bool                           `json:"archived"`
	Properties     map[string]PropertyValueObject `json:"properties"` // Key is property name/id
	URL            string                         `json:"url"`
	PublicURL      *string                        `json:"public_url,omitempty"` // Nullable, only if public
}

// PropertySchemaObject defines the schema (type and settings) for a Database property.
type PropertySchemaObject struct {
	ID   string `json:"id"`   // Property ID
	Name string `json:"name"` // Property Name
	Type string `json:"type"` // "title", "rich_text", "number", etc.

	Title          map[string]interface{} `json:"title,omitempty"`     // CHANGE: Often {}, use map
	RichText       map[string]interface{} `json:"rich_text,omitempty"` // CHANGE: Often {}, use map
	Number         *NumberSchema          `json:"number,omitempty"`
	Select         *SelectSchema          `json:"select,omitempty"`
	MultiSelect    *SelectSchema          `json:"multi_select,omitempty"` // Uses same options structure as Select
	Date           map[string]interface{} `json:"date,omitempty"`         // CHANGE: Often {}, use map
	People         map[string]interface{} `json:"people,omitempty"`       // CHANGE: Often {}, use map
	Files          map[string]interface{} `json:"files,omitempty"`        // CHANGE: Often {}, use map
	Checkbox       map[string]interface{} `json:"checkbox,omitempty"`     // CHANGE: Often {}, use map
	URL            map[string]interface{} `json:"url,omitempty"`          // CHANGE: Often {}, use map
	Email          map[string]interface{} `json:"email,omitempty"`        // CHANGE: Often {}, use map
	PhoneNumber    map[string]interface{} `json:"phone_number,omitempty"` // CHANGE: Often {}, use map
	Formula        *FormulaSchema         `json:"formula,omitempty"`
	Relation       *RelationSchema        `json:"relation,omitempty"`
	Rollup         *RollupSchema          `json:"rollup,omitempty"`
	CreatedTime    map[string]interface{} `json:"created_time,omitempty"`     // CHANGE: Often {}, use map
	CreatedBy      map[string]interface{} `json:"created_by,omitempty"`       // CHANGE: Often {}, use map
	LastEditedTime map[string]interface{} `json:"last_edited_time,omitempty"` // CHANGE: Often {}, use map
	LastEditedBy   map[string]interface{} `json:"last_edited_by,omitempty"`   // CHANGE: Often {}, use map
	Status         *SelectSchema          `json:"status,omitempty"`           // Status schema is like Select
	UniqueID       map[string]interface{} `json:"unique_id,omitempty"`        // CHANGE: Added UniqueID Schema {}, use map
	Verification   map[string]interface{} `json:"verification,omitempty"`     // CHANGE: Added Verification Schema {}, use map
}

// NumberSchema defines options for a Number property.
type NumberSchema struct {
	Format string `json:"format"` // "number", "number_with_commas", "percent", "dollar", etc.
}

// SelectSchema defines options for Select or MultiSelect properties.
type SelectSchema struct {
	Options []SelectOption `json:"options"`
}

// FormulaSchema defines the expression for a Formula property.
type FormulaSchema struct {
	Expression string `json:"expression"`
}

// RelationSchema defines options for a Relation property.
type RelationSchema struct {
	DatabaseID     string                `json:"database_id"`
	Type           string                `json:"type"` // "single_property" or "dual_property"
	SingleProperty *SinglePropertySchema `json:"single_property,omitempty"`
	DualProperty   *DualPropertySchema   `json:"dual_property,omitempty"`
}

// SinglePropertySchema is currently empty for Relation creation/retrieval.
type SinglePropertySchema struct {
	// Empty object {}
}

// DualPropertySchema defines the related property in a two-way relation.
type DualPropertySchema struct {
	SyncedPropertyName string `json:"synced_property_name"`
	SyncedPropertyID   string `json:"synced_property_id"`
}

// RollupSchema defines options for a Rollup property.
type RollupSchema struct {
	RelationPropertyName string `json:"relation_property_name"`
	RelationPropertyID   string `json:"relation_property_id"`
	RollupPropertyName   string `json:"rollup_property_name"`
	RollupPropertyID     string `json:"rollup_property_id"`
	Function             string `json:"function"` // e.g., "count", "sum", "average"
}

// Database represents a Notion database object.
type Database struct {
	Object         string                          `json:"object"` // "database"
	ID             string                          `json:"id"`
	CreatedTime    string                          `json:"created_time"`
	LastEditedTime string                          `json:"last_edited_time"`
	CreatedBy      PartialUser                     `json:"created_by"`
	LastEditedBy   PartialUser                     `json:"last_edited_by"`
	Title          []RichText                      `json:"title"`
	Description    []RichText                      `json:"description"`
	Icon           *Icon                           `json:"icon"`       // Nullable
	Cover          *Cover                          `json:"cover"`      // Nullable
	Properties     map[string]PropertySchemaObject `json:"properties"` // Key is property name
	Parent         Parent                          `json:"parent"`
	URL            string                          `json:"url"`
	Archived       bool                            `json:"archived"`
	IsInline       bool                            `json:"is_inline"`
	PublicURL      *string                         `json:"public_url,omitempty"` // Nullable
}

// Block represents a generic Notion block object.
// Type-specific fields contain the actual content.
type Block struct {
	Object         string      `json:"object"` // "block"
	ID             string      `json:"id"`
	Parent         Parent      `json:"parent"`
	CreatedTime    string      `json:"created_time"`
	LastEditedTime string      `json:"last_edited_time"`
	CreatedBy      PartialUser `json:"created_by"`
	LastEditedBy   PartialUser `json:"last_edited_by"`
	HasChildren    bool        `json:"has_children"`
	Archived       bool        `json:"archived"`
	Type           string      `json:"type"` // "paragraph", "heading_1", "to_do", etc.

	Paragraph        *ParagraphBlock        `json:"paragraph,omitempty"`
	Heading1         *HeadingBlock          `json:"heading_1,omitempty"`
	Heading2         *HeadingBlock          `json:"heading_2,omitempty"`
	Heading3         *HeadingBlock          `json:"heading_3,omitempty"`
	BulletedListItem *ListItemBlock         `json:"bulleted_list_item,omitempty"`
	NumberedListItem *ListItemBlock         `json:"numbered_list_item,omitempty"`
	ToDo             *ToDoBlock             `json:"to_do,omitempty"`
	Toggle           *ToggleBlock           `json:"toggle,omitempty"`
	ChildPage        *ChildPageBlock        `json:"child_page,omitempty"`
	ChildDatabase    *ChildDatabaseBlock    `json:"child_database,omitempty"`
	Embed            *EmbedBlock            `json:"embed,omitempty"`
	Image            *FileBlock             `json:"image,omitempty"` // Uses FileBlock structure
	Video            *FileBlock             `json:"video,omitempty"`
	File             *FileBlock             `json:"file,omitempty"`
	PDF              *FileBlock             `json:"pdf,omitempty"`
	Bookmark         *BookmarkBlock         `json:"bookmark,omitempty"`
	Callout          *CalloutBlock          `json:"callout,omitempty"`
	Quote            *QuoteBlock            `json:"quote,omitempty"`
	Equation         *EquationBlock         `json:"equation,omitempty"`
	Divider          map[string]interface{} `json:"divider,omitempty"` // CHANGE: Often {}, use map for robustness with mapstructure
	TableOfContents  *TableOfContentsBlock  `json:"table_of_contents,omitempty"`
	Code             *CodeBlock             `json:"code,omitempty"`
	Table            *TableBlock            `json:"table,omitempty"`
	TableRow         *TableRowBlock         `json:"table_row,omitempty"`
	ColumnList       map[string]interface{} `json:"column_list,omitempty"` // CHANGE: No unique fields, use map
	Column           map[string]interface{} `json:"column,omitempty"`      // CHANGE: No unique fields, use map
	LinkToPage       *LinkToPageBlock       `json:"link_to_page,omitempty"`
	SyncedBlock      *SyncedBlock           `json:"synced_block,omitempty"`
	Template         *TemplateBlock         `json:"template,omitempty"`
	LinkPreview      *LinkPreviewBlock      `json:"link_preview,omitempty"`
	// ... other block types ...
}

type ParagraphBlock struct {
	RichText []RichText `json:"rich_text"`
	Color    string     `json:"color"`
	Children *[]Block   `json:"children,omitempty"` // Sometimes blocks can have children
}

type HeadingBlock struct { // Used for heading_1, heading_2, heading_3
	RichText     []RichText `json:"rich_text"`
	Color        string     `json:"color"`
	IsToggleable bool       `json:"is_toggleable"`
	Children     *[]Block   `json:"children,omitempty"` // Only if toggleable
}

type ListItemBlock struct { // Used for bulleted_list_item, numbered_list_item
	RichText []RichText `json:"rich_text"`
	Color    string     `json:"color"`
	Children *[]Block   `json:"children,omitempty"`
}

type ToDoBlock struct {
	RichText []RichText `json:"rich_text"`
	Checked  *bool      `json:"checked"` // Use pointer for explicit false vs. not set
	Color    string     `json:"color"`
	Children *[]Block   `json:"children,omitempty"`
}

type ToggleBlock struct {
	RichText []RichText `json:"rich_text"`
	Color    string     `json:"color"`
	Children *[]Block   `json:"children,omitempty"` // Content inside the toggle
}

type ChildPageBlock struct {
	Title string `json:"title"`
}

type ChildDatabaseBlock struct {
	Title string `json:"title"`
}

type EmbedBlock struct {
	URL     string     `json:"url"`
	Caption []RichText `json:"caption"`
}

type FileBlock struct { // Used for image, video, file, pdf
	Type     string      `json:"type"` // "external" or "file"
	External *External   `json:"external,omitempty"`
	File     *FileObject `json:"file,omitempty"` // Note: API uses FileObject structure here
	Caption  []RichText  `json:"caption"`
	Name     *string     `json:"name,omitempty"` // For type: "file"
}

type BookmarkBlock struct {
	URL     string     `json:"url"`
	Caption []RichText `json:"caption"`
}

type CalloutBlock struct {
	RichText []RichText `json:"rich_text"`
	Icon     *Icon      `json:"icon"`
	Color    string     `json:"color"`
}

type QuoteBlock struct {
	RichText []RichText `json:"rich_text"`
	Color    string     `json:"color"`
	Children *[]Block   `json:"children,omitempty"`
}

type EquationBlock struct {
	Expression string `json:"expression"`
}

type TableOfContentsBlock struct {
	Color string `json:"color"`
}

type CodeBlock struct {
	Caption  []RichText `json:"caption"`
	RichText []RichText `json:"rich_text"`
	Language string     `json:"language"`
}

type TableBlock struct {
	TableWidth      int      `json:"table_width"`
	HasColumnHeader bool     `json:"has_column_header"`
	HasRowHeader    bool     `json:"has_row_header"`
	Children        *[]Block `json:"children,omitempty"` // Contains table_row blocks
}

type TableRowBlock struct {
	Cells [][]RichText `json:"cells"` // Array of cells, each cell is an array of rich text
}

type LinkToPageBlock struct {
	Type       string `json:"type"` // "page_id" or "database_id"
	PageID     string `json:"page_id,omitempty"`
	DatabaseID string `json:"database_id,omitempty"`
}

type SyncedBlock struct {
	SyncedFrom *SyncedFromInfo `json:"synced_from"`        // Null if this is the original block
	Children   *[]Block        `json:"children,omitempty"` // Content of the synced block
}

type SyncedFromInfo struct {
	Type    string `json:"type"` // "block_id"
	BlockID string `json:"block_id"`
}

type TemplateBlock struct {
	RichText []RichText `json:"rich_text"`
	Children *[]Block   `json:"children,omitempty"`
}

type LinkPreviewBlock struct {
	URL string `json:"url"`
}

// Comment represents a comment on a page or discussion.
type Comment struct {
	Object         string      `json:"object"` // "comment"
	ID             string      `json:"id"`
	Parent         Parent      `json:"parent"` // Page or block parent
	DiscussionID   string      `json:"discussion_id"`
	CreatedTime    string      `json:"created_time"`
	LastEditedTime string      `json:"last_edited_time"`
	CreatedBy      PartialUser `json:"created_by"`
	RichText       []RichText  `json:"rich_text"`
}

// PropertyItem represents the value of a single property when retrieved individually.
// Structure varies greatly depending on the property type.
type PropertyItem struct {
	Object string `json:"object"` // "property_item"
	Type   string `json:"type"`   // "title", "rich_text", "number", etc.
	ID     string `json:"id"`     // Only for 'relation' type usually?

	Title          *[]RichText                `json:"title,omitempty"`
	RichText       *[]RichText                `json:"rich_text,omitempty"`
	Number         *float64                   `json:"number,omitempty"`
	Select         *SelectOption              `json:"select,omitempty"`
	MultiSelect    *[]SelectOption            `json:"multi_select,omitempty"`
	Date           *DateValue                 `json:"date,omitempty"`
	People         *[]PartialUser             `json:"people,omitempty"`
	Files          *[]FileObject              `json:"files,omitempty"`
	Checkbox       *bool                      `json:"checkbox,omitempty"`
	URL            *string                    `json:"url,omitempty"`
	Email          *string                    `json:"email,omitempty"`
	PhoneNumber    *string                    `json:"phone_number,omitempty"`
	Formula        *FormulaValue              `json:"formula,omitempty"`
	Relation       *[]RelationValue           `json:"relation,omitempty"` // May be paginated within a ListResponse
	Rollup         *RollupValue               `json:"rollup,omitempty"`
	CreatedTime    *string                    `json:"created_time,omitempty"`
	CreatedBy      *PartialUser               `json:"created_by,omitempty"`
	LastEditedTime *string                    `json:"last_edited_time,omitempty"`
	LastEditedBy   *PartialUser               `json:"last_edited_by,omitempty"`
	Status         *SelectOption              `json:"status,omitempty"`
	UniqueID       *UniqueIDPropertyValue     `json:"unique_id,omitempty"`
	Verification   *VerificationPropertyValue `json:"verification,omitempty"`
}

// ListResponse is a generic structure for Notion API list responses.
type ListResponse struct {
	Object     string        `json:"object"`      // "list"
	Results    []interface{} `json:"results"`     // CHANGE: Array of objects (Page, Database, User, Block, Comment, PropertyItem) - Decoded as []interface{}
	NextCursor *string       `json:"next_cursor"` // Nullable
	HasMore    bool          `json:"has_more"`
	Type       string        `json:"type"` // "user", "page_or_database", "block", "comment", "property_item"
	// REMOVED Specific type fields as Type + Results provide the info
	// User           *json.RawMessage `json:"user,omitempty"`
	// PageOrDatabase *json.RawMessage `json:"page_or_database,omitempty"`
	// Block          *json.RawMessage `json:"block,omitempty"`
	// Comment        *json.RawMessage `json:"comment,omitempty"`
	// PropertyItem   *json.RawMessage `json:"property_item,omitempty"`
}

// You can define specific list types for better type safety if you know the expected result type.

type UserListResponse struct {
	Object     string   `json:"object"` // "list"
	Results    []User   `json:"results"`
	NextCursor *string  `json:"next_cursor"`
	HasMore    bool     `json:"has_more"`
	Type       string   `json:"type"` // "user"
	User       struct{} `json:"user"`
}

type SearchListResponse struct {
	Object         string        `json:"object"`  // "list"
	Results        []interface{} `json:"results"` // CHANGE: Contains Page or Database objects (map[string]interface{})
	NextCursor     *string       `json:"next_cursor"`
	HasMore        bool          `json:"has_more"`
	Type           string        `json:"type"` // "page_or_database"
	PageOrDatabase struct{}      `json:"page_or_database"`
}

type DatabaseQueryListResponse struct {
	Object         string   `json:"object"`  // "list"
	Results        []Page   `json:"results"` // Contains Page objects conforming to DB schema
	NextCursor     *string  `json:"next_cursor"`
	HasMore        bool     `json:"has_more"`
	Type           string   `json:"type"` // "page_or_database"
	PageOrDatabase struct{} `json:"page_or_database"`
}

type BlockListResponse struct {
	Object     string   `json:"object"`  // "list"
	Results    []Block  `json:"results"` // Contains Block objects
	NextCursor *string  `json:"next_cursor"`
	HasMore    bool     `json:"has_more"`
	Type       string   `json:"type"` // "block"
	Block      struct{} `json:"block"`
}

type CommentListResponse struct {
	Object     string    `json:"object"` // "list"
	Results    []Comment `json:"results"`
	NextCursor *string   `json:"next_cursor"`
	HasMore    bool      `json:"has_more"`
	Type       string    `json:"type"` // "comment"
	Comment    struct{}  `json:"comment"`
}

// PropertyItemListResponse used for retrieving paginated property items
type PropertyItemListResponse struct {
	Object       string         `json:"object"` // "list"
	Results      []PropertyItem `json:"results"`
	NextCursor   *string        `json:"next_cursor"`
	HasMore      bool           `json:"has_more"`
	Type         string         `json:"type"` // "property_item"
	PropertyItem struct{}       `json:"property_item"`
}
