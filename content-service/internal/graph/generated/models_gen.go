package generated

type Menu struct {
	ID        uint           `json:"id"`
	CreatedAt string         `json:"created_at"`
	UpdatedAt string         `json:"updated_at"`
	Name      string         `json:"name"`
	Slug      string         `json:"slug"`
	Channel   string         `json:"channel"`
	Items     []*MenuItem    `json:"items"`
}

type MenuItem struct {
	ID           uint         `json:"id"`
	CreatedAt    string       `json:"created_at"`
	UpdatedAt    string       `json:"updated_at"`
	MenuID       uint         `json:"menu_id"`
	Name         string       `json:"name"`
	Level        int          `json:"level"`
	URL          string       `json:"url"`
	CategoryID   *uint        `json:"category_id"`
	Category     *Category    `json:"category,omitempty"`
	CollectionID *uint        `json:"collection_id"`
	Collection   *Collection  `json:"collection,omitempty"`
	PageID       *uint        `json:"page_id"`
	Page         *Page        `json:"page,omitempty"`
	ParentID     *uint        `json:"parent_id"`
	Children     []*MenuItem  `json:"children"`
}

type Category struct {
	ID          uint   `json:"id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
}

type Collection struct {
	ID          uint   `json:"id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
}

type Page struct {
	ID          uint   `json:"id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	Title       string `json:"title"`
	Slug        string `json:"slug"`
	Content     string `json:"content"`
	Description string `json:"description"`
} 