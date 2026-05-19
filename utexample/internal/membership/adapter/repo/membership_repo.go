package repo

import (
	"context"

	"github.com/herocwhsu/training/utexample/internal/membership/domain"
	"github.com/herocwhsu/training/utexample/internal/membership/interfaces"
)

type membershipDoc struct {
	ID        string
	CompanyID string
	UserID    string
	Role      string
}

func docToEntity(doc *membershipDoc) *domain.Membership {
	return &domain.Membership{ID: doc.ID, CompanyID: doc.CompanyID, UserID: doc.UserID, Role: doc.Role}
}

type MembershipRepository struct {
	dao interfaces.MembershipDAO
}

func NewMembershipRepository(dao interfaces.MembershipDAO) *MembershipRepository {
	return &MembershipRepository{dao: dao}
}

func (r *MembershipRepository) Save(ctx context.Context, m *domain.Membership) error {
	id, err := r.dao.Insert(ctx, m.CompanyID, m.UserID, m.Role)
	if err != nil {
		return err
	}
	m.ID = id
	return nil
}

func (r *MembershipRepository) FindByID(ctx context.Context, id string) (*domain.Membership, error) {
	companyID, userID, role, err := r.dao.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return docToEntity(&membershipDoc{ID: id, CompanyID: companyID, UserID: userID, Role: role}), nil
}

func (r *MembershipRepository) FindByCompanyID(ctx context.Context, companyID string) ([]*domain.Membership, error) {
	rows, err := r.dao.FindByCompanyID(ctx, companyID)
	if err != nil {
		return nil, err
	}
	result := make([]*domain.Membership, 0, len(rows))
	for _, row := range rows {
		result = append(result, docToEntity(&membershipDoc{ID: row.ID, CompanyID: row.CompanyID, UserID: row.UserID, Role: row.Role}))
	}
	return result, nil
}

func (r *MembershipRepository) FindByUserID(ctx context.Context, userID string) ([]*domain.Membership, error) {
	rows, err := r.dao.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	result := make([]*domain.Membership, 0, len(rows))
	for _, row := range rows {
		result = append(result, docToEntity(&membershipDoc{ID: row.ID, CompanyID: row.CompanyID, UserID: row.UserID, Role: row.Role}))
	}
	return result, nil
}

func (r *MembershipRepository) ExistsByCompanyAndUser(ctx context.Context, companyID, userID string) (bool, error) {
	return r.dao.ExistsByCompanyAndUser(ctx, companyID, userID)
}

func (r *MembershipRepository) Remove(ctx context.Context, id string) error {
	return r.dao.DeleteByID(ctx, id)
}
