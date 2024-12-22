-- name: CreateTransaction :exec
INSERT INTO invoice_transaction (
    invoice_id,
    hash,
    amount,
    received_address,
    confirmation_status
) VALUES (
    @invoice_id,
    @hash,
    @amount,
    @received_address,
    @confirmation_status
);

-- name: SaveTransaction :exec
INSERT INTO invoice_transaction (
    invoice_id,
    hash,
    amount,
    received_address,
    confirmation_status
) VALUES (
    @invoice_id,
    @hash,
    @amount,
    @received_address,
    @confirmation_status
) ON CONFLICT (hash) DO UPDATE SET
    invoice_id = EXCLUDED.invoice_id,
    amount = EXCLUDED.amount,
    received_address = EXCLUDED.received_address,
    confirmation_status = EXCLUDED.confirmation_status;

-- name: GetInvoiceTransactions :many
SELECT *
FROM invoice_transaction
WHERE invoice_id = @invoice_id
ORDER BY created_at; 