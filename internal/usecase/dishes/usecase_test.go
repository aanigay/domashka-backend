package dishes

import (
	"context"
	entity "domashka-backend/internal/entity/dishes"
	"github.com/golang/mock/gomock"
	"mime/multipart"
	"reflect"
	"testing"
)

func TestUsecase_Create(t *testing.T) {
	type args struct {
		ctx         context.Context
		dish        *entity.Dish
		nutrition   *entity.Nutrition
		sizes       []entity.Size
		ingredients []entity.Ingredient
	}
	tests := []struct {
		name     string
		dishRepo func(ctrl *gomock.Controller) dishRepo
		s3Client func(ctrl *gomock.Controller) s3Client
		args     args
		want     int64
		wantErr  bool
	}{
		{
			name: "success",
			args: args{},
			dishRepo: func(ctrl *gomock.Controller) dishRepo {
				m := NewMockdishRepo(ctrl)
				return m
			},
			s3Client: func(ctrl *gomock.Controller) s3Client {
				m := NewMocks3Client(ctrl)
				return m
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(tt.dishRepo(ctrl), tt.s3Client(ctrl))
			got, err := u.Create(tt.args.ctx, tt.args.dish, tt.args.nutrition, tt.args.sizes, tt.args.ingredients)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Create() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsecase_Delete(t *testing.T) {
	type args struct {
		ctx    context.Context
		dishID int64
	}
	tests := []struct {
		name     string
		dishRepo func(ctrl *gomock.Controller) dishRepo
		s3Client func(ctrl *gomock.Controller) s3Client
		args     args
		wantErr  bool
	}{
		{
			name: "success",
			args: args{},
			dishRepo: func(ctrl *gomock.Controller) dishRepo {
				m := NewMockdishRepo(ctrl)
				return m
			},
			s3Client: func(ctrl *gomock.Controller) s3Client {
				m := NewMocks3Client(ctrl)
				return m
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(tt.dishRepo(ctrl), tt.s3Client(ctrl))
			if err := u.Delete(tt.args.ctx, tt.args.dishID); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUsecase_GetAll(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name     string
		dishRepo func(ctrl *gomock.Controller) dishRepo
		s3Client func(ctrl *gomock.Controller) s3Client
		args     args
		want     []entity.Dish
		wantErr  bool
	}{
		{
			name: "success",
			args: args{},
			dishRepo: func(ctrl *gomock.Controller) dishRepo {
				m := NewMockdishRepo(ctrl)
				return m
			},
			s3Client: func(ctrl *gomock.Controller) s3Client {
				m := NewMocks3Client(ctrl)
				return m
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(tt.dishRepo(ctrl), tt.s3Client(ctrl))
			got, err := u.GetAll(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAll() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsecase_GetAllCategories(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name     string
		dishRepo func(ctrl *gomock.Controller) dishRepo
		s3Client func(ctrl *gomock.Controller) s3Client
		args     args
		want     []entity.Category
		wantErr  bool
	}{
		{
			name: "success",
			args: args{},
			dishRepo: func(ctrl *gomock.Controller) dishRepo {
				m := NewMockdishRepo(ctrl)
				return m
			},
			s3Client: func(ctrl *gomock.Controller) s3Client {
				m := NewMocks3Client(ctrl)
				return m
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(tt.dishRepo(ctrl), tt.s3Client(ctrl))
			got, err := u.GetAllCategories(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllCategories() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllCategories() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsecase_GetAllDishesByChefID(t *testing.T) {
	type args struct {
		ctx    context.Context
		chefID int64
	}
	tests := []struct {
		name     string
		dishRepo func(ctrl *gomock.Controller) dishRepo
		s3Client func(ctrl *gomock.Controller) s3Client
		args     args
		want     []entity.Dish
		wantErr  bool
	}{
		{
			name: "success",
			args: args{},
			dishRepo: func(ctrl *gomock.Controller) dishRepo {
				m := NewMockdishRepo(ctrl)
				return m
			},
			s3Client: func(ctrl *gomock.Controller) s3Client {
				m := NewMocks3Client(ctrl)
				return m
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(tt.dishRepo(ctrl), tt.s3Client(ctrl))
			got, err := u.GetAllDishesByChefID(tt.args.ctx, tt.args.chefID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllDishesByChefID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllDishesByChefID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsecase_GetAllIngredients(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name     string
		dishRepo func(ctrl *gomock.Controller) dishRepo
		s3Client func(ctrl *gomock.Controller) s3Client
		args     args
		want     []entity.Ingredient
		wantErr  bool
	}{
		{
			name: "success",
			args: args{},
			dishRepo: func(ctrl *gomock.Controller) dishRepo {
				m := NewMockdishRepo(ctrl)
				return m
			},
			s3Client: func(ctrl *gomock.Controller) s3Client {
				m := NewMocks3Client(ctrl)
				return m
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(tt.dishRepo(ctrl), tt.s3Client(ctrl))
			got, err := u.GetAllIngredients(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllIngredients() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllIngredients() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsecase_GetCategoryTitleByDishID(t *testing.T) {
	type args struct {
		ctx    context.Context
		dishID int64
	}
	tests := []struct {
		name     string
		dishRepo func(ctrl *gomock.Controller) dishRepo
		s3Client func(ctrl *gomock.Controller) s3Client
		args     args
		want     string
		wantErr  bool
	}{
		{
			name: "success",
			args: args{},
			dishRepo: func(ctrl *gomock.Controller) dishRepo {
				m := NewMockdishRepo(ctrl)
				return m
			},
			s3Client: func(ctrl *gomock.Controller) s3Client {
				m := NewMocks3Client(ctrl)
				return m
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(tt.dishRepo(ctrl), tt.s3Client(ctrl))
			got, err := u.GetCategoryTitleByDishID(tt.args.ctx, tt.args.dishID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCategoryTitleByDishID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetCategoryTitleByDishID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsecase_GetDishByID(t *testing.T) {
	type args struct {
		ctx    context.Context
		dishID int64
	}
	tests := []struct {
		name     string
		dishRepo func(ctrl *gomock.Controller) dishRepo
		s3Client func(ctrl *gomock.Controller) s3Client
		args     args
		want     *entity.Dish
		wantErr  bool
	}{
		{
			name: "success",
			args: args{},
			dishRepo: func(ctrl *gomock.Controller) dishRepo {
				m := NewMockdishRepo(ctrl)
				return m
			},
			s3Client: func(ctrl *gomock.Controller) s3Client {
				m := NewMocks3Client(ctrl)
				return m
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(tt.dishRepo(ctrl), tt.s3Client(ctrl))
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

func TestUsecase_GetDishRatingByID(t *testing.T) {
	type args struct {
		ctx    context.Context
		dishID int64
	}
	tests := []struct {
		name     string
		dishRepo func(ctrl *gomock.Controller) dishRepo
		s3Client func(ctrl *gomock.Controller) s3Client
		args     args
		want     *entity.Dish
		wantErr  bool
	}{
		{
			name: "success",
			args: args{},
			dishRepo: func(ctrl *gomock.Controller) dishRepo {
				m := NewMockdishRepo(ctrl)
				return m
			},
			s3Client: func(ctrl *gomock.Controller) s3Client {
				m := NewMocks3Client(ctrl)
				return m
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(tt.dishRepo(ctrl), tt.s3Client(ctrl))
			got, err := u.GetDishRatingByID(tt.args.ctx, tt.args.dishID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDishRatingByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetDishRatingByID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsecase_GetDishSizesByDishID(t *testing.T) {
	type args struct {
		ctx    context.Context
		dishID int64
	}
	tests := []struct {
		name     string
		dishRepo func(ctrl *gomock.Controller) dishRepo
		s3Client func(ctrl *gomock.Controller) s3Client
		args     args
		want     []entity.Size
		wantErr  bool
	}{
		{
			name: "success",
			args: args{},
			dishRepo: func(ctrl *gomock.Controller) dishRepo {
				m := NewMockdishRepo(ctrl)
				return m
			},
			s3Client: func(ctrl *gomock.Controller) s3Client {
				m := NewMocks3Client(ctrl)
				return m
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(tt.dishRepo(ctrl), tt.s3Client(ctrl))
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
	type args struct {
		ctx    context.Context
		chefID int64
		limit  int
	}
	tests := []struct {
		name     string
		dishRepo func(ctrl *gomock.Controller) dishRepo
		s3Client func(ctrl *gomock.Controller) s3Client
		args     args
		want     []entity.Dish
		wantErr  bool
	}{
		{
			name: "success",
			args: args{},
			dishRepo: func(ctrl *gomock.Controller) dishRepo {
				m := NewMockdishRepo(ctrl)
				return m
			},
			s3Client: func(ctrl *gomock.Controller) s3Client {
				m := NewMocks3Client(ctrl)
				return m
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(tt.dishRepo(ctrl), tt.s3Client(ctrl))
			got, err := u.GetDishesByChefID(tt.args.ctx, tt.args.chefID, tt.args.limit)
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
	type args struct {
		ctx    context.Context
		dishID int64
	}
	tests := []struct {
		name     string
		dishRepo func(ctrl *gomock.Controller) dishRepo
		s3Client func(ctrl *gomock.Controller) s3Client
		args     args
		want     []entity.Ingredient
		wantErr  bool
	}{
		{
			name: "success",
			args: args{},
			dishRepo: func(ctrl *gomock.Controller) dishRepo {
				m := NewMockdishRepo(ctrl)
				return m
			},
			s3Client: func(ctrl *gomock.Controller) s3Client {
				m := NewMocks3Client(ctrl)
				return m
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(tt.dishRepo(ctrl), tt.s3Client(ctrl))
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
	type args struct {
		ctx    context.Context
		dishID int64
	}
	tests := []struct {
		name     string
		dishRepo func(ctrl *gomock.Controller) dishRepo
		s3Client func(ctrl *gomock.Controller) s3Client
		args     args
		want     *entity.Price
		wantErr  bool
	}{
		{
			name: "success",
			args: args{},
			dishRepo: func(ctrl *gomock.Controller) dishRepo {
				m := NewMockdishRepo(ctrl)
				return m
			},
			s3Client: func(ctrl *gomock.Controller) s3Client {
				m := NewMocks3Client(ctrl)
				return m
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(tt.dishRepo(ctrl), tt.s3Client(ctrl))
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
	type args struct {
		ctx    context.Context
		dishID int64
	}
	tests := []struct {
		name     string
		dishRepo func(ctrl *gomock.Controller) dishRepo
		s3Client func(ctrl *gomock.Controller) s3Client
		args     args
		want     *entity.Nutrition
		wantErr  bool
	}{
		{
			name: "success",
			args: args{},
			dishRepo: func(ctrl *gomock.Controller) dishRepo {
				m := NewMockdishRepo(ctrl)
				return m
			},
			s3Client: func(ctrl *gomock.Controller) s3Client {
				m := NewMocks3Client(ctrl)
				return m
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(tt.dishRepo(ctrl), tt.s3Client(ctrl))
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
	type args struct {
		ctx   context.Context
		limit int
	}
	tests := []struct {
		name     string
		dishRepo func(ctrl *gomock.Controller) dishRepo
		s3Client func(ctrl *gomock.Controller) s3Client
		args     args
		want     []entity.Dish
		wantErr  bool
	}{
		{
			name: "success",
			args: args{},
			dishRepo: func(ctrl *gomock.Controller) dishRepo {
				m := NewMockdishRepo(ctrl)
				return m
			},
			s3Client: func(ctrl *gomock.Controller) s3Client {
				m := NewMocks3Client(ctrl)
				return m
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(tt.dishRepo(ctrl), tt.s3Client(ctrl))
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

func TestUsecase_SetDishImage(t *testing.T) {
	type args struct {
		ctx    context.Context
		dishID int64
		image  *multipart.FileHeader
	}
	tests := []struct {
		name     string
		dishRepo func(ctrl *gomock.Controller) dishRepo
		s3Client func(ctrl *gomock.Controller) s3Client
		args     args
		want     string
		wantErr  bool
	}{
		{
			name: "success",
			args: args{},
			dishRepo: func(ctrl *gomock.Controller) dishRepo {
				m := NewMockdishRepo(ctrl)
				return m
			},
			s3Client: func(ctrl *gomock.Controller) s3Client {
				m := NewMocks3Client(ctrl)
				return m
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(tt.dishRepo(ctrl), tt.s3Client(ctrl))
			got, err := u.SetDishImage(tt.args.ctx, tt.args.dishID, tt.args.image)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetDishImage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SetDishImage() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsecase_SetIngredientImage(t *testing.T) {
	type fields struct {
		dishRepo dishRepo
		s3Client s3Client
	}
	type args struct {
		ctx          context.Context
		ingredientID int64
		image        *multipart.FileHeader
	}
	tests := []struct {
		name     string
		dishRepo func(ctrl *gomock.Controller) dishRepo
		s3Client func(ctrl *gomock.Controller) s3Client
		args     args
		want     string
		wantErr  bool
	}{
		{
			name: "success",
			args: args{},
			dishRepo: func(ctrl *gomock.Controller) dishRepo {
				m := NewMockdishRepo(ctrl)
				return m
			},
			s3Client: func(ctrl *gomock.Controller) s3Client {
				m := NewMocks3Client(ctrl)
				return m
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(tt.dishRepo(ctrl), tt.s3Client(ctrl))
			got, err := u.SetIngredientImage(tt.args.ctx, tt.args.ingredientID, tt.args.image)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetIngredientImage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SetIngredientImage() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsecase_Update(t *testing.T) {
	type args struct {
		ctx         context.Context
		dish        *entity.Dish
		nutrition   *entity.Nutrition
		sizes       []entity.Size
		ingredients []entity.Ingredient
	}
	tests := []struct {
		name     string
		dishRepo func(ctrl *gomock.Controller) dishRepo
		s3Client func(ctrl *gomock.Controller) s3Client
		args     args
		want     int64
		wantErr  bool
	}{
		{
			name: "success",
			args: args{},
			dishRepo: func(ctrl *gomock.Controller) dishRepo {
				m := NewMockdishRepo(ctrl)
				return m
			},
			s3Client: func(ctrl *gomock.Controller) s3Client {
				m := NewMocks3Client(ctrl)
				return m
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(tt.dishRepo(ctrl), tt.s3Client(ctrl))
			got, err := u.Update(tt.args.ctx, tt.args.dish, tt.args.nutrition, tt.args.sizes, tt.args.ingredients)
			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Update() got = %v, want %v", got, tt.want)
			}
		})
	}
}
