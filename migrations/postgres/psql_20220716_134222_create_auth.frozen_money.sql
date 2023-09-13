-- +migrate Up
CREATE TABLE IF NOT EXISTS auth.frozen_money
(
    id         uuid      NOT NULL PRIMARY KEY,
    wallet_id  UUID      NOT NULL,
    amount     BIGINT    NOT NULL,
    lottery_id UUID      NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    CONSTRAINT fk_frozen_money_lottery FOREIGN KEY (lottery_id) REFERENCES bc.lotteries (id),
    CONSTRAINT fk_frozen_money_wallets FOREIGN KEY (wallet_id) REFERENCES auth.wallet (id)
);

-- +migrate Down
DROP TABLE IF EXISTS auth.frozen_money;
