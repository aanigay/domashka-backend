CREATE TABLE IF NOT EXISTS ordered_items
(
    id             BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    user_id        BIGINT       NOT NULL,
    dish_id        BIGINT       NOT NULL,
    chef_id        BIGINT       NOT NULL,
    order_id       BIGINT       NOT NULL,
    dish_size_id   BIGINT       NOT NULL,
    quantity       INT          NOT NULL,
    customer_notes VARCHAR(255) NOT NULL DEFAULT '',
    added_at       TIMESTAMP    NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS ordered_item_added_ingredients
(
    ordered_item_id  BIGINT NOT NULL,
    ingredient_id BIGINT NOT NULL,
    PRIMARY KEY (ordered_item_id, ingredient_id)
);

CREATE TABLE IF NOT EXISTS ordered_item_removed_ingredients
(
    ordered_item_id  BIGINT NOT NULL,
    ingredient_id BIGINT NOT NULL,
    PRIMARY KEY (ordered_item_id, ingredient_id)
);
