-- +migrate Up
CREATE TABLE IF NOT EXISTS cfg.messages(
    id INT  NOT NULL PRIMARY KEY,
    spa VARCHAR (100) NOT NULL,
    eng VARCHAR (100) NOT NULL,
    type_message INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);
