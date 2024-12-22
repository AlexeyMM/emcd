CREATE TABLE profile (
    user_id integer,
    fee numeric,
    whitelabel_fee numeric,
    reward numeric,
    coin character varying,
    PRIMARY KEY(user_id, coin)
);

CREATE TABLE profile_history (
                                 user_id integer,
                                 fee numeric,
                                 whitelabel_fee numeric,
                                 reward numeric,
                                 coin character varying,
                                 created_at timestamp,
                                 PRIMARY KEY(user_id, coin, created_at)
);