ALTER TABLE coins DROP COLUMN sort_priority;

ALTER TABLE coins ADD COLUMN sort_priority_mining SMALLINT NOT NULL DEFAULT 0;
ALTER TABLE coins ADD COLUMN sort_priority_wallet SMALLINT NOT NULL DEFAULT 0;

UPDATE coins SET sort_priority_mining = 32767;
UPDATE coins SET sort_priority_wallet = 32767;

UPDATE coins SET sort_priority_mining = 10 WHERE id = 'btc';
UPDATE coins SET sort_priority_mining = 20 WHERE id = 'bch';
UPDATE coins SET sort_priority_mining = 30 WHERE id = 'dash';
UPDATE coins SET sort_priority_mining = 40 WHERE id = 'doge';
UPDATE coins SET sort_priority_mining = 50 WHERE id = 'etc';
UPDATE coins SET sort_priority_mining = 60 WHERE id = 'ltc';

UPDATE coins SET sort_priority_wallet = 10 WHERE id = 'usdt';
UPDATE coins SET sort_priority_wallet = 20 WHERE id = 'btc';
UPDATE coins SET sort_priority_wallet = 30 WHERE id = 'eth';
UPDATE coins SET sort_priority_wallet = 40 WHERE id = 'bch';
UPDATE coins SET sort_priority_wallet = 50 WHERE id = 'dash';
UPDATE coins SET sort_priority_wallet = 60 WHERE id = 'doge';
UPDATE coins SET sort_priority_wallet = 70 WHERE id = 'etc';
UPDATE coins SET sort_priority_wallet = 80 WHERE id = 'ltc';
UPDATE coins SET sort_priority_wallet = 90 WHERE id = 'ton';
UPDATE coins SET sort_priority_wallet = 100 WHERE id = 'usdc';
