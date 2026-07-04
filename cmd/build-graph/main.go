package main

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	ctx := context.Background()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	radiusM := 20.0
	if v := os.Getenv("LAMP_RADIUS_METERS"); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			radiusM = f
		}
	}

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	log.Printf("Building lamp graph, radius=%.0fm ...", radiusM)

	if _, err := pool.Exec(ctx, `TRUNCATE lamp_edges`); err != nil {
		log.Fatalf("truncate lamp_edges: %v", err)
	}

	tag, err := pool.Exec(ctx, `
		INSERT INTO lamp_edges (from_lamp, to_lamp, distance_m)
		SELECT a.id, b.id, ST_Distance(a.geog, b.geog)
		FROM lamps a
		JOIN lamps b
		  ON a.id <> b.id
		 AND ST_DWithin(a.geog, b.geog, $1)
	`, radiusM)
	if err != nil {
		log.Fatalf("build graph: %v", err)
	}

	log.Printf("Done. Inserted %d directed edges.", tag.RowsAffected())
}
