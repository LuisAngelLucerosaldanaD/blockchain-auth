
-- +migrate Up
CREATE TABLE IF NOT EXISTS cfg.dictionaries(
    id INT  NOT NULL PRIMARY KEY,
    name VARCHAR (50) NOT NULL,
    value VARCHAR (100) NOT NULL,
    description VARCHAR (200) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

-- +migrate Down

