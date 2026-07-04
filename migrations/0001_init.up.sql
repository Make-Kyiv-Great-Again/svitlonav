CREATE EXTENSION IF NOT EXISTS postgis;

CREATE TABLE lamps (
    id           BIGSERIAL PRIMARY KEY,
    geog         GEOGRAPHY(Point, 4326) NOT NULL,
    osm_id       BIGINT,
    queue_group  TEXT,
    lit          BOOLEAN DEFAULT TRUE,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX lamps_geog_gist_idx ON lamps USING GIST (geog);
CREATE INDEX lamps_queue_group_idx ON lamps (queue_group);
