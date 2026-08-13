package user

import (
	"context"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type repo struct {
	driver neo4j.DriverWithContext
}

func newRepo(driver neo4j.DriverWithContext) *repo {
	return &repo{driver: driver}
}

func (r *repo) ListUsers(ctx context.Context) ([]User, error) {
	result, err := neo4j.ExecuteQuery(ctx, r.driver,
		`MATCH (u:User) RETURN u.id AS id, u.name AS name ORDER BY u.name LIMIT 50`,
		nil, neo4j.EagerResultTransformer)
	if err != nil {
		return nil, fmt.Errorf("listing users: %w", err)
	}
	out := make([]User, 0, len(result.Records))
	for _, rec := range result.Records {
		id, _ := rec.Get("id")
		name, _ := rec.Get("name")
		out = append(out, User{ID: id.(string), Name: name.(string)})
	}
	return out, nil
}

func (r *repo) TrustPath(ctx context.Context, meID, recommenderID string) ([]PathStep, error) {
	result, err := neo4j.ExecuteQuery(ctx, r.driver,
		`MATCH p = shortestPath((me:User {id:$meId})-[:TRUSTS*..4]->(rec:User {id:$recommenderId}))
		 RETURN [n IN nodes(p) | {id: n.id, name: n.name}] AS steps`,
		map[string]any{"meId": meID, "recommenderId": recommenderID},
		neo4j.EagerResultTransformer)
	if err != nil {
		return nil, fmt.Errorf("finding trust path: %w", err)
	}
	if len(result.Records) == 0 {
		return nil, nil
	}
	raw, _ := result.Records[0].Get("steps")
	list, _ := raw.([]any)
	steps := make([]PathStep, 0, len(list))
	for _, item := range list {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		id, _ := m["id"].(string)
		name, _ := m["name"].(string)
		steps = append(steps, PathStep{ID: id, Name: name})
	}
	return steps, nil
}