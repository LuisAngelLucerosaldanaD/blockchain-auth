-- +migrate Up
ALTER TABLE cfg.messages
    ADD CONSTRAINT FK_cfg_messages_type FOREIGN KEY (type_message)
        REFERENCES cfg.dictionaries (id);

ALTER TABLE auth.accounting
    ADD CONSTRAINT fk_accounting_wallets FOREIGN KEY (wallet_id)
        REFERENCES auth.wallet (id);


-- +migrate Down
