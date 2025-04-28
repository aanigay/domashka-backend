package chefs

import (
	"encoding/json"
	"time"
)

// ChefExperience содержит информацию о профессиональном опыте повара.
type ChefExperience struct {
	// ID - Уникальный идентификатор записи.
	ID int64 `db:"id" json:"id"`

	// ChefID - Ссылка на идентификатор повара.
	ChefID int64 `db:"chef_id" json:"chef_id"`

	// TypeExperience - Вид опыта: personal_exp – личный опыт, education_exp – образовательный опыт, work_exp – рабочий опыт.
	TypeExperience string `db:"type_experience" json:"type_experience"`

	// Status - Статус: 0 - дефолт (нет данных), 1 - документ отправлен и требует проверки, 2 - документ проверен и подтвержден.
	Status int `db:"status" json:"status"`

	// URLPhotoExperienceID - Идентификатор фото, подтверждающих опыт.
	URLPhotoExperienceID int `db:"url_photo_experience_id" json:"url_photo_experience_id"`

	// ExperienceYears - Количество лет опыта.
	ExperienceYears int `db:"experience_years" json:"experience_years"`

	// MetaData - Дополнительные метаданные в формате JSON.
	MetaData json.RawMessage `db:"meta_data" json:"meta_data"`

	// CreatedAt - Время создания записи.
	CreatedAt time.Time `db:"created_at" json:"created_at"`

	// UpdatedAt - Время последнего обновления записи.
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
