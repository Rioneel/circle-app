package trust

import (
	"context"
	"fmt"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type repo struct {
	driver neo4j.DriverWithContext
}

func newRepo(driver neo4j.DriverWithContext) *repo {
	return &repo{driver : driver}
}
func(r *repo) exists(ctx context.Context, fromID, toID string)(bool, error) {
	res, err := neo4j.ExecuteQuery(ctx, r.driver, 
		`OPTIONAL MATCH (f: User {id: $fromID}) - [k:TRUSTS] -> (t: User{id: $toID}) 
		RETURN k IS NOT Null AS exists`, map[string]any{"fromID" : fromID, "toID" : toID},neo4j.EagerResultTransformer)
		if err != nil {
			return false, fmt.Errorf("checking trust existence: %w",err)
		}
		// apprently cant assume how many rows something returns eventhough its guarnateed
		if len(res.Records) == 0 {
		// No row at all means nothing matched — treat as "doesn't exist,"
		// don't crash the process over it.
		return false, nil
	}
	val, found := res.Records[0].Get("exists")
	if !found {
		return false, fmt.Errorf("unexpected query result shape")
	}
	return val.(bool), nil
}
func(r *repo) CreateTrust(ctx context.Context, edge TrustEdge) (error){
		already, err1 := r.exists(ctx, edge.FromId, edge.ToId)
	if err1 != nil {
		return err1
	}
	if already {
		return r.UpdateTrust(ctx, edge)
	}

		_, err := neo4j.ExecuteQuery(ctx, r.driver, 
		`MATCH
		 (f:User {id : $fromID}),
		 (t:User {id : $toID })
		 CREATE 
		 (f)-[:TRUSTS {weight: $weight, createdAt : $createdAt}] ->(t)
		 `,
		 map[string]any{
			"fromID": edge.FromId,
			"toID" : edge.ToId,
			"weight" : edge.Weight,
			"createdAt": edge.CreatedAt.Format(time.RFC3339),
		 }, 
		 neo4j.EagerResultTransformer,
		)
		if err != nil {
			return fmt.Errorf("query failed : %w", err)
		}
		
		return nil 
	}


func(r *repo) UpdateTrust(ctx context.Context,edge TrustEdge) error{
	
	_, err := neo4j.ExecuteQuery(ctx, r.driver, 
	`MATCH (f: User {id: $fromID}) - [r:TRUSTS] -> (t : User {id: $toID}) 
	SET r.weight = $weight`,
	map[string]any{"fromID": edge.FromId, "toID": edge.ToId, "weight": edge.Weight},neo4j.EagerResultTransformer, 
)

	if err != nil {
		return fmt.Errorf("query failed : %w", err)
	}

	return nil
}

func(r *repo) RemoveTrust(ctx context.Context, fromID string, toID string) error {
	
	_, err := neo4j.ExecuteQuery(ctx, r.driver, 
	`MATCH (f :User {id: $fromID}) - [r:TRUSTS] -> (t:User {id: $toID})
	DELETE r`, map[string]any{"fromID": fromID, "toID": toID}, neo4j.EagerResultTransformer)
	if err != nil {
		return fmt.Errorf("query failed : %w", err)
	}

	return nil
}

func(r *repo) ListTrusted(ctx context.Context, userID string) ([]TrustedUser ,error) {

	result , err := neo4j.ExecuteQuery(ctx, r.driver, 
	`MATCH (f : User {id: $userId}) - [t:TRUSTS] -> (x: User)
	 RETURN x.id AS id, x.name AS name, t.weight AS weight`, map[string]any{"userId": userID},neo4j.EagerResultTransformer,
	)

	 if err != nil {
		return nil , fmt.Errorf("query failed : %w", err)
	 }

	 users := make([]TrustedUser, 0 , len(result.Records))
	 for _, rec := range result.Records {
		IdVal, _ := rec.Get("id")
		nameVal, _ := rec.Get("name")
		weightVal, _ := rec.Get("weight")
		users = append(users, TrustedUser{
			ID : IdVal.(string),
			Name : nameVal.(string),
			Weight : weightVal.(float64),
		})
	 }
	 return users, nil
}
