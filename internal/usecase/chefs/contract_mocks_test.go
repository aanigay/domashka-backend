// Code generated by MockGen. DO NOT EDIT.
// Source: contract.go

// Package chefs is a generated GoMock package.
package chefs

import (
	context "context"
	chefs "domashka-backend/internal/entity/chefs"
	dishes "domashka-backend/internal/entity/dishes"
	multipart "mime/multipart"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockchefRepo is a mock of chefRepo interface.
type MockchefRepo struct {
	ctrl     *gomock.Controller
	recorder *MockchefRepoMockRecorder
}

// MockchefRepoMockRecorder is the mock recorder for MockchefRepo.
type MockchefRepoMockRecorder struct {
	mock *MockchefRepo
}

// NewMockchefRepo creates a new mock instance.
func NewMockchefRepo(ctrl *gomock.Controller) *MockchefRepo {
	mock := &MockchefRepo{ctrl: ctrl}
	mock.recorder = &MockchefRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockchefRepo) EXPECT() *MockchefRepoMockRecorder {
	return m.recorder
}

// GetAll mocks base method.
func (m *MockchefRepo) GetAll(ctx context.Context) ([]chefs.Chef, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", ctx)
	ret0, _ := ret[0].([]chefs.Chef)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockchefRepoMockRecorder) GetAll(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockchefRepo)(nil).GetAll), ctx)
}

// GetChefAvatarURLByChefID mocks base method.
func (m *MockchefRepo) GetChefAvatarURLByChefID(ctx context.Context, chefID int64) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetChefAvatarURLByChefID", ctx, chefID)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetChefAvatarURLByChefID indicates an expected call of GetChefAvatarURLByChefID.
func (mr *MockchefRepoMockRecorder) GetChefAvatarURLByChefID(ctx, chefID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetChefAvatarURLByChefID", reflect.TypeOf((*MockchefRepo)(nil).GetChefAvatarURLByChefID), ctx, chefID)
}

// GetChefAvatarURLByDishID mocks base method.
func (m *MockchefRepo) GetChefAvatarURLByDishID(ctx context.Context, dishID int64) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetChefAvatarURLByDishID", ctx, dishID)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetChefAvatarURLByDishID indicates an expected call of GetChefAvatarURLByDishID.
func (mr *MockchefRepoMockRecorder) GetChefAvatarURLByDishID(ctx, dishID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetChefAvatarURLByDishID", reflect.TypeOf((*MockchefRepo)(nil).GetChefAvatarURLByDishID), ctx, dishID)
}

