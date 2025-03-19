CREATE TABLE IF NOT EXISTS carts
(
    user_id    UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS cart_items
(
    id                         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id                    UUID NOT NULL,
    dish_id                    UUID NOT NULL,
    chef_id                    UUID NOT NULL,
    additional_ingredients_ids BIGINT ARRAY NOT NULL,
    removed_ingredients_ids    BIGINT ARRAY NOT NULL ,
    added_at                   TIMESTAMP NOT NULL DEFAULT now(),
    customer_notes             VARCHAR(255) NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_cart_item_user_id ON cart_items (user_id);