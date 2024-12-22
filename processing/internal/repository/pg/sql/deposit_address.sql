-- name: SaveDepositAddress :exec
insert into deposit_address (address, network_id, merchant_id, available)
values (@address, @network_id, @merchant_id, @available)
on conflict (address)
    do update set network_id  = excluded.network_id,
                  merchant_id = excluded.merchant_id,
                  available   = excluded.available;

-- name: OccupyDepositAddress :one
UPDATE deposit_address
SET available = false
WHERE network_id = @network_id
  AND merchant_id = @merchant_id
  AND available
RETURNING address;
