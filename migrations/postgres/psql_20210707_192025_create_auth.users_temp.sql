-- +migrate Up
CREATE TABLE auth.users_temp
(
    id            uuid         NOT NULL,
    nickname      varchar(50)  NOT NULL,
    email         varchar(150) NOT NULL,
    "password"    varchar(150) NOT NULL,
    "name"        varchar(100) NULL,
    lastname      varchar(150) NULL,
    id_type       int4 NULL,
    id_number     varchar(15) NULL,
    cellphone     varchar(20) NULL,
    birth_date    timestamp NULL,
    verified_code varchar(150) NOT NULL,
    is_deleted    bool NULL,
    id_user       uuid         NOT NULL,
    deleted_at    timestamp    NOT NULL,
    created_at    timestamp    NOT NULL DEFAULT now(),
    updated_at    timestamp    NOT NULL DEFAULT now(),
    CONSTRAINT users_temp_email_key UNIQUE (email),
    CONSTRAINT users_temp_nickname_key UNIQUE (nickname),
    CONSTRAINT users_temp_pkey PRIMARY KEY (id)
);

-- +migrate Down

