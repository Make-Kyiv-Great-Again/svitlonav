package main

import (
	"context"
	"log"
	"os"
	"runtime"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/paulmach/osm"
	"github.com/paulmach/osm/osmpbf"
)

func main() {
	_ = godotenv.Load()

	ctx := context.Background()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	file, err := os.Open("lamps.osm.pbf")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := osmpbf.New(context.Background(), file, runtime.NumCPU())

	log.Println("Import started...")

	tx, err := pool.Begin(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback(ctx)

	imported := 0
	total := 0

	for scanner.Scan() {
		total++

		if total%1000 == 0 {
			log.Printf("Processed %d objects...", total)
		}

		obj := scanner.Object()

		node, ok := obj.(*osm.Node)
		if !ok {
			continue
		}

		if node.Tags.Find("highway") != "street_lamp" {
			continue
		}

		_, err := tx.Exec(ctx, `
			INSERT INTO lamps (
				osm_id,
				geog,
				lit
			)
			VALUES (
				$1,
				ST_SetSRID(
					ST_Point($2,$3),
					4326
				)::geography,
				true
			)
			ON CONFLICT (osm_id) DO NOTHING
		`,
			int64(node.ID),
			node.Lon,
			node.Lat,
		)

		if err != nil {
			log.Println(err)
			continue
		}

		imported++

		if imported%500 == 0 {
			log.Printf("Imported %d lamps", imported)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	if err := tx.Commit(ctx); err != nil {
		log.Fatal(err)
	}

	log.Printf("Done!")
	log.Printf("Objects processed: %d", total)
	log.Printf("Lamps imported: %d", imported)
}