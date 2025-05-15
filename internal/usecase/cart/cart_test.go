package cart

import (
	"context"
	cartentity "domashka-backend/internal/entity/cart"
	"domashka-backend/internal/entity/dishes"
	"github.com/golang/mock/gomock"
	"reflect"
	"testing"
)

func TestUsecase_AddItem(t *testing.T) {
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
		name     string
		cartRepo func(ctrl *gomock.Controller) CartRepository
		args     args
		want     int64
		wantErr  bool
	}{
		{
			name: "success",
			cartRepo: func(ctrl *gomock.Controller) CartRepository {
				m := NewMockCartRepository(ctrl)
				m.EXPECT().AddItem(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(int64(1), nil)
				return m
			},
			args: args{
				ctx:    context.Background(),
				userID: 1,
				dish: dishes.Dish{
					ID:           1,
					Name:         "test",
					Description:  "test",
					ChefID:       1,
					ImageURL:     "example.com",
					Rating:       nil,
					ReviewsCount: nil,
					CategoryID:   1,
					IsDeleted:    false,
				},
			},
			want: int64(1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(tt.cartRepo(ctrl))
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
	type args struct {
		ctx    context.Context
		userID int64
	}
	tests := []struct {
		name     string
		cartRepo func(ctrl *gomock.Controller) CartRepository
		args     args
		wantErr  bool
	}{
		{
			name: "success",
			args: args{
				ctx:    context.Background(),
				userID: 1,
			},
			cartRepo: func(ctrl *gomock.Controller) CartRepository {
				m := NewMockCartRepository(ctrl)
				m.EXPECT().Clear(gomock.Any(), gomock.Any()).Return(nil)
				return m
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(tt.cartRepo(ctrl))
			if err := u.ClearCart(tt.args.ctx, tt.args.userID); (err != nil) != tt.wantErr {
				t.Errorf("ClearCart() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUsecase_DecrementCartItem(t *testing.T) {
	type args struct {
		ctx        context.Context
		cartItemID int64
	}
	tests := []struct {
		name            string
		cartRepo        func(ctrl *gomock.Controller) CartRepository
		args            args
		wantNewQuantity int32
		wantErr         bool
	}{
		{
			name: "success",
			cartRepo: func(ctrl *gomock.Controller) CartRepository {
				m := NewMockCartRepository(ctrl)
				m.EXPECT().DecrementCartItemQuantity(gomock.Any(), int64(1)).Return(int32(2), nil)
				return m
			},
			args: args{
				ctx:        context.Background(),
				cartItemID: 1,
			},
			wantNewQuantity: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(tt.cartRepo(ctrl))
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
	type args struct {
		ctx    context.Context
		userID int64
	}
	tests := []struct {
		name     string
		cartRepo func(ctrl *gomock.Controller) CartRepository
		args     args
		want     []cartentity.CartItem
		wantErr  bool
	}{
		{
			name: "success",
			args: args{
				ctx:    context.Background(),
				userID: 1,
			},
			cartRepo: func(ctrl *gomock.Controller) CartRepository {
				m := NewMockCartRepository(ctrl)
				m.EXPECT().GetCartItems(gomock.Any(), int64(1)).Return([]cartentity.CartItem{
					{
						ID: 1,
						Dish: dishes.Dish{
							ID:           0,
							Name:         "",
							Description:  "",
							ChefID:       0,
							ImageURL:     "",
							Rating:       nil,
							ReviewsCount: nil,
							CategoryID:   0,
							IsDeleted:    false,
						},
						Quantity:           0,
						AddedIngredients:   nil,
						RemovedIngredients: nil,
						Size: dishes.Size{
							ID:            0,
							DishID:        0,
							Label:         "",
							WeightValue:   0,
							WeightUnit:    "",
							PriceValue:    0,
							PriceCurrency: "",
						},
						Notes: "",
					},
				}, nil)
				return m
			},
			want: []cartentity.CartItem{
				{
					ID: 1,
					Dish: dishes.Dish{
						ID:           0,
						Name:         "",
						Description:  "",
						ChefID:       0,
						ImageURL:     "",
						Rating:       nil,
						ReviewsCount: nil,
						CategoryID:   0,
						IsDeleted:    false,
					},
					Quantity:           0,
					AddedIngredients:   nil,
					RemovedIngredients: nil,
					Size: dishes.Size{
						ID:            0,
						DishID:        0,
						Label:         "",
						WeightValue:   0,
						WeightUnit:    "",
						PriceValue:    0,
						PriceCurrency: "",
					},
					Notes: "",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(tt.cartRepo(ctrl))
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
	type args struct {
		ctx     context.Context
		orderID int64
	}
	tests := []struct {
		name     string
		cartRepo func(ctrl *gomock.Controller) CartRepository
		args     args
		want     []cartentity.CartItem
		wantErr  bool
	}{
		{
			name: "success",
			cartRepo: func(ctrl *gomock.Controller) CartRepository {
				m := NewMockCartRepository(ctrl)
				m.EXPECT().GetCartItemsByOrderID(gomock.Any(), int64(1)).Return([]cartentity.CartItem{
					{
						ID: 1,
						Dish: dishes.Dish{
							ID:           0,
							Name:         "",
							Description:  "",
							ChefID:       0,
							ImageURL:     "",
							Rating:       nil,
							ReviewsCount: nil,
							CategoryID:   0,
							IsDeleted:    false,
						},
						Quantity:           0,
						AddedIngredients:   nil,
						RemovedIngredients: nil,
						Size: dishes.Size{
							ID:            0,
							DishID:        0,
							Label:         "",
							WeightValue:   0,
							WeightUnit:    "",
							PriceValue:    0,
							PriceCurrency: "",
						},
						Notes: "",
					},
				}, nil)
				return m
			},
			want: []cartentity.CartItem{

				{
					ID: 1,
					Dish: dishes.Dish{
						ID:           0,
						Name:         "",
						Description:  "",
						ChefID:       0,
						ImageURL:     "",
						Rating:       nil,
						ReviewsCount: nil,
						CategoryID:   0,
						IsDeleted:    false,
					},
					Quantity:           0,
					AddedIngredients:   nil,
					RemovedIngredients: nil,
					Size: dishes.Size{
						ID:            0,
						DishID:        0,
						Label:         "",
						WeightValue:   0,
						WeightUnit:    "",
						PriceValue:    0,
						PriceCurrency: "",
					},
					Notes: "",
				},
			},
			args: args{
				ctx:     context.Background(),
				orderID: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(tt.cartRepo(ctrl))
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
	type args struct {
		ctx        context.Context
		cartItemID int64
	}
	tests := []struct {
		name            string
		cartRepo        func(ctrl *gomock.Controller) CartRepository
		args            args
		wantNewQuantity int32
		wantErr         bool
	}{
		{
			name: "success",
			args: args{
				ctx:        context.Background(),
				cartItemID: 1,
			},
			wantNewQuantity: 2,
			cartRepo: func(ctrl *gomock.Controller) CartRepository {
				m := NewMockCartRepository(ctrl)
				m.EXPECT().IncrementCartItemQuantity(gomock.Any(), int64(1)).Return(int32(2), nil)
				return m
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(tt.cartRepo(ctrl))
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
	type args struct {
		ctx        context.Context
		cartItemID int64
	}
	tests := []struct {
		name     string
		cartRepo func(ctrl *gomock.Controller) CartRepository
		args     args
		wantErr  bool
	}{
		{
			name: "success",
			cartRepo: func(ctrl *gomock.Controller) CartRepository {
				m := NewMockCartRepository(ctrl)
				m.EXPECT().RemoveItem(gomock.Any(), int64(1)).Return(nil)
				return m
			},
			args: args{
				ctx:        context.Background(),
				cartItemID: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(tt.cartRepo(ctrl))
			if err := u.RemoveItem(tt.args.ctx, tt.args.cartItemID); (err != nil) != tt.wantErr {
				t.Errorf("RemoveItem() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
