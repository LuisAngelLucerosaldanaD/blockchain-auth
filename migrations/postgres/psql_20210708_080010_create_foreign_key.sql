-- +migrate Up
ALTER TABLE cfg.messages
    ADD CONSTRAINT FK_cfg_messages_type FOREIGN KEY (type_message)
        REFERENCES cfg.dictionaries (id);

ALTER TABLE auth.users_temp
    ADD CONSTRAINT FK_auth_users_temp_id_type FOREIGN KEY (id_type)
        REFERENCES cfg.dictionaries (id);

ALTER TABLE auth.users
    ADD CONSTRAINT FK_auth_users_id_type FOREIGN KEY (id_type)
        REFERENCES cfg.dictionaries (id);

ALTER TABLE auth.wallets_temp
    ADD CONSTRAINT FK_auth_wallets_temp_id_user FOREIGN KEY (id_user)
        REFERENCES auth.users (id);

ALTER TABLE auth.user_wallet
    ADD CONSTRAINT fk_user_wallet_users FOREIGN KEY (id_user)
        REFERENCES auth.users (id);

ALTER TABLE auth.user_wallet
    ADD CONSTRAINT fk_user_wallet_wallets FOREIGN KEY (id_wallet)
        REFERENCES auth.wallets (id);

ALTER TABLE auth.accounting
    ADD CONSTRAINT fk_accounting_users FOREIGN KEY (id_user)
        REFERENCES auth.users (id);

ALTER TABLE auth.accounting
    ADD CONSTRAINT fk_accounting_wallets FOREIGN KEY (id_wallet)
        REFERENCES auth.wallets (id);


-- +migrate Down
