ALTER TABLE profile RENAME TO commissions;
ALTER TABLE profile_history RENAME TO commissions_history;

CREATE TABLE users (
    id uuid,
    username text,
    vip bool,
    segment_id int,
    ref_id int,
    PRIMARY KEY (id)
);

ALTER TABLE commissions
    ADD CONSTRAINT fk_id
        FOREIGN KEY (user_id)
            REFERENCES users(id);