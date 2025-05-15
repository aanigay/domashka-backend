CREATE TABLE IF NOT EXISTS user_favorite_dishes (
                                                    user_id BIGINT NOT NULL,
                                                    dish_id BIGINT NOT NULL ,
                                                    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                                                    PRIMARY KEY (user_id, dish_id)
);


COMMENT ON COLUMN user_favorite_dishes.created_at IS 'Время добавления записи о любимом блюде пользователя';
COMMENT ON COLUMN user_favorite_dishes.user_id IS 'Идентификатор пользователя';
COMMENT ON COLUMN user_favorite_dishes.dish_id IS 'Идентификатор блюда';


CREATE TABLE IF NOT EXISTS user_favorite_chefs (
                                                   user_id BIGINT NOT NULL,
                                                   chef_id BIGINT NOT NULL,
                                                   created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                                                   PRIMARY KEY (user_id, chef_id)
);

COMMENT ON COLUMN user_favorite_chefs.created_at IS 'Время добавления записи о любимом шеф-поваре пользователя';
COMMENT ON COLUMN user_favorite_chefs.user_id IS 'Идентификатор пользователя';
COMMENT ON COLUMN user_favorite_chefs.chef_id IS 'Идентификатор шефа';