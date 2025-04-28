CREATE TABLE IF NOT EXISTS chefs_experience
(
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    chef_id BIGINT NOT NULL,
    type_experience TEXT NOT NULL,
    status int NOT NULL,
    url_photo_experience_id int NOT NULL,
    meta_data json,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

COMMENT ON TABLE chefs_experience IS 'Содержит информацию о профессиональном опыте повара';

COMMENT ON COLUMN chefs_experience.id IS 'Уникальный идентификатор записи';
COMMENT ON COLUMN chefs_experience.chef_id IS 'Ссылка на идентификатор повара';
COMMENT ON COLUMN chefs_experience.type_experience IS 'Вид опыта: personal_exp – личный опыт, education_exp – образовательный опыт, work_exp – рабочий опыт';
COMMENT ON COLUMN chefs_experience.status IS 'Статус: 0 - дефолт (нет данных), 1 - документ отправлен и требует проверки, 2 - документ проверен и подтвержден';
COMMENT ON COLUMN chefs_experience.url_photo_experience_id IS 'Идентификатор фото, подтверждающих опыт';
COMMENT ON COLUMN chefs_experience.meta_data IS 'Дополнительные метаданные';
COMMENT ON COLUMN chefs_experience.created_at IS 'Время создания записи';
COMMENT ON COLUMN chefs_experience.updated_at IS 'Время последнего обновления записи';