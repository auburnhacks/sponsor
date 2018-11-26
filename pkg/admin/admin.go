// Package admin deals with all the functionality needed by an admin to modify
// state in the application
package admin

import (
	"fmt"
	"time"

	"github.com/auburnhacks/sponsor/pkg/db"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// DefaultACL is a variables that is used for all admins if no ACL list is
// provided during signup
var DefaultACL = "read,update"

// ErrInvalidAuth is an error that is returns when there is a failed login attempt
var ErrInvalidAuth = errors.New("pkg/admin: invalid credentials provided, please try again")

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
	q := `
	UPDATE admins 
	SET name = :name, email = :email, password = :password, acl = :acl
	WHERE id = :id
	RETURNING id`
	stmt, err := db.Conn.PrepareNamed(q)
	if err != nil {
		return err
	}
	_ = stmt.QueryRow(map[string]interface{}{
		"id":       a.ID,
		"name":     a.Name,
		"email":    a.Email,
		"password": a.Password,
		"acl":      a.ACL,
	})
	return nil
}

// Register is only called once when the admin first signs up
func (a *Admin) Register() error {
	// hash password and save it to the database
	pwdHash, err := bcrypt.GenerateFromPassword([]byte(a.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	// update the struct with the secure password
	a.Password = fmt.Sprintf("%s", pwdHash)
	query := `
	INSERT INTO admins
	(name, email, password, acl)
	VALUES(:name, :email, :password, :acl)
	RETURNING id`
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
		return nil, ErrInvalidAuth
	}
	// using bcrypt to match hashed password
	err = bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(password))
	if err != nil {
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
