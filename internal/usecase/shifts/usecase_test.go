package shifts

import (
	"context"
	"domashka-backend/internal/entity/shifts"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		repo ShiftsRepo
	}
	tests := []struct {
		name string
		args args
		want *Usecase
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.repo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsecase_GetActiveShiftByChefID(t *testing.T) {
	type fields struct {
		repo ShiftsRepo
	}
	type args struct {
		ctx    context.Context
		chefID int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *shifts.Shift
		wantErr bool
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &Usecase{
				repo: tt.fields.repo,
			}
			got, err := u.GetActiveShiftByChefID(tt.args.ctx, tt.args.chefID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetActiveShiftByChefID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetActiveShiftByChefID() got = %v, want %v", got, tt.want)
			}
		})
	}
}
