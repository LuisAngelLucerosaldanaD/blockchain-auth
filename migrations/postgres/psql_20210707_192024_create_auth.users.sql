-- +migrate Up
CREATE TABLE auth.users
(
    id                   uuid          NOT NULL,
    nickname             varchar(50)   NOT NULL,
    email                varchar(150)  NOT NULL,
    "password"           varchar(150)  NOT NULL,
    "name"               varchar(100)  NOT NULL,
    lastname             varchar(150)  NOT NULL,
    id_type              int4 NULL,
    id_number            varchar(15) NULL,
    cellphone            varchar(20) NULL,
    status_id            int4          NOT NULL,
    failed_attempts      int4          NOT NULL,
    block_date           timestamp NULL,
    disabled_date        timestamp NULL,
    last_login           timestamp NULL,
    last_change_password timestamp NULL,
    birth_date           timestamp NULL,
    verified_code        varchar(150)  NOT NULL,
    verified_at          timestamp     NOT NULL,
    is_deleted           bool NULL,
    id_user              uuid          NOT NULL,
    id_role              int4 NULL,
    full_path_photo      varchar(250)  NOT NULL,
    rsa_private          varchar(3500) NOT NULL,
    rsa_public           varchar(3500) NOT NULL,
    recovery_account_at  timestamp NULL,
    deleted_at           timestamp     NULL,
    created_at           timestamp     NOT NULL DEFAULT now(),
    updated_at           timestamp     NOT NULL DEFAULT now(),
    CONSTRAINT users_email_key UNIQUE (email),
    CONSTRAINT users_nickname_key UNIQUE (nickname),
    CONSTRAINT users_pkey PRIMARY KEY (id)
);
-- +migrate Down

