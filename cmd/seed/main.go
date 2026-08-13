package main

import (
	"context"
	"log"
	"math/rand"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"

	"circle/internal/config"
	"circle/internal/db"
)

var names = []string{
	"Aarav", "Priya", "Rohan", "Ananya", "Kiran",
	"Meera", "Sanjay", "Divya", "Arjun", "Neha",
}

var categories = []string{"Plumber", "Electrician", "Dentist", "AC Repair", "Carpenter", "Tutor", "Pest Control", "House Cleaning"}
var areas = []string{"Hanamkonda", "Warangal", "Hyderabad"}
var providerNames = []string{
	"Sri Balaji Plumbing", "Warangal Electricals", "Smile Care Dental",
	"CoolFix AC Services", "Master Carpentry Works", "BrightMinds Tuition",
	"SafeHome Pest Control", "SparkleClean Services", "Om Sai Plumbing",
	"PowerLine Electricals",
}

func main() {
	ctx := context.Background()
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}
	driver, err := db.NewDriver(ctx, cfg)
	if err != nil {
		log.Fatalf("could not connect to CognoDB: %v", err)
	}
	defer driver.Close(ctx)

	rng := rand.New(rand.NewSource(42))

	log.Println("seeding users...")
	userRows := make([]map[string]any, len(names))
	for i, name := range names {
		userRows[i] = map[string]any{"id": "user-" + itoa(i+1), "name": name}
	}
	mustRun(ctx, driver, `UNWIND $rows AS row MERGE (u:User {id: row.id}) SET u.name = row.name`,
		map[string]any{"rows": userRows})

	log.Println("wiring a denser trust network...")
	var edges []map[string]any
	for i := range names {
		fromID := "user-" + itoa(i+1)
		for j := 0; j < 5; j++ { // denser: 5 outbound trust edges per user, more real multi-hop paths
			toID := "user-" + itoa(rng.Intn(len(names))+1)
			if toID == fromID {
				continue
			}
			edges = append(edges, map[string]any{"from": fromID, "to": toID, "weight": 0.5 + rng.Float64()*0.5})
		}
	}
	mustRun(ctx, driver, `UNWIND $edges AS e MATCH (a:User {id:e.from}),(b:User {id:e.to}) MERGE (a)-[t:TRUSTS]->(b) SET t.weight=e.weight`,
		map[string]any{"edges": edges})

	log.Println("seeding categories and areas...")
	mustRun(ctx, driver, `UNWIND $cats AS c MERGE (:ServiceCategory {name: c})`, map[string]any{"cats": categories})
	mustRun(ctx, driver, `UNWIND $areas AS a MERGE (:Area {name: a})`, map[string]any{"areas": areas})

	log.Println("seeding providers...")
	provRows := make([]map[string]any, len(providerNames))
	for i, name := range providerNames {
		provRows[i] = map[string]any{
			"id": "provider-" + itoa(i+1), "name": name,
			"category": categories[i%len(categories)],
			"area":     areas[rng.Intn(len(areas))],
		}
	}
	mustRun(ctx, driver, `
		UNWIND $rows AS row
		MERGE (p:ServiceProvider {id: row.id}) SET p.name = row.name
		WITH p, row MATCH (c:ServiceCategory {name: row.category}) MERGE (p)-[:PROVIDES]->(c)
		WITH p, row MATCH (a:Area {name: row.area}) MERGE (p)-[:LOCATED_IN]->(a)`,
		map[string]any{"rows": provRows})

	log.Println("seeding baseline recommendations...")
	var recs []map[string]any
	for i := range names {
		userID := "user-" + itoa(i+1)
		for j := 0; j < 1+rng.Intn(2); j++ {
			recs = append(recs, map[string]any{
				"userId": userID, "providerId": "provider-" + itoa(rng.Intn(len(providerNames))+1),
				"rating": float64(2 + rng.Intn(4)),
			})
		}
	}
	mustRun(ctx, driver, `
		UNWIND $rows AS row
		MATCH (u:User {id: row.userId}), (p:ServiceProvider {id: row.providerId})
		MERGE (u)-[r:RECOMMENDS]->(p) SET r.rating = row.rating`,
		map[string]any{"rows": recs})

	log.Println("planting deliberate multi-hop demo scenarios across categories...")

	// Scenario 1: user-1 -> user-6 -> user-9 -> recommends provider-1 (Plumber)
	mustRun(ctx, driver, `MATCH (a:User {id:"user-1"}),(b:User {id:"user-6"}) MERGE (a)-[t:TRUSTS]->(b) SET t.weight=0.9`, nil)
	mustRun(ctx, driver, `MATCH (a:User {id:"user-6"}),(b:User {id:"user-9"}) MERGE (a)-[t:TRUSTS]->(b) SET t.weight=0.9`, nil)
	mustRun(ctx, driver, `MATCH (u:User {id:"user-9"}),(p:ServiceProvider {id:"provider-1"}) MERGE (u)-[r:RECOMMENDS]->(p) SET r.rating=5.0`, nil)

	// Scenario 2: user-2 -> user-7 -> user-10 -> recommends provider-2 (Electrician)
	mustRun(ctx, driver, `MATCH (a:User {id:"user-2"}),(b:User {id:"user-7"}) MERGE (a)-[t:TRUSTS]->(b) SET t.weight=0.85`, nil)
	mustRun(ctx, driver, `MATCH (a:User {id:"user-7"}),(b:User {id:"user-10"}) MERGE (a)-[t:TRUSTS]->(b) SET t.weight=0.85`, nil)
	mustRun(ctx, driver, `MATCH (u:User {id:"user-10"}),(p:ServiceProvider {id:"provider-2"}) MERGE (u)-[r:RECOMMENDS]->(p) SET r.rating=5.0`, nil)

	// Scenario 3: user-4 -> user-8 -> user-5 -> recommends provider-6 (Tutor), 3-hop chain
	mustRun(ctx, driver, `MATCH (a:User {id:"user-4"}),(b:User {id:"user-8"}) MERGE (a)-[t:TRUSTS]->(b) SET t.weight=0.8`, nil)
	mustRun(ctx, driver, `MATCH (a:User {id:"user-8"}),(b:User {id:"user-5"}) MERGE (a)-[t:TRUSTS]->(b) SET t.weight=0.8`, nil)
	mustRun(ctx, driver, `MATCH (u:User {id:"user-5"}),(p:ServiceProvider {id:"provider-6"}) MERGE (u)-[r:RECOMMENDS]->(p) SET r.rating=4.0`, nil)

	// Scenario 4: mixed-reputation provider-4 (AC Repair) — one great review, one poor one,
	// so the weighted-average query visibly does something interesting in a demo.
	mustRun(ctx, driver, `MATCH (u:User {id:"user-3"}),(p:ServiceProvider {id:"provider-4"}) MERGE (u)-[r:RECOMMENDS]->(p) SET r.rating=5.0`, nil)
	mustRun(ctx, driver, `MATCH (u:User {id:"user-4"}),(p:ServiceProvider {id:"provider-4"}) MERGE (u)-[r:RECOMMENDS]->(p) SET r.rating=2.0`, nil)
	log.Println("removing any random shortcuts that would defeat the planted scenarios...")
	mustRun(ctx, driver, `MATCH (a:User {id:"user-1"})-[t:TRUSTS]->(b:User {id:"user-9"}) DELETE t`, nil)
	mustRun(ctx, driver, `MATCH (a:User {id:"user-2"})-[t:TRUSTS]->(b:User {id:"user-10"}) DELETE t`, nil)
	mustRun(ctx, driver, `MATCH (a:User {id:"user-4"})-[t:TRUSTS]->(b:User {id:"user-5"}) DELETE t`, nil)
	log.Println("done. demo users: user-1 through user-10.")
	log.Println("try: user-1 searching Plumber, user-2 searching Electrician, user-4 searching Tutor.")
}

func mustRun(ctx context.Context, driver neo4j.DriverWithContext, cypher string, params map[string]any) {
	if _, err := neo4j.ExecuteQuery(ctx, driver, cypher, params, neo4j.EagerResultTransformer); err != nil {
		log.Fatalf("seed query failed: %v", err)
	}
}

func itoa(n int) string {
	if n < 10 {
		return string(rune('0' + n))
	}
	return string(rune('0'+n/10)) + string(rune('0'+n%10))
}