package favorites

import (
	"context"
	"github.com/golang/mock/gomock"
	"testing"
)

func TestUsecase_AddFavoriteChef(t *testing.T) {
	type args struct {
		ctx    context.Context
		userID int64
		chefID int64
	}
	tests := []struct {
		name    string
		favRepo func(ctrl *gomock.Controller) favRepo
		args    args
		wantErr bool
	}{
		{
			name: "success",
			favRepo: func(ctrl *gomock.Controller) favRepo {
				mock := NewMockfavRepo(ctrl)
				mock.EXPECT().AddFavoriteChef(gomock.Any(), int64(1), int64(1)).Return(nil)
				return mock
			},
			args: args{
				ctx:    context.Background(),
				userID: 1,
				chefID: 1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(tt.favRepo(ctrl))
			if err := u.AddFavoriteChef(tt.args.ctx, tt.args.userID, tt.args.chefID); (err != nil) != tt.wantErr {
				t.Errorf("AddFavoriteChef() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUsecase_AddFavoriteDish(t *testing.T) {
	type args struct {
		ctx    context.Context
		userID int64
		dishID int64
	}
	tests := []struct {
		name    string
		favRepo func(ctrl *gomock.Controller) favRepo
		args    args
		wantErr bool
	}{
		{
			name: "success",
			favRepo: func(ctrl *gomock.Controller) favRepo {
				mock := NewMockfavRepo(ctrl)
				mock.EXPECT().AddFavoriteDish(gomock.Any(), int64(1), int64(1)).Return(nil)
				return mock
			},
			args: args{
				ctx:    context.Background(),
				userID: 1,
				dishID: 1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(tt.favRepo(ctrl))
			if err := u.AddFavoriteDish(tt.args.ctx, tt.args.userID, tt.args.dishID); (err != nil) != tt.wantErr {
				t.Errorf("AddFavoriteDish() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUsecase_RemoveFavoriteChef(t *testing.T) {
	type args struct {
		ctx    context.Context
		userID int64
		chefID int64
	}
	tests := []struct {
		name    string
		favRepo func(ctrl *gomock.Controller) favRepo
		args    args
		wantErr bool
	}{
		{
			name: "success",
			favRepo: func(ctrl *gomock.Controller) favRepo {
				mock := NewMockfavRepo(ctrl)
				mock.EXPECT().RemoveFavoriteChef(gomock.Any(), int64(1), int64(1)).Return(nil)
				return mock
			},
			args: args{
				ctx:    context.Background(),
				userID: 1,
				chefID: 1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(tt.favRepo(ctrl))
			if err := u.RemoveFavoriteChef(tt.args.ctx, tt.args.userID, tt.args.chefID); (err != nil) != tt.wantErr {
				t.Errorf("RemoveFavoriteChef() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUsecase_RemoveFavoriteDish(t *testing.T) {
	type args struct {
		ctx    context.Context
		userID int64
		dishID int64
	}
	tests := []struct {
		name    string
		favRepo func(ctrl *gomock.Controller) favRepo
		args    args
		wantErr bool
	}{
		{
			name: "success",
			favRepo: func(ctrl *gomock.Controller) favRepo {
				mock := NewMockfavRepo(ctrl)
				mock.EXPECT().RemoveFavoriteDish(gomock.Any(), int64(1), int64(1)).Return(nil)
				return mock
			},
			args: args{
				ctx:    context.Background(),
				userID: 1,
				dishID: 1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(tt.favRepo(ctrl))
			if err := u.RemoveFavoriteDish(tt.args.ctx, tt.args.userID, tt.args.dishID); (err != nil) != tt.wantErr {
				t.Errorf("RemoveFavoriteDish() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
