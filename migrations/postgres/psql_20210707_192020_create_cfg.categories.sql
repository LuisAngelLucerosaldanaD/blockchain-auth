-- +migrate Up
CREATE TABLE IF NOT EXISTS cfg.categories
(
    id         UUID         NOT NULL PRIMARY KEY,
    name       VARCHAR(100) NOT NULL,
    icon       VARCHAR(100) NOT NULL,
    color      VARCHAR(100) NOT NULL,
    created_at TIMESTAMP    NOT NULL DEFAULT now(),
    updated_at TIMESTAMP    NOT NULL DEFAULT now()
);

-- +migrate Down
