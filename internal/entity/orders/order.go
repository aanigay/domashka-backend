package orders

type Order struct {
	ID     int64 `db:"id"`
	Status int32 `db:"status"`
}
