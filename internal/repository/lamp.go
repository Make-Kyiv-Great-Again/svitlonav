package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type LampRepository struct {
	pool *pgxpool.Pool
}

func NewLampRepository(pool *pgxpool.Pool) *LampRepository {
	return &LampRepository{pool: pool}
}

type NearbyLamp struct {
	ID  int64
	Lat float64
	Lon float64
}

func (r *LampRepository) LampsNearLine(ctx context.Context, coords [][2]float64, bufferM float64) ([]NearbyLamp, error) {
	if len(coords) < 2 {
		return nil, nil
	}

	var wkt strings.Builder
	wkt.WriteString("LINESTRING(")
	for i, c := range coords {
		if i > 0 {
			wkt.WriteString(",")
		}
		fmt.Fprintf(&wkt, "%f %f", c[1], c[0])
	}
	wkt.WriteString(")")

	const query = `
		SELECT id, ST_Y(geog::geometry), ST_X(geog::geometry)
		FROM lamps
		WHERE ST_DWithin(
    geog,
    ST_SetSRID(ST_GeomFromText($1), 4326)::geography,
    $2
		)
		`
	rows, err := r.pool.Query(ctx, query, wkt.String(), bufferM)
	if err != nil {
		return nil, fmt.Errorf("lamps near line: %w", err)
	}
	defer rows.Close()

	var lamps []NearbyLamp
	for rows.Next() {
		var l NearbyLamp
		if err := rows.Scan(&l.ID, &l.Lat, &l.Lon); err != nil {
			return nil, fmt.Errorf("scan lamp: %w", err)
		}
		lamps = append(lamps, l)
	}
	return lamps, rows.Err()
}
