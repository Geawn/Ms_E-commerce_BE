package graphql

// Page represents pagination information
type Page struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}
