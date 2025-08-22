package domain

import "errors"

type Company struct {
	ID    string
	Email string
	Name  string
}

// Simple domain validation/business rule example
func (c *Company) Validate() error {
	if c.Email == "" {
		return errors.New("email required")
	}
	if c.Name == "" {
		return errors.New("name required")
	}
	return nil
}
