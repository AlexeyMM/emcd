CREATE TABLE referral.promocodes_updates (
    user_id uuid,
    action_id uuid,
    coin text,
    created_at timestamp,
    PRIMARY KEY (user_id, action_id, coin)
);