package domain

type Role string

const (
	RoleMember Role = "member"
	RoleAdmin  Role = "admin"
)

type Membership struct {
	ID        string
	CompanyID string
	UserID    string
	Role      Role
}

func (m *Membership) Validate() error {
	if m.Role != RoleMember && m.Role != RoleAdmin {
		return ErrInvalidRole
	}
	return nil
}
