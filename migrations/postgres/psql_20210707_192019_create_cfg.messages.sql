-- +migrate Up
CREATE TABLE IF NOT EXISTS cfg.messages(
    id INT  NOT NULL PRIMARY KEY,
    spa VARCHAR (100) NOT NULL,
    eng VARCHAR (100) NOT NULL,
    type_message INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

-- +migrate Down
DROP TABLE IF EXISTS auth.accounting;
DROP TABLE IF EXISTS auth.user_wallet;
DROP TABLE IF EXISTS auth.wallets;
DROP TABLE IF EXISTS auth.wallets_temp;
DROP TABLE IF EXISTS auth.users;
DROP TABLE IF EXISTS auth.users_temp;

DROP TABLE IF EXISTS bc.blocks;
DROP TABLE IF EXISTS bc.blocks_tmp;
DROP TABLE IF EXISTS bc.transactions;

DROP TABLE IF EXISTS cfg.messages;
DROP TABLE IF EXISTS cfg.dictionaries;
