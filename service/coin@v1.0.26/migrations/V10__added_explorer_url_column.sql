ALTER TABLE networks
	ADD COLUMN explorer_url VARCHAR;

UPDATE networks
SET explorer_url = 'https://blockchair.com/bitcoin/transaction/%txid%'
WHERE id = 'btc';

UPDATE networks
SET explorer_url = 'https://blockchair.com/bitcoin/transaction/%txid%'
WHERE id = 'btc';

UPDATE networks
SET explorer_url = 'https://blockchair.com/bitcoin/transaction/%txid%'
WHERE id = 'bch';

UPDATE networks
SET explorer_url = 'https://blockchair.com/bitcoin-cash/transaction/%txid%'
WHERE id = 'bch';

UPDATE networks
SET explorer_url = 'https://blockchair.com/litecoin/transaction/%txid%'
WHERE id = 'ltc';

UPDATE networks
SET explorer_url = 'https://blockchair.com/dash/transaction/%txid%'
WHERE id = 'dash';

UPDATE networks
SET explorer_url = 'https://blockchair.com/dogecoin/transaction/%txid%'
WHERE id = 'doge';

UPDATE networks
SET explorer_url = 'https://etherscan.io/tx/%txid%'
WHERE id = 'erc20';

UPDATE networks
SET explorer_url = 'https://blockscout.com/etc/mainnet/tx/%txid%'
WHERE id = 'etc';

UPDATE networks
SET explorer_url = 'https://bscscan.com/tx/%txid%'
WHERE id = 'bep20';

UPDATE networks
SET explorer_url = 'https://tronscan.org/#/transaction/%txid%'
WHERE id = 'trc20';

UPDATE networks
SET explorer_url = 'https://tonscan.org/tx/%txid%'
WHERE id = 'ton';

UPDATE networks
SET explorer_url = 'https://kas.fyi/transaction/%txid%'
WHERE id = 'kas';
