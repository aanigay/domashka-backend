package users

import "time"

type User struct {
	ID               string     `json:"id" db:"id"`                                 // UUID
	Username         string     `json:"username" db:"username"`                     // Уникальный идентификатор
	Alias            string     `json:"alias" db:"alias"`                           // Полное имя (ФИО)
	FirstName        string     `json:"first_name" db:"first_name"`                 // Имя
	SecondName       *string    `json:"second_name,omitempty" db:"second_name"`     // Отчество (опционально)
	LastName         *string    `json:"last_name,omitempty" db:"last_name"`         // Фамилия (опционально)
	Email            *string    `json:"email,omitempty" db:"email"`                 // Электронная почта (опционально)
	NumberPhone      *string    `json:"number_phone,omitempty" db:"number_phone"`   // Телефон (опционально)
	Status           int        `json:"status" db:"status"`                         // Состояние (0 = не бан, 1 = бан)
	ExternalType     int        `json:"external_type" db:"external_type"`           // Внешний статус (0, 1, 2)
	TelegramName     *string    `json:"telegram_name,omitempty" db:"telegram_name"` // Telegram username (опционально)
	ExternalID       *string    `json:"external_id,omitempty" db:"external_id"`     // Внешний ID (опционально)
	NotificationFlag int        `json:"notification_flag" db:"notification_flag"`   // Уведомления (0 = включены, 1 = выключены)
	Role             string     `json:"role" db:"role"`                             // Роль (client, chef, admin)
	Birthday         *time.Time `json:"birthday,omitempty" db:"birthday"`           // День рождения (опционально)
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`                 // Время создания
	UpdatedAt        time.Time  `json:"updated_at" db:"updated_at"`                 // Время последнего обновления
}
