package users

import (
	"context"
	chefsEntity "domashka-backend/internal/entity/chefs"
	dishesEntity "domashka-backend/internal/entity/dishes"
	"domashka-backend/internal/entity/users"
	"domashka-backend/internal/utils/pointers"
	"github.com/golang/mock/gomock"
	"reflect"
	"testing"
)

func TestUseCase_CheckIfUserIsChef(t *testing.T) {
	type fields struct {
		repo repo
	}
	type args struct {
		ctx    context.Context
		userID int64
	}
	tests := []struct {
		name    string
		repo    func(ctrl *gomock.Controller) repo
		args    args
		want    *int64
		want1   bool
		wantErr bool
	}{
		{
			name: "success",
			repo: func(ctrl *gomock.Controller) repo {
				m := NewMockrepo(ctrl)
				m.EXPECT().CheckIfUserIsChef(gomock.Any(), gomock.Any()).Return(pointers.To(int64(1)), true, nil)
				return m
			},
			want:  pointers.To(int64(1)),
			want1: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(tt.repo(ctrl))
			got, got1, err := u.CheckIfUserIsChef(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckIfUserIsChef() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CheckIfUserIsChef() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("CheckIfUserIsChef() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestUseCase_Create(t *testing.T) {
	type fields struct {
		repo repo
	}
	type args struct {
		ctx  context.Context
		user *users.User
	}
	tests := []struct {
		name    string
		repo    func(ctrl *gomock.Controller) repo
		args    args
		wantErr bool
	}{
		{
			name: "success",
			repo: func(ctrl *gomock.Controller) repo {
				m := NewMockrepo(ctrl)
				m.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
				return m
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(tt.repo(ctrl))
			if err := u.Create(tt.args.ctx, tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUseCase_Delete(t *testing.T) {
	type fields struct {
		repo repo
	}
	type args struct {
		ctx context.Context
		id  int64
	}
	tests := []struct {
		name    string
		repo    func(ctrl *gomock.Controller) repo
		args    args
		wantErr bool
	}{
		{
			name: "success",
			repo: func(ctrl *gomock.Controller) repo {
				m := NewMockrepo(ctrl)
				m.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil)
				return m
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(tt.repo(ctrl))
			if err := u.Delete(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUseCase_GetByID(t *testing.T) {
	type fields struct {
		repo repo
	}
	type args struct {
		ctx context.Context
		id  int64
	}
	tests := []struct {
		name    string
		repo    func(ctrl *gomock.Controller) repo
		args    args
		want    *users.User
		wantErr bool
	}{
		{
			name: "success",
			repo: func(ctrl *gomock.Controller) repo {
				m := NewMockrepo(ctrl)
				m.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(&users.User{}, nil)
				return m
			},
			want: &users.User{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(tt.repo(ctrl))
			got, err := u.GetByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetByID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCase_GetFavoritesChefsByUserID(t *testing.T) {
	type fields struct {
		repo repo
	}
	type args struct {
		ctx    context.Context
		userID int64
	}
	tests := []struct {
		name    string
		repo    func(ctrl *gomock.Controller) repo
		args    args
		want    []chefsEntity.Chef
		wantErr bool
	}{
		{
			name: "success",
			repo: func(ctrl *gomock.Controller) repo {
				m := NewMockrepo(ctrl)
				m.EXPECT().GetFavoritesChefsByUserID(gomock.Any(), gomock.Any()).Return([]chefsEntity.Chef{}, nil)
				return m
			},
			want: []chefsEntity.Chef{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(tt.repo(ctrl))
			got, err := u.GetFavoritesChefsByUserID(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFavoritesChefsByUserID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFavoritesChefsByUserID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCase_GetFavoritesDishesByUserID(t *testing.T) {
	type fields struct {
		repo repo
	}
	type args struct {
		ctx    context.Context
		userID int64
	}
	tests := []struct {
		name    string
		repo    func(ctrl *gomock.Controller) repo
		args    args
		want    []dishesEntity.Dish
		wantErr bool
	}{
		{
			name: "success",
			repo: func(ctrl *gomock.Controller) repo {
				m := NewMockrepo(ctrl)
				m.EXPECT().GetFavoritesDishesByUserID(gomock.Any(), gomock.Any()).Return([]dishesEntity.Dish{}, nil)
				return m
			},
			want: []dishesEntity.Dish{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(tt.repo(ctrl))
			got, err := u.GetFavoritesDishesByUserID(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFavoritesDishesByUserID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFavoritesDishesByUserID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCase_Update(t *testing.T) {
	type fields struct {
		repo repo
	}
	type args struct {
		ctx  context.Context
		id   int64
		user users.User
	}
	tests := []struct {
		name    string
		repo    func(ctrl *gomock.Controller) repo
		args    args
		wantErr bool
	}{
		{
			name: "success",
			repo: func(ctrl *gomock.Controller) repo {
				m := NewMockrepo(ctrl)
				m.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return m
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(tt.repo(ctrl))
			if err := u.Update(tt.args.ctx, tt.args.id, tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
