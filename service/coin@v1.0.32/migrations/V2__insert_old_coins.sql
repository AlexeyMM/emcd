INSERT INTO s_coin.t_wallet_coins VALUES
     ('btc', 'BTC', 'Bitcoin', 'https://emcd.io/static/coins/btc.svg', 8, 0.0005, 0.0001, true, now(), now()),
     ('bch', 'BCH', 'Bitcoin Cash', 'https://emcd.io/static/coins/bch.svg', 8, 0.001, 0.01, true, now(), now()),
     ('ltc', 'LTC', 'Litecoin', 'https://emcd.io/static/coins/ltc.svg', 8, 0.0017, 0.1, true, now(), now()),
     ('dash', 'DASH', 'Dash', 'https://emcd.io/static/coins/dash.svg', 8, 0.002, 0.0001, true, now(), now()),
     ('eth', 'ETH', 'Ethereum', 'https://emcd.io/static/coins/eth.svg', 8, 0, 0, true, now(), now()),
     ('etc', 'ETC', 'Ethereum Classic', 'https://emcd.io/static/coins/etc.svg', 8, 0.01, 0.1, true, now(), now()),
     ('doge', 'DOGE', 'Dogecoin', 'https://emcd.io/static/coins/doge.svg', 8, 5, 1, true, now(), now()),
     ('usdt', 'USDT', 'USDT', 'https://emcd.io/static/coins/usdt.svg', 8, 1, 1, true, now(), now()),
     ('usdc', 'USDC', 'USDC', 'https://emcd.io/static/coins/usdc.svg', 8, 1, 1, true, now(), now()),
     ('ton', 'TON', 'Toncoin', 'https://emcd.io/static/coins/ton.svg', 8, 0.1, 1, true, now(), now());

INSERT INTO s_coin.t_wallet_coin_networks VALUES
    ('usdt', 'usdt-bep20', 'USDT-BEP20', 'BEP20 USDT token, contract address', true, now(), now()),
    ('usdc', 'usdc-bep20', 'USDC-BEP20', 'BEP20 USDC token, contract address', true, now(), now()),
    ('usdt', 'usdt-trc20', 'USDT-TRC20', 'TRC20 USDT token, contract address', true, now(), now()),
    ('usdc', 'usdc-trc20', 'USDC-TRC20', 'TRC20 USDC token, contract address', true, now(), now());



INSERT INTO s_coin.t_mining_coins VALUES
    ('btc', 'BTC', 'Bitcoin', 0, '', true, now(), now()),
    ('bch', 'BCH', 'Bitcoin Cash', 0, '', true, now(), now()),
    ('ltc', 'LTC', 'Litecoin', 0, 'doge', true, now(), now()),
    ('dash', 'DASH', 'Dash', 0, '', true, now(), now()),
    ('etc', 'ETC', 'Ethereum Classic', 0, '', true, now(), now()),
    ('doge', 'DOGE', 'Dogecoin', 0, 'ltc', true, now(), now());