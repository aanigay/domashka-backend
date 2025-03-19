CREATE TABLE IF NOT EXISTS dishes_categories
(
    id          BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name        VARCHAR(255) NOT NULL UNIQUE,
    description TEXT         NOT NULL
);
INSERT INTO dishes_categories (name, description)
VALUES ('Итальянская кухня', 'Блюда из Италии, такие как паста, пицца, ризотто и другие.'),
       ('Французская кухня', 'Классические французские блюда, такие как крепс, кюри, соус бешамель и др.'),
       ('Азийская кухня', 'Блюда из Азии, включая японскую, китайскую, индийскую и другие.'),
       ('Мексиканская кухня', 'Традиционные мексиканские блюда, такие как тако, чипсы, гуакамоле и др.'),
       ('Русская кухня', 'Классические русские блюда, такие как борщ, котлеты по-киевски, супы и др.'),
       ('Вегетарианская кухня', 'Блюда без мяса, включая салаты, супы, пасту и другие.'),
       ('Веганская кухня',
        'Блюда без мяса и продуктов животного происхождения, такие как салаты, супы, пасты и другие.'),
       ('Острые блюда', 'Блюда поострее.'),
       ('Неострые блюда', 'Блюда без острых специй.'),
       ('Здоровая кухня', 'Блюда, богатые полезными питательными веществами, с низким содержанием жиров и калорий.'),
       ('Морепродукты', 'Блюда из морепродуктов, такие как рыба, креветки и др.'),
       ('Десерты', 'Сладкие блюда, такие как торты, пирожные, мороженое и др.');


CREATE TABLE IF NOT EXISTS dishes
(
    id          UUID PRIMARY KEY      DEFAULT gen_random_uuid(),
    chef_id     UUID         NOT NULL,
    name        VARCHAR(255) NOT NULL UNIQUE,
    description TEXT         NOT NULL,
    price       BIGINT       NOT NULL,
    stock       BIGINT       NOT NULL,
    created_at  TIMESTAMP    NOT NULL DEFAULT current_timestamp,
    updated_at  TIMESTAMP    NOT NULL DEFAULT current_timestamp
);

CREATE TABLE IF NOT EXISTS dishes_categories_dishes
(
    id                   uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    dishes_categories_id BIGINT NOT NULL,
    dishes_id            UUID   NOT NULL
);

CREATE TABLE IF NOT EXISTS dishes_ingredients
(
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    dishes_id      UUID   NOT NULL,
    ingredients_id BIGINT NOT NULL
)