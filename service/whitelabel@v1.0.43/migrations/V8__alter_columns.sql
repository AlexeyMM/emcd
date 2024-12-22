ALTER TABLE whitelabel
    DROP COLUMN logo,
    DROP COLUMN description;

ALTER TABLE whitelabel
    ADD COLUMN domain text,
    ADD COLUMN media_id uuid;


ALTER TABLE whitelabel_owner RENAME TO whitelabel_roles;

ALTER TABLE whitelabel_roles
    ADD COLUMN role text;
