// Code generated by MockGen. DO NOT EDIT.
// Source: contract.go

// Package users is a generated GoMock package.
package users

import (
	context "context"
	chefs "domashka-backend/internal/entity/chefs"
	dishes "domashka-backend/internal/entity/dishes"
	users "domashka-backend/internal/entity/users"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// Mockrepo is a mock of repo interface.
type Mockrepo struct {
	ctrl     *gomock.Controller
	recorder *MockrepoMockRecorder
}

// MockrepoMockRecorder is the mock recorder for Mockrepo.
type MockrepoMockRecorder struct {
	mock *Mockrepo
}

// NewMockrepo creates a new mock instance.
func NewMockrepo(ctrl *gomock.Controller) *Mockrepo {
	mock := &Mockrepo{ctrl: ctrl}
	mock.recorder = &MockrepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mockrepo) EXPECT() *MockrepoMockRecorder {
	return m.recorder
}

// CheckIfUserIsChef mocks base method.
func (m *Mockrepo) CheckIfUserIsChef(ctx context.Context, userID int64) (*int64, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckIfUserIsChef", ctx, userID)
	ret0, _ := ret[0].(*int64)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// CheckIfUserIsChef indicates an expected call of CheckIfUserIsChef.
func (mr *MockrepoMockRecorder) CheckIfUserIsChef(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckIfUserIsChef", reflect.TypeOf((*Mockrepo)(nil).CheckIfUserIsChef), ctx, userID)
}

// Create mocks base method.
func (m *Mockrepo) Create(ctx context.Context, user *users.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockrepoMockRecorder) Create(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*Mockrepo)(nil).Create), ctx, user)
}

// Delete mocks base method.
func (m *Mockrepo) Delete(ctx context.Context, id int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockrepoMockRecorder) Delete(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*Mockrepo)(nil).Delete), ctx, id)
}

// GetByID mocks base method.
func (m *Mockrepo) GetByID(ctx context.Context, id int64) (*users.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, id)
	ret0, _ := ret[0].(*users.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockrepoMockRecorder) GetByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*Mockrepo)(nil).GetByID), ctx, id)
}

// GetFavoritesChefsByUserID mocks base method.
func (m *Mockrepo) GetFavoritesChefsByUserID(ctx context.Context, userID int64) ([]chefs.Chef, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFavoritesChefsByUserID", ctx, userID)
	ret0, _ := ret[0].([]chefs.Chef)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFavoritesChefsByUserID indicates an expected call of GetFavoritesChefsByUserID.
func (mr *MockrepoMockRecorder) GetFavoritesChefsByUserID(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFavoritesChefsByUserID", reflect.TypeOf((*Mockrepo)(nil).GetFavoritesChefsByUserID), ctx, userID)
}

// GetFavoritesDishesByUserID mocks base method.
func (m *Mockrepo) GetFavoritesDishesByUserID(ctx context.Context, userID int64) ([]dishes.Dish, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFavoritesDishesByUserID", ctx, userID)
	ret0, _ := ret[0].([]dishes.Dish)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFavoritesDishesByUserID indicates an expected call of GetFavoritesDishesByUserID.
func (mr *MockrepoMockRecorder) GetFavoritesDishesByUserID(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFavoritesDishesByUserID", reflect.TypeOf((*Mockrepo)(nil).GetFavoritesDishesByUserID), ctx, userID)
}

// Update mocks base method.
func (m *Mockrepo) Update(ctx context.Context, id int64, user users.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, id, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockrepoMockRecorder) Update(ctx, id, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*Mockrepo)(nil).Update), ctx, id, user)
}
