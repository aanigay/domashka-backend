package geo

import (
	"context"
	"domashka-backend/internal/entity/geo"
	"github.com/golang/mock/gomock"
	"reflect"
	"testing"
)

func TestUseCase_AddChefAddress(t *testing.T) {
	type args struct {
		ctx     context.Context
		chefID  int
		address geo.Address
	}
	tests := []struct {
		name    string
		geoRepo func(ctrl *gomock.Controller) *MockGeoRepository
		args    args
		wantErr bool
	}{
		{
			name: "success",
			geoRepo: func(ctrl *gomock.Controller) *MockGeoRepository {
				m := NewMockGeoRepository(ctrl)
				m.EXPECT().AddChefAddress(gomock.Any(), 1, geo.Address{
					ID:        1,
					Longitude: 37,
					Latitude:  55,
				}).Return(nil)
				return m
			},
			args: args{
				ctx:    context.Background(),
				chefID: 1,
				address: geo.Address{
					ID:        1,
					Longitude: 37,
					Latitude:  55,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			s := New(tt.geoRepo(ctrl))
			if err := s.AddChefAddress(tt.args.ctx, tt.args.chefID, tt.args.address); (err != nil) != tt.wantErr {
				t.Errorf("AddChefAddress() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUseCase_AddClientAddress(t *testing.T) {
	type args struct {
		ctx      context.Context
		clientID int
		address  geo.Address
	}
	tests := []struct {
		name    string
		geoRepo func(ctrl *gomock.Controller) *MockGeoRepository
		args    args
		wantErr bool
	}{
		{
			name: "success",
			geoRepo: func(ctrl *gomock.Controller) *MockGeoRepository {
				m := NewMockGeoRepository(ctrl)
				m.EXPECT().AddClientAddress(gomock.Any(), 1, geo.Address{
					ID:        1,
					Longitude: 37,
					Latitude:  55,
				}).Return(nil)
				return m
			},
			args: args{
				ctx:      context.Background(),
				clientID: 1,
				address: geo.Address{
					ID:        1,
					Longitude: 37,
					Latitude:  55,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			s := New(tt.geoRepo(ctrl))
			if err := s.AddClientAddress(tt.args.ctx, tt.args.clientID, tt.args.address); (err != nil) != tt.wantErr {
				t.Errorf("AddClientAddress() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUseCase_FindChefsNearAddress(t *testing.T) {
	type args struct {
		ctx             context.Context
		clientAddressID int
		radius          float64
	}
	tests := []struct {
		name    string
		geoRepo func(ctrl *gomock.Controller) *MockGeoRepository
		args    args
		want    []geo.Address
		wantErr bool
	}{
		{
			name: "success",
			geoRepo: func(ctrl *gomock.Controller) *MockGeoRepository {
				m := NewMockGeoRepository(ctrl)
				m.EXPECT().GetChefsAddrByRange(gomock.Any(), 1, float64(100)).Return([]geo.Address{}, nil)
				return m
			},
			args: args{
				ctx:             context.Background(),
				clientAddressID: 1,
				radius:          100,
			},
			want: []geo.Address{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			s := New(tt.geoRepo(ctrl))
			got, err := s.FindChefsNearAddress(tt.args.ctx, tt.args.clientAddressID, tt.args.radius)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindChefsNearAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindChefsNearAddress() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCase_FindClientsNearAddress(t *testing.T) {
	type args struct {
		ctx    context.Context
		chefID int
		radius float64
	}
	tests := []struct {
		name    string
		geoRepo func(ctrl *gomock.Controller) *MockGeoRepository
		args    args
		want    []geo.Address
		wantErr bool
	}{
		{
			name: "success",
			geoRepo: func(ctrl *gomock.Controller) *MockGeoRepository {
				m := NewMockGeoRepository(ctrl)
				m.EXPECT().GetClientsAddrByRange(gomock.Any(), 1, float64(100)).Return([]geo.Address{}, nil)
				return m
			},
			args: args{
				ctx:    context.Background(),
				chefID: 1,
				radius: 100,
			},
			want: []geo.Address{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			s := New(tt.geoRepo(ctrl))
			got, err := s.FindClientsNearAddress(tt.args.ctx, tt.args.chefID, tt.args.radius)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindClientsNearAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindClientsNearAddress() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCase_GetAddressByID(t *testing.T) {
	type args struct {
		ctx context.Context
		id  int64
	}
	tests := []struct {
		name    string
		geoRepo func(ctrl *gomock.Controller) *MockGeoRepository
		args    args
		want    *geo.Address
		wantErr bool
	}{
		{
			name: "success",
			geoRepo: func(ctrl *gomock.Controller) *MockGeoRepository {
				m := NewMockGeoRepository(ctrl)
				m.EXPECT().GetAddressByID(gomock.Any(), gomock.Any()).Return(&geo.Address{}, nil)
				return m
			},
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			want: &geo.Address{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(tt.geoRepo(ctrl))
			got, err := u.GetAddressByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAddressByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAddressByID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCase_GetChefAddress(t *testing.T) {
	type args struct {
		ctx    context.Context
		chefID int
	}
	tests := []struct {
		name    string
		geoRepo func(ctrl *gomock.Controller) *MockGeoRepository
		args    args
		want    geo.Address
		wantErr bool
	}{
		{
			name: "success",
			geoRepo: func(ctrl *gomock.Controller) *MockGeoRepository {
				m := NewMockGeoRepository(ctrl)
				m.EXPECT().GetChefAddress(gomock.Any(), gomock.Any()).Return(geo.Address{}, nil)
				return m
			},
			args: args{
				ctx:    context.Background(),
				chefID: 1,
			},
			want:    geo.Address{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			s := New(tt.geoRepo(ctrl))
			got, err := s.GetChefAddress(tt.args.ctx, tt.args.chefID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetChefAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetChefAddress() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCase_GetClientAddresses(t *testing.T) {
	type args struct {
		ctx      context.Context
		clientID int
	}
	tests := []struct {
		name    string
		geoRepo func(ctrl *gomock.Controller) *MockGeoRepository
		args    args
		want    []geo.Address
		wantErr bool
	}{
		{
			name: "success",
			geoRepo: func(ctrl *gomock.Controller) *MockGeoRepository {
				m := NewMockGeoRepository(ctrl)
				m.EXPECT().GetClientAddresses(gomock.Any(), gomock.Any()).Return([]geo.Address{}, nil)
				return m
			},
			args: args{
				ctx:      context.Background(),
				clientID: 1,
			},
			want: []geo.Address{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			s := New(tt.geoRepo(ctrl))
			got, err := s.GetClientAddresses(tt.args.ctx, tt.args.clientID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetClientAddresses() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetClientAddresses() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCase_GetLastUpdatedClientAddress(t *testing.T) {
	type args struct {
		ctx      context.Context
		clientID int64
	}
	tests := []struct {
		name    string
		geoRepo func(ctrl *gomock.Controller) *MockGeoRepository
		args    args
		want    *geo.Address
		wantErr bool
	}{
		{
			name: "success",
			geoRepo: func(ctrl *gomock.Controller) *MockGeoRepository {
				m := NewMockGeoRepository(ctrl)
				m.EXPECT().GetLastUpdatedClientAddress(gomock.Any(), gomock.Any()).Return(&geo.Address{}, nil)
				return m
			},
			args: args{
				ctx:      context.Background(),
				clientID: 1,
			},
			want: &geo.Address{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			s := New(tt.geoRepo(ctrl))
			got, err := s.GetLastUpdatedClientAddress(tt.args.ctx, tt.args.clientID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLastUpdatedClientAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetLastUpdatedClientAddress() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCase_PushClientAddress(t *testing.T) {
	type args struct {
		ctx       context.Context
		addressID int64
	}
	tests := []struct {
		name    string
		geoRepo func(ctrl *gomock.Controller) *MockGeoRepository
		args    args
		wantErr bool
	}{
		{
			name: "success",
			geoRepo: func(ctrl *gomock.Controller) *MockGeoRepository {
				m := NewMockGeoRepository(ctrl)
				m.EXPECT().PushClientAddress(gomock.Any(), gomock.Any()).Return(nil)
				return m
			},
			args: args{
				ctx:       context.Background(),
				addressID: 1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(tt.geoRepo(ctrl))
			if err := u.PushClientAddress(tt.args.ctx, tt.args.addressID); (err != nil) != tt.wantErr {
				t.Errorf("PushClientAddress() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUseCase_UpdateChefAddress(t *testing.T) {
	type args struct {
		ctx     context.Context
		chefID  int
		address geo.Address
	}
	tests := []struct {
		name    string
		geoRepo func(ctrl *gomock.Controller) *MockGeoRepository
		args    args
		wantErr bool
	}{
		{
			name: "success",
			geoRepo: func(ctrl *gomock.Controller) *MockGeoRepository {
				m := NewMockGeoRepository(ctrl)
				m.EXPECT().UpdateChefAddress(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return m
			},
			args: args{
				ctx:    context.Background(),
				chefID: 1,
				address: geo.Address{
					Longitude: 37,
					Latitude:  55,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			s := New(tt.geoRepo(ctrl))
			if err := s.UpdateChefAddress(tt.args.ctx, tt.args.chefID, tt.args.address); (err != nil) != tt.wantErr {
				t.Errorf("UpdateChefAddress() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUseCase_UpdateClientAddress(t *testing.T) {
	type args struct {
		ctx       context.Context
		clientID  int
		addressID int
		address   geo.Address
	}
	tests := []struct {
		name    string
		geoRepo func(ctrl *gomock.Controller) *MockGeoRepository
		args    args
		wantErr bool
	}{
		{
			name: "success",
			geoRepo: func(ctrl *gomock.Controller) *MockGeoRepository {
				m := NewMockGeoRepository(ctrl)
				m.EXPECT().UpdateClientAddress(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return m
			},
			args: args{
				ctx:       context.Background(),
				clientID:  1,
				addressID: 1,
				address: geo.Address{
					Longitude: 37,
					Latitude:  55,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			s := New(tt.geoRepo(ctrl))
			if err := s.UpdateClientAddress(tt.args.ctx, tt.args.clientID, tt.args.addressID, tt.args.address); (err != nil) != tt.wantErr {
				t.Errorf("UpdateClientAddress() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
