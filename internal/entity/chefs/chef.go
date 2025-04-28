package chefs

import (
	"database/sql"
	"time"
)

// Chef представляет данные шеф-повара.
// Поля RoleID, IsSelfEmployed, IsArchive, IsBlock соответствуют описанию в DDL.
type Chef struct {
	ID             int64     `db:"id"`
	Name           string    `db:"name"`
	ImageURL       string    `db:"image_url"`
	RoleID         *int      `db:"role_id"`          // Идентификатор роли шефа (связь с таблицей roles)
	IsSelfEmployed *bool     `db:"is_self_employed"` // Флаг, указывающий на то, что шеф-повар работает как самозанятый
	IsArchive      *bool     `db:"is_archive"`       // Флаг, указывающий на активность профиля шеф-повара
	IsBlock        *bool     `db:"is_block"`         // Флаг, определяющий, заблокирован профиль шеф-повара
	Rating         *float32  `db:"rating"`
	ReviewsCount   *int32    `db:"reviews_count"`
	CreatedAt      time.Time `db:"created_at"`  // Дата и время создания записи
	UpdatedAt      time.Time `db:"updated_at"`  // Дата и время последнего обновления записи
	Description    string    `db:"description"` // Описание шефа BIO
	LegalInfo      string    `db:"legal_info"`  // Юридическая информация
}

// Certification — запись в таблице certifications.
type Certification struct {
	ID          int            `db:"id"`
	Name        string         `db:"name"`
	Description sql.NullString `db:"description"`
	CreatedAt   time.Time      `db:"created_at"`
}

// ChefCertification — запись в связующей таблице chef_certifications.
type ChefCertification struct {
	ChefID          int64        `db:"chef_id"`
	CertificationID int          `db:"certification_id"`
	IssuedAt        sql.NullTime `db:"issued_at"`
}
