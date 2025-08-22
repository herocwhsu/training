package companydao

import "context"

//go:generate mockgen -destination=../../../mocks/mock_dao.go -package=mocks github.com/herocwhsu/training/utexample/internal/dao/company CompanyDAO

type CompanyDAO interface {
	Insert(ctx context.Context, email, name string) (string, error)
	FindByID(ctx context.Context, id string) (email string, name string, err error)
}

// RDS implementation (placeholder)
type RDSCompanyDAO struct {
	// db *sql.DB // inject in real code
}

func NewRDSCompanyDAO() *RDSCompanyDAO { return &RDSCompanyDAO{} }

func (d *RDSCompanyDAO) Insert(ctx context.Context, email, name string) (string, error) {
	// TODO: implement DB insert, return generated id
	return "cmp_123", nil
}

func (d *RDSCompanyDAO) FindByID(ctx context.Context, id string) (string, string, error) {
	// TODO: implement DB select
	if id == "cmp_404" {
		return "", "", context.Canceled // just a stub error
	}
	return "team@example.com", "Example Inc", nil
}