// GetChefByDishID mocks base method.
func (m *MockchefRepo) GetChefByDishID(ctx context.Context, dishID int64) (*chefs.Chef, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetChefByDishID", ctx, dishID)
	ret0, _ := ret[0].(*chefs.Chef)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetChefByDishID indicates an expected call of GetChefByDishID.
func (mr *MockchefRepoMockRecorder) GetChefByDishID(ctx, dishID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetChefByDishID", reflect.TypeOf((*MockchefRepo)(nil).GetChefByDishID), ctx, dishID)
}

// GetChefByID mocks base method.
func (m *MockchefRepo) GetChefByID(ctx context.Context, chefID int64) (*chefs.Chef, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetChefByID", ctx, chefID)
	ret0, _ := ret[0].(*chefs.Chef)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetChefByID indicates an expected call of GetChefByID.
func (mr *MockchefRepoMockRecorder) GetChefByID(ctx, chefID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetChefByID", reflect.TypeOf((*MockchefRepo)(nil).GetChefByID), ctx, chefID)
}

// GetChefCertifications mocks base method.
func (m *MockchefRepo) GetChefCertifications(ctx context.Context, chefID int64) ([]chefs.Certification, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetChefCertifications", ctx, chefID)
	ret0, _ := ret[0].([]chefs.Certification)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetChefCertifications indicates an expected call of GetChefCertifications.
func (mr *MockchefRepoMockRecorder) GetChefCertifications(ctx, chefID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetChefCertifications", reflect.TypeOf((*MockchefRepo)(nil).GetChefCertifications), ctx, chefID)
}

// GetChefExperienceYears mocks base method.
func (m *MockchefRepo) GetChefExperienceYears(ctx context.Context, chefID int64) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetChefExperienceYears", ctx, chefID)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetChefExperienceYears indicates an expected call of GetChefExperienceYears.
func (mr *MockchefRepoMockRecorder) GetChefExperienceYears(ctx, chefID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetChefExperienceYears", reflect.TypeOf((*MockchefRepo)(nil).GetChefExperienceYears), ctx, chefID)
}

// GetChefRatingByChefID mocks base method.
func (m *MockchefRepo) GetChefRatingByChefID(ctx context.Context, chefID int64) (*chefs.Chef, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetChefRatingByChefID", ctx, chefID)
	ret0, _ := ret[0].(*chefs.Chef)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetChefRatingByChefID indicates an expected call of GetChefRatingByChefID.
func (mr *MockchefRepoMockRecorder) GetChefRatingByChefID(ctx, chefID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetChefRatingByChefID", reflect.TypeOf((*MockchefRepo)(nil).GetChefRatingByChefID), ctx, chefID)
}

// GetDishesByChefID mocks base method.
func (m *MockchefRepo) GetDishesByChefID(ctx context.Context, chefID int64) ([]dishes.Dish, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDishesByChefID", ctx, chefID)
	ret0, _ := ret[0].([]dishes.Dish)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDishesByChefID indicates an expected call of GetDishesByChefID.
func (mr *MockchefRepoMockRecorder) GetDishesByChefID(ctx, chefID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDishesByChefID", reflect.TypeOf((*MockchefRepo)(nil).GetDishesByChefID), ctx, chefID)
}

// GetNearestChefs mocks base method.
func (m *MockchefRepo) GetNearestChefs(ctx context.Context, lat, long float64, distance, limit int) ([]chefs.Chef, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNearestChefs", ctx, lat, long, distance, limit)
	ret0, _ := ret[0].([]chefs.Chef)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNearestChefs indicates an expected call of GetNearestChefs.
func (mr *MockchefRepoMockRecorder) GetNearestChefs(ctx, lat, long, distance, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNearestChefs", reflect.TypeOf((*MockchefRepo)(nil).GetNearestChefs), ctx, lat, long, distance, limit)
}

// GetTopChefs mocks base method.
func (m *MockchefRepo) GetTopChefs(ctx context.Context, limit int) ([]chefs.Chef, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTopChefs", ctx, limit)
	ret0, _ := ret[0].([]chefs.Chef)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTopChefs indicates an expected call of GetTopChefs.
func (mr *MockchefRepoMockRecorder) GetTopChefs(ctx, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTopChefs", reflect.TypeOf((*MockchefRepo)(nil).GetTopChefs), ctx, limit)
}

// SaveChefAvatar mocks base method.
func (m *MockchefRepo) SaveChefAvatar(ctx context.Context, chefID int64, publicURL string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveChefAvatar", ctx, chefID, publicURL)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveChefAvatar indicates an expected call of SaveChefAvatar.
func (mr *MockchefRepoMockRecorder) SaveChefAvatar(ctx, chefID, publicURL interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveChefAvatar", reflect.TypeOf((*MockchefRepo)(nil).SaveChefAvatar), ctx, chefID, publicURL)
}

// SetSmallAvatar mocks base method.
func (m *MockchefRepo) SetSmallAvatar(ctx context.Context, chefID int64, publicURL string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetSmallAvatar", ctx, chefID, publicURL)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetSmallAvatar indicates an expected call of SetSmallAvatar.
func (mr *MockchefRepoMockRecorder) SetSmallAvatar(ctx, chefID, publicURL interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetSmallAvatar", reflect.TypeOf((*MockchefRepo)(nil).SetSmallAvatar), ctx, chefID, publicURL)
}

// MockgeoRepo is a mock of geoRepo interface.
type MockgeoRepo struct {
	ctrl     *gomock.Controller
	recorder *MockgeoRepoMockRecorder
}

// MockgeoRepoMockRecorder is the mock recorder for MockgeoRepo.
type MockgeoRepoMockRecorder struct {
	mock *MockgeoRepo
}

// NewMockgeoRepo creates a new mock instance.
func NewMockgeoRepo(ctrl *gomock.Controller) *MockgeoRepo {
	mock := &MockgeoRepo{ctrl: ctrl}
	mock.recorder = &MockgeoRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockgeoRepo) EXPECT() *MockgeoRepoMockRecorder {
	return m.recorder
}

// GetDistanceToChef mocks base method.
func (m *MockgeoRepo) GetDistanceToChef(ctx context.Context, lat, long float64, chefID int64) (float64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDistanceToChef", ctx, lat, long, chefID)
	ret0, _ := ret[0].(float64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDistanceToChef indicates an expected call of GetDistanceToChef.
func (mr *MockgeoRepoMockRecorder) GetDistanceToChef(ctx, lat, long, chefID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDistanceToChef", reflect.TypeOf((*MockgeoRepo)(nil).GetDistanceToChef), ctx, lat, long, chefID)
}

// Mocks3Client is a mock of s3Client interface.
type Mocks3Client struct {
	ctrl     *gomock.Controller
	recorder *Mocks3ClientMockRecorder
}

// Mocks3ClientMockRecorder is the mock recorder for Mocks3Client.
type Mocks3ClientMockRecorder struct {
	mock *Mocks3Client
}

// NewMocks3Client creates a new mock instance.
func NewMocks3Client(ctrl *gomock.Controller) *Mocks3Client {
	mock := &Mocks3Client{ctrl: ctrl}
	mock.recorder = &Mocks3ClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mocks3Client) EXPECT() *Mocks3ClientMockRecorder {
	return m.recorder
}

// UploadPicture mocks base method.
func (m *Mocks3Client) UploadPicture(ctx context.Context, filePrefix string, fileHeader *multipart.FileHeader) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UploadPicture", ctx, filePrefix, fileHeader)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UploadPicture indicates an expected call of UploadPicture.
func (mr *Mocks3ClientMockRecorder) UploadPicture(ctx, filePrefix, fileHeader interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadPicture", reflect.TypeOf((*Mocks3Client)(nil).UploadPicture), ctx, filePrefix, fileHeader)
}
