CREATE TABLE IF NOT EXISTS ingredients_categories
(
    id        BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name      VARCHAR(255) NOT NULL UNIQUE,
    image_url VARCHAR(255) NOT NULL
);

INSERT INTO ingredients_categories (name, image_url)
VALUES ('Овощи', ''),    -- 1
       ('Фрукты', ''),   -- 2
       ('Орехи', ''),    -- 3
       ('Мучное', ''),   -- 4
       ('Молочное', ''), -- 5
       ('Мясо', ''); -- 6

CREATE TABLE IF NOT EXISTS ingredients
(
    id          BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name        VARCHAR(255) NOT NULL UNIQUE,
    image_url   VARCHAR(255) NOT NULL,
    category_id BIGINT       NOT NULL
);

CREATE INDEX idx_ingredients_category_id ON ingredients (category_id);
INSERT INTO ingredients (name, category_id, image_url)
VALUES ('Картофель', 1, ''),
       ('Морковь', 1, ''),
       ('Лук', 1, ''),
       ('Петрушка', 1, ''),
       ('Свекла', 1, ''),
       ('Огурец', 1, ''),
       ('Баклажан', 1, ''),
       ('Помидор', 1, ''),
       ('Редис', 1, ''),
       ('Сельдерей', 1, '');

INSERT INTO ingredients (name, category_id, image_url)
VALUES ('Яблоко', 2, ''),
       ('Банан', 2, ''),
       ('Груша', 2, ''),
       ('Апельсин', 2, ''),
       ('Манго', 2, ''),
       ('Папайя', 2, ''),
       ('Киви', 2, ''),
       ('Лимон', 2, ''),
       ('Грейпфрут', 2, ''),
       ('Ананас', 2, '');
INSERT INTO ingredients (name, category_id, image_url)
VALUES ('Грецкий орех', 3, ''),
       ('Фисташки', 3, ''),
       ('Макадамия', 3, ''),
       ('Кешью', 3, ''),
       ('Пекан', 3, '');
INSERT INTO ingredients (name, category_id, image_url)
VALUES ('Пшеничная мука', 4, ''),
       ('Макароны', 4, ''),
       ('Хлеб', 4, '');
INSERT INTO ingredients (name, category_id, image_url)
VALUES ('Молоко', 5, ''),
       ('Сыр', 5, ''),
       ('Творог', 5, ''),
       ('Кефир', 5, ''),
       ('Йогурт', 5, ''),
       ('Сметана', 5, ''),
       ('Ряженка', 5, ''),
       ('Масло', 5, ''),
       ('Кисломолочные продукты', 5, ''),
       ('Соус из сливок', 5, '');
INSERT INTO ingredients (name, category_id, image_url)
VALUES ('Говядина', 6, ''),
       ('Свинина', 6, ''),
       ('Курица', 6, ''),
       ('Гусь', 6, ''),
       ('Индейка', 6, ''),
       ('Лосятина', 6, ''),
       ('Мраморная телятина', 6, ''),
       ('Сердце свинины', 6, ''),
       ('Шея говядина', 6, ''),
       ('Печень курицы', 6, '');
