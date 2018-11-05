// Package admin deals with all the functionality needed by an admin to modify
// state in the application
package admin

import (
	"time"

	"github.com/auburnhacks/sponsor/pkg/db"
	"github.com/pkg/errors"
)

// DefaultACL is a variables that is used for all admins if no ACL list is
// provided during signup
var DefaultACL = "read,update"

// ErrInvalidAuth is an error that is returns when there is a failed login attempt
var ErrInvalidAuth = errors.New("admin: invalud credentials provided")

// Admin is a struct that is the highest entity in the system
// It has write and read access to any thing in the database.
// More fine grained controlled can be given based on the ACL
// property
type Admin struct {
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	ACL       string    `db:"acl"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
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
func NewWithACL(name, email, password string, acl string) *Admin {
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
	query := `INSERT INTO admins(name, email, password, acl) VALUES(:name, :email, :password, :acl) RETURNING id`
	stmt, err := db.Conn.PrepareNamed(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	var id string
	err = stmt.QueryRow(map[string]interface{}{
		"name":     a.Name,
		"email":    a.Email,
		"password": a.Password,
		"acl":      a.ACL,
	}).Scan(&id)
	if err != nil {
		return err
	}
	a.ID = id
	return nil
}

// Login is a function that returns an instance of an admin is
// successcful or returns an error if there was some of error
func Login(email, password string) (*Admin, error) {
	a, err := ByEmail(email)
	if err != nil {
		return nil, err
	}
	if a.Password != password {
		return nil, ErrInvalidAuth
	}
	return a, nil
}

// ByID is a function that gets an admin from the given ID
func ByID(adminID string) (*Admin, error) {
	query := `SELECT * FROM admins WHERE id=:id`
	stmt, err := db.Conn.PrepareNamed(query)
	if err != nil {
		return nil, err
	}
	a := new(Admin)
	err = stmt.Get(a, map[string]interface{}{
		"id": adminID,
	})
	if err != nil {
		return nil, err
	}
	return a, nil
}

// ByEmail is a function that gets an admin from the given email
func ByEmail(adminEmail string) (*Admin, error) {
	query := `SELECT * FROM admins WHERE email=:email`
	stmt, err := db.Conn.PrepareNamed(query)
	if err != nil {
		return nil, err
	}
	a := new(Admin)
	err = stmt.Get(a, map[string]interface{}{
		"email": adminEmail,
	})
	if err != nil {
		return nil, err
	}
	return a, nil
}
