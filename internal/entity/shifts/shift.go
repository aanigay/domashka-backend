package shifts

import "time"

type Shift struct {
	ID          int64
	ChefID      int64
	IsActive    bool
	CreatedAt   time.Time
	ClosedAt    *time.Time
	TotalProfit float64
}
