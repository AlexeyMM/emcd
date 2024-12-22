ALTER TABLE profile DROP CONSTRAINT profile_pkey;
ALTER TABLE profile DROP COLUMN user_id;
ALTER TABLE profile ADD COLUMN user_id uuid;
ALTER TABLE profile ADD PRIMARY KEY(user_id,coin);

ALTER TABLE profile_history DROP CONSTRAINT profile_history_pkey;
ALTER TABLE profile_history DROP COLUMN user_id;
ALTER TABLE profile_history ADD COLUMN user_id uuid;
ALTER TABLE profile_history ADD PRIMARY KEY(user_id,coin,created_at);

