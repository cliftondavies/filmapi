CREATE INDEX IF NOT EXISTS films_title_idx ON films USING GIN (to_tsvector('simple', title));
CREATE INDEX IF NOT EXISTS films_genres_idx ON films USING GIN (genres);