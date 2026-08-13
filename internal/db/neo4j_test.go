package db_test


import (
	"context"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)


func TestCognoDBconnectivity(t *testing.T){
	_ = godotenv.Load("../../.env")

	uri := os.Getenv("COGNODB_URI")
	user := os.Getenv("COGNODB_USER")
	pass := os.Getenv("COGNODB_PASSWORD")

	driver, err := neo4j.NewDriverWithContext(uri, neo4j.BasicAuth(user,pass,""))
	if err != nil {
		t.Fatalf("could not connect: %v", err)
	}
	defer driver.Close(context.Background())

	ctx := context.Background()
	if err := driver.VerifyConnectivity(ctx); err != nil {
		t.Fatalf("could not connect:%v", err)
	}

	result, err := neo4j.ExecuteQuery(ctx, driver, "RETURN 1 AS n", nil, neo4j.EagerResultTransformer)
	if err != nil {
		t.Fatalf("query failed: %v", err)
	}

	got, _ := result.Records[0].Get("n")
	if got != int64(1){
		t.Fatalf("expected 1, got %v", got)
	}
}
