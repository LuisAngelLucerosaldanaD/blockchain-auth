-- +migrate Up
CREATE TABLE auth.accounting
(
    id         uuid      NOT NULL,
    id_wallet  uuid      NOT NULL,
    amount     float8    NOT NULL,
    id_user    uuid      NOT NULL,
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp NOT NULL DEFAULT now(),
    CONSTRAINT accounting_id_wallet_key UNIQUE (id_wallet),
    CONSTRAINT accounting_pkey PRIMARY KEY (id)
);

-- +migrate Down

