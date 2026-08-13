package db

import (
	"context"
	"fmt"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"

	"circle/internal/config"
)

func NewDriver(ctx context.Context, cfg *config.Config) (neo4j.DriverWithContext, error) {
	driver, err := neo4j.NewDriverWithContext(
		cfg.CognoURI,
		neo4j.BasicAuth(cfg.CognoUser, cfg.CognoPassword, ""),
	)
	if err != nil {
		return nil, fmt.Errorf("creating driver: %w", err)
	}

	verifyCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := driver.VerifyConnectivity(verifyCtx); err != nil {
		return nil, fmt.Errorf("could not reach CognoDB at startup: %w", err)
	}
	return driver, nil
}

func HealthCheck(ctx context.Context, driver neo4j.DriverWithContext) error {
	checkCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	return driver.VerifyConnectivity(checkCtx)
}