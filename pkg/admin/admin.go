package admin

// SuperUser is an interface that has to be satisfied by other types
type Admin interface {
	Email() string
	Password() string
	// Save is a function that should save the details of a admin to the database
	Save() bool
	//CanAddAdmins is a function that can be used to check if a certain admin can
	// add other admins to the system
	CanAddAdmins() bool
	// CanAddSponsors is method that can be used to specify if an admin can add
	// sponsors to the system
	CanAddSponsors() bool
}

func NewAdmin(email, password string, isSuperAdmin bool) Admin {
	if isSuperAdmin {
		return &SuperAdmin{
			email:    email,
			password: password,
		}
	}
	return &RegularAdmin{
		email:    email,
		password: password,
	}
}

// Admin is a type that is a SuperUser and can edit crucial
// elements in the system
type SuperAdmin struct {
	email    string
	password string
}

func (sa *SuperAdmin) Email() string {
	return sa.email
}

func (sa *SuperAdmin) Password() string {
	return sa.password
}

func (sa *SuperAdmin) CanAddAdmins() bool {
	return true
}
func (sa *SuperAdmin) CanAddSponsors() bool {
	return true
}

func (sa *SuperAdmin) Save() bool {
	return true
}

type RegularAdmin struct {
	email    string
	password string
}

func (ra *RegularAdmin) Email() string {
	return ra.email
}

func (ra *RegularAdmin) Password() string {
	return ra.password
}

func (ra *RegularAdmin) CanAddAdmins() bool {
	return false
}
func (ra *RegularAdmin) CanAddSponsors() bool {
	return true
}

func (sa *RegularAdmin) Save() bool {
	return true
}

type Sponsor struct {
	Name        string   `json:"name,omitempty"`
	Company     string   `json:"company,omitempty"`
	Email       string   `json:"email,omitempty"`
	Username    string   `json:"username,omitempty"`
	ACL         []string `json:"acl,omitempty"`
	Password    string   `json:"password,omitempty"`
	CanAddUsers string   `json:"can_add_users,omitempty"`
}
