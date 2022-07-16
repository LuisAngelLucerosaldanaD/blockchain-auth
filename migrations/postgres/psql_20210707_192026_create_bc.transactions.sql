-- +migrate Up
CREATE TABLE bc.transactions
(
    id         uuid      NOT NULL,
    from_id    uuid      NOT NULL,
    to_id      uuid      NOT NULL,
    amount     float8    NOT NULL,
    type_id    int4      NOT NULL,
    "data"     varchar   NOT NULL,
    block      int8      NOT NULL,
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp NOT NULL DEFAULT now(),
    CONSTRAINT transactions_pkey PRIMARY KEY (id)
);
-- +migrate Down
