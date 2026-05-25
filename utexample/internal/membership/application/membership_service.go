package application

import (
	"context"
	"fmt"

	"github.com/herocwhsu/training/utexample/internal/membership/domain"
	"github.com/herocwhsu/training/utexample/internal/membership/interfaces"
)

type MembershipService struct {
	repo          interfaces.MembershipRepository
	companyReader interfaces.CompanyReader
	userReader    interfaces.UserReader
}

func NewMembershipService(
	repo interfaces.MembershipRepository,
	companyReader interfaces.CompanyReader,
	userReader interfaces.UserReader,
) *MembershipService {
	return &MembershipService{repo: repo, companyReader: companyReader, userReader: userReader}
}

func (s *MembershipService) Add(ctx context.Context, companyID, userID, role string) (string, error) {
	if _, err := s.companyReader.FindByID(ctx, companyID); err != nil {
		return "", fmt.Errorf("find company: %w", err)
	}
	if _, err := s.userReader.FindByID(ctx, userID); err != nil {
		return "", fmt.Errorf("find user: %w", err)
	}
	exists, err := s.repo.ExistsByCompanyAndUser(ctx, companyID, userID)
	if err != nil {
		return "", fmt.Errorf("check membership: %w", err)
	}
	if exists {
		return "", domain.ErrAlreadyMember
	}
	m := &domain.Membership{CompanyID: companyID, UserID: userID, Role: role}
	if err := m.Validate(); err != nil {
		return "", err
	}
	saved, err := s.repo.Save(ctx, m)
	if err != nil {
		return "", fmt.Errorf("save membership: %w", err)
	}
	return saved.ID, nil
}

func (s *MembershipService) Remove(ctx context.Context, membershipID string) error {
	if err := s.repo.Remove(ctx, membershipID); err != nil {
		return fmt.Errorf("remove membership: %w", err)
	}
	return nil
}

func (s *MembershipService) ListByCompany(ctx context.Context, companyID string) ([]*domain.Membership, error) {
	memberships, err := s.repo.FindByCompanyID(ctx, companyID)
	if err != nil {
		return nil, fmt.Errorf("list by company: %w", err)
	}
	return memberships, nil
}

func (s *MembershipService) ListByUser(ctx context.Context, userID string) ([]*domain.Membership, error) {
	memberships, err := s.repo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list by user: %w", err)
	}
	return memberships, nil
}
