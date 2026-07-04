CREATE TABLE lamp_edges (
    from_lamp  BIGINT NOT NULL REFERENCES lamps(id) ON DELETE CASCADE,
    to_lamp    BIGINT NOT NULL REFERENCES lamps(id) ON DELETE CASCADE,
    distance_m DOUBLE PRECISION NOT NULL,
    PRIMARY KEY (from_lamp, to_lamp)
);

CREATE INDEX idx_lamp_edges_from ON lamp_edges (from_lamp);
