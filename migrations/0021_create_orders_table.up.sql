CREATE TABLE IF NOT EXISTS orders
(
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    chef_id BIGINT NOT NULL,
    shift_id BIGINT NOT NULL,
    status INT NOT NULL DEFAULT 1,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    total_cost NUMERIC NOT NULL,
    leave_by_the_door BOOLEAN NOT NULL DEFAULT false,
    call_beforehand BOOLEAN NOT NULL DEFAULT false,
    client_address_id BIGINT NOT NULL
);

CREATE TABLE IF NOT EXISTS orders_cart_items
(
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    cart_item_id BIGINT NOT NULL,
    order_id BIGINT NOT NULL
);