package deps

import "database/sql"

type Client struct {
	baseURL string
	db      *sql.DB
}

func NewClient(db *sql.DB) *Client {
	return &Client{
		baseURL: "https:
		db:      db,
	}
}
