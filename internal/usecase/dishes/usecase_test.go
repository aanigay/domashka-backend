package dishes

import (
	"context"
	entity "domashka-backend/internal/entity/dishes"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		dishRepo dishRepo
	}
	tests := []struct {
		name string
		args args
		want *Usecase
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.dishRepo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsecase_GetDishByID(t *testing.T) {
	type fields struct {
		dishRepo dishRepo
	}
	type args struct {
		ctx    context.Context
		dishID int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *entity.Dish
		wantErr bool
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &Usecase{
				dishRepo: tt.fields.dishRepo,
			}
			got, err := u.GetDishByID(tt.args.ctx, tt.args.dishID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDishByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetDishByID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsecase_GetDishSizesByDishID(t *testing.T) {
	type fields struct {
		dishRepo dishRepo
	}
	type args struct {
		ctx    context.Context
		dishID int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []entity.Size
		wantErr bool
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &Usecase{
				dishRepo: tt.fields.dishRepo,
			}
			got, err := u.GetDishSizesByDishID(tt.args.ctx, tt.args.dishID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDishSizesByDishID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetDishSizesByDishID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsecase_GetDishesByChefID(t *testing.T) {
	type fields struct {
		dishRepo dishRepo
	}
	type args struct {
		ctx    context.Context
		chefID int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []entity.Dish
		wantErr bool
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &Usecase{
				dishRepo: tt.fields.dishRepo,
			}
			got, err := u.GetDishesByChefID(tt.args.ctx, tt.args.chefID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDishesByChefID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetDishesByChefID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsecase_GetIngredientsByDishID(t *testing.T) {
	type fields struct {
		dishRepo dishRepo
	}
	type args struct {
		ctx    context.Context
		dishID int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []entity.Ingredient
		wantErr bool
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &Usecase{
				dishRepo: tt.fields.dishRepo,
			}
			got, err := u.GetIngredientsByDishID(tt.args.ctx, tt.args.dishID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetIngredientsByDishID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetIngredientsByDishID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsecase_GetMinimalPriceByDishID(t *testing.T) {
	type fields struct {
		dishRepo dishRepo
	}
	type args struct {
		ctx    context.Context
		dishID int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *entity.Price
		wantErr bool
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &Usecase{
				dishRepo: tt.fields.dishRepo,
			}
			got, err := u.GetMinimalPriceByDishID(tt.args.ctx, tt.args.dishID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMinimalPriceByDishID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetMinimalPriceByDishID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsecase_GetNutritionByDishID(t *testing.T) {
	type fields struct {
		dishRepo dishRepo
	}
	type args struct {
		ctx    context.Context
		dishID int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *entity.Nutrition
		wantErr bool
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &Usecase{
				dishRepo: tt.fields.dishRepo,
			}
			got, err := u.GetNutritionByDishID(tt.args.ctx, tt.args.dishID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNutritionByDishID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetNutritionByDishID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsecase_GetTopDishes(t *testing.T) {
	type fields struct {
		dishRepo dishRepo
	}
	type args struct {
		ctx   context.Context
		limit int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []entity.Dish
		wantErr bool
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &Usecase{
				dishRepo: tt.fields.dishRepo,
			}
			got, err := u.GetTopDishes(tt.args.ctx, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTopDishes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTopDishes() got = %v, want %v", got, tt.want)
			}
		})
	}
}
