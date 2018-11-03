// package db provides a global interface for the database connection
// to all packages
package db

import (
	"database/sql"
)

var Conn *sql.DB
