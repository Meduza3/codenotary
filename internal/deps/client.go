package deps

import "database/sql"

type client struct {
	baseURL string
	db      *sql.DB
}

func NewClient(db *sql.DB) *client {
	return &client{
		baseURL: "https://api.deps.dev/v3",
		db:      db,
	}
}
