-- +migrate Up
CREATE TABLE auth.accounting
(
    id         uuid      NOT NULL,
    wallet_id  uuid      NOT NULL unique,
    amount     float8    NOT NULL,
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp NOT NULL DEFAULT now(),
    CONSTRAINT accounting_wallet_id_key UNIQUE (wallet_id),
    CONSTRAINT accounting_pkey PRIMARY KEY (id),
    CONSTRAINT FK_wallet_id foreign key(wallet_id) references auth.wallet(id)
);

-- +migrate Down

