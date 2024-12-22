ALTER TABLE coins
	ADD COLUMN mining_reward_type VARCHAR(20) NOT NULL DEFAULT '';

UPDATE coins SET mining_reward_type = 'FPPS' WHERE id = 'btc';
UPDATE coins SET mining_reward_type = 'FPPS' WHERE id = 'bch';
UPDATE coins SET mining_reward_type = 'FPPS' WHERE id = 'dash';
UPDATE coins SET mining_reward_type = 'FPPS' WHERE id = 'bsv';
UPDATE coins SET mining_reward_type = 'FPPS' WHERE id = 'eth';
UPDATE coins SET mining_reward_type = 'FPPS' WHERE id = 'etc';
UPDATE coins SET mining_reward_type = 'PPS+' WHERE id = 'ltc';
UPDATE coins SET mining_reward_type = 'PPLNS' WHERE id = 'doge';
