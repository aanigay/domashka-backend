package order

import (
	"context"
	cartentity "domashka-backend/internal/entity/cart"
	chefEntity "domashka-backend/internal/entity/chefs"
	dishEntity "domashka-backend/internal/entity/dishes"
	geoEntity "domashka-backend/internal/entity/geo"
	"domashka-backend/internal/entity/orders"
	"github.com/golang/mock/gomock"
	"reflect"
	"testing"
)

func TestUsecase_Accept(t *testing.T) {
	type fields struct {
		geoUsecase    geoUsecase
		cartUsecase   cartUsecase
		dishesUsecase dishesUsecase
		chefsUsecase  chefsUsecase
		reviewUsecase reviewUsecase
		shiftsRepo    shiftsRepo
		ordersRepo    ordersRepo
	}
	type args struct {
		ctx     context.Context
		orderID int64
	}
	tests := []struct {
		name          string
		geoUsecase    func(ctrl *gomock.Controller) geoUsecase
		cartUsecase   func(ctrl *gomock.Controller) cartUsecase
		dishesUsecase func(ctrl *gomock.Controller) dishesUsecase
		chefsUsecase  func(ctrl *gomock.Controller) chefsUsecase
		reviewUsecase func(ctrl *gomock.Controller) reviewUsecase
		shiftsRepo    func(ctrl *gomock.Controller) shiftsRepo
		ordersRepo    func(ctrl *gomock.Controller) ordersRepo
		args          args
		wantErr       bool
	}{
		{
			name: "success",
			geoUsecase: func(ctrl *gomock.Controller) geoUsecase {
				m := NewMockgeoUsecase(ctrl)
				return m
			},
			cartUsecase: func(ctrl *gomock.Controller) cartUsecase {
				m := NewMockcartUsecase(ctrl)
				return m
			},
			dishesUsecase: func(ctrl *gomock.Controller) dishesUsecase {
				m := NewMockdishesUsecase(ctrl)
				return m
			},
			chefsUsecase: func(ctrl *gomock.Controller) chefsUsecase {
				m := NewMockchefsUsecase(ctrl)
				return m
			},
			reviewUsecase: func(ctrl *gomock.Controller) reviewUsecase {
				m := NewMockreviewUsecase(ctrl)
				return m
			},
			shiftsRepo: func(ctrl *gomock.Controller) shiftsRepo {
				m := NewMockshiftsRepo(ctrl)
				return m
			},
			ordersRepo: func(ctrl *gomock.Controller) ordersRepo {
				m := NewMockordersRepo(ctrl)
				return m
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(
				tt.geoUsecase(ctrl),
				tt.cartUsecase(ctrl),
				tt.shiftsRepo(ctrl),
				tt.ordersRepo(ctrl),
				tt.dishesUsecase(ctrl),
				tt.chefsUsecase(ctrl),
				tt.reviewUsecase(ctrl),
			)
			if err := u.Accept(tt.args.ctx, tt.args.orderID); (err != nil) != tt.wantErr {
				t.Errorf("Accept() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUsecase_CallDelivery(t *testing.T) {
	type fields struct {
		geoUsecase    geoUsecase
		cartUsecase   cartUsecase
		dishesUsecase dishesUsecase
		chefsUsecase  chefsUsecase
		reviewUsecase reviewUsecase
		shiftsRepo    shiftsRepo
		ordersRepo    ordersRepo
	}
	type args struct {
		ctx     context.Context
		orderID int64
	}
	tests := []struct {
		name          string
		geoUsecase    func(ctrl *gomock.Controller) geoUsecase
		cartUsecase   func(ctrl *gomock.Controller) cartUsecase
		dishesUsecase func(ctrl *gomock.Controller) dishesUsecase
		chefsUsecase  func(ctrl *gomock.Controller) chefsUsecase
		reviewUsecase func(ctrl *gomock.Controller) reviewUsecase
		shiftsRepo    func(ctrl *gomock.Controller) shiftsRepo
		ordersRepo    func(ctrl *gomock.Controller) ordersRepo
		args          args
		wantErr       bool
	}{
		{
			name: "success",
			geoUsecase: func(ctrl *gomock.Controller) geoUsecase {
				m := NewMockgeoUsecase(ctrl)
				return m
			},
			cartUsecase: func(ctrl *gomock.Controller) cartUsecase {
				m := NewMockcartUsecase(ctrl)
				return m
			},
			dishesUsecase: func(ctrl *gomock.Controller) dishesUsecase {
				m := NewMockdishesUsecase(ctrl)
				return m
			},
			chefsUsecase: func(ctrl *gomock.Controller) chefsUsecase {
				m := NewMockchefsUsecase(ctrl)
				return m
			},
			reviewUsecase: func(ctrl *gomock.Controller) reviewUsecase {
				m := NewMockreviewUsecase(ctrl)
				return m
			},
			shiftsRepo: func(ctrl *gomock.Controller) shiftsRepo {
				m := NewMockshiftsRepo(ctrl)
				return m
			},
			ordersRepo: func(ctrl *gomock.Controller) ordersRepo {
				m := NewMockordersRepo(ctrl)
				return m
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(
				tt.geoUsecase(ctrl),
				tt.cartUsecase(ctrl),
				tt.shiftsRepo(ctrl),
				tt.ordersRepo(ctrl),
				tt.dishesUsecase(ctrl),
				tt.chefsUsecase(ctrl),
				tt.reviewUsecase(ctrl),
			)
			if err := u.CallDelivery(tt.args.ctx, tt.args.orderID); (err != nil) != tt.wantErr {
				t.Errorf("CallDelivery() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUsecase_CreateOrder(t *testing.T) {
	type fields struct {
		geoUsecase    geoUsecase
		cartUsecase   cartUsecase
		dishesUsecase dishesUsecase
		chefsUsecase  chefsUsecase
		reviewUsecase reviewUsecase
		shiftsRepo    shiftsRepo
		ordersRepo    ordersRepo
	}
	type args struct {
		ctx            context.Context
		userID         int64
		leaveByTheDoor bool
		callBeforehand bool
	}
	tests := []struct {
		name          string
		geoUsecase    func(ctrl *gomock.Controller) geoUsecase
		cartUsecase   func(ctrl *gomock.Controller) cartUsecase
		dishesUsecase func(ctrl *gomock.Controller) dishesUsecase
		chefsUsecase  func(ctrl *gomock.Controller) chefsUsecase
		reviewUsecase func(ctrl *gomock.Controller) reviewUsecase
		shiftsRepo    func(ctrl *gomock.Controller) shiftsRepo
		ordersRepo    func(ctrl *gomock.Controller) ordersRepo
		args          args
		want          int64
		wantErr       bool
	}{
		{
			name: "success",
			geoUsecase: func(ctrl *gomock.Controller) geoUsecase {
				m := NewMockgeoUsecase(ctrl)
				m.EXPECT().GetLastUpdatedClientAddress(gomock.Any(), gomock.Any()).Return(&geoEntity.Address{}, nil)
				return m
			},
			cartUsecase: func(ctrl *gomock.Controller) cartUsecase {
				m := NewMockcartUsecase(ctrl)
				m.EXPECT().GetCartItems(gomock.Any(), gomock.Any()).Return([]cartentity.CartItem{{
					ID: 1,
					Dish: dishEntity.Dish{
						ID: 1,
					},
					Quantity:           1,
					AddedIngredients:   []dishEntity.Ingredient{},
					RemovedIngredients: []dishEntity.Ingredient{},
					Size:               dishEntity.Size{},
					Notes:              "",
				}}, nil)
				return m
			},
			dishesUsecase: func(ctrl *gomock.Controller) dishesUsecase {
				m := NewMockdishesUsecase(ctrl)
				return m
			},
			chefsUsecase: func(ctrl *gomock.Controller) chefsUsecase {
				m := NewMockchefsUsecase(ctrl)
				return m
			},
			reviewUsecase: func(ctrl *gomock.Controller) reviewUsecase {
				m := NewMockreviewUsecase(ctrl)
				return m
			},
			shiftsRepo: func(ctrl *gomock.Controller) shiftsRepo {
				m := NewMockshiftsRepo(ctrl)
				m.EXPECT().GetActiveShiftIDByChefID(gomock.Any(), gomock.Any()).Return(int64(1), nil)
				return m
			},
			ordersRepo: func(ctrl *gomock.Controller) ordersRepo {
				m := NewMockordersRepo(ctrl)
				m.EXPECT().CreateOrder(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return(int64(1), nil)
				m.EXPECT().AddCartItemToOrder(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return m
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(
				tt.geoUsecase(ctrl),
				tt.cartUsecase(ctrl),
				tt.shiftsRepo(ctrl),
				tt.ordersRepo(ctrl),
				tt.dishesUsecase(ctrl),
				tt.chefsUsecase(ctrl),
				tt.reviewUsecase(ctrl),
			)
			got, err := u.CreateOrder(tt.args.ctx, tt.args.userID, tt.args.leaveByTheDoor, tt.args.callBeforehand)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateOrder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CreateOrder() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsecase_Deliver(t *testing.T) {
	type fields struct {
		geoUsecase    geoUsecase
		cartUsecase   cartUsecase
		dishesUsecase dishesUsecase
		chefsUsecase  chefsUsecase
		reviewUsecase reviewUsecase
		shiftsRepo    shiftsRepo
		ordersRepo    ordersRepo
	}
	type args struct {
		ctx     context.Context
		orderID int64
	}
	tests := []struct {
		name          string
		geoUsecase    func(ctrl *gomock.Controller) geoUsecase
		cartUsecase   func(ctrl *gomock.Controller) cartUsecase
		dishesUsecase func(ctrl *gomock.Controller) dishesUsecase
		chefsUsecase  func(ctrl *gomock.Controller) chefsUsecase
		reviewUsecase func(ctrl *gomock.Controller) reviewUsecase
		shiftsRepo    func(ctrl *gomock.Controller) shiftsRepo
		ordersRepo    func(ctrl *gomock.Controller) ordersRepo
		args          args
		wantErr       bool
	}{
		{
			name: "success",
			geoUsecase: func(ctrl *gomock.Controller) geoUsecase {
				m := NewMockgeoUsecase(ctrl)
				return m
			},
			cartUsecase: func(ctrl *gomock.Controller) cartUsecase {
				m := NewMockcartUsecase(ctrl)
				return m
			},
			dishesUsecase: func(ctrl *gomock.Controller) dishesUsecase {
				m := NewMockdishesUsecase(ctrl)
				return m
			},
			chefsUsecase: func(ctrl *gomock.Controller) chefsUsecase {
				m := NewMockchefsUsecase(ctrl)
				return m
			},
			reviewUsecase: func(ctrl *gomock.Controller) reviewUsecase {
				m := NewMockreviewUsecase(ctrl)
				return m
			},
			shiftsRepo: func(ctrl *gomock.Controller) shiftsRepo {
				m := NewMockshiftsRepo(ctrl)
				return m
			},
			ordersRepo: func(ctrl *gomock.Controller) ordersRepo {
				m := NewMockordersRepo(ctrl)
				return m
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(
				tt.geoUsecase(ctrl),
				tt.cartUsecase(ctrl),
				tt.shiftsRepo(ctrl),
				tt.ordersRepo(ctrl),
				tt.dishesUsecase(ctrl),
				tt.chefsUsecase(ctrl),
				tt.reviewUsecase(ctrl),
			)
			if err := u.Deliver(tt.args.ctx, tt.args.orderID); (err != nil) != tt.wantErr {
				t.Errorf("Deliver() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUsecase_GetActiveOrdersByUserID(t *testing.T) {
	type fields struct {
		geoUsecase    geoUsecase
		cartUsecase   cartUsecase
		dishesUsecase dishesUsecase
		chefsUsecase  chefsUsecase
		reviewUsecase reviewUsecase
		shiftsRepo    shiftsRepo
		ordersRepo    ordersRepo
	}
	type args struct {
		ctx    context.Context
		userID int64
	}
	tests := []struct {
		name          string
		geoUsecase    func(ctrl *gomock.Controller) geoUsecase
		cartUsecase   func(ctrl *gomock.Controller) cartUsecase
		dishesUsecase func(ctrl *gomock.Controller) dishesUsecase
		chefsUsecase  func(ctrl *gomock.Controller) chefsUsecase
		reviewUsecase func(ctrl *gomock.Controller) reviewUsecase
		shiftsRepo    func(ctrl *gomock.Controller) shiftsRepo
		ordersRepo    func(ctrl *gomock.Controller) ordersRepo
		args          args
		want          []orders.Order
		wantErr       bool
	}{

		{
			name: "success",
			geoUsecase: func(ctrl *gomock.Controller) geoUsecase {
				m := NewMockgeoUsecase(ctrl)
				return m
			},
			cartUsecase: func(ctrl *gomock.Controller) cartUsecase {
				m := NewMockcartUsecase(ctrl)
				return m
			},
			dishesUsecase: func(ctrl *gomock.Controller) dishesUsecase {
				m := NewMockdishesUsecase(ctrl)
				return m
			},
			chefsUsecase: func(ctrl *gomock.Controller) chefsUsecase {
				m := NewMockchefsUsecase(ctrl)
				return m
			},
			reviewUsecase: func(ctrl *gomock.Controller) reviewUsecase {
				m := NewMockreviewUsecase(ctrl)
				return m
			},
			shiftsRepo: func(ctrl *gomock.Controller) shiftsRepo {
				m := NewMockshiftsRepo(ctrl)
				return m
			},
			ordersRepo: func(ctrl *gomock.Controller) ordersRepo {
				m := NewMockordersRepo(ctrl)
				return m
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(
				tt.geoUsecase(ctrl),
				tt.cartUsecase(ctrl),
				tt.shiftsRepo(ctrl),
				tt.ordersRepo(ctrl),
				tt.dishesUsecase(ctrl),
				tt.chefsUsecase(ctrl),
				tt.reviewUsecase(ctrl),
			)
			got, err := u.GetActiveOrdersByUserID(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetActiveOrdersByUserID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetActiveOrdersByUserID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsecase_GetCartItemsByOrderID(t *testing.T) {
	type fields struct {
		geoUsecase    geoUsecase
		cartUsecase   cartUsecase
		dishesUsecase dishesUsecase
		chefsUsecase  chefsUsecase
		reviewUsecase reviewUsecase
		shiftsRepo    shiftsRepo
		ordersRepo    ordersRepo
	}
	type args struct {
		ctx     context.Context
		orderID int64
	}
	tests := []struct {
		name          string
		geoUsecase    func(ctrl *gomock.Controller) geoUsecase
		cartUsecase   func(ctrl *gomock.Controller) cartUsecase
		dishesUsecase func(ctrl *gomock.Controller) dishesUsecase
		chefsUsecase  func(ctrl *gomock.Controller) chefsUsecase
		reviewUsecase func(ctrl *gomock.Controller) reviewUsecase
		shiftsRepo    func(ctrl *gomock.Controller) shiftsRepo
		ordersRepo    func(ctrl *gomock.Controller) ordersRepo
		args          args
		want          []cartentity.CartItem
		wantErr       bool
	}{
		{
			name: "success",
			geoUsecase: func(ctrl *gomock.Controller) geoUsecase {
				m := NewMockgeoUsecase(ctrl)
				return m
			},
			cartUsecase: func(ctrl *gomock.Controller) cartUsecase {
				m := NewMockcartUsecase(ctrl)
				return m
			},
			dishesUsecase: func(ctrl *gomock.Controller) dishesUsecase {
				m := NewMockdishesUsecase(ctrl)
				return m
			},
			chefsUsecase: func(ctrl *gomock.Controller) chefsUsecase {
				m := NewMockchefsUsecase(ctrl)
				return m
			},
			reviewUsecase: func(ctrl *gomock.Controller) reviewUsecase {
				m := NewMockreviewUsecase(ctrl)
				return m
			},
			shiftsRepo: func(ctrl *gomock.Controller) shiftsRepo {
				m := NewMockshiftsRepo(ctrl)
				return m
			},
			ordersRepo: func(ctrl *gomock.Controller) ordersRepo {
				m := NewMockordersRepo(ctrl)
				return m
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(
				tt.geoUsecase(ctrl),
				tt.cartUsecase(ctrl),
				tt.shiftsRepo(ctrl),
				tt.ordersRepo(ctrl),
				tt.dishesUsecase(ctrl),
				tt.chefsUsecase(ctrl),
				tt.reviewUsecase(ctrl),
			)
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

func TestUsecase_GetOrderByID(t *testing.T) {
	type fields struct {
		geoUsecase    geoUsecase
		cartUsecase   cartUsecase
		dishesUsecase dishesUsecase
		chefsUsecase  chefsUsecase
		reviewUsecase reviewUsecase
		shiftsRepo    shiftsRepo
		ordersRepo    ordersRepo
	}
	type args struct {
		ctx     context.Context
		orderID int64
	}
	tests := []struct {
		name          string
		geoUsecase    func(ctrl *gomock.Controller) geoUsecase
		cartUsecase   func(ctrl *gomock.Controller) cartUsecase
		dishesUsecase func(ctrl *gomock.Controller) dishesUsecase
		chefsUsecase  func(ctrl *gomock.Controller) chefsUsecase
		reviewUsecase func(ctrl *gomock.Controller) reviewUsecase
		shiftsRepo    func(ctrl *gomock.Controller) shiftsRepo
		ordersRepo    func(ctrl *gomock.Controller) ordersRepo
		args          args
		want          *orders.Order
		wantErr       bool
	}{
		{
			name: "success",
			geoUsecase: func(ctrl *gomock.Controller) geoUsecase {
				m := NewMockgeoUsecase(ctrl)
				return m
			},
			cartUsecase: func(ctrl *gomock.Controller) cartUsecase {
				m := NewMockcartUsecase(ctrl)
				return m
			},
			dishesUsecase: func(ctrl *gomock.Controller) dishesUsecase {
				m := NewMockdishesUsecase(ctrl)
				return m
			},
			chefsUsecase: func(ctrl *gomock.Controller) chefsUsecase {
				m := NewMockchefsUsecase(ctrl)
				return m
			},
			reviewUsecase: func(ctrl *gomock.Controller) reviewUsecase {
				m := NewMockreviewUsecase(ctrl)
				return m
			},
			shiftsRepo: func(ctrl *gomock.Controller) shiftsRepo {
				m := NewMockshiftsRepo(ctrl)
				return m
			},
			ordersRepo: func(ctrl *gomock.Controller) ordersRepo {
				m := NewMockordersRepo(ctrl)
				return m
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(
				tt.geoUsecase(ctrl),
				tt.cartUsecase(ctrl),
				tt.shiftsRepo(ctrl),
				tt.ordersRepo(ctrl),
				tt.dishesUsecase(ctrl),
				tt.chefsUsecase(ctrl),
				tt.reviewUsecase(ctrl),
			)
			got, err := u.GetOrderByID(tt.args.ctx, tt.args.orderID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetOrderByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetOrderByID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsecase_GetOrderedDishesAndChefsByUserID(t *testing.T) {
	type fields struct {
		geoUsecase    geoUsecase
		cartUsecase   cartUsecase
		dishesUsecase dishesUsecase
		chefsUsecase  chefsUsecase
		reviewUsecase reviewUsecase
		shiftsRepo    shiftsRepo
		ordersRepo    ordersRepo
	}
	type args struct {
		ctx    context.Context
		userID int64
	}
	tests := []struct {
		name          string
		geoUsecase    func(ctrl *gomock.Controller) geoUsecase
		cartUsecase   func(ctrl *gomock.Controller) cartUsecase
		dishesUsecase func(ctrl *gomock.Controller) dishesUsecase
		chefsUsecase  func(ctrl *gomock.Controller) chefsUsecase
		reviewUsecase func(ctrl *gomock.Controller) reviewUsecase
		shiftsRepo    func(ctrl *gomock.Controller) shiftsRepo
		ordersRepo    func(ctrl *gomock.Controller) ordersRepo
		args          args
		want          []dishEntity.Dish
		want1         []chefEntity.Chef
		wantErr       bool
	}{
		{
			name: "success",
			geoUsecase: func(ctrl *gomock.Controller) geoUsecase {
				m := NewMockgeoUsecase(ctrl)
				return m
			},
			cartUsecase: func(ctrl *gomock.Controller) cartUsecase {
				m := NewMockcartUsecase(ctrl)
				return m
			},
			dishesUsecase: func(ctrl *gomock.Controller) dishesUsecase {
				m := NewMockdishesUsecase(ctrl)
				return m
			},
			chefsUsecase: func(ctrl *gomock.Controller) chefsUsecase {
				m := NewMockchefsUsecase(ctrl)
				return m
			},
			reviewUsecase: func(ctrl *gomock.Controller) reviewUsecase {
				m := NewMockreviewUsecase(ctrl)
				return m
			},
			shiftsRepo: func(ctrl *gomock.Controller) shiftsRepo {
				m := NewMockshiftsRepo(ctrl)
				return m
			},
			ordersRepo: func(ctrl *gomock.Controller) ordersRepo {
				m := NewMockordersRepo(ctrl)
				return m
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(
				tt.geoUsecase(ctrl),
				tt.cartUsecase(ctrl),
				tt.shiftsRepo(ctrl),
				tt.ordersRepo(ctrl),
				tt.dishesUsecase(ctrl),
				tt.chefsUsecase(ctrl),
				tt.reviewUsecase(ctrl),
			)
			got, got1, err := u.GetOrderedDishesAndChefsByUserID(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetOrderedDishesAndChefsByUserID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetOrderedDishesAndChefsByUserID() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("GetOrderedDishesAndChefsByUserID() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestUsecase_GetOrdersByShiftID(t *testing.T) {
	type fields struct {
		geoUsecase    geoUsecase
		cartUsecase   cartUsecase
		dishesUsecase dishesUsecase
		chefsUsecase  chefsUsecase
		reviewUsecase reviewUsecase
		shiftsRepo    shiftsRepo
		ordersRepo    ordersRepo
	}
	type args struct {
		ctx     context.Context
		shiftID int64
	}
	tests := []struct {
		name          string
		geoUsecase    func(ctrl *gomock.Controller) geoUsecase
		cartUsecase   func(ctrl *gomock.Controller) cartUsecase
		dishesUsecase func(ctrl *gomock.Controller) dishesUsecase
		chefsUsecase  func(ctrl *gomock.Controller) chefsUsecase
		reviewUsecase func(ctrl *gomock.Controller) reviewUsecase
		shiftsRepo    func(ctrl *gomock.Controller) shiftsRepo
		ordersRepo    func(ctrl *gomock.Controller) ordersRepo
		args          args
		want          []orders.Order
		wantErr       bool
	}{
		{
			name: "success",
			geoUsecase: func(ctrl *gomock.Controller) geoUsecase {
				m := NewMockgeoUsecase(ctrl)
				return m
			},
			cartUsecase: func(ctrl *gomock.Controller) cartUsecase {
				m := NewMockcartUsecase(ctrl)
				return m
			},
			dishesUsecase: func(ctrl *gomock.Controller) dishesUsecase {
				m := NewMockdishesUsecase(ctrl)
				return m
			},
			chefsUsecase: func(ctrl *gomock.Controller) chefsUsecase {
				m := NewMockchefsUsecase(ctrl)
				return m
			},
			reviewUsecase: func(ctrl *gomock.Controller) reviewUsecase {
				m := NewMockreviewUsecase(ctrl)
				return m
			},
			shiftsRepo: func(ctrl *gomock.Controller) shiftsRepo {
				m := NewMockshiftsRepo(ctrl)
				return m
			},
			ordersRepo: func(ctrl *gomock.Controller) ordersRepo {
				m := NewMockordersRepo(ctrl)
				return m
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(
				tt.geoUsecase(ctrl),
				tt.cartUsecase(ctrl),
				tt.shiftsRepo(ctrl),
				tt.ordersRepo(ctrl),
				tt.dishesUsecase(ctrl),
				tt.chefsUsecase(ctrl),
				tt.reviewUsecase(ctrl),
			)
			got, err := u.GetOrdersByShiftID(tt.args.ctx, tt.args.shiftID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetOrdersByShiftID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetOrdersByShiftID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsecase_GetOrdersByUserID(t *testing.T) {
	type fields struct {
		geoUsecase    geoUsecase
		cartUsecase   cartUsecase
		dishesUsecase dishesUsecase
		chefsUsecase  chefsUsecase
		reviewUsecase reviewUsecase
		shiftsRepo    shiftsRepo
		ordersRepo    ordersRepo
	}
	type args struct {
		ctx    context.Context
		userID int64
	}
	tests := []struct {
		name          string
		geoUsecase    func(ctrl *gomock.Controller) geoUsecase
		cartUsecase   func(ctrl *gomock.Controller) cartUsecase
		dishesUsecase func(ctrl *gomock.Controller) dishesUsecase
		chefsUsecase  func(ctrl *gomock.Controller) chefsUsecase
		reviewUsecase func(ctrl *gomock.Controller) reviewUsecase
		shiftsRepo    func(ctrl *gomock.Controller) shiftsRepo
		ordersRepo    func(ctrl *gomock.Controller) ordersRepo
		args          args
		want          []orders.OrderProfile
		wantErr       bool
	}{
		{
			name: "success",
			geoUsecase: func(ctrl *gomock.Controller) geoUsecase {
				m := NewMockgeoUsecase(ctrl)
				return m
			},
			cartUsecase: func(ctrl *gomock.Controller) cartUsecase {
				m := NewMockcartUsecase(ctrl)
				return m
			},
			dishesUsecase: func(ctrl *gomock.Controller) dishesUsecase {
				m := NewMockdishesUsecase(ctrl)
				return m
			},
			chefsUsecase: func(ctrl *gomock.Controller) chefsUsecase {
				m := NewMockchefsUsecase(ctrl)
				return m
			},
			reviewUsecase: func(ctrl *gomock.Controller) reviewUsecase {
				m := NewMockreviewUsecase(ctrl)
				return m
			},
			shiftsRepo: func(ctrl *gomock.Controller) shiftsRepo {
				m := NewMockshiftsRepo(ctrl)
				return m
			},
			ordersRepo: func(ctrl *gomock.Controller) ordersRepo {
				m := NewMockordersRepo(ctrl)
				return m
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(
				tt.geoUsecase(ctrl),
				tt.cartUsecase(ctrl),
				tt.shiftsRepo(ctrl),
				tt.ordersRepo(ctrl),
				tt.dishesUsecase(ctrl),
				tt.chefsUsecase(ctrl),
				tt.reviewUsecase(ctrl),
			)
			got, err := u.GetOrdersByUserID(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetOrdersByUserID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetOrdersByUserID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsecase_GetStatus(t *testing.T) {
	type fields struct {
		geoUsecase    geoUsecase
		cartUsecase   cartUsecase
		dishesUsecase dishesUsecase
		chefsUsecase  chefsUsecase
		reviewUsecase reviewUsecase
		shiftsRepo    shiftsRepo
		ordersRepo    ordersRepo
	}
	type args struct {
		ctx     context.Context
		orderID int64
	}
	tests := []struct {
		name          string
		geoUsecase    func(ctrl *gomock.Controller) geoUsecase
		cartUsecase   func(ctrl *gomock.Controller) cartUsecase
		dishesUsecase func(ctrl *gomock.Controller) dishesUsecase
		chefsUsecase  func(ctrl *gomock.Controller) chefsUsecase
		reviewUsecase func(ctrl *gomock.Controller) reviewUsecase
		shiftsRepo    func(ctrl *gomock.Controller) shiftsRepo
		ordersRepo    func(ctrl *gomock.Controller) ordersRepo
		args          args
		want          int32
		wantErr       bool
	}{
		{
			name: "success",
			geoUsecase: func(ctrl *gomock.Controller) geoUsecase {
				m := NewMockgeoUsecase(ctrl)
				return m
			},
			cartUsecase: func(ctrl *gomock.Controller) cartUsecase {
				m := NewMockcartUsecase(ctrl)
				return m
			},
			dishesUsecase: func(ctrl *gomock.Controller) dishesUsecase {
				m := NewMockdishesUsecase(ctrl)
				return m
			},
			chefsUsecase: func(ctrl *gomock.Controller) chefsUsecase {
				m := NewMockchefsUsecase(ctrl)
				return m
			},
			reviewUsecase: func(ctrl *gomock.Controller) reviewUsecase {
				m := NewMockreviewUsecase(ctrl)
				return m
			},
			shiftsRepo: func(ctrl *gomock.Controller) shiftsRepo {
				m := NewMockshiftsRepo(ctrl)
				return m
			},
			ordersRepo: func(ctrl *gomock.Controller) ordersRepo {
				m := NewMockordersRepo(ctrl)
				return m
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(
				tt.geoUsecase(ctrl),
				tt.cartUsecase(ctrl),
				tt.shiftsRepo(ctrl),
				tt.ordersRepo(ctrl),
				tt.dishesUsecase(ctrl),
				tt.chefsUsecase(ctrl),
				tt.reviewUsecase(ctrl),
			)
			got, err := u.GetStatus(tt.args.ctx, tt.args.orderID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetStatus() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsecase_PickUp(t *testing.T) {
	type fields struct {
		geoUsecase    geoUsecase
		cartUsecase   cartUsecase
		dishesUsecase dishesUsecase
		chefsUsecase  chefsUsecase
		reviewUsecase reviewUsecase
		shiftsRepo    shiftsRepo
		ordersRepo    ordersRepo
	}
	type args struct {
		ctx     context.Context
		orderID int64
	}
	tests := []struct {
		name          string
		geoUsecase    func(ctrl *gomock.Controller) geoUsecase
		cartUsecase   func(ctrl *gomock.Controller) cartUsecase
		dishesUsecase func(ctrl *gomock.Controller) dishesUsecase
		chefsUsecase  func(ctrl *gomock.Controller) chefsUsecase
		reviewUsecase func(ctrl *gomock.Controller) reviewUsecase
		shiftsRepo    func(ctrl *gomock.Controller) shiftsRepo
		ordersRepo    func(ctrl *gomock.Controller) ordersRepo
		args          args
		wantErr       bool
	}{

		{
			name: "success",
			geoUsecase: func(ctrl *gomock.Controller) geoUsecase {
				m := NewMockgeoUsecase(ctrl)
				return m
			},
			cartUsecase: func(ctrl *gomock.Controller) cartUsecase {
				m := NewMockcartUsecase(ctrl)
				return m
			},
			dishesUsecase: func(ctrl *gomock.Controller) dishesUsecase {
				m := NewMockdishesUsecase(ctrl)
				return m
			},
			chefsUsecase: func(ctrl *gomock.Controller) chefsUsecase {
				m := NewMockchefsUsecase(ctrl)
				return m
			},
			reviewUsecase: func(ctrl *gomock.Controller) reviewUsecase {
				m := NewMockreviewUsecase(ctrl)
				return m
			},
			shiftsRepo: func(ctrl *gomock.Controller) shiftsRepo {
				m := NewMockshiftsRepo(ctrl)
				return m
			},
			ordersRepo: func(ctrl *gomock.Controller) ordersRepo {
				m := NewMockordersRepo(ctrl)
				return m
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(
				tt.geoUsecase(ctrl),
				tt.cartUsecase(ctrl),
				tt.shiftsRepo(ctrl),
				tt.ordersRepo(ctrl),
				tt.dishesUsecase(ctrl),
				tt.chefsUsecase(ctrl),
				tt.reviewUsecase(ctrl),
			)
			if err := u.PickUp(tt.args.ctx, tt.args.orderID); (err != nil) != tt.wantErr {
				t.Errorf("PickUp() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUsecase_Reject(t *testing.T) {
	type fields struct {
		geoUsecase    geoUsecase
		cartUsecase   cartUsecase
		dishesUsecase dishesUsecase
		chefsUsecase  chefsUsecase
		reviewUsecase reviewUsecase
		shiftsRepo    shiftsRepo
		ordersRepo    ordersRepo
	}
	type args struct {
		ctx     context.Context
		orderID int64
	}
	tests := []struct {
		name          string
		geoUsecase    func(ctrl *gomock.Controller) geoUsecase
		cartUsecase   func(ctrl *gomock.Controller) cartUsecase
		dishesUsecase func(ctrl *gomock.Controller) dishesUsecase
		chefsUsecase  func(ctrl *gomock.Controller) chefsUsecase
		reviewUsecase func(ctrl *gomock.Controller) reviewUsecase
		shiftsRepo    func(ctrl *gomock.Controller) shiftsRepo
		ordersRepo    func(ctrl *gomock.Controller) ordersRepo
		args          args
		wantErr       bool
	}{
		{
			name: "success",
			geoUsecase: func(ctrl *gomock.Controller) geoUsecase {
				m := NewMockgeoUsecase(ctrl)
				return m
			},
			cartUsecase: func(ctrl *gomock.Controller) cartUsecase {
				m := NewMockcartUsecase(ctrl)
				return m
			},
			dishesUsecase: func(ctrl *gomock.Controller) dishesUsecase {
				m := NewMockdishesUsecase(ctrl)
				return m
			},
			chefsUsecase: func(ctrl *gomock.Controller) chefsUsecase {
				m := NewMockchefsUsecase(ctrl)
				return m
			},
			reviewUsecase: func(ctrl *gomock.Controller) reviewUsecase {
				m := NewMockreviewUsecase(ctrl)
				return m
			},
			shiftsRepo: func(ctrl *gomock.Controller) shiftsRepo {
				m := NewMockshiftsRepo(ctrl)
				return m
			},
			ordersRepo: func(ctrl *gomock.Controller) ordersRepo {
				m := NewMockordersRepo(ctrl)
				return m
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(
				tt.geoUsecase(ctrl),
				tt.cartUsecase(ctrl),
				tt.shiftsRepo(ctrl),
				tt.ordersRepo(ctrl),
				tt.dishesUsecase(ctrl),
				tt.chefsUsecase(ctrl),
				tt.reviewUsecase(ctrl),
			)
			if err := u.Reject(tt.args.ctx, tt.args.orderID); (err != nil) != tt.wantErr {
				t.Errorf("Reject() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUsecase_SetStatus(t *testing.T) {
	type fields struct {
		geoUsecase    geoUsecase
		cartUsecase   cartUsecase
		dishesUsecase dishesUsecase
		chefsUsecase  chefsUsecase
		reviewUsecase reviewUsecase
		shiftsRepo    shiftsRepo
		ordersRepo    ordersRepo
	}
	type args struct {
		ctx     context.Context
		orderID int64
		status  int32
	}
	tests := []struct {
		name          string
		geoUsecase    func(ctrl *gomock.Controller) geoUsecase
		cartUsecase   func(ctrl *gomock.Controller) cartUsecase
		dishesUsecase func(ctrl *gomock.Controller) dishesUsecase
		chefsUsecase  func(ctrl *gomock.Controller) chefsUsecase
		reviewUsecase func(ctrl *gomock.Controller) reviewUsecase
		shiftsRepo    func(ctrl *gomock.Controller) shiftsRepo
		ordersRepo    func(ctrl *gomock.Controller) ordersRepo
		args          args
		wantErr       bool
	}{
		{
			name: "success",
			geoUsecase: func(ctrl *gomock.Controller) geoUsecase {
				m := NewMockgeoUsecase(ctrl)
				return m
			},
			cartUsecase: func(ctrl *gomock.Controller) cartUsecase {
				m := NewMockcartUsecase(ctrl)
				return m
			},
			dishesUsecase: func(ctrl *gomock.Controller) dishesUsecase {
				m := NewMockdishesUsecase(ctrl)
				return m
			},
			chefsUsecase: func(ctrl *gomock.Controller) chefsUsecase {
				m := NewMockchefsUsecase(ctrl)
				return m
			},
			reviewUsecase: func(ctrl *gomock.Controller) reviewUsecase {
				m := NewMockreviewUsecase(ctrl)
				return m
			},
			shiftsRepo: func(ctrl *gomock.Controller) shiftsRepo {
				m := NewMockshiftsRepo(ctrl)
				return m
			},
			ordersRepo: func(ctrl *gomock.Controller) ordersRepo {
				m := NewMockordersRepo(ctrl)
				return m
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(
				tt.geoUsecase(ctrl),
				tt.cartUsecase(ctrl),
				tt.shiftsRepo(ctrl),
				tt.ordersRepo(ctrl),
				tt.dishesUsecase(ctrl),
				tt.chefsUsecase(ctrl),
				tt.reviewUsecase(ctrl),
			)
			if err := u.SetStatus(tt.args.ctx, tt.args.orderID, tt.args.status); (err != nil) != tt.wantErr {
				t.Errorf("SetStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
