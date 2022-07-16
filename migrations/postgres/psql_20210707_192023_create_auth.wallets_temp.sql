-- +migrate Up
CREATE TABLE auth.wallets_temp
(
    id         uuid         NOT NULL,
    mnemonic   varchar(250) NOT NULL,
    id_user    uuid         NOT NULL,
    created_at timestamp    NOT NULL DEFAULT now(),
    updated_at timestamp    NOT NULL DEFAULT now(),
    CONSTRAINT wallets_temp_pkey PRIMARY KEY (id)
);
-- +migrate Down

