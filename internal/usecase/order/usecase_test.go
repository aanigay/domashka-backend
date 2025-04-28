package order

import (
	"context"
	"domashka-backend/internal/entity/orders"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		geoUsecase  geoUsecase
		cartUsecase cartUsecase
		shiftsRepo  shiftsRepo
		ordersRepo  ordersRepo
	}
	tests := []struct {
		name string
		args args
		want *Usecase
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.geoUsecase, tt.args.cartUsecase, tt.args.shiftsRepo, tt.args.ordersRepo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsecase_CreateOrder(t *testing.T) {
	type fields struct {
		geoUsecase  geoUsecase
		cartUsecase cartUsecase
		shiftsRepo  shiftsRepo
		ordersRepo  ordersRepo
	}
	type args struct {
		ctx            context.Context
		userID         int64
		leaveByTheDoor bool
		callBeforehand bool
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
				geoUsecase:  tt.fields.geoUsecase,
				cartUsecase: tt.fields.cartUsecase,
				shiftsRepo:  tt.fields.shiftsRepo,
				ordersRepo:  tt.fields.ordersRepo,
			}
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

func TestUsecase_GetOrdersByShiftID(t *testing.T) {
	type fields struct {
		geoUsecase  geoUsecase
		cartUsecase cartUsecase
		shiftsRepo  shiftsRepo
		ordersRepo  ordersRepo
	}
	type args struct {
		ctx     context.Context
		shiftID int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []orders.Order
		wantErr bool
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &Usecase{
				geoUsecase:  tt.fields.geoUsecase,
				cartUsecase: tt.fields.cartUsecase,
				shiftsRepo:  tt.fields.shiftsRepo,
				ordersRepo:  tt.fields.ordersRepo,
			}
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
