package chefs

import (
	"context"
	entity "domashka-backend/internal/entity/chefs"
	dish "domashka-backend/internal/entity/dishes"
	"domashka-backend/internal/utils/pointers"
	"github.com/golang/mock/gomock"
	"mime/multipart"
	"reflect"
	"testing"
)

func TestUsecase_GetAll(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name     string
		chefRepo func(ctrl *gomock.Controller) chefRepo
		geoRepo  func(ctrl *gomock.Controller) geoRepo
		s3Client func(ctrl *gomock.Controller) s3Client
		args     args
		want     []entity.Chef
		wantErr  bool
	}{
		{
			name: "success",
			chefRepo: func(ctrl *gomock.Controller) chefRepo {
				m := NewMockchefRepo(ctrl)
				m.EXPECT().GetAll(gomock.Any()).Return([]entity.Chef{
					{
						ID:   1,
						Name: "test 1",
					},
					{
						ID:   2,
						Name: "test 2",
					},
				}, nil)
				m.EXPECT().GetChefRatingByChefID(gomock.Any(), int64(1)).Return(&entity.Chef{
					Rating:       pointers.To(float32(5)),
					ReviewsCount: pointers.To(int32(5)),
				}, nil)
				m.EXPECT().GetChefRatingByChefID(gomock.Any(), int64(2)).Return(&entity.Chef{
					Rating:       pointers.To(float32(5)),
					ReviewsCount: pointers.To(int32(5)),
				}, nil)
				return m
			},
			geoRepo: func(ctrl *gomock.Controller) geoRepo {
				m := NewMockgeoRepo(ctrl)
				return m
			},
			s3Client: func(ctrl *gomock.Controller) s3Client {
				m := NewMocks3Client(ctrl)
				return m
			},
			args: args{
				ctx: context.Background(),
			},
			want: []entity.Chef{
				{
					ID:           1,
					Name:         "test 1",
					Rating:       pointers.To(float32(5)),
					ReviewsCount: pointers.To(int32(5)),
				},
				{
					ID:           2,
					Name:         "test 2",
					Rating:       pointers.To(float32(5)),
					ReviewsCount: pointers.To(int32(5)),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(tt.chefRepo(ctrl), tt.geoRepo(ctrl), tt.s3Client(ctrl))
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

func TestUsecase_GetChefAvatarURLByChefID(t *testing.T) {
	type args struct {
		ctx    context.Context
		chefID int64
	}
	tests := []struct {
		name     string
		chefRepo func(ctrl *gomock.Controller) chefRepo
		geoRepo  func(ctrl *gomock.Controller) geoRepo
		s3Client func(ctrl *gomock.Controller) s3Client
		args     args
		want     string
		wantErr  bool
	}{
		{
			name: "success",
			chefRepo: func(ctrl *gomock.Controller) chefRepo {
				m := NewMockchefRepo(ctrl)
				m.EXPECT().GetChefAvatarURLByChefID(gomock.Any(), int64(1)).Return("example.com", nil)
				return m
			},
			geoRepo: func(ctrl *gomock.Controller) geoRepo {
				m := NewMockgeoRepo(ctrl)
				return m
			},
			s3Client: func(ctrl *gomock.Controller) s3Client {
				m := NewMocks3Client(ctrl)
				return m
			},
			want: "example.com",
			args: args{
				ctx:    context.Background(),
				chefID: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(tt.chefRepo(ctrl), tt.geoRepo(ctrl), tt.s3Client(ctrl))
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
	type args struct {
		ctx    context.Context
		dishID int64
	}
	tests := []struct {
		name     string
		chefRepo func(ctrl *gomock.Controller) chefRepo
		geoRepo  func(ctrl *gomock.Controller) geoRepo
		s3Client func(ctrl *gomock.Controller) s3Client
		args     args
		want     string
		wantErr  bool
	}{
		{
			name: "success",
			chefRepo: func(ctrl *gomock.Controller) chefRepo {
				m := NewMockchefRepo(ctrl)
				m.EXPECT().GetChefAvatarURLByDishID(gomock.Any(), int64(1)).Return("example.com", nil)
				return m
			},
			geoRepo: func(ctrl *gomock.Controller) geoRepo {
				m := NewMockgeoRepo(ctrl)
				return m
			},
			s3Client: func(ctrl *gomock.Controller) s3Client {
				m := NewMocks3Client(ctrl)
				return m
			},
			want: "example.com",
			args: args{
				ctx:    context.Background(),
				dishID: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(tt.chefRepo(ctrl), tt.geoRepo(ctrl), tt.s3Client(ctrl))
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
	type args struct {
		ctx    context.Context
		dishID int64
	}
	tests := []struct {
		name     string
		chefRepo func(ctrl *gomock.Controller) chefRepo
		geoRepo  func(ctrl *gomock.Controller) geoRepo
		s3Client func(ctrl *gomock.Controller) s3Client
		args     args
		want     *entity.Chef
		wantErr  bool
	}{
		{
			name: "success",
			chefRepo: func(ctrl *gomock.Controller) chefRepo {
				m := NewMockchefRepo(ctrl)
				m.EXPECT().GetChefByDishID(gomock.Any(), int64(1)).Return(&entity.Chef{
					ID:   2,
					Name: "Chef",
				}, nil)
				m.EXPECT().GetChefRatingByChefID(gomock.Any(), int64(2)).Return(&entity.Chef{Rating: pointers.To(float32(2)), ReviewsCount: pointers.To(int32(5))}, nil)
				return m
			},
			geoRepo: func(ctrl *gomock.Controller) geoRepo {
				m := NewMockgeoRepo(ctrl)
				return m
			},
			s3Client: func(ctrl *gomock.Controller) s3Client {
				m := NewMocks3Client(ctrl)
				return m
			},
			args: args{
				ctx:    context.Background(),
				dishID: 1,
			},
			want: &entity.Chef{
				ID:           2,
				Name:         "Chef",
				Rating:       pointers.To(float32(2)),
				ReviewsCount: pointers.To(int32(5)),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(tt.chefRepo(ctrl), tt.geoRepo(ctrl), tt.s3Client(ctrl))
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
	type args struct {
		ctx    context.Context
		chefID int64
	}
	tests := []struct {
		name     string
		chefRepo func(ctrl *gomock.Controller) chefRepo
		geoRepo  func(ctrl *gomock.Controller) geoRepo
		s3Client func(ctrl *gomock.Controller) s3Client
		args     args
		want     *entity.Chef
		wantErr  bool
	}{
		{
			name: "success",
			chefRepo: func(ctrl *gomock.Controller) chefRepo {
				m := NewMockchefRepo(ctrl)
				m.EXPECT().GetChefByID(gomock.Any(), int64(1)).Return(&entity.Chef{
					ID:   1,
					Name: "Chef",
				}, nil)
				m.EXPECT().GetChefRatingByChefID(gomock.Any(), int64(1)).Return(&entity.Chef{
					Rating:       pointers.To(float32(1)),
					ReviewsCount: pointers.To(int32(1)),
				}, nil)
				return m
			},
			geoRepo: func(ctrl *gomock.Controller) geoRepo {
				m := NewMockgeoRepo(ctrl)
				return m
			},
			s3Client: func(ctrl *gomock.Controller) s3Client {
				m := NewMocks3Client(ctrl)
				return m
			},
			args: args{
				ctx:    context.Background(),
				chefID: 1,
			},
			want: &entity.Chef{
				ID:           1,
				Name:         "Chef",
				Rating:       pointers.To(float32(1)),
				ReviewsCount: pointers.To(int32(1)),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(tt.chefRepo(ctrl), tt.geoRepo(ctrl), tt.s3Client(ctrl))
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
	type args struct {
		ctx    context.Context
		chefID int64
	}
	tests := []struct {
		name     string
		chefRepo func(ctrl *gomock.Controller) chefRepo
		geoRepo  func(ctrl *gomock.Controller) geoRepo
		s3Client func(ctrl *gomock.Controller) s3Client
		args     args
		want     []entity.Certification
		wantErr  bool
	}{
		{
			name: "success",
			chefRepo: func(ctrl *gomock.Controller) chefRepo {
				m := NewMockchefRepo(ctrl)
				m.EXPECT().GetChefCertifications(gomock.Any(), int64(1)).Return([]entity.Certification{
					{
						ID:   1,
						Name: "бабушка",
					},
					{
						ID:   2,
						Name: "дедушка",
					},
				}, nil)
				return m
			},
			geoRepo: func(ctrl *gomock.Controller) geoRepo {
				m := NewMockgeoRepo(ctrl)
				return m
			},
			s3Client: func(ctrl *gomock.Controller) s3Client {
				m := NewMocks3Client(ctrl)
				return m
			},
			args: args{
				ctx:    context.Background(),
				chefID: 1,
			},
			want: []entity.Certification{
				{
					ID:   1,
					Name: "бабушка",
				},
				{
					ID:   2,
					Name: "дедушка",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(tt.chefRepo(ctrl), tt.geoRepo(ctrl), tt.s3Client(ctrl))
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
	type args struct {
		ctx    context.Context
		chefID int64
	}
	tests := []struct {
		name     string
		chefRepo func(ctrl *gomock.Controller) chefRepo
		geoRepo  func(ctrl *gomock.Controller) geoRepo
		s3Client func(ctrl *gomock.Controller) s3Client
		args     args
		want     int
		wantErr  bool
	}{
		{
			name: "success",
			chefRepo: func(ctrl *gomock.Controller) chefRepo {
				m := NewMockchefRepo(ctrl)
				m.EXPECT().GetChefExperienceYears(gomock.Any(), int64(1)).Return(2, nil)
				return m
			},
			geoRepo: func(ctrl *gomock.Controller) geoRepo {
				m := NewMockgeoRepo(ctrl)
				return m
			},
			s3Client: func(ctrl *gomock.Controller) s3Client {
				m := NewMocks3Client(ctrl)
				return m
			},
			args: args{
				ctx:    context.Background(),
				chefID: 1,
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(tt.chefRepo(ctrl), tt.geoRepo(ctrl), tt.s3Client(ctrl))
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
	type args struct {
		ctx    context.Context
		chefID int64
	}
	tests := []struct {
		name     string
		chefRepo func(ctrl *gomock.Controller) chefRepo
		geoRepo  func(ctrl *gomock.Controller) geoRepo
		s3Client func(ctrl *gomock.Controller) s3Client
		args     args
		want     []dish.Dish
		wantErr  bool
	}{
		{
			name: "success",
			chefRepo: func(ctrl *gomock.Controller) chefRepo {
				m := NewMockchefRepo(ctrl)
				m.EXPECT().GetDishesByChefID(gomock.Any(), int64(1)).Return([]dish.Dish{
					{
						ID:   1,
						Name: "asd",
					},
				}, nil)
				return m
			},
			geoRepo: func(ctrl *gomock.Controller) geoRepo {
				m := NewMockgeoRepo(ctrl)
				return m
			},
			s3Client: func(ctrl *gomock.Controller) s3Client {
				m := NewMocks3Client(ctrl)
				return m
			},
			args: args{
				ctx:    context.Background(),
				chefID: 1,
			},
			want: []dish.Dish{
				{
					ID:   1,
					Name: "asd",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(tt.chefRepo(ctrl), tt.geoRepo(ctrl), tt.s3Client(ctrl))
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

func TestUsecase_GetDistanceToChef(t *testing.T) {
	type args struct {
		ctx  context.Context
		lat  float64
		long float64
		id   int64
	}
	tests := []struct {
		name     string
		chefRepo func(ctrl *gomock.Controller) chefRepo
		geoRepo  func(ctrl *gomock.Controller) geoRepo
		s3Client func(ctrl *gomock.Controller) s3Client
		args     args
		want     float64
		wantErr  bool
	}{
		{
			name: "success",
			chefRepo: func(ctrl *gomock.Controller) chefRepo {
				m := NewMockchefRepo(ctrl)
				return m
			},
			geoRepo: func(ctrl *gomock.Controller) geoRepo {
				m := NewMockgeoRepo(ctrl)
				m.EXPECT().GetDistanceToChef(gomock.Any(), float64(1), float64(2), int64(1)).Return(float64(1), nil)
				return m
			},
			s3Client: func(ctrl *gomock.Controller) s3Client {
				m := NewMocks3Client(ctrl)
				return m
			},
			want: float64(1),
			args: args{
				ctx:  context.Background(),
				lat:  float64(1),
				long: float64(2),
				id:   int64(1),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(tt.chefRepo(ctrl), tt.geoRepo(ctrl), tt.s3Client(ctrl))
			got, err := u.GetDistanceToChef(tt.args.ctx, tt.args.lat, tt.args.long, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDistanceToChef() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetDistanceToChef() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsecase_GetNearestChefs(t *testing.T) {
	type args struct {
		ctx      context.Context
		lat      float64
		long     float64
		distance int
		limit    int
	}
	tests := []struct {
		name     string
		chefRepo func(ctrl *gomock.Controller) chefRepo
		geoRepo  func(ctrl *gomock.Controller) geoRepo
		s3Client func(ctrl *gomock.Controller) s3Client
		args     args
		want     []entity.Chef
		wantErr  bool
	}{
		{
			name: "success",
			chefRepo: func(ctrl *gomock.Controller) chefRepo {
				m := NewMockchefRepo(ctrl)
				m.EXPECT().GetNearestChefs(gomock.Any(), float64(1), float64(2), 1, 2).Return([]entity.Chef{
					{
						ID:   1,
						Name: "asd",
					},
					{
						ID:   2,
						Name: "asd",
					},
				}, nil)
				m.EXPECT().GetChefRatingByChefID(gomock.Any(), int64(1)).Return(&entity.Chef{
					Rating:       pointers.To(float32(1)),
					ReviewsCount: pointers.To(int32(1)),
				}, nil)
				m.EXPECT().GetChefRatingByChefID(gomock.Any(), int64(2)).Return(&entity.Chef{
					Rating:       pointers.To(float32(2)),
					ReviewsCount: pointers.To(int32(2)),
				}, nil)
				return m
			},
			geoRepo: func(ctrl *gomock.Controller) geoRepo {
				m := NewMockgeoRepo(ctrl)
				return m
			},
			s3Client: func(ctrl *gomock.Controller) s3Client {
				m := NewMocks3Client(ctrl)
				return m
			},
			want: []entity.Chef{
				{
					ID:           1,
					Name:         "asd",
					Rating:       pointers.To(float32(1)),
					ReviewsCount: pointers.To(int32(1)),
				},
				{
					ID:           2,
					Name:         "asd",
					Rating:       pointers.To(float32(2)),
					ReviewsCount: pointers.To(int32(2)),
				},
			},
			args: args{
				ctx:      context.Background(),
				lat:      float64(1),
				long:     float64(2),
				limit:    2,
				distance: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(tt.chefRepo(ctrl), tt.geoRepo(ctrl), tt.s3Client(ctrl))
			got, err := u.GetNearestChefs(tt.args.ctx, tt.args.lat, tt.args.long, tt.args.distance, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNearestChefs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetNearestChefs() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsecase_GetTopChefs(t *testing.T) {
	type args struct {
		ctx   context.Context
		limit int
	}
	tests := []struct {
		name     string
		chefRepo func(ctrl *gomock.Controller) chefRepo
		geoRepo  func(ctrl *gomock.Controller) geoRepo
		s3Client func(ctrl *gomock.Controller) s3Client
		args     args
		want     []entity.Chef
		wantErr  bool
	}{
		{
			name: "success",
			chefRepo: func(ctrl *gomock.Controller) chefRepo {
				m := NewMockchefRepo(ctrl)
				m.EXPECT().GetTopChefs(gomock.Any(), 2).Return([]entity.Chef{
					{
						ID:   1,
						Name: "asd",
					},
					{
						ID:   2,
						Name: "asd",
					},
				}, nil)
				return m
			},
			geoRepo: func(ctrl *gomock.Controller) geoRepo {
				m := NewMockgeoRepo(ctrl)
				return m
			},
			s3Client: func(ctrl *gomock.Controller) s3Client {
				m := NewMocks3Client(ctrl)
				return m
			},
			args: args{
				ctx:   context.Background(),
				limit: 2,
			},
			want: []entity.Chef{
				{
					ID:   1,
					Name: "asd",
				},
				{
					ID:   2,
					Name: "asd",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(tt.chefRepo(ctrl), tt.geoRepo(ctrl), tt.s3Client(ctrl))
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
	type args struct {
		ctx        context.Context
		chefID     int64
		fileHeader *multipart.FileHeader
	}
	tests := []struct {
		name     string
		chefRepo func(ctrl *gomock.Controller) chefRepo
		geoRepo  func(ctrl *gomock.Controller) geoRepo
		s3Client func(ctrl *gomock.Controller) s3Client
		args     args
		want     string
		wantErr  bool
	}{
		{
			name: "success",
			chefRepo: func(ctrl *gomock.Controller) chefRepo {
				m := NewMockchefRepo(ctrl)
				m.EXPECT().SaveChefAvatar(gomock.Any(), int64(1), "url").Return(nil)
				return m
			},
			geoRepo: func(ctrl *gomock.Controller) geoRepo {
				m := NewMockgeoRepo(ctrl)
				return m
			},
			s3Client: func(ctrl *gomock.Controller) s3Client {
				m := NewMocks3Client(ctrl)
				m.EXPECT().UploadPicture(gomock.Any(), "avatar/m/chef_1", gomock.Any()).Return("url", nil)
				return m
			},
			args: args{
				ctx:        context.Background(),
				chefID:     1,
				fileHeader: nil,
			},
			want: "url",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(tt.chefRepo(ctrl), tt.geoRepo(ctrl), tt.s3Client(ctrl))
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

func TestUsecase_UploadSmallAvatar(t *testing.T) {
	type args struct {
		ctx        context.Context
		chefID     int64
		fileHeader *multipart.FileHeader
	}
	tests := []struct {
		name     string
		chefRepo func(ctrl *gomock.Controller) chefRepo
		geoRepo  func(ctrl *gomock.Controller) geoRepo
		s3Client func(ctrl *gomock.Controller) s3Client
		args     args
		want     string
		wantErr  bool
	}{
		{
			name: "success",
			chefRepo: func(ctrl *gomock.Controller) chefRepo {
				m := NewMockchefRepo(ctrl)
				m.EXPECT().SetSmallAvatar(gomock.Any(), int64(1), "url").Return(nil)
				return m
			},
			geoRepo: func(ctrl *gomock.Controller) geoRepo {
				m := NewMockgeoRepo(ctrl)
				return m
			},
			s3Client: func(ctrl *gomock.Controller) s3Client {
				m := NewMocks3Client(ctrl)
				m.EXPECT().UploadPicture(gomock.Any(), "avatar/s/chef_1", gomock.Any()).Return("url", nil)
				return m
			},
			args: args{
				ctx:        context.Background(),
				chefID:     1,
				fileHeader: nil,
			},
			want: "url",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(tt.chefRepo(ctrl), tt.geoRepo(ctrl), tt.s3Client(ctrl))
			got, err := u.UploadSmallAvatar(tt.args.ctx, tt.args.chefID, tt.args.fileHeader)
			if (err != nil) != tt.wantErr {
				t.Errorf("UploadSmallAvatar() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UploadSmallAvatar() got = %v, want %v", got, tt.want)
			}
		})
	}
}
