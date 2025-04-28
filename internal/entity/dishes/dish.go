package dishes

type Dish struct {
	ID           int64  `db:"id"`
	Name         string `db:"name"`
	Description  string `db:"description"`
	ChefID       int64  `db:"chef_id"`
	ImageURL     string `db:"image_url"`
	Rating       *float32
	ReviewsCount *int32
}

type Nutrition struct {
	DishID        int64 `db:"dish_id"`
	Calories      int   `db:"calories"`
	Protein       int   `db:"protein"`
	Fat           int   `db:"fat"`
	Carbohydrates int   `db:"carbohydrates"`
}

type Ingredient struct {
	ID          int64  `db:"id"`
	Name        string `db:"name"`
	ImageURL    string `db:"image_url"`
	IsAllergen  bool   `db:"is_allergen"`
	CategoryID  int64  `db:"category_id"`
	IsRemovable bool
}

type Size struct {
	ID            int64   `db:"id"`
	DishID        int64   `db:"dish_id"`
	Label         string  `db:"label"`
	WeightValue   float32 `db:"weight_value"`
	WeightUnit    string  `db:"weight_unit"`
	PriceValue    float32 `db:"price_value"`
	PriceCurrency string  `db:"price_currency"`
}

type Price struct {
	Value    float32
	Currency string
}
