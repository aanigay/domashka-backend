package chefs

type ChefRating struct {
	ChefID       int64   `db:"chef_id"`
	Rating       float32 `db:"rating"`
	ReviewsCount int32   `db:"reviews_count"`
}
