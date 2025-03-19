package carts

import "time"

type CartItem struct {
	ID                       string    `json:"id" db:"id"`
	UserID                   string    `json:"user_id" db:"user_id"`
	DishID                   string    `json:"dish_id" db:"dish_id"`
	ChefID                   string    `json:"chef_id" db:"chef_id"`
	AdditionalIngredientsIDs []int64   `json:"additional_ingredients_ids" db:"additional_ingredients_ids"`
	RemovedIngredientsIDs    []int64   `json:"removed_ingredients_ids" db:"removed_ingredients_ids"`
	AddedAt                  time.Time `json:"added_at" db:"added_at"`
	CustomerNotes            string    `json:"customer_notes" db:"customer_notes"`
}
