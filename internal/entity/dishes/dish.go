package dishes

type Dish struct {
	ID          string `json:"id" db:"id"`
	ChefID      string `json:"chef_id" db:"chef_id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	Price       int64  `json:"price" db:"price"`
	Stock       int64  `json:"stock" db:"stock"`
	CreatedAt   string `json:"created_at" db:"created_at"`
	UpdatedAt   string `json:"updated_at" db:"updated_at"`
}

type DishCategory struct {
	ID          string `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
}
