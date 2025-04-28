package cart

import (
	"context"
	cartentity "domashka-backend/internal/entity/cart"
	"domashka-backend/internal/entity/dishes"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		cartRepo CartRepository
	}
	tests := []struct {
		name string
		args args
		want *Usecase
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.cartRepo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsecase_AddItem(t *testing.T) {
	type fields struct {
		cartRepo CartRepository
	}
	type args struct {
		ctx                context.Context
		userID             int64
		dish               dishes.Dish
		sizeID             int64
		addedIngredients   []int64
		removedIngredients []int64
		notes              string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &Usecase{
				cartRepo: tt.fields.cartRepo,
			}
			got, err := u.AddItem(tt.args.ctx, tt.args.userID, tt.args.dish, tt.args.sizeID, tt.args.addedIngredients, tt.args.removedIngredients, tt.args.notes)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AddItem() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsecase_ClearCart(t *testing.T) {
	type fields struct {
		cartRepo CartRepository
	}
	type args struct {
		ctx    context.Context
		userID int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &Usecase{
				cartRepo: tt.fields.cartRepo,
			}
			if err := u.ClearCart(tt.args.ctx, tt.args.userID); (err != nil) != tt.wantErr {
				t.Errorf("ClearCart() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUsecase_DecrementCartItem(t *testing.T) {
	type fields struct {
		cartRepo CartRepository
	}
	type args struct {
		ctx        context.Context
		cartItemID int64
	}
	tests := []struct {
		name            string
		fields          fields
		args            args
		wantNewQuantity int32
		wantErr         bool
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &Usecase{
				cartRepo: tt.fields.cartRepo,
			}
			gotNewQuantity, err := u.DecrementCartItem(tt.args.ctx, tt.args.cartItemID)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecrementCartItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotNewQuantity != tt.wantNewQuantity {
				t.Errorf("DecrementCartItem() gotNewQuantity = %v, want %v", gotNewQuantity, tt.wantNewQuantity)
			}
		})
	}
}

func TestUsecase_GetCartItems(t *testing.T) {
	type fields struct {
		cartRepo CartRepository
	}
	type args struct {
		ctx    context.Context
		userID int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []cartentity.CartItem
		wantErr bool
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &Usecase{
				cartRepo: tt.fields.cartRepo,
			}
			got, err := u.GetCartItems(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCartItems() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCartItems() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsecase_GetCartItemsByOrderID(t *testing.T) {
	type fields struct {
		cartRepo CartRepository
	}
	type args struct {
		ctx     context.Context
		orderID int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []cartentity.CartItem
		wantErr bool
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &Usecase{
				cartRepo: tt.fields.cartRepo,
			}
			got, err := u.GetCartItemsByOrderID(tt.args.ctx, tt.args.orderID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCartItemsByOrderID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCartItemsByOrderID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsecase_IncrementCartItem(t *testing.T) {
	type fields struct {
		cartRepo CartRepository
	}
	type args struct {
		ctx        context.Context
		cartItemID int64
	}
	tests := []struct {
		name            string
		fields          fields
		args            args
		wantNewQuantity int32
		wantErr         bool
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &Usecase{
				cartRepo: tt.fields.cartRepo,
			}
			gotNewQuantity, err := u.IncrementCartItem(tt.args.ctx, tt.args.cartItemID)
			if (err != nil) != tt.wantErr {
				t.Errorf("IncrementCartItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotNewQuantity != tt.wantNewQuantity {
				t.Errorf("IncrementCartItem() gotNewQuantity = %v, want %v", gotNewQuantity, tt.wantNewQuantity)
			}
		})
	}
}

func TestUsecase_RemoveItem(t *testing.T) {
	type fields struct {
		cartRepo CartRepository
	}
	type args struct {
		ctx        context.Context
		cartItemID int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &Usecase{
				cartRepo: tt.fields.cartRepo,
			}
			if err := u.RemoveItem(tt.args.ctx, tt.args.cartItemID); (err != nil) != tt.wantErr {
				t.Errorf("RemoveItem() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
