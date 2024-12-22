-- name: GetInvoice :one
SELECT i.*,
       COALESCE(
               (SELECT SUM(t.amount)
                FROM invoice_transaction t
                WHERE t.invoice_id = i.id
                  AND t.confirmation_status = 'confirmed'),
               0
       )::decimal as paid_amount
FROM invoice i
WHERE i.id = @id;

-- name: CreateInvoice :exec
insert into invoice (id,
                     merchant_id,
                     coin_id,
                     network_id,
                     deposit_address,
                     amount,
                     buyer_fee,
                     merchant_fee,
                     title,
                     description,
                     checkout_url,
                     status,
                     expires_at,
                     external_id,
                     buyer_email,
                     required_payment)
values (@id,
        @merchant_id,
        @coin_id,
        @network_id,
        @deposit_address,
        @amount,
        @buyer_fee,
        @merchant_fee,
        @title,
        @description,
        @checkout_url,
        @status,
        @expires_at,
        @external_id,
        @buyer_email,
        @required_payment);

-- name: GetActiveInvoiceByDepositAddressForUpdate :one
SELECT i.*,
       COALESCE(
               (SELECT SUM(t.amount)
                FROM invoice_transaction t
                WHERE t.invoice_id = i.id
                  AND t.confirmation_status = 'confirmed'),
               0
       )::decimal as paid_amount
FROM invoice i
WHERE i.deposit_address = @deposit_address
  AND i.status NOT IN ('finished', 'expired')
    FOR UPDATE;

-- name: UpdateStatus :exec
UPDATE invoice
SET status      = @status::invoice_status,
    finished_at = CASE
                      WHEN @status IN ('finished',
                                       'cancelled',
                                       'expired') THEN NOW()
                      ELSE finished_at
        END
WHERE id = @id;

-- name: SetInvoicesExpired :execrows
UPDATE invoice
SET status      = 'expired',
    finished_at = NOW()
WHERE status = 'waiting_for_deposit'
  AND expires_at < NOW()
  AND NOT EXISTS ( -- sanity check
    SELECT 1
    FROM invoice_transaction t
    WHERE t.invoice_id = invoice.id);

-- name: CountInvoiceByStatus :many
SELECT status, COUNT(1) AS invoice_count
FROM invoice
GROUP BY status;
