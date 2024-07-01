-- +goose Up
-- +goose StatementBegin

CREATE TYPE resource_type AS ENUM ('BINARY', 'TEXT', 'CARD', 'PASSWORD');

CREATE TABLE IF NOT EXISTS users(
                                    id UUID DEFAULT gen_random_uuid(),
                                    name VARCHAR(255),
                                    password VARCHAR(255),
                                    registered_at timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                    UNIQUE(name),
                                    PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS pieces(
                                     id uuid DEFAULT gen_random_uuid(),
                                     content BYTEA,
                                     salt BYTEA,
                                     iv BYTEA,
                                     PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS resources(
                                      id serial PRIMARY KEY ,
                                      piece_uuid UUID,
                                      user_id UUID,
                                      data BYTEA,
                                      type resource_type,
                                      meta TEXT,
                                      CONSTRAINT fk_user
                                          FOREIGN KEY(user_id)
                                              REFERENCES users(id)
                                              ON DELETE CASCADE,
                                      CONSTRAINT fk_piece
                                          FOREIGN KEY (piece_uuid)
                                              REFERENCES pieces(id)

);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS resources;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS pieces;
DROP TYPE resource_type;
-- +goose StatementEnd
