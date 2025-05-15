CREATE TABLE IF NOT EXISTS users_chefs
(
    user_id BIGINT NOT NULL,
    chef_id BIGINT NOT NULL,
    PRIMARY KEY (user_id, chef_id)
)