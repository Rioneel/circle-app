package user

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

func (s *Service) ListUsers(ctx context.Context) ([]User, error) {
	return s.repo.ListUsers(ctx)
}

func (s *Service) TrustPath(ctx context.Context, meID, recommenderID string) ([]PathStep, error) {
	if meID == "" || recommenderID == "" {
		return nil, fmt.Errorf("meId and recommenderId are required")
	}
	return s.repo.TrustPath(ctx, meID, recommenderID)
}