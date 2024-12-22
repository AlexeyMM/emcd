INSERT INTO networks (id, is_active)
VALUES ('btc', TRUE),
	   ('bch', TRUE),
	   ('ltc', TRUE),
	   ('dash', TRUE),
	   ('doge', TRUE),

	   ('erc20', TRUE),
	   ('etc', TRUE),

	   ('bep20', TRUE),

	   ('trc20', TRUE),

	   ('ton', TRUE),
	   ('kas', TRUE)
;

INSERT INTO coins (id, sort_priority, is_active, is_withdrawals_disabled, media_url)
VALUES ('btc', 1, TRUE, FALSE, 'https://emcd.io/static/coins/btc.svg'),
	   ('bch', 2, TRUE, FALSE, 'https://emcd.io/static/coins/bch.svg'),
	   ('ltc', 3, TRUE, FALSE, 'https://emcd.io/static/coins/ltc.svg'),
	   ('dash', 4, TRUE, FALSE, 'https://emcd.io/static/coins/dash.svg'),
	   ('doge', 5, TRUE, FALSE, 'https://emcd.io/static/coins/doge.svg'),

	   ('eth', 6, TRUE, FALSE, 'https://emcd.io/static/coins/eth.svg'),
	   ('etc', 7, TRUE, FALSE, 'https://emcd.io/static/coins/etc.svg'),

	   ('usdt', 8, TRUE, FALSE, 'https://emcd.io/static/coins/usdt.svg'),
	   ('usdc', 9, TRUE, FALSE, 'https://emcd.io/static/coins/usdс.svg'),

	   ('ton', 10, TRUE, FALSE, 'https://emcd.io/static/coins/ton.svg'),
	   ('kas', 11, FALSE, TRUE, ''),
	   ('trx', 999999, FALSE, TRUE, ''),
	   ('bnb', 999999, FALSE, TRUE, '')
;

INSERT INTO coins_networks (coin_id, network_id, is_active, title, description, contract_address, decimals, is_wallet,
							withdrawal_fee, withdrawal_min_limit, withdrawal_fee_ttl_seconds, is_mining,
							is_free_withdraw, mining_fee, is_withdrawals_disabled, withdrawals_disabled_description)
VALUES ('btc', 'btc', TRUE, 'btc', 'btc', NULL, 8, TRUE, 0.0005, 0.0001, 0, TRUE, TRUE, 0.015, FALSE, NULL),
	   ('bch', 'bch', TRUE, 'bch', 'bch', NULL, 8, TRUE, 0.001, 0.01, 0, TRUE, TRUE, 0.015, FALSE, NULL),
	   ('ltc', 'ltc', TRUE, 'ltc', 'ltc', NULL, 8, TRUE, 0.0017, 0.1, 0, TRUE, TRUE, 0.015, FALSE, NULL),
	   ('dash', 'dash', TRUE, 'dash', 'dash', NULL, 8, TRUE, 0.002, 0.0001, 0, TRUE, TRUE, 0.015, FALSE, NULL),
	   ('doge', 'doge', TRUE, 'doge', 'doge', NULL, 8, TRUE, 5, 1, 0, TRUE, TRUE, 0.015, FALSE, NULL),

	   ('kas', 'kas', TRUE, 'kaspa', 'kaspa', NULL, 8, TRUE, 30, 60, 0, TRUE, TRUE, 0.015, FALSE, NULL),

	   ('eth', 'erc20', TRUE, 'eth', 'eth', NULL, 18, TRUE, 0, 0, 30, FALSE, FALSE, NULL, FALSE, FALSE),
	   ('etc', 'etc', TRUE, 'etc', 'etc', NULL, 18, TRUE, 0.01, 0.1, 0, TRUE, FALSE, 0.015, FALSE, FALSE),

	   ('ton', 'ton', TRUE, 'ton', 'ton', NULL, 9, TRUE, 0.1, 1, 0, FALSE, FALSE, NULL, FALSE, NULL),

	   ('usdt', 'bep20', TRUE, 'usdt-bep20', 'usdt-bep20', '0x55d398326f99059ff775485246999027b3197955', 18, TRUE, 1, 1,
		0, FALSE, FALSE, NULL, FALSE, NULL),
	   ('usdc', 'bep20', TRUE, 'usdc-bep20', 'usdc-bep20', '0x8ac76a51cc950d9822d68b83fe1ad97b32cd580d', 18, TRUE, 1, 1,
		0, FALSE, FALSE, NULL, FALSE, NULL),

	   ('usdt', 'trc20', TRUE, 'usdt-trc20', 'usdt-trc20', 'TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t', 6, TRUE, 1, 1, 0,
		FALSE, FALSE, NULL, FALSE, NULL),
	   ('usdc', 'trc20', TRUE, 'usdc-trc20', 'usdc-trc20', 'TEkxiTehnzSmSe2XqrBj4w32RUN966rdz8', 6, TRUE, 1, 1, 0,
		FALSE, FALSE, NULL, FALSE, NULL)
;
