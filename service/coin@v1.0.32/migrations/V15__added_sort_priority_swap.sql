ALTER TABLE public.coins
    ADD COLUMN sort_priority_swap SMALLINT DEFAULT 0;

UPDATE public.coins SET sort_priority_swap = 10 WHERE id = 'ton';
UPDATE public.coins SET sort_priority_swap = 20 WHERE id = 'usdt';
UPDATE public.coins SET sort_priority_swap = 30 WHERE id = 'btc';
UPDATE public.coins SET sort_priority_swap = 40 WHERE id = 'eth';
UPDATE public.coins SET sort_priority_swap = 50 WHERE id = 'bnb';
UPDATE public.coins SET sort_priority_swap = 60 WHERE id = 'usdc';
UPDATE public.coins SET sort_priority_swap = 70 WHERE id = 'doge';
UPDATE public.coins SET sort_priority_swap = 80 WHERE id = 'matic';
UPDATE public.coins SET sort_priority_swap = 90 WHERE id = 'bel';
UPDATE public.coins SET sort_priority_swap = 100 WHERE id = 'avax';
UPDATE public.coins SET sort_priority_swap = 110 WHERE id = 'trx';
UPDATE public.coins SET sort_priority_swap = 120 WHERE id = 'ltc';
UPDATE public.coins SET sort_priority_swap = 130 WHERE id = 'bch';
UPDATE public.coins SET sort_priority_swap = 140 WHERE id = 'arb';
UPDATE public.coins SET sort_priority_swap = 150 WHERE id = 'op';
UPDATE public.coins SET sort_priority_swap = 160 WHERE id = 'etc';
UPDATE public.coins SET sort_priority_swap = 170 WHERE id = 'dash';
UPDATE public.coins SET sort_priority_swap = 180 WHERE id = 'bsv';
UPDATE public.coins SET sort_priority_swap = 190 WHERE id = 'kas';
UPDATE public.coins SET sort_priority_swap = 200 WHERE id = 'not';
UPDATE public.coins SET sort_priority_swap = 210 WHERE id = 'fb';