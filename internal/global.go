package internal

import (
	"codenotary/internal/deps"
	"database/sql"
)

const Database string = "codenotary.db"

var Db *sql.DB
var Client *deps.Client
