package chefs

import (
	"context"
	entity "domashka-backend/internal/entity/chefs"
	dish "domashka-backend/internal/entity/dishes"
	"mime/multipart"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		chefRepo chefRepo
	}
	tests := []struct {
		name string
		args args
		want *Usecase
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.chefRepo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsecase_GetChefAvatarURLByChefID(t *testing.T) {
	type fields struct {
		chefRepo chefRepo
	}
	type args struct {
		ctx    context.Context
		chefID int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &Usecase{
				chefRepo: tt.fields.chefRepo,
			}
			got, err := u.GetChefAvatarURLByChefID(tt.args.ctx, tt.args.chefID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetChefAvatarURLByChefID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetChefAvatarURLByChefID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsecase_GetChefAvatarURLByDishID(t *testing.T) {
	type fields struct {
		chefRepo chefRepo
	}
	type args struct {
		ctx    context.Context
		dishID int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &Usecase{
				chefRepo: tt.fields.chefRepo,
			}
			got, err := u.GetChefAvatarURLByDishID(tt.args.ctx, tt.args.dishID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetChefAvatarURLByDishID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetChefAvatarURLByDishID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsecase_GetChefByDishID(t *testing.T) {
	type fields struct {
		chefRepo chefRepo
	}
	type args struct {
		ctx    context.Context
		dishID int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *entity.Chef
		wantErr bool
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &Usecase{
				chefRepo: tt.fields.chefRepo,
			}
			got, err := u.GetChefByDishID(tt.args.ctx, tt.args.dishID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetChefByDishID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetChefByDishID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsecase_GetChefByID(t *testing.T) {
	type fields struct {
		chefRepo chefRepo
	}
	type args struct {
		ctx    context.Context
		chefID int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *entity.Chef
		wantErr bool
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &Usecase{
				chefRepo: tt.fields.chefRepo,
			}
			got, err := u.GetChefByID(tt.args.ctx, tt.args.chefID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetChefByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetChefByID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsecase_GetChefCertifications(t *testing.T) {
	type fields struct {
		chefRepo chefRepo
	}
	type args struct {
		ctx    context.Context
		chefID int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []entity.Certification
		wantErr bool
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &Usecase{
				chefRepo: tt.fields.chefRepo,
			}
			got, err := u.GetChefCertifications(tt.args.ctx, tt.args.chefID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetChefCertifications() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetChefCertifications() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsecase_GetChefExperienceYears(t *testing.T) {
	type fields struct {
		chefRepo chefRepo
	}
	type args struct {
		ctx    context.Context
		chefID int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &Usecase{
				chefRepo: tt.fields.chefRepo,
			}
			got, err := u.GetChefExperienceYears(tt.args.ctx, tt.args.chefID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetChefExperienceYears() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetChefExperienceYears() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsecase_GetDishesByChefID(t *testing.T) {
	type fields struct {
		chefRepo chefRepo
	}
	type args struct {
		ctx    context.Context
		chefID int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []dish.Dish
		wantErr bool
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &Usecase{
				chefRepo: tt.fields.chefRepo,
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

func TestUsecase_GetTopChefs(t *testing.T) {
	type fields struct {
		chefRepo chefRepo
	}
	type args struct {
		ctx   context.Context
		limit int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []entity.Chef
		wantErr bool
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &Usecase{
				chefRepo: tt.fields.chefRepo,
			}
			got, err := u.GetTopChefs(tt.args.ctx, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTopChefs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTopChefs() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsecase_UploadAvatar(t *testing.T) {
	type fields struct {
		chefRepo chefRepo
	}
	type args struct {
		ctx        context.Context
		chefID     int64
		fileHeader *multipart.FileHeader
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &Usecase{
				chefRepo: tt.fields.chefRepo,
			}
			got, err := u.UploadAvatar(tt.args.ctx, tt.args.chefID, tt.args.fileHeader)
			if (err != nil) != tt.wantErr {
				t.Errorf("UploadAvatar() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UploadAvatar() got = %v, want %v", got, tt.want)
			}
		})
	}
}
