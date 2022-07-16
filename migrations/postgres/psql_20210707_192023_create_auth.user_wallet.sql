-- +migrate Up
CREATE TABLE auth.user_wallet
(
    id         uuid      NOT NULL,
    id_user    uuid      NOT NULL,
    id_wallet  uuid      NOT NULL,
    is_delete  bool      NOT NULL DEFAULT false,
    created_at timestamp NOT NULL DEFAULT now(),
    deleted_at timestamp NULL,
    updated_at timestamp NOT NULL DEFAULT now(),
    CONSTRAINT user_wallet_id_wallet_key UNIQUE (id_wallet),
    CONSTRAINT user_wallet_pkey PRIMARY KEY (id)
);

-- +migrate Down

