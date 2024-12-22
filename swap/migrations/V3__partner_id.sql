UPDATE swap.swaps
SET partner_id = ''
WHERE partner_id IS NULL;

ALTER TABLE swap.swaps
    ALTER COLUMN partner_id SET DEFAULT '',
    ALTER COLUMN partner_id SET NOT NULL;
