-- +migrate Up
CREATE TABLE auth.wallets
(
    id                 uuid          NOT NULL,
    mnemonic           varchar(250)  NOT NULL,
    rsa_public         varchar(3500) NOT NULL,
    ip_device          varchar(50)   NOT NULL,
    status_id          int4          NOT NULL,
    identity_number    varchar(250)  NOT NULL,
    created_at         timestamp     NOT NULL DEFAULT now(),
    updated_at         timestamp     NOT NULL DEFAULT now(),
    CONSTRAINT wallets_identity_number_key UNIQUE (identity_number),
    CONSTRAINT wallets_pkey PRIMARY KEY (id)
);

-- +migrate Down

