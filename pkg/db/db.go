// Package db provides a global interface for the database connection
// to all packages
package db

import "github.com/jmoiron/sqlx"

// Conn is a global variables that is used by all package to access a SQL database
var Conn *sqlx.DB
