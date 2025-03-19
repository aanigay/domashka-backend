package notifications

import (
	"database/sql"
	"time"
)

const (
	ChannelEmail        = "email"
	ChannelSMS          = "sms"
	ChannelPush         = "push"
	ScenarioOrderStatus = "order_status"
	ScenarioMarketing   = "marketing"
	ScenarioSystem      = "system"
	StatusSent          = "sent"
	StatusError         = "error"
)

type Notification struct {
	ID           int
	UserID       sql.NullInt64 // Может быть NULL
	Channel      string
	Scenario     string
	Subject      sql.NullString
	Message      string
	Recipient    string
	Status       string
	ErrorMessage sql.NullString
	CreatedAt    time.Time
	UpdatedAt    string
	SendAttempts int
	Metadata     sql.NullString
}

type EmailData struct {
	Title    string
	UserName string
	Body     string
	Footer   string
}
