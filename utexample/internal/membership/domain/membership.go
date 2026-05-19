package domain

type Membership struct {
	ID        string
	CompanyID string
	UserID    string
	Role      string // "member" | "admin"
}

func (m *Membership) Validate() error {
	if m.Role != "member" && m.Role != "admin" {
		return ErrInvalidRole
	}
	return nil
}
