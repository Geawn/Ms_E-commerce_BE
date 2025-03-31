package model

type Country struct {
	Country string `json:"country"`
	Code    string `json:"code"`
}

type Channel struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Slug         string    `json:"slug"`
	IsActive     bool      `json:"isActive"`
	CurrencyCode string    `json:"currencyCode"`
	Countries    []Country `json:"countries"`
}
