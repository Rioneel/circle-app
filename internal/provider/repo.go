package provider

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

func toFloat64(v any) float64 {
	switch n := v.(type) {
	case float64:
		return n
	case int64:
		return float64(n)
	default:
		return 0
	}
}
func (r *repo) exists(ctx context.Context, userID, providerID string) (bool, error) {
	res, err := neo4j.ExecuteQuery(ctx, r.driver,
		`OPTIONAL MATCH (u:User {id:$userId})-[rec:RECOMMENDS]->(p:ServiceProvider {id:$providerId})
		 RETURN rec IS NOT NULL AS exists`,
		map[string]any{"userId": userID, "providerId": providerID},
		neo4j.EagerResultTransformer,
	)
	if err != nil {
		return false, fmt.Errorf("checking recommendation existence: %w", err)
	}
	if len(res.Records) == 0 {
		return false, nil
	}
	val, found := res.Records[0].Get("exists")
	if !found {
		return false, fmt.Errorf("unexpected query result shape")
	}
	return val.(bool), nil
}

func (r *repo) Search(ctx context.Context, userID, category, area string) ([]SearchResult, error) {
	cypher := `
		MATCH (cat:ServiceCategory {name:$category})<-[:PROVIDES]-(p:ServiceProvider)`
	params := map[string]any{"userId": userID, "category": category}
	if area != "" {
		cypher += ` MATCH (p)-[:LOCATED_IN]->(:Area {name:$area})`
		params["area"] = area
	}
	cypher += `
		MATCH (p)<-[rec:RECOMMENDS]-(recommender:User)
		MATCH path = (recommender)<-[:TRUSTS*1..3]-(me:User {id:$userId})
		WITH p, recommender, rec, min(length(path)) AS hops
		WITH p, sum(rec.rating * (1.0/hops)) / sum(1.0/hops) AS ws, count(rec) AS voices
		RETURN p.id AS id, p.name AS name, ws AS weightedScore, voices
		ORDER BY ws DESC LIMIT 20`

	result, err := neo4j.ExecuteQuery(ctx, r.driver, cypher, params, neo4j.EagerResultTransformer)
	if err != nil {
		return nil, fmt.Errorf("searching providers: %w", err)
	}

	out := make([]SearchResult, 0, len(result.Records))
	for _, rec := range result.Records {
		idVal, _ := rec.Get("id")
		nameVal, _ := rec.Get("name")
		wsVal, _ := rec.Get("weightedScore")
		voicesVal, _ := rec.Get("voices")
		out = append(out, SearchResult{
			ID: idVal.(string), Name: nameVal.(string),
			WeightedScore: toFloat64(wsVal), Voices: voicesVal.(int64),
		})
	}
	return out, nil
}

func (r *repo) GetProvider(ctx context.Context, id string) (*Provider, error) {
	result, err := neo4j.ExecuteQuery(ctx, r.driver,
		`MATCH (p:ServiceProvider {id:$id})
		 OPTIONAL MATCH (p)-[:PROVIDES]->(c:ServiceCategory)
		 OPTIONAL MATCH (p)-[:LOCATED_IN]->(a:Area)
		 RETURN p.id AS id, p.name AS name, c.name AS category, a.name AS area`,
		map[string]any{"id": id}, neo4j.EagerResultTransformer)
	if err != nil {
		return nil, fmt.Errorf("getting provider: %w", err)
	}
	if len(result.Records) == 0 {
		return nil, nil
	}
	rec := result.Records[0]
	idVal, _ := rec.Get("id")
	nameVal, _ := rec.Get("name")
	catVal, _ := rec.Get("category")
	areaVal, _ := rec.Get("area")
	p := &Provider{ID: idVal.(string), Name: nameVal.(string)}
	if catVal != nil {
		p.Category, _ = catVal.(string)
	}
	if areaVal != nil {
		p.Area, _ = areaVal.(string)
	}
	return p, nil
}

func (r *repo) Recommendations(ctx context.Context, providerID string) ([]Recommendation, error) {
	result, err := neo4j.ExecuteQuery(ctx, r.driver,
		`MATCH (u:User)-[r:RECOMMENDS]->(p:ServiceProvider {id:$providerId})
		 RETURN u.id AS userId, u.name AS userName, r.rating AS rating
		 ORDER BY r.rating DESC`,
		map[string]any{"providerId": providerID}, neo4j.EagerResultTransformer)
	if err != nil {
		return nil, fmt.Errorf("loading recommendations: %w", err)
	}
	out := make([]Recommendation, 0, len(result.Records))
	for _, rec := range result.Records {
		uid, _ := rec.Get("userId")
		uname, _ := rec.Get("userName")
		rating, _ := rec.Get("rating")
		out = append(out, Recommendation{UserID: uid.(string), UserName: uname.(string), Rating: toFloat64(rating)})
	}
	return out, nil
}

func (r *repo) CreateRecommendation(ctx context.Context, userID, providerID string, rating float64) error {
	already, err := r.exists(ctx, userID, providerID)
	if err != nil {
		return err
	}
	if already {
		_, err := neo4j.ExecuteQuery(ctx, r.driver,
			`MATCH (u:User {id:$userId})-[r:RECOMMENDS]->(p:ServiceProvider {id:$providerId}) SET r.rating=$rating`,
			map[string]any{"userId": userID, "providerId": providerID, "rating": rating},
			neo4j.EagerResultTransformer)
		return err
	}
	_, err = neo4j.ExecuteQuery(ctx, r.driver,
		`MATCH (u:User {id:$userId}), (p:ServiceProvider {id:$providerId})
		 CREATE (u)-[:RECOMMENDS {rating:$rating}]->(p)`,
		map[string]any{"userId": userID, "providerId": providerID, "rating": rating},
		neo4j.EagerResultTransformer)
	if err != nil {
		return fmt.Errorf("creating recommendation: %w", err)
	}
	return nil
}