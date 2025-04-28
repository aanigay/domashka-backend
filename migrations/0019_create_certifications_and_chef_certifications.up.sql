-- 1) Создаём таблицу всех доступных сертификатов
CREATE TABLE certifications (
                                id SERIAL PRIMARY KEY,                  -- Уникальный идентификатор сертификата
                                name VARCHAR(255) NOT NULL UNIQUE,      -- Название сертификата
                                description TEXT,                       -- Описание сертификата
                                created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()  -- Время создания записи
);

-- 2) Создаём связующую таблицу многие-ко-многим между chefs и certifications
CREATE TABLE chef_certifications (
                                     chef_id BIGINT NOT NULL,                -- Ссылка на повара (chefs.id)
                                     certification_id INTEGER NOT NULL,      -- Ссылка на сертификат (certifications.id)
                                     issued_at DATE,                         -- Дата выдачи сертификата повару
                                     PRIMARY KEY (chef_id, certification_id),-- Композитный первичный ключ
                                     CONSTRAINT fk_chef
                                         FOREIGN KEY (chef_id)
                                             REFERENCES chefs(id)
                                             ON DELETE CASCADE,                   -- При удалении повара удалить свои сертификаты
                                     CONSTRAINT fk_certification
                                         FOREIGN KEY (certification_id)
                                             REFERENCES certifications(id)
                                             ON DELETE CASCADE                    -- При удалении сертификата убрать все связи
);

-- 3) Добавляем комментарии к таблицам
COMMENT ON TABLE certifications IS 'Таблица всех доступных сертификатов';
COMMENT ON TABLE chef_certifications IS 'Связующая таблица многие-ко-многим между chefs и certifications';

-- 4) Добавляем комментарии к колонкам certifications
COMMENT ON COLUMN certifications.id          IS 'Уникальный идентификатор сертификата';
COMMENT ON COLUMN certifications.name        IS 'Название сертификата';
COMMENT ON COLUMN certifications.description IS 'Описание сертификата';
COMMENT ON COLUMN certifications.created_at  IS 'Время создания записи';

-- 5) Добавляем комментарии к колонкам chef_certifications
COMMENT ON COLUMN chef_certifications.chef_id          IS 'Ссылка на повара (chefs.id)';
COMMENT ON COLUMN chef_certifications.certification_id IS 'Ссылка на сертификат (certifications.id)';
COMMENT ON COLUMN chef_certifications.issued_at        IS 'Дата выдачи сертификата повару';
