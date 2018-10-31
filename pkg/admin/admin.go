package admin

// SuperUser is an interface that has to be satisfied by other types
type Admin interface {
	CanAddAdmins() bool
	CanAddSponsors() bool
}

// Admin is a type that is a SuperUser and can edit crucial
// elements in the system
type SuperAdmin struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

func (sa *SuperAdmin) CanAddAdmins() bool {
	return true
}
func (sa *SuperAdmin) CanAddSponsors() bool {
	return true
}

type RegularAdmin struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

func (sa *RegularAdmin) CanAddAdmins() bool {
	return false
}
func (sa *RegularAdmin) CanAddSponsors() bool {
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
