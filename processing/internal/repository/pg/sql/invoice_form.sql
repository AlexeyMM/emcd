-- name: CreateInvoiceForm :exec
INSERT INTO invoice_form (
    id,
    merchant_id,
    coin_id,
    network_id,
    amount,
    title,
    description,
    buyer_email,
    checkout_url,
    external_id,
    expires_at
) VALUES (
    @id,
    @merchant_id,
    @coin_id,
    @network_id,
    @amount,
    @title,
    @description,
    @buyer_email,
    @checkout_url,
    @external_id,
    @expires_at
);

-- name: GetInvoiceForm :one
SELECT * FROM invoice_form WHERE id = @id;