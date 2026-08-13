package trust

import (
	"context"
	"fmt"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Service struct {
	repo *repo
}

func NewService(driver neo4j.DriverWithContext) *Service {
	return &Service{repo: newRepo(driver)}
}

func (s *Service) CreateTrust(ctx context.Context, fromID, toID string, weight float64) error {
	if fromID == "" || toID == "" {
		return fmt.Errorf("fromId and toId are required")
	}
	if fromID == toID {
		return fmt.Errorf("a user cannot trust themselves")
	}
	return s.repo.CreateTrust(ctx, TrustEdge{FromId: fromID, ToId: toID, Weight: weight, CreatedAt: time.Now()})
}

func (s *Service) UpdateTrust(ctx context.Context, fromID, toID string, weight float64) error {
	if fromID == "" || toID == "" {
		return fmt.Errorf("fromId and toId are required")
	}
	return s.repo.UpdateTrust(ctx, TrustEdge{FromId: fromID, ToId: toID, Weight: weight})
}

func (s *Service) RemoveTrust(ctx context.Context, fromID, toID string) error {
	if fromID == "" || toID == "" {
		return fmt.Errorf("fromId and toId are required")
	}
	return s.repo.RemoveTrust(ctx, fromID, toID)
}

func (s *Service) ListTrusted(ctx context.Context, userID string) ([]TrustedUser, error) {
	if userID == "" {
		return nil, fmt.Errorf("userId is required")
	}
	return s.repo.ListTrusted(ctx, userID)
}