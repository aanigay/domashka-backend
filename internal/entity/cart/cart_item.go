package cart

import (
	dishentity "domashka-backend/internal/entity/dishes"
	"fmt"
	"strings"
)

type CartItem struct {
	ID                 int64
	Dish               dishentity.Dish
	Quantity           int32
	AddedIngredients   []dishentity.Ingredient
	RemovedIngredients []dishentity.Ingredient
	Size               dishentity.Size
	Notes              string
}

func (c *CartItem) GetTotalPrice() dishentity.Price {
	return dishentity.Price{
		Value:    float32(c.Quantity) * c.Size.PriceValue,
		Currency: c.Size.PriceCurrency,
	}
}

func GetTotalCartPrice(cartItems []CartItem) dishentity.Price {
	totalPrice := dishentity.Price{Value: 0, Currency: "RUB"}
	for _, item := range cartItems {
		totalPrice.Value += float32(item.Quantity) * item.Size.PriceValue
		totalPrice.Currency = item.Size.PriceCurrency
	}
	return totalPrice
}

func (c *CartItem) GetDetailsString() string {
	addIngredientsLabels := make([]string, 0, len(c.AddedIngredients))
	for _, ingr := range c.AddedIngredients {
		addIngredientsLabels = append(addIngredientsLabels, ingr.Name)
	}
	removeIngredientsLabels := make([]string, 0, len(c.RemovedIngredients))
	for _, ingr := range c.RemovedIngredients {
		removeIngredientsLabels = append(removeIngredientsLabels, ingr.Name)
	}
	details := fmt.Sprintf("%s", c.Size.Label)
	if len(addIngredientsLabels) > 0 {
		details = fmt.Sprintf("%s, Добавить: %s", details, strings.Join(addIngredientsLabels, ", "))
	}
	if len(removeIngredientsLabels) > 0 {
		details = fmt.Sprintf("%s, Убрать: %s", details, strings.Join(removeIngredientsLabels, ", "))
	}
	return details
}
