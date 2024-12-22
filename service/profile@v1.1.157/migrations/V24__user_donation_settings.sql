CREATE TABLE user_donation_settings (
    user_id uuid,
    back_donation_on bool,
    front_donation_on bool,
    percent numeric,
    PRIMARY KEY (user_id)
);