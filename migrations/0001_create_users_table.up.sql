CREATE TABLE IF NOT EXISTS users
(
    id                UUID PRIMARY KEY                                 DEFAULT gen_random_uuid(),                   -- Уникальный идентификатор
    username          VARCHAR(256) UNIQUE                     NOT NULL,                                             -- Уникальный логин
    alias             VARCHAR(256)                            NOT NULL,                                             -- Полное имя (ФИО)
    first_name        VARCHAR(256)                            NOT NULL,                                             -- Имя
    second_name       VARCHAR(256),                                                                                 -- Отчество
    last_name         VARCHAR(256),                                                                                 -- Фамилия
    email             VARCHAR(256),                                                                                 -- Электронная почта
    number_phone      VARCHAR(256),                                                                                 -- Телефон
    status            INT CHECK (status IN (0, 1))            NOT NULL,                                             -- Состояние (0 = не бан, 1 = бан)
    external_type     INT CHECK (external_type IN (0, 1, 2))  NOT NULL,                                             -- Внешний статус
    telegram_name     VARCHAR(256),                                                                                 -- Telegram username
    external_id       VARCHAR(256),                                                                                 -- Внешний идентификатор (тип данных зависит от контекста)
    notification_flag INT CHECK (notification_flag IN (0, 1)) NOT NULL,                                             -- Отключение уведомлений
    role              VARCHAR(256)                            NOT NULL CHECK (role IN ('client', 'chef', 'admin')), -- Роль (ENUM)
    birthday          DATE,                                                                                         -- День рождения
    created_at        TIMESTAMP                               NOT NULL DEFAULT now(),                               -- Время создания записи
    updated_at        TIMESTAMP                               NOT NULL DEFAULT now()                                -- Время последнего обновления записи
);

CREATE INDEX idx_users_username ON users (username);
CREATE INDEX idx_users_email ON users (email);
