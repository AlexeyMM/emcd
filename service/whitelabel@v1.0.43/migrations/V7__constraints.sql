ALTER TABLE whitelabel
    ADD CONSTRAINT name_unique UNIQUE (name);

ALTER TABLE whitelabel_owner
    ADD CONSTRAINT fk_wl_id
        FOREIGN KEY (whitelabel_id)
            REFERENCES whitelabel(id);