package dishes

type DishRating struct {
	DishID       int64   `db:"dish_id"`
	Rating       float32 `db:"rating"`
	ReviewsCount int32   `db:"reviews_count"`
}

type DishIngredient struct {
	DishID       int64 `db:"dish_id"`
	IngredientID int64 `db:"ingredient_id"`
	IsRemovable  bool  `db:"is_removable"`
}
