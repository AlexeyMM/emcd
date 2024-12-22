insert into emcd.transaction_types (id, description, account_type_id)
values (95, '[P2P2] заморозка крипты при создании ордера (wallet -> p2p)', 1),
       (96, '[P2P2] разморозка крипты при отмене ордера (p2p -> wallet)', 7),
       (97, '[P2P2] комиссия в пользу emcd за успешную p2p-сделку', 1),
       (98, '[P2P2] комиссия в пользу рефовода за успешную p2p-сделку', 1) on conflict do nothing;