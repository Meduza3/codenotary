package internal

import (
	"codenotary/internal/deps"
	"database/sql"
)

const Database string = "codenotarydatabase.db"

var Db *sql.DB
var Client *deps.Client
