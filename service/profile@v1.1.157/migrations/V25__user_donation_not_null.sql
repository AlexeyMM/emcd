ALTER TABLE user_donation_settings ALTER COLUMN back_donation_on SET NOT NULL;

ALTER TABLE user_donation_settings ALTER COLUMN front_donation_on SET NOT NULL;

ALTER TABLE user_donation_settings ALTER COLUMN percent SET NOT NULL;
