// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package mock_service is a generated GoMock package.
package mock_service

import (
	context "context"
	reflect "reflect"
	entity "rodnik/internal/entity"
	service "rodnik/internal/service"

	gomock "github.com/golang/mock/gomock"
)

// MockToken is a mock of Token interface.
type MockToken struct {
	ctrl     *gomock.Controller
	recorder *MockTokenMockRecorder
}

// MockTokenMockRecorder is the mock recorder for MockToken.
type MockTokenMockRecorder struct {
	mock *MockToken
}

// NewMockToken creates a new mock instance.
func NewMockToken(ctrl *gomock.Controller) *MockToken {
	mock := &MockToken{ctrl: ctrl}
	mock.recorder = &MockTokenMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockToken) EXPECT() *MockTokenMockRecorder {
	return m.recorder
}

// DeleteUserTokens mocks base method.
func (m *MockToken) DeleteUserTokens(ctx context.Context, userId string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUserTokens", ctx, userId)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUserTokens indicates an expected call of DeleteUserTokens.
func (mr *MockTokenMockRecorder) DeleteUserTokens(ctx, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUserTokens", reflect.TypeOf((*MockToken)(nil).DeleteUserTokens), ctx, userId)
}

// GetTokenPair mocks base method.
func (m *MockToken) GetTokenPair(ctx context.Context, userId string) (*service.TokenPair, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTokenPair", ctx, userId)
	ret0, _ := ret[0].(*service.TokenPair)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTokenPair indicates an expected call of GetTokenPair.
func (mr *MockTokenMockRecorder) GetTokenPair(ctx, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTokenPair", reflect.TypeOf((*MockToken)(nil).GetTokenPair), ctx, userId)
}

// ParseToken mocks base method.
func (m *MockToken) ParseToken(tokenString string) (*service.CustomClaims, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ParseToken", tokenString)
	ret0, _ := ret[0].(*service.CustomClaims)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ParseToken indicates an expected call of ParseToken.
func (mr *MockTokenMockRecorder) ParseToken(tokenString interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ParseToken", reflect.TypeOf((*MockToken)(nil).ParseToken), tokenString)
}

// RefreshToken mocks base method.
func (m *MockToken) RefreshToken(ctx context.Context, refreshToken string) (*service.TokenPair, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RefreshToken", ctx, refreshToken)
	ret0, _ := ret[0].(*service.TokenPair)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RefreshToken indicates an expected call of RefreshToken.
func (mr *MockTokenMockRecorder) RefreshToken(ctx, refreshToken interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RefreshToken", reflect.TypeOf((*MockToken)(nil).RefreshToken), ctx, refreshToken)
}

// MockUsers is a mock of Users interface.
type MockUsers struct {
	ctrl     *gomock.Controller
	recorder *MockUsersMockRecorder
}

// MockUsersMockRecorder is the mock recorder for MockUsers.
type MockUsersMockRecorder struct {
	mock *MockUsers
}

// NewMockUsers creates a new mock instance.
func NewMockUsers(ctrl *gomock.Controller) *MockUsers {
	mock := &MockUsers{ctrl: ctrl}
	mock.recorder = &MockUsersMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUsers) EXPECT() *MockUsersMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockUsers) Create(ctx context.Context, user *entity.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockUsersMockRecorder) Create(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockUsers)(nil).Create), ctx, user)
}

// Login mocks base method.
func (m *MockUsers) Login(ctx context.Context, user *entity.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login", ctx, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// Login indicates an expected call of Login.
func (mr *MockUsersMockRecorder) Login(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockUsers)(nil).Login), ctx, user)
}

// MockTasks is a mock of Tasks interface.
type MockTasks struct {
	ctrl     *gomock.Controller
	recorder *MockTasksMockRecorder
}

// MockTasksMockRecorder is the mock recorder for MockTasks.
type MockTasksMockRecorder struct {
	mock *MockTasks
}

// NewMockTasks creates a new mock instance.
func NewMockTasks(ctrl *gomock.Controller) *MockTasks {
	mock := &MockTasks{ctrl: ctrl}
	mock.recorder = &MockTasksMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTasks) EXPECT() *MockTasksMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockTasks) Create(ctx context.Context, task *entity.Task) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, task)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockTasksMockRecorder) Create(ctx, task interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockTasks)(nil).Create), ctx, task)
}
