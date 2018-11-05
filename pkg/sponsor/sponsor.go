// Package sponsor contains all the types and methods related to a sponsor
// in the api
package sponsor

import (
	"time"

	"github.com/auburnhacks/sponsor/pkg/db"
	"github.com/lib/pq"
)

// DefaultACL is a variables that is used for all admins if no ACL list is
// provided during signup
var DefaultACL = []string{"read"}

// Sponsor is a struc that repesents a sponsor in the system and the database
type Sponsor struct {
	ID        int
	Name      string
	Email     string
	Password  string
	Company   string
	ACL       []string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// New returns a instance of a Sponsor it uses the input parameters but
// if no ACL list is provided then it defaults to the DefaultACL list in the
// package definition
func New(name, email, password, company string, acl []string) *Sponsor {
	s := &Sponsor{
		Name:     name,
		Password: password,
		Company:  company,
	}
	if acl == nil {
		s.ACL = DefaultACL
	} else if len(acl) == 0 {
		s.ACL = DefaultACL
	}
	return s
}

// Register is a function that is used when a new instance of a sponsor has to
// be saved to the database and the in-memory instance has to be updated with
// the lastInsertedID
// NOTE: use this only when New Sponsors have to be created
// Use the Save method on the Sponsor type for all other subsequent calls
func (s *Sponsor) Register() error {
	query := `INSERT INTO sponsors(name, email, password, company, acl) VALUES($1, $2, $3, $4, $5) RETURNING id`
	stmt, err := db.Conn.Prepare(query)
	if err != nil {
		return err
	}
	var sponsorID int
	err = stmt.QueryRow(s.Name, s.Email, s.Password, s.Company, pq.Array(s.ACL)).Scan(&sponsorID)
	if err != nil {
		return err
	}
	s.ID = sponsorID
	return nil
}
