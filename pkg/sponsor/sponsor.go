package sponsor

import (
	"time"

	"github.com/auburnhacks/sponsor/pkg/db"
	"github.com/lib/pq"
)

// DefaultACL is a variables that is used for all admins if no ACL list is
// provided during signup
var DefaultACL []string = []string{"read"}

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

func (s *Sponsor) Register() error {
	query := `INSERT INTO sponsors(name, email, password, company, acl)
	VALUES($1, $2, $3, $4, $5) RETURNING id`
	stmt, err := db.Conn.Prepare(query)
	if err != nil {
		return err
	}
	var sponsorId int
	err = stmt.QueryRow(s.Name, s.Email, s.Password, s.Company, pq.Array(s.ACL)).Scan(&sponsorId)
	if err != nil {
		return err
	}
	s.ID = sponsorId
	return nil
}
