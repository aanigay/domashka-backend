package shifts

import (
	"context"
	"domashka-backend/internal/entity/shifts"
	"github.com/golang/mock/gomock"
	"reflect"
	"testing"
)

func TestUsecase_CloseShift(t *testing.T) {
	type fields struct {
		repo ShiftsRepo
	}
	type args struct {
		ctx    context.Context
		chefID int64
	}
	tests := []struct {
		name       string
		shiftsRepo func(ctrl *gomock.Controller) ShiftsRepo
		args       args
		wantErr    bool
	}{
		{
			name: "success",
			shiftsRepo: func(ctrl *gomock.Controller) ShiftsRepo {
				m := NewMockShiftsRepo(ctrl)
				m.EXPECT().CloseActiveShiftByChefID(gomock.Any(), gomock.Any()).Return(nil)
				return m
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(tt.shiftsRepo(ctrl))
			if err := u.CloseShift(tt.args.ctx, tt.args.chefID); (err != nil) != tt.wantErr {
				t.Errorf("CloseShift() error = %v, wantErr %v", err, tt.wantErr)
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
		name       string
		shiftsRepo func(ctrl *gomock.Controller) ShiftsRepo
		args       args
		want       *shifts.Shift
		wantErr    bool
	}{
		{
			name: "success",
			shiftsRepo: func(ctrl *gomock.Controller) ShiftsRepo {
				m := NewMockShiftsRepo(ctrl)
				m.EXPECT().GetActiveShiftByChefID(gomock.Any(), gomock.Any()).Return(&shifts.Shift{}, nil)
				return m
			},
			want: &shifts.Shift{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(tt.shiftsRepo(ctrl))
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

func TestUsecase_OpenShift(t *testing.T) {
	type fields struct {
		repo ShiftsRepo
	}
	type args struct {
		ctx    context.Context
		chefID int64
	}
	tests := []struct {
		name       string
		shiftsRepo func(ctrl *gomock.Controller) ShiftsRepo
		args       args
		wantErr    bool
	}{
		{
			name: "success",
			shiftsRepo: func(ctrl *gomock.Controller) ShiftsRepo {
				m := NewMockShiftsRepo(ctrl)
				m.EXPECT().OpenShift(gomock.Any(), gomock.Any()).Return(nil)
				return m
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(tt.shiftsRepo(ctrl))
			if err := u.OpenShift(tt.args.ctx, tt.args.chefID); (err != nil) != tt.wantErr {
				t.Errorf("OpenShift() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
