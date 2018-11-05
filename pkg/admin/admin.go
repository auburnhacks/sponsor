// Package admin deals with all the functionality needed by an admin to modify
// state in the application
package admin

import (
	"database/sql"
	"time"

	"github.com/auburnhacks/sponsor/pkg/db"
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

// DefaultACL is a variables that is used for all admins if no ACL list is
// provided during signup
var DefaultACL = []string{"read", "update"}

// Admin is a struct that is the highest entity in the system
// It has write and read access to any thing in the database.
// More fine grained controlled can be given based on the ACL
// property
type Admin struct {
	ID        int
	Name      string
	Email     string
	Password  string
	ACL       []string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// New is a function that returns an instance of admin
// based on the name, email and password
// NOTE: using this function will use the DefaultACL
func New(name, email, password string) *Admin {
	a := &Admin{
		Name:     name,
		Email:    email,
		Password: password,
		ACL:      DefaultACL,
	}
	return a
}

// NewWithACL is a function that return an instance of an admin
// based on the name, email, password and ACL
func NewWithACL(name, email, password string, acl []string) *Admin {
	a := &Admin{
		Name:     name,
		Email:    email,
		Password: password,
		ACL:      acl,
	}
	return a
}

// Save saves the instance of an admin to the database
func (a *Admin) Save() error {
	query := `INSERT INTO admins(name, email, password, acl) VALUES($1, $2, $3, $4)`
	_, err := db.Conn.Exec(query, a.Name, a.Email, a.Password, a.ACL)
	if err != nil {
		return errors.Wrap(err, "cloud not save user to the database")
	}
	return nil
}

// Register is only called once when the admin first signs up
func (a *Admin) Register() error {
	query := `INSERT INTO admins(name, email, password, acl) VALUES($1, $2, $3, $4) RETURNING id`
	stmt, err := db.Conn.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	var adminID int
	err = stmt.QueryRow(a.Name, a.Email, a.Password, pq.Array(a.ACL)).Scan(&adminID)
	if err != nil {
		return err
	}
	// set the adminId to the instance
	a.ID = adminID
	return nil
}

// Login is a function that returns an instance of an admin is
// successcful or returns an error if there was some of error
func Login(email, password string) (*Admin, error) {
	query := `SELECT * FROM admins WHERE email=$1`
	stmt, err := db.Conn.Prepare(query)
	if err != nil {
		return nil, err
	}
	var a Admin
	err = stmt.QueryRow(email).Scan(&a.ID, &a.Name, &a.Email, &a.Password, pq.Array(&a.ACL), &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.Wrapf(err, "user with email %s not found", email)
		}
		return nil, err
	}
	return &a, nil
}

// GetAdminByID is a function that gets an admin from the given ID
func GetAdminByID(adminID int) (*Admin, error) {
	query := `SELECT * FROM admins WHERE id=$1`
	stmt, err := db.Conn.Prepare(query)
	if err != nil {
		return nil, err
	}
	var a Admin
	err = stmt.QueryRow(adminID).Scan(&a.ID, &a.Name, &a.Email, &a.Password, pq.Array(&a.ACL), &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.Wrapf(err, "admin: error getting admin with id %d", adminID)
		}
		return nil, err
	}
	return &a, nil
}
