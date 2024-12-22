ALTER TABLE smtp_settings ADD COLUMN  from_address_displayed_as text default '';

ALTER TABLE mailgun_settings ADD COLUMN  from_address_displayed_as text default '';