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
var DefaultACL []string = []string{"read", "update"}

type Admin struct {
	ID        int      `db:"id,omitempty"`
	Name      string   `db:"name,omitempty"`
	Email     string   `db:"email,omitempty"`
	Password  string   `db:"password,omitempty"`
	ACL       []string `db:"acl,omitempty"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func New(name, email, password string) *Admin {
	a := &Admin{
		Name:     name,
		Email:    email,
		Password: password,
		ACL:      DefaultACL,
	}
	return a
}

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
	var adminId int
	err = stmt.QueryRow(a.Name, a.Email, a.Password, pq.Array(a.ACL)).Scan(&adminId)
	if err != nil {
		return err
	}
	// set the adminId to the instance
	a.ID = adminId
	return nil
}

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
