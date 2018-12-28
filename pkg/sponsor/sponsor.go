// Package sponsor contains all the types and methods related to a sponsor
// in the api
package sponsor

import (
	"fmt"
	"time"

	"github.com/auburnhacks/sponsor/pkg/db"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// DefaultACL is a variables that is used for all admins if no ACL list is
// provided during signup
var DefaultACL = "read"

// ErrInvalidAuth is an error returned when there are invalid credentials
var ErrInvalidAuth = errors.New("pkg/sponsor: invalid auth credentials")

// Sponsor is a struc that repesents a sponsor in the system and the database
type Sponsor struct {
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	CompanyID string    `db:"company_id"`
	ACL       string    `db:"acl"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// New returns a instance of a Sponsor it uses the input parameters but
// if no ACL list is provided then it defaults to the DefaultACL list in the
// package definition
func New(name, email, password, companyID string, acl string) *Sponsor {
	s := &Sponsor{
		Name:      name,
		Email:     email,
		Password:  password,
		CompanyID: companyID,
	}
	if acl == "" {
		s.ACL = DefaultACL
	} else if len(acl) == 0 {
		s.ACL = DefaultACL
	} else {
		s.ACL = acl
	}
	return s
}

// ByID is a function that gets an admin from the given ID
func ByID(sponsorID string) (*Sponsor, error) {
	query := `SELECT * FROM sponsors WHERE id = :id LIMIT 1`
	stmt, err := db.Conn.PrepareNamed(query)
	if err != nil {
		return nil, err
	}
	s := new(Sponsor)
	err = stmt.Get(s, map[string]interface{}{
		"id": sponsorID,
	})
	if err != nil {
		return nil, err
	}
	return s, nil
}

// ByEmail is a function that gets an admin from the given email
func ByEmail(sponsorEmail string) (*Sponsor, error) {
	query := `SELECT * FROM sponsors WHERE email=:email`
	stmt, err := db.Conn.PrepareNamed(query)
	if err != nil {
		return nil, err
	}
	s := new(Sponsor)
	err = stmt.Get(s, map[string]interface{}{
		"email": sponsorEmail,
	})
	if err != nil {
		return nil, err
	}
	return s, nil
}

// Register is a function that is used when a new instance of a sponsor has to
// be saved to the database and the in-memory instance has to be updated with
// the lastInsertedID
// NOTE: use this only when New Sponsors have to be created
// Use the Save method on the Sponsor type for all other subsequent calls
func (s *Sponsor) Register() error {
	// hash password before storing it
	pwdHash, err := bcrypt.GenerateFromPassword([]byte(s.Password),
		bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	s.Password = fmt.Sprintf("%s", pwdHash)
	query := `
	INSERT INTO sponsors(name, email, password, company_id, acl)
	VALUES(:name, :email, :password, :company, :acl) RETURNING id`
	stmt, err := db.Conn.PrepareNamed(query)
	if err != nil {
		return err
	}
	var id string
	err = stmt.QueryRow(map[string]interface{}{
		"name":     s.Name,
		"email":    s.Email,
		"password": s.Password,
		"company":  s.CompanyID,
		"acl":      s.ACL,
	}).Scan(&id)
	if err != nil {
		return err
	}
	s.ID = id
	return nil
}

// Login is a function that takes an email and plain text password
// and returns a Sponsor if the credentials are valid
func Login(email, password string) (*Sponsor, error) {
	sp, err := ByEmail(email)
	if err != nil {
		return nil, ErrInvalidAuth
	}
	err = bcrypt.CompareHashAndPassword([]byte(sp.Password), []byte(password))
	if err != nil {
		return nil, ErrInvalidAuth
	}
	return sp, nil
}

// Save is method on sponsor to save the state of a sponsor to the db
func (s *Sponsor) Save() error {
	query := `
	UPDATE sponsors
	SET name = :name, email = :email, password = :password, acl = :acl
	WHERE id = :id
	RETURNING id`
	stmt, err := db.Conn.PrepareNamed(query)
	if err != nil {
		return err
	}
	_ = stmt.QueryRow(map[string]interface{}{
		"id":       s.ID,
		"name":     s.Name,
		"email":    s.Email,
		"password": s.Password,
		"acl":      s.ACL,
	})
	return nil
}

// Company is a struct that represents all the parameters required by a company
type Company struct {
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	Logo      string    `db:"logo"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// NewCompany is a constructor that returns an instance of a Company
func NewCompany(name, logo string) *Company {
	return &Company{
		Name: name,
		Logo: logo,
	}
}

// CompanyByID fetches a given company by the ID
func CompanyByID(ID string) (*Company, error) {
	query := `SELECT * FROM company WHERE id=:id LIMIT 1`
	stmt, err := db.Conn.PrepareNamed(query)
	if err != nil {
		return nil, err
	}
	c := new(Company)
	err = stmt.Get(c, map[string]interface{}{
		"id": ID,
	})
	if err != nil {
		return nil, err
	}
	return c, nil
}

// ListCompanies fetches all the list of companies from the database
func ListCompanies() ([]*Company, error) {
	query := `SELECT * FROM company ORDER BY id DESC`
	rows, err := db.Conn.Query(query)
	if err != nil {
		return nil, err
	}
	var companies []*Company
	for rows.Next() {
		c := new(Company)
		err := rows.Scan(&c.ID, &c.Name, &c.Logo, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			return nil, errors.Wrap(err, "sponsor: error while scanning rows for companies")
		}
		companies = append(companies, c)
	}
	return companies, nil
}

// Save saves an instance of the company to the database
func (c *Company) Save() error {
	q := `INSERT INTO company(name, logo) VALUES(:name, :logo) RETURNING id`
	stmt, err := db.Conn.PrepareNamed(q)
	if err != nil {
		return err
	}
	var id string
	err = stmt.QueryRowx(c).Scan(&id)
	if err != nil {
		return err
	}
	c.ID = id
	return nil
}
