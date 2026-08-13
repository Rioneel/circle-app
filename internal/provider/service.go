package provider

import (
	"context"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Service struct {
	repo *repo
}

func NewService(driver neo4j.DriverWithContext) *Service {
	return &Service{repo: newRepo(driver)}
}

func (s *Service) Search(ctx context.Context, userID, category, area string) ([]SearchResult, error) {
	if userID == "" || category == "" {
		return nil, fmt.Errorf("userId and category are required")
	}
	return s.repo.Search(ctx, userID, category, area)
}

func (s *Service) GetProvider(ctx context.Context, id string) (*Provider, error) {
	if id == "" {
		return nil, fmt.Errorf("id is required")
	}
	return s.repo.GetProvider(ctx, id)
}

func (s *Service) Recommendations(ctx context.Context, providerID string) ([]Recommendation, error) {
	if providerID == "" {
		return nil, fmt.Errorf("providerId is required")
	}
	return s.repo.Recommendations(ctx, providerID)
}

func (s *Service) Recommend(ctx context.Context, userID, providerID string, rating float64) error {
	if userID == "" || providerID == "" {
		return fmt.Errorf("userId and providerId are required")
	}
	if rating < 1 || rating > 5 {
		return fmt.Errorf("rating must be between 1 and 5")
	}
	return s.repo.CreateRecommendation(ctx, userID, providerID, rating)
}