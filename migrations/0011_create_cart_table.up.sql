CREATE TABLE IF NOT EXISTS carts
(
    user_id    BIGINT PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS cart_items
(
    id             BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    user_id        BIGINT       NOT NULL,
    dish_id        BIGINT       NOT NULL,
    chef_id        BIGINT       NOT NULL,
    dish_size_id   BIGINT       NOT NULL,
    quantity       INT          NOT NULL,
    customer_notes VARCHAR(255) NOT NULL DEFAULT '',
    added_at       TIMESTAMP    NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS cart_item_added_ingredients
(
    cart_item_id  BIGINT NOT NULL,
    ingredient_id BIGINT NOT NULL,
    PRIMARY KEY (cart_item_id, ingredient_id)
);

CREATE TABLE IF NOT EXISTS cart_item_removed_ingredients
(
    cart_item_id  BIGINT NOT NULL,
    ingredient_id BIGINT NOT NULL,
    PRIMARY KEY (cart_item_id, ingredient_id)
);

CREATE INDEX IF NOT EXISTS idx_cart_item_user_id ON cart_items (user_id);