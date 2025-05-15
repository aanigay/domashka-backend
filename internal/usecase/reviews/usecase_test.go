package reviews

import (
	"context"
	cartentity "domashka-backend/internal/entity/cart"
	"domashka-backend/internal/entity/reviews"
	"domashka-backend/internal/entity/users"
	"fmt"
	"github.com/golang/mock/gomock"
	"reflect"
	"testing"
	"time"
)

func TestUsecase_Create(t *testing.T) {
	type fields struct {
		repo                reviewRepo
		userRepo            userRepo
		orderRepo           orderRepo
		dishReviewsProducer KafkaWriter
		chefReviewsProducer KafkaWriter
	}
	type args struct {
		ctx context.Context
		rv  *reviews.Review
	}
	tests := []struct {
		name                string
		repo                func(ctrl *gomock.Controller) reviewRepo
		userRepo            func(ctrl *gomock.Controller) userRepo
		orderRepo           func(ctrl *gomock.Controller) orderRepo
		dishReviewsProducer func(ctrl *gomock.Controller) KafkaWriter
		chefReviewsProducer func(ctrl *gomock.Controller) KafkaWriter
		args                args
		wantErr             bool
	}{
		{
			name: "success",
			repo: func(ctrl *gomock.Controller) reviewRepo {
				m := NewMockreviewRepo(ctrl)
				m.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
				return m
			},
			userRepo: func(ctrl *gomock.Controller) userRepo {
				m := NewMockuserRepo(ctrl)
				return m
			},
			orderRepo: func(ctrl *gomock.Controller) orderRepo {
				m := NewMockorderRepo(ctrl)
				m.EXPECT().GetCartItemsByOrderID(gomock.Any(), gomock.Any()).Return([]cartentity.CartItem{}, nil)
				return m
			},
			dishReviewsProducer: func(ctrl *gomock.Controller) KafkaWriter {
				m := NewMockKafkaWriter(ctrl)
				m.EXPECT().WriteMessages(gomock.Any(), gomock.Any()).Return(nil)
				return m
			},
			chefReviewsProducer: func(ctrl *gomock.Controller) KafkaWriter {
				m := NewMockKafkaWriter(ctrl)
				m.EXPECT().WriteMessages(gomock.Any(), gomock.Any()).Return(nil)
				return m
			},
			args: args{
				ctx: context.Background(),
				rv: &reviews.Review{
					ChefID:          1,
					UserID:          1,
					Stars:           3,
					Comment:         "WOW",
					IsVerified:      false,
					IncludeInRating: false,
					IsDeleted:       false,
					CreatedAt:       time.Time{},
					UpdatedAt:       time.Time{},
					OrderID:         1,
				},
			},
		},
		{
			name: "rv nil",
			repo: func(ctrl *gomock.Controller) reviewRepo {
				m := NewMockreviewRepo(ctrl)
				return m
			},
			userRepo: func(ctrl *gomock.Controller) userRepo {
				m := NewMockuserRepo(ctrl)
				return m
			},
			orderRepo: func(ctrl *gomock.Controller) orderRepo {
				m := NewMockorderRepo(ctrl)
				return m
			},
			dishReviewsProducer: func(ctrl *gomock.Controller) KafkaWriter {
				m := NewMockKafkaWriter(ctrl)
				return m
			},
			chefReviewsProducer: func(ctrl *gomock.Controller) KafkaWriter {
				m := NewMockKafkaWriter(ctrl)
				return m
			},
			args: args{
				ctx: context.Background(),
				rv:  nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(
				tt.repo(ctrl),
				tt.userRepo(ctrl),
				tt.orderRepo(ctrl),
				tt.dishReviewsProducer(ctrl),
				tt.chefReviewsProducer(ctrl),
			)
			if err := u.Create(tt.args.ctx, tt.args.rv); (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUsecase_CreateReview(t *testing.T) {
	type fields struct {
		repo                reviewRepo
		userRepo            userRepo
		orderRepo           orderRepo
		dishReviewsProducer KafkaWriter
		chefReviewsProducer KafkaWriter
	}
	type args struct {
		ctx    context.Context
		review reviews.Review
	}
	tests := []struct {
		name                string
		repo                func(ctrl *gomock.Controller) reviewRepo
		userRepo            func(ctrl *gomock.Controller) userRepo
		orderRepo           func(ctrl *gomock.Controller) orderRepo
		dishReviewsProducer func(ctrl *gomock.Controller) KafkaWriter
		chefReviewsProducer func(ctrl *gomock.Controller) KafkaWriter
		args                args
		wantErr             bool
	}{
		{
			name: "success",
			repo: func(ctrl *gomock.Controller) reviewRepo {
				m := NewMockreviewRepo(ctrl)
				m.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
				return m
			},
			userRepo: func(ctrl *gomock.Controller) userRepo {
				m := NewMockuserRepo(ctrl)
				return m
			},
			orderRepo: func(ctrl *gomock.Controller) orderRepo {
				m := NewMockorderRepo(ctrl)
				m.EXPECT().GetCartItemsByOrderID(gomock.Any(), gomock.Any()).Return([]cartentity.CartItem{}, nil)
				return m
			},
			dishReviewsProducer: func(ctrl *gomock.Controller) KafkaWriter {
				m := NewMockKafkaWriter(ctrl)
				m.EXPECT().WriteMessages(gomock.Any(), gomock.Any()).Return(nil)
				return m
			},
			chefReviewsProducer: func(ctrl *gomock.Controller) KafkaWriter {
				m := NewMockKafkaWriter(ctrl)
				m.EXPECT().WriteMessages(gomock.Any(), gomock.Any()).Return(nil)
				return m
			},
			args: args{
				ctx: context.Background(),
				review: reviews.Review{
					ChefID:          1,
					UserID:          1,
					Stars:           3,
					Comment:         "WOW",
					IsVerified:      false,
					IncludeInRating: false,
					IsDeleted:       false,
					CreatedAt:       time.Time{},
					UpdatedAt:       time.Time{},
					OrderID:         1,
				},
			},
		},
		{
			name: "rv nil",
			repo: func(ctrl *gomock.Controller) reviewRepo {
				m := NewMockreviewRepo(ctrl)
				return m
			},
			userRepo: func(ctrl *gomock.Controller) userRepo {
				m := NewMockuserRepo(ctrl)
				return m
			},
			orderRepo: func(ctrl *gomock.Controller) orderRepo {
				m := NewMockorderRepo(ctrl)
				return m
			},
			dishReviewsProducer: func(ctrl *gomock.Controller) KafkaWriter {
				m := NewMockKafkaWriter(ctrl)
				return m
			},
			chefReviewsProducer: func(ctrl *gomock.Controller) KafkaWriter {
				m := NewMockKafkaWriter(ctrl)
				return m
			},
			args: args{
				ctx: context.Background(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(
				tt.repo(ctrl),
				tt.userRepo(ctrl),
				tt.orderRepo(ctrl),
				tt.dishReviewsProducer(ctrl),
				tt.chefReviewsProducer(ctrl),
			)
			if err := u.CreateReview(tt.args.ctx, tt.args.review); (err != nil) != tt.wantErr {
				t.Errorf("CreateReview() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUsecase_GetByID(t *testing.T) {
	type fields struct {
		repo                reviewRepo
		userRepo            userRepo
		orderRepo           orderRepo
		dishReviewsProducer KafkaWriter
		chefReviewsProducer KafkaWriter
	}
	type args struct {
		ctx context.Context
		id  int64
	}
	tests := []struct {
		name                string
		repo                func(ctrl *gomock.Controller) reviewRepo
		userRepo            func(ctrl *gomock.Controller) userRepo
		orderRepo           func(ctrl *gomock.Controller) orderRepo
		dishReviewsProducer func(ctrl *gomock.Controller) KafkaWriter
		chefReviewsProducer func(ctrl *gomock.Controller) KafkaWriter
		args                args
		want                *reviews.Review
		wantErr             bool
	}{
		{
			name: "success",
			repo: func(ctrl *gomock.Controller) reviewRepo {
				m := NewMockreviewRepo(ctrl)
				m.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(&reviews.Review{}, nil)
				return m
			},
			userRepo: func(ctrl *gomock.Controller) userRepo {
				m := NewMockuserRepo(ctrl)
				return m
			},
			orderRepo: func(ctrl *gomock.Controller) orderRepo {
				m := NewMockorderRepo(ctrl)
				return m
			},
			dishReviewsProducer: func(ctrl *gomock.Controller) KafkaWriter {
				m := NewMockKafkaWriter(ctrl)
				return m
			},
			chefReviewsProducer: func(ctrl *gomock.Controller) KafkaWriter {
				m := NewMockKafkaWriter(ctrl)
				return m
			},
			want: &reviews.Review{},
		},
		{
			name: "error",
			repo: func(ctrl *gomock.Controller) reviewRepo {
				m := NewMockreviewRepo(ctrl)
				m.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error"))
				return m
			},
			userRepo: func(ctrl *gomock.Controller) userRepo {
				m := NewMockuserRepo(ctrl)
				return m
			},
			orderRepo: func(ctrl *gomock.Controller) orderRepo {
				m := NewMockorderRepo(ctrl)
				return m
			},
			dishReviewsProducer: func(ctrl *gomock.Controller) KafkaWriter {
				m := NewMockKafkaWriter(ctrl)
				return m
			},
			chefReviewsProducer: func(ctrl *gomock.Controller) KafkaWriter {
				m := NewMockKafkaWriter(ctrl)
				return m
			},
			wantErr: true,
		},
		{
			name: "deleted",
			repo: func(ctrl *gomock.Controller) reviewRepo {
				m := NewMockreviewRepo(ctrl)
				m.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(&reviews.Review{IsDeleted: true}, nil)
				return m
			},
			userRepo: func(ctrl *gomock.Controller) userRepo {
				m := NewMockuserRepo(ctrl)
				return m
			},
			orderRepo: func(ctrl *gomock.Controller) orderRepo {
				m := NewMockorderRepo(ctrl)
				return m
			},
			dishReviewsProducer: func(ctrl *gomock.Controller) KafkaWriter {
				m := NewMockKafkaWriter(ctrl)
				return m
			},
			chefReviewsProducer: func(ctrl *gomock.Controller) KafkaWriter {
				m := NewMockKafkaWriter(ctrl)
				return m
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(
				tt.repo(ctrl),
				tt.userRepo(ctrl),
				tt.orderRepo(ctrl),
				tt.dishReviewsProducer(ctrl),
				tt.chefReviewsProducer(ctrl),
			)
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

func TestUsecase_GetFullReviewsByChefID(t *testing.T) {
	type fields struct {
		repo                reviewRepo
		userRepo            userRepo
		orderRepo           orderRepo
		dishReviewsProducer KafkaWriter
		chefReviewsProducer KafkaWriter
	}
	type args struct {
		ctx    context.Context
		chefID int64
		limit  int
	}
	tests := []struct {
		name                string
		repo                func(ctrl *gomock.Controller) reviewRepo
		userRepo            func(ctrl *gomock.Controller) userRepo
		orderRepo           func(ctrl *gomock.Controller) orderRepo
		dishReviewsProducer func(ctrl *gomock.Controller) KafkaWriter
		chefReviewsProducer func(ctrl *gomock.Controller) KafkaWriter
		args                args
		want                []reviews.ReviewWithUser
		wantErr             bool
	}{
		{
			name: "success",
			repo: func(ctrl *gomock.Controller) reviewRepo {
				m := NewMockreviewRepo(ctrl)
				m.EXPECT().ListByChef(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]reviews.Review{
					{
						UserID: 1,
					},
				}, nil)
				return m
			},
			userRepo: func(ctrl *gomock.Controller) userRepo {
				m := NewMockuserRepo(ctrl)
				m.EXPECT().GetByID(gomock.Any(), int64(1)).Return(&users.User{
					FirstName: "kek",
				}, nil)
				return m
			},
			orderRepo: func(ctrl *gomock.Controller) orderRepo {
				m := NewMockorderRepo(ctrl)
				return m
			},
			dishReviewsProducer: func(ctrl *gomock.Controller) KafkaWriter {
				m := NewMockKafkaWriter(ctrl)
				return m
			},
			chefReviewsProducer: func(ctrl *gomock.Controller) KafkaWriter {
				m := NewMockKafkaWriter(ctrl)
				return m
			},
			want: []reviews.ReviewWithUser{
				{
					Review:   reviews.Review{UserID: 1},
					UserName: "kek",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(
				tt.repo(ctrl),
				tt.userRepo(ctrl),
				tt.orderRepo(ctrl),
				tt.dishReviewsProducer(ctrl),
				tt.chefReviewsProducer(ctrl),
			)
			got, err := u.GetFullReviewsByChefID(tt.args.ctx, tt.args.chefID, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFullReviewsByChefID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFullReviewsByChefID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsecase_GetReviewByOrderAndUserID(t *testing.T) {
	type fields struct {
		repo                reviewRepo
		userRepo            userRepo
		orderRepo           orderRepo
		dishReviewsProducer KafkaWriter
		chefReviewsProducer KafkaWriter
	}
	type args struct {
		ctx    context.Context
		chefID int64
		userID int64
	}
	tests := []struct {
		name                string
		repo                func(ctrl *gomock.Controller) reviewRepo
		userRepo            func(ctrl *gomock.Controller) userRepo
		orderRepo           func(ctrl *gomock.Controller) orderRepo
		dishReviewsProducer func(ctrl *gomock.Controller) KafkaWriter
		chefReviewsProducer func(ctrl *gomock.Controller) KafkaWriter
		args                args
		want                *reviews.Review
		wantErr             bool
	}{
		{
			name: "success",
			repo: func(ctrl *gomock.Controller) reviewRepo {
				m := NewMockreviewRepo(ctrl)
				m.EXPECT().GetReviewByOrderAndUserID(gomock.Any(), gomock.Any(), gomock.Any()).Return(&reviews.Review{}, nil)
				return m
			},
			userRepo: func(ctrl *gomock.Controller) userRepo {
				m := NewMockuserRepo(ctrl)
				return m
			},
			orderRepo: func(ctrl *gomock.Controller) orderRepo {
				m := NewMockorderRepo(ctrl)
				return m
			},
			dishReviewsProducer: func(ctrl *gomock.Controller) KafkaWriter {
				m := NewMockKafkaWriter(ctrl)
				return m
			},
			chefReviewsProducer: func(ctrl *gomock.Controller) KafkaWriter {
				m := NewMockKafkaWriter(ctrl)
				return m
			},
			want: &reviews.Review{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(
				tt.repo(ctrl),
				tt.userRepo(ctrl),
				tt.orderRepo(ctrl),
				tt.dishReviewsProducer(ctrl),
				tt.chefReviewsProducer(ctrl),
			)
			got, err := u.GetReviewByOrderAndUserID(tt.args.ctx, tt.args.chefID, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetReviewByOrderAndUserID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetReviewByOrderAndUserID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsecase_GetReviewsByChefID(t *testing.T) {
	type fields struct {
		repo                reviewRepo
		userRepo            userRepo
		orderRepo           orderRepo
		dishReviewsProducer KafkaWriter
		chefReviewsProducer KafkaWriter
	}
	type args struct {
		ctx    context.Context
		chefID int64
	}
	tests := []struct {
		name                string
		repo                func(ctrl *gomock.Controller) reviewRepo
		userRepo            func(ctrl *gomock.Controller) userRepo
		orderRepo           func(ctrl *gomock.Controller) orderRepo
		dishReviewsProducer func(ctrl *gomock.Controller) KafkaWriter
		chefReviewsProducer func(ctrl *gomock.Controller) KafkaWriter
		args                args
		want                []reviews.Review
		wantErr             bool
	}{
		{
			name: "success",
			repo: func(ctrl *gomock.Controller) reviewRepo {
				m := NewMockreviewRepo(ctrl)
				m.EXPECT().ListByChef(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]reviews.Review{
					{
						UserID: 1,
					},
				}, nil)
				return m
			},
			userRepo: func(ctrl *gomock.Controller) userRepo {
				m := NewMockuserRepo(ctrl)
				return m
			},
			orderRepo: func(ctrl *gomock.Controller) orderRepo {
				m := NewMockorderRepo(ctrl)
				return m
			},
			dishReviewsProducer: func(ctrl *gomock.Controller) KafkaWriter {
				m := NewMockKafkaWriter(ctrl)
				return m
			},
			chefReviewsProducer: func(ctrl *gomock.Controller) KafkaWriter {
				m := NewMockKafkaWriter(ctrl)
				return m
			},
			want: []reviews.Review{
				{
					UserID: 1,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(
				tt.repo(ctrl),
				tt.userRepo(ctrl),
				tt.orderRepo(ctrl),
				tt.dishReviewsProducer(ctrl),
				tt.chefReviewsProducer(ctrl),
			)
			got, err := u.GetReviewsByChefID(tt.args.ctx, tt.args.chefID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetReviewsByChefID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetReviewsByChefID() got = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestUsecase_ListByChef(t *testing.T) {
	type fields struct {
		repo                reviewRepo
		userRepo            userRepo
		orderRepo           orderRepo
		dishReviewsProducer KafkaWriter
		chefReviewsProducer KafkaWriter
	}
	type args struct {
		ctx         context.Context
		chefID      int64
		includeOnly bool
		limit       *int
	}
	tests := []struct {
		name                string
		repo                func(ctrl *gomock.Controller) reviewRepo
		userRepo            func(ctrl *gomock.Controller) userRepo
		orderRepo           func(ctrl *gomock.Controller) orderRepo
		dishReviewsProducer func(ctrl *gomock.Controller) KafkaWriter
		chefReviewsProducer func(ctrl *gomock.Controller) KafkaWriter
		args                args
		want                []reviews.Review
		wantErr             bool
	}{
		{
			name: "success",
			repo: func(ctrl *gomock.Controller) reviewRepo {
				m := NewMockreviewRepo(ctrl)
				m.EXPECT().ListByChef(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]reviews.Review{
					{
						UserID: 1,
					},
				}, nil)
				return m
			},
			userRepo: func(ctrl *gomock.Controller) userRepo {
				m := NewMockuserRepo(ctrl)
				return m
			},
			orderRepo: func(ctrl *gomock.Controller) orderRepo {
				m := NewMockorderRepo(ctrl)
				return m
			},
			dishReviewsProducer: func(ctrl *gomock.Controller) KafkaWriter {
				m := NewMockKafkaWriter(ctrl)
				return m
			},
			chefReviewsProducer: func(ctrl *gomock.Controller) KafkaWriter {
				m := NewMockKafkaWriter(ctrl)
				return m
			},
			want: []reviews.Review{
				{
					UserID: 1,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(
				tt.repo(ctrl),
				tt.userRepo(ctrl),
				tt.orderRepo(ctrl),
				tt.dishReviewsProducer(ctrl),
				tt.chefReviewsProducer(ctrl),
			)
			got, err := u.ListByChef(tt.args.ctx, tt.args.chefID, tt.args.includeOnly, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListByChef() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListByChef() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsecase_SoftDelete(t *testing.T) {
	type fields struct {
		repo                reviewRepo
		userRepo            userRepo
		orderRepo           orderRepo
		dishReviewsProducer KafkaWriter
		chefReviewsProducer KafkaWriter
	}
	type args struct {
		ctx context.Context
		id  int64
	}
	tests := []struct {
		name                string
		repo                func(ctrl *gomock.Controller) reviewRepo
		userRepo            func(ctrl *gomock.Controller) userRepo
		orderRepo           func(ctrl *gomock.Controller) orderRepo
		dishReviewsProducer func(ctrl *gomock.Controller) KafkaWriter
		chefReviewsProducer func(ctrl *gomock.Controller) KafkaWriter
		args                args
		wantErr             bool
	}{
		{
			name: "success",
			repo: func(ctrl *gomock.Controller) reviewRepo {
				m := NewMockreviewRepo(ctrl)
				m.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(&reviews.Review{}, nil)
				m.EXPECT().SoftDelete(gomock.Any(), gomock.Any()).Return(nil)
				return m
			},
			userRepo: func(ctrl *gomock.Controller) userRepo {
				m := NewMockuserRepo(ctrl)
				return m
			},
			orderRepo: func(ctrl *gomock.Controller) orderRepo {
				m := NewMockorderRepo(ctrl)
				return m
			},
			dishReviewsProducer: func(ctrl *gomock.Controller) KafkaWriter {
				m := NewMockKafkaWriter(ctrl)
				return m
			},
			chefReviewsProducer: func(ctrl *gomock.Controller) KafkaWriter {
				m := NewMockKafkaWriter(ctrl)
				return m
			},
		},
		{
			name: "error",
			repo: func(ctrl *gomock.Controller) reviewRepo {
				m := NewMockreviewRepo(ctrl)
				m.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("some error"))
				return m
			},
			userRepo: func(ctrl *gomock.Controller) userRepo {
				m := NewMockuserRepo(ctrl)
				return m
			},
			orderRepo: func(ctrl *gomock.Controller) orderRepo {
				m := NewMockorderRepo(ctrl)
				return m
			},
			dishReviewsProducer: func(ctrl *gomock.Controller) KafkaWriter {
				m := NewMockKafkaWriter(ctrl)
				return m
			},
			chefReviewsProducer: func(ctrl *gomock.Controller) KafkaWriter {
				m := NewMockKafkaWriter(ctrl)
				return m
			},
			wantErr: true,
		},
		{
			name: "error delete",
			repo: func(ctrl *gomock.Controller) reviewRepo {
				m := NewMockreviewRepo(ctrl)
				m.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(&reviews.Review{}, nil)
				m.EXPECT().SoftDelete(gomock.Any(), gomock.Any()).Return(fmt.Errorf("some error"))
				return m
			},
			userRepo: func(ctrl *gomock.Controller) userRepo {
				m := NewMockuserRepo(ctrl)
				return m
			},
			orderRepo: func(ctrl *gomock.Controller) orderRepo {
				m := NewMockorderRepo(ctrl)
				return m
			},
			dishReviewsProducer: func(ctrl *gomock.Controller) KafkaWriter {
				m := NewMockKafkaWriter(ctrl)
				return m
			},
			chefReviewsProducer: func(ctrl *gomock.Controller) KafkaWriter {
				m := NewMockKafkaWriter(ctrl)
				return m
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(
				tt.repo(ctrl),
				tt.userRepo(ctrl),
				tt.orderRepo(ctrl),
				tt.dishReviewsProducer(ctrl),
				tt.chefReviewsProducer(ctrl),
			)
			if err := u.SoftDelete(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("SoftDelete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUsecase_Update(t *testing.T) {
	type fields struct {
		repo                reviewRepo
		userRepo            userRepo
		orderRepo           orderRepo
		dishReviewsProducer KafkaWriter
		chefReviewsProducer KafkaWriter
	}
	type args struct {
		ctx context.Context
		rv  *reviews.Review
	}
	tests := []struct {
		name                string
		repo                func(ctrl *gomock.Controller) reviewRepo
		userRepo            func(ctrl *gomock.Controller) userRepo
		orderRepo           func(ctrl *gomock.Controller) orderRepo
		dishReviewsProducer func(ctrl *gomock.Controller) KafkaWriter
		chefReviewsProducer func(ctrl *gomock.Controller) KafkaWriter
		args                args
		wantErr             bool
	}{
		{
			name: "success",
			repo: func(ctrl *gomock.Controller) reviewRepo {
				m := NewMockreviewRepo(ctrl)
				m.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(&reviews.Review{IsDeleted: false}, nil)
				m.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
				return m
			},
			userRepo: func(ctrl *gomock.Controller) userRepo {
				m := NewMockuserRepo(ctrl)
				return m
			},
			orderRepo: func(ctrl *gomock.Controller) orderRepo {
				m := NewMockorderRepo(ctrl)
				return m
			},
			dishReviewsProducer: func(ctrl *gomock.Controller) KafkaWriter {
				m := NewMockKafkaWriter(ctrl)
				return m
			},
			chefReviewsProducer: func(ctrl *gomock.Controller) KafkaWriter {
				m := NewMockKafkaWriter(ctrl)
				return m
			},
			args: args{
				ctx: context.Background(),
				rv:  &reviews.Review{IsDeleted: false, ID: 1, Stars: 5},
			},
		},
		{
			name: "deleted",
			repo: func(ctrl *gomock.Controller) reviewRepo {
				m := NewMockreviewRepo(ctrl)
				m.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(&reviews.Review{IsDeleted: true}, nil)
				return m
			},
			userRepo: func(ctrl *gomock.Controller) userRepo {
				m := NewMockuserRepo(ctrl)
				return m
			},
			orderRepo: func(ctrl *gomock.Controller) orderRepo {
				m := NewMockorderRepo(ctrl)
				return m
			},
			dishReviewsProducer: func(ctrl *gomock.Controller) KafkaWriter {
				m := NewMockKafkaWriter(ctrl)
				return m
			},
			chefReviewsProducer: func(ctrl *gomock.Controller) KafkaWriter {
				m := NewMockKafkaWriter(ctrl)
				return m
			},
			args: args{
				ctx: context.Background(),
				rv:  &reviews.Review{IsDeleted: false, ID: 1, Stars: 5},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(
				tt.repo(ctrl),
				tt.userRepo(ctrl),
				tt.orderRepo(ctrl),
				tt.dishReviewsProducer(ctrl),
				tt.chefReviewsProducer(ctrl),
			)
			if err := u.Update(tt.args.ctx, tt.args.rv); (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
