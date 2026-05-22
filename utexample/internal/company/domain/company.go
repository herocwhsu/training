package domain

type Company struct {
	ID    string
	Email string
	Name  string
}

func (c *Company) Validate() error {
	if c.Email == "" {
		return ErrInvalidEmail
	}
	if c.Name == "" {
		return ErrInvalidName
	}
	return nil
}
