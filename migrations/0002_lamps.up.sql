ALTER TABLE lamps
ADD CONSTRAINT lamps_osm_id_unique
UNIQUE(osm_id);