TRUNCATE TABLE users;
INSERT INTO public.users (username, alias, first_name, second_name, last_name, email, number_phone,
                          status, external_type, telegram_name, external_id, notification_flag,
                          role, birthday, is_spam, sms_attempts, last_sms_request, name, chat_id)
VALUES ('ivan_ivanov', 'ivan_i', 'Иван', 'Сергеевич', 'Иванов', 'ivan.ivanov@example.com', '+79161234567',
        1, 0, 'ivan_telega', 'ext001', 1, 'client', '1988-03-15', 0, 0, NULL, 'Иван Иванов', 'chat001'),
       ('anna_petrova', 'anna_p', 'Анна', 'Владимировна', 'Петрова', 'anna.petrova@example.com', '+79262345678',
        1, 1, 'anna_telega', 'ext002', 1, 'chef', '1992-07-22', 1, 2, '2023-10-05 12:30:00', 'Анна Петрова', 'chat002'),
       ('sergey_smirnov', 'sergey_s', 'Сергей', NULL, 'Смирнов', 'sergey.smirnov@example.com', '+79373456789',
        0, 2, NULL, 'ext003', 0, 'client', '1985-11-10', 0, 1, '2023-09-20 15:45:00', 'Сергей Смирнов', NULL),
       ('elena_kuznetsova', 'elena_k', 'Елена', 'Александровна', 'Кузнецова', 'elena.kuznetsova@example.com',
        '+79484567890',
        1, 0, 'elena_telega', 'ext004', 1, 'admin', '1979-05-28', 0, 0, NULL, 'Елена Кузнецова', 'chat004'),
       ('dmitry_voronov', 'dmitry_v', 'Дмитрий', 'Игоревич', 'Воронов', 'dmitry.voronov@example.com', '+79595678901',
        1, 1, 'dmitry_telega', 'ext005', 1, 'chef', '1990-09-03', 1, 0, NULL, 'Дмитрий Воронов', 'chat005');

TRUNCATE TABLE carts;
INSERT INTO carts (user_id)
VALUES (1),
       (2),
       (3),
       (4),
       (5);

TRUNCATE TABLE chefs_experience RESTART IDENTITY CASCADE;
TRUNCATE TABLE chef_certifications RESTART IDENTITY CASCADE;
TRUNCATE TABLE chefs RESTART IDENTITY CASCADE;
INSERT INTO chefs (name, image_url)
VALUES ('Иван Петров', 'https://example.com/images/ivan_petrov.jpg');
INSERT INTO chefs (name, image_url)
VALUES ('Анна Сидорова', 'https://example.com/images/anna_sidorova.jpg');
INSERT INTO chefs (name, image_url)
VALUES ('Сергей Иванов', 'https://example.com/images/sergei_ivanov.jpg');
INSERT INTO chefs (name, image_url)
VALUES ('Мария Козлова', 'https://example.com/images/maria_kozlova.jpg');

TRUNCATE TABLE chef_ratings;
INSERT INTO chef_ratings (chef_id, rating, reviews_count)
VALUES (1, 4.5, 10),
       (2, 4.2, 15),
       (3, 4.8, 20),
       (4, 4.0, 10),
       (5, 4.3, 18);

TRUNCATE TABLE chefs_experience;
INSERT INTO chefs_experience (chef_id, type_experience, status, url_photo_experience_id, meta_data, created_at,
                              updated_at, experience_years)
VALUES (1, 'Опыт работы', 1, 1, null, '2022-01-01 00:00:00', '2022-01-01 00:00:00', 5),
       (2, 'Опыт работы', 1, 1, null, '2022-01-01 00:00:00', '2022-01-01 00:00:00', 5);

TRUNCATE TABLE dishes RESTART IDENTITY CASCADE;
INSERT INTO dishes (chef_id, name, description, image_url)
VALUES (1, 'Цезарь с курицей', 'Салат Цезарь с курицей на гриле.', 'https://example.com/images/caesar_chicken.jpg'),
       (2, 'Паста Болоньезе', 'Классическая итальянская паста с мясным соусом.',
        'https://example.com/images/pasta_bolognese.jpg'),
       (3, 'Суп Харчо', 'Острый грузинский суп с говядиной и пряностями.',
        'https://example.com/images/soup_harcho.jpg'),
       (1, 'Борщ', 'Традиционный украинский борщ со свеклой и мясом.', 'https://example.com/images/borscht.jpg'),
       (4, 'Суши ассорти', 'Набор различных суши с рыбой и овощами.', 'https://example.com/images/sushi_assort.jpg'),
       (2, 'Стейк Рибай', 'Сочный стейк рибай с овощным гарниром.', 'https://example.com/images/ribeye_steak.jpg')
ON CONFLICT DO NOTHING;

TRUNCATE TABLE ingredients_categories RESTART IDENTITY CASCADE;
INSERT INTO ingredients_categories (name, image_url)
VALUES ('Овощи', ''),    -- 1
       ('Фрукты', ''),   -- 2
       ('Орехи', ''),    -- 3
       ('Мучное', ''),   -- 4
       ('Молочное', ''), -- 5
       ('Мясо', ''),-- 6
       ('Крупа', ''),    -- 7
       ('Другое', '')    -- 8
ON CONFLICT DO NOTHING;

-------

TRUNCATE TABLE ingredients RESTART IDENTITY CASCADE;
INSERT INTO ingredients (name, category_id, image_url)
VALUES ('Картофель', 1, 'https://s3.twcstorage.ru/1bd3f8a3-s3-public/ingredients/potato.jpg'),
       ('Морковь', 1, 'https://s3.twcstorage.ru/1bd3f8a3-s3-public/ingredients/carrot.jpg'),
       ('Лук', 1, 'https://s3.twcstorage.ru/1bd3f8a3-s3-public/ingredients/onion.jpg'),
       ('Петрушка', 1, 'https://s3.twcstorage.ru/1bd3f8a3-s3-public/ingredients/greens.png'),
       ('Свекла', 1, 'https://s3.twcstorage.ru/1bd3f8a3-s3-public/ingredients/beetroot.jpg'),
       ('Огурец', 1, 'https://s3.twcstorage.ru/1bd3f8a3-s3-public/ingredients/'),
       ('Баклажан', 1, 'https://s3.twcstorage.ru/1bd3f8a3-s3-public/ingredients/'),
       ('Томат', 1, 'https://s3.twcstorage.ru/1bd3f8a3-s3-public/ingredients/tomato.jpg'),
       ('Редис', 1, 'https://s3.twcstorage.ru/1bd3f8a3-s3-public/ingredients/'),
       ('Сельдерей', 1, 'https://s3.twcstorage.ru/1bd3f8a3-s3-public/ingredients/'),
       ('Чеснок', 1, 'https://s3.twcstorage.ru/1bd3f8a3-s3-public/ingredients/garlic.jpg'),
       ('Лавровый лист', 1, 'https://s3.twcstorage.ru/1bd3f8a3-s3-public/ingredients/bay_leaf.jpg'),
       ('Капуста', 1, 'https://s3.twcstorage.ru/1bd3f8a3-s3-public/ingredients/cabbage.jpeg'),
       ('Салат', 1, 'https://s3.twcstorage.ru/1bd3f8a3-s3-public/ingredients/romano_salad.jpg')
ON CONFLICT DO NOTHING;

-------

INSERT INTO ingredients (name, category_id, image_url)
VALUES ('Яблоко', 2, 'https://s3.twcstorage.ru/1bd3f8a3-s3-public/ingredients/'),
       ('Банан', 2, 'https://s3.twcstorage.ru/1bd3f8a3-s3-public/ingredients/'),
       ('Груша', 2, 'https://s3.twcstorage.ru/1bd3f8a3-s3-public/ingredients/'),
       ('Апельсин', 2, 'https://s3.twcstorage.ru/1bd3f8a3-s3-public/ingredients/'),
       ('Манго', 2, 'https://s3.twcstorage.ru/1bd3f8a3-s3-public/ingredients/'),
       ('Папайя', 2, 'https://s3.twcstorage.ru/1bd3f8a3-s3-public/ingredients/'),
       ('Киви', 2, 'https://s3.twcstorage.ru/1bd3f8a3-s3-public/ingredients/'),
       ('Лимон', 2, 'https://s3.twcstorage.ru/1bd3f8a3-s3-public/ingredients/'),
       ('Грейпфрут', 2, 'https://s3.twcstorage.ru/1bd3f8a3-s3-public/ingredients/'),
       ('Ананас', 2, 'https://s3.twcstorage.ru/1bd3f8a3-s3-public/ingredients/')
ON CONFLICT DO NOTHING;

-------

INSERT INTO ingredients (name, category_id, image_url)
VALUES ('Грецкий орех', 3, 'https://s3.twcstorage.ru/1bd3f8a3-s3-public/ingredients/'),
       ('Фисташки', 3, 'https://s3.twcstorage.ru/1bd3f8a3-s3-public/ingredients/'),
       ('Макадамия', 3, 'https://s3.twcstorage.ru/1bd3f8a3-s3-public/ingredients/'),
       ('Кешью', 3, 'https://s3.twcstorage.ru/1bd3f8a3-s3-public/ingredients/'),
       ('Пекан', 3, 'https://s3.twcstorage.ru/1bd3f8a3-s3-public/ingredients/')
ON CONFLICT DO NOTHING;

-------

INSERT INTO ingredients (name, category_id, image_url)
VALUES ('Пшеничная мука', 4, 'https://s3.twcstorage.ru/1bd3f8a3-s3-public/ingredients/'),
       ('Макароны', 4, 'https://s3.twcstorage.ru/1bd3f8a3-s3-public/ingredients/macaroni.jpg'),
       ('Хлеб', 4, 'https://s3.twcstorage.ru/1bd3f8a3-s3-public/ingredients/bread.jpg')
ON CONFLICT DO NOTHING;

-------

INSERT INTO ingredients (name, category_id, image_url)
VALUES ('Молоко', 5, 'https://s3.twcstorage.ru/1bd3f8a3-s3-public/ingredients/'),
       ('Сыр', 5, 'https://s3.twcstorage.ru/1bd3f8a3-s3-public/ingredients/'),
       ('Творог', 5, 'https://s3.twcstorage.ru/1bd3f8a3-s3-public/ingredients/'),
       ('Кефир', 5, 'https://s3.twcstorage.ru/1bd3f8a3-s3-public/ingredients/'),
       ('Йогурт', 5, 'https://s3.twcstorage.ru/1bd3f8a3-s3-public/ingredients/'),
       ('Сметана', 5, 'https://s3.twcstorage.ru/1bd3f8a3-s3-public/ingredients/'),
       ('Ряженка', 5, 'https://s3.twcstorage.ru/1bd3f8a3-s3-public/ingredients/'),
       ('Масло', 5, 'https://s3.twcstorage.ru/1bd3f8a3-s3-public/ingredients/'),
       ('Кисломолочные продукты', 5, 'https://s3.twcstorage.ru/1bd3f8a3-s3-public/ingredients/'),
       ('Соус из сливок', 5, 'https://s3.twcstorage.ru/1bd3f8a3-s3-public/ingredients/'),
       ('Сливочное масло', 5, 'https://s3.twcstorage.ru/1bd3f8a3-s3-public/ingredients/butter.jpg'),
       ('Сыр пармезан', 5, 'https://s3.twcstorage.ru/1bd3f8a3-s3-public/ingredients/parmesan.png')
ON CONFLICT DO NOTHING;

-------

INSERT INTO ingredients (name, category_id, image_url)
VALUES ('Говядина', 6, 'https://s3.twcstorage.ru/1bd3f8a3-s3-public/ingredients/beef.jpg'),
       ('Свинина', 6, 'https://s3.twcstorage.ru/1bd3f8a3-s3-public/ingredients/'),
       ('Курица', 6, 'https://s3.twcstorage.ru/1bd3f8a3-s3-public/ingredients/chicken.png'),
       ('Гусь', 6, 'https://s3.twcstorage.ru/1bd3f8a3-s3-public/ingredients/'),
       ('Индейка', 6, 'https://s3.twcstorage.ru/1bd3f8a3-s3-public/ingredients/'),
       ('Лосятина', 6, 'https://s3.twcstorage.ru/1bd3f8a3-s3-public/ingredients/'),
       ('Мраморная телятина', 6, 'https://s3.twcstorage.ru/1bd3f8a3-s3-public/ingredients/'),
       ('Сердце свинины', 6, 'https://s3.twcstorage.ru/1bd3f8a3-s3-public/ingredients/'),
       ('Шея говядины', 6, 'https://s3.twcstorage.ru/1bd3f8a3-s3-public/ingredients/'),
       ('Печень курицы', 6, 'https://s3.twcstorage.ru/1bd3f8a3-s3-public/ingredients/'),
       ('Морепродукты', 6, 'https://s3.twcstorage.ru/1bd3f8a3-s3-public/ingredients/seafood.jpg'),
       ('Рыба', 6, 'https://s3.twcstorage.ru/1bd3f8a3-s3-public/ingredients/fish.jpg')
ON CONFLICT DO NOTHING;

-------

INSERT INTO ingredients (name, category_id, image_url)
VALUES ('Рис', 7, 'https://s3.twcstorage.ru/1bd3f8a3-s3-public/ingredients/rice.jpg')
ON CONFLICT DO NOTHING;

-------

INSERT INTO ingredients (name, category_id, image_url)
VALUES ('Васаби', 8, 'https://s3.twcstorage.ru/1bd3f8a3-s3-public/ingredients/vasabi.jpg'),
       ('Специи', 8, 'https://s3.twcstorage.ru/1bd3f8a3-s3-public/ingredients/spices.png'),
       ('Нори', 8, 'https://s3.twcstorage.ru/1bd3f8a3-s3-public/ingredients/nori.jpg'),
       ('Соевый соус', 8, 'https://s3.twcstorage.ru/1bd3f8a3-s3-public/ingredients/soy_souce.png'),
       ('Соус для цезаря', 8, 'https://s3.twcstorage.ru/1bd3f8a3-s3-public/ingredients/cezar_souce.jpg')
ON CONFLICT DO NOTHING;

TRUNCATE TABLE dish_ratings;
INSERT INTO dish_ratings (dish_id, rating, reviews_count)
VALUES (1, 4.5, 20),
       (2, 4.2, 15),
       (3, 4.8, 25),
       (4, 4.0, 10),
       (5, 4.3, 18),
       (6, 4.7, 22)
ON CONFLICT DO NOTHING;


TRUNCATE TABLE dish_sizes;
INSERT INTO dish_sizes (dish_id, label, weight_value, weight_unit, price_value, price_currency)
VALUES (1, 'S', 250, 'г', 350, 'RUB'),
       (2, 'L', 300, 'г', 400, 'RUB'),
       (3, 'S', 350, 'г', 300, 'RUB'),
       (4, 'S', 400, 'г', 450, 'RUB'),
       (5, 'L', 200, 'г', 500, 'RUB'),
       (6, 'M', 350, 'г', 700, 'RUB');

TRUNCATE TABLE dish_ingredients RESTART IDENTITY CASCADE;
INSERT INTO dish_ingredients (dish_id, ingredient_id, is_removable)
VALUES
    -- Блюдо 1: Салат Цезарь
    (1, 14, false), -- Салат
    (1, 47, true),  -- Курица
    (1, 32, true),  -- Хлеб
    (1, 44, true),  -- Сыр пармезан
    (1, 62, true),  -- Соус для цезаря
    (1, 8, true),   -- Томат (помидоры)

    -- Блюдо 2: Паста Болоньезе
    (2, 31, false), -- Макароны (паста)
    (2, 45, false), -- Говядина
    (2, 8, false),  -- Томат
    (2, 3, true),   -- Лук
    (2, 2, true),   -- Морковь
    (2, 11, true),  -- Чеснок
    (2, 59, true),  -- Специи

    -- Блюдо 3: Суп Харчо
    (3, 45, false), -- Говядина
    (3, 57, false), -- Рис
    (3, 3, true),   -- Лук
    (3, 2, true),   -- Морковь
    (3, 8, false),  -- Томат
    (3, 11, true),  -- Чеснок
    (3, 12, true),  -- Лавровый лист
    (3, 59, true),  -- Специи
    (3, 4, true),   -- Петрушка (зелень)

    -- Блюдо 4: Борщ
    (4, 45, false),
    (4, 5, false),
    (4, 13, true),
    (4, 2, true),
    (4, 3, true),
    (4, 8, true),
    (4, 11, true),
    (4, 59, true),

    -- Блюдо 5: Суши ассорти
    (5, 57, false), -- Рис
    (5, 60, false), -- Нори
    (5, 56, true),  -- Рыба
    (5, 55, true),  -- Морепродукты
    (5, 61, true),  -- Соевый соус
    (5, 58, true),  -- Васаби

    -- Блюдо 6: Стейк Рибай
    (6, 45, false), -- Говядина
    (6, 43, false), -- Сливочное масло
    (6, 1, true),   -- Картофель
    (6, 11, true),  -- Чеснок
    (6, 59, true),  -- Специи
    (6, 4, true); -- Петрушка (зелень)


TRUNCATE TABLE nutritions;
INSERT INTO nutritions (dish_id, calories, protein, fat, carbohydrates)
VALUES (1, 600, 30, 20, 60),
       (2, 800, 35, 25, 90),
       (3, 500, 25, 15, 70),
       (4, 550, 28, 18, 80),
       (5, 400, 20, 10, 50),
       (6, 750, 45, 30, 50);


TRUNCATE TABLE client_addresses;
INSERT INTO client_addresses (client_id, address_type, name, full_address, comment, geom)
VALUES (1, 'домашний', 'Главный дом', 'ул. Ленина, д. 10, Москва', 'Тестовый адрес',
        ST_SetSRID(ST_MakePoint(37.6173, 55.7558), 4326)),
       (2, 'рабочий', 'Офис', 'пр. Мира, д. 15, Москва', NULL, ST_SetSRID(ST_MakePoint(37.6225, 55.7522), 4326));


TRUNCATE TABLE notifications;
INSERT INTO notifications (user_id, channel, scenario, subject, message, recipient, status, created_at, send_attempts)
VALUES (1, 'email', 'order_confirmation', 'Заказ подтвержден', 'Ваш заказ принят и обрабатывается.', 'ivan@example.com',
        'created', CURRENT_TIMESTAMP, 0),
       (2, 'sms', 'delivery', 'Заказ отправлен', 'Ваш заказ в пути.', '+70000000000', 'sent', CURRENT_TIMESTAMP, 1);

TRUNCATE TABLE chef_certifications RESTART IDENTITY CASCADE;
TRUNCATE TABLE certifications RESTART IDENTITY CASCADE;

INSERT INTO certifications (name, description, created_at)
VALUES ('Мед. книжка', 'Действующая медицинская книжка для работы с продуктами', NOW()),
       ('Работа в ресторане', 'Прохождение санитарно-гигиенического минимума', NOW()),
       ('Образование', 'Сертификат о прохождении курса по приготовлению суши', NOW()),
       ('Бабушка', 'Курс по приготовлению кофе и напитков на его основе', NOW());

-- Привязка сертификатов к шефам
INSERT INTO chef_certifications (chef_id, certification_id, issued_at)
VALUES (1, 1, '2024-01-15'),
       (1, 3, '2024-02-10'),

       (2, 2, '2023-11-05'),
       (2, 4, '2023-12-20'),

       (3, 1, '2024-03-01'),
       (3, 2, '2024-03-05'),
       (3, 3, '2024-03-10'),
       (3, 4, '2024-03-15');

TRUNCATE chef_addresses RESTART IDENTITY CASCADE;
INSERT INTO chef_addresses (chef_id, geom, full_address)
VALUES (1, ST_Makepoint(37.595329, 55.752248), 'улица Новый Арбат, 11с1'),
       (2, ST_Makepoint(37.631092, 55.765357), 'улица Большая Лубянка, 24/15с1'),
       (3, ST_Makepoint(37.656676, 55.736164), 'улица Большие Каменщики, 19'),
       (4, ST_Makepoint(37.594872, 55.708900), 'улица Орджоникидзе, 11');


TRUNCATE dish_categories RESTART IDENTITY CASCADE;
INSERT INTO dish_categories (title)
VALUES ('Супы'),
       ('Салаты'),
       ('Второе');

TRUNCATE dishes_dish_categories RESTART IDENTITY CASCADE;
INSERT INTO dishes_dish_categories (dish_id, category_id)
VALUES (1, 2),
       (2, 3),
       (3, 1),
       (4, 1),
       (5, 3),
       (6, 3);