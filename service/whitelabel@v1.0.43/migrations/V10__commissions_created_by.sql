ALTER TABLE whitelabel_commissions
    ADD COLUMN created_by uuid;

ALTER TABLE whitelabel_history
    ADD COLUMN created_by uuid;

