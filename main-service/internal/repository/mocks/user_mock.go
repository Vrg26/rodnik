// Code generated by MockGen. DO NOT EDIT.
// Source: users.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	context "context"
	entity "main-service/internal/entity"
	http "net/http"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
)

// MockUsersRepo is a mock of UsersRepo interface.
type MockUsersRepo struct {
	ctrl     *gomock.Controller
	recorder *MockUsersRepoMockRecorder
}

// MockUsersRepoMockRecorder is the mock recorder for MockUsersRepo.
type MockUsersRepoMockRecorder struct {
	mock *MockUsersRepo
}

// NewMockUsersRepo creates a new mock instance.
func NewMockUsersRepo(ctrl *gomock.Controller) *MockUsersRepo {
	mock := &MockUsersRepo{ctrl: ctrl}
	mock.recorder = &MockUsersRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUsersRepo) EXPECT() *MockUsersRepoMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockUsersRepo) Create(ctx context.Context, user *entity.User) (*entity.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, user)
	ret0, _ := ret[0].(*entity.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockUsersRepoMockRecorder) Create(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockUsersRepo)(nil).Create), ctx, user)
}

// FindById mocks base method.
func (m *MockUsersRepo) FindById(ctx context.Context, userID uuid.UUID) (*entity.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindById", ctx, userID)
	ret0, _ := ret[0].(*entity.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindById indicates an expected call of FindById.
func (mr *MockUsersRepoMockRecorder) FindById(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindById", reflect.TypeOf((*MockUsersRepo)(nil).FindById), ctx, userID)
}

// FindByPhone mocks base method.
func (m *MockUsersRepo) FindByPhone(ctx context.Context, phone string) (*entity.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByPhone", ctx, phone)
	ret0, _ := ret[0].(*entity.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByPhone indicates an expected call of FindByPhone.
func (mr *MockUsersRepoMockRecorder) FindByPhone(ctx, phone interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByPhone", reflect.TypeOf((*MockUsersRepo)(nil).FindByPhone), ctx, phone)
}

// SetAvatar mocks base method.
func (m *MockUsersRepo) SetAvatar(ctx context.Context, userID string, avatarName string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetAvatar", ctx, userID, avatarName)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetAvatar indicates an expected call of SetAvatar.
func (mr *MockUsersRepoMockRecorder) SetAvatar(ctx, userID, avatarName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetAvatar", reflect.TypeOf((*MockUsersRepo)(nil).SetAvatar), ctx, userID, avatarName)
}

// UpdateUserBalance mocks base method.
func (m *MockUsersRepo) UpdateUserBalance(ctx context.Context, user *entity.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserBalance", ctx, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUserBalance indicates an expected call of UpdateUserBalance.
func (mr *MockUsersRepoMockRecorder) UpdateUserBalance(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserBalance", reflect.TypeOf((*MockUsersRepo)(nil).UpdateUserBalance), ctx, user)
}

// MockclientImageService is a mock of clientImageService interface.
type MockclientImageService struct {
	ctrl     *gomock.Controller
	recorder *MockclientImageServiceMockRecorder
}

// MockclientImageServiceMockRecorder is the mock recorder for MockclientImageService.
type MockclientImageServiceMockRecorder struct {
	mock *MockclientImageService
}

// NewMockclientImageService creates a new mock instance.
func NewMockclientImageService(ctrl *gomock.Controller) *MockclientImageService {
	mock := &MockclientImageService{ctrl: ctrl}
	mock.recorder = &MockclientImageServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockclientImageService) EXPECT() *MockclientImageServiceMockRecorder {
	return m.recorder
}

// GetURL mocks base method.
func (m *MockclientImageService) GetURL() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetURL")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetURL indicates an expected call of GetURL.
func (mr *MockclientImageServiceMockRecorder) GetURL() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetURL", reflect.TypeOf((*MockclientImageService)(nil).GetURL))
}

// Upload mocks base method.
func (m *MockclientImageService) Upload(ctx context.Context, image []byte) (*http.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Upload", ctx, image)
	ret0, _ := ret[0].(*http.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Upload indicates an expected call of Upload.
func (mr *MockclientImageServiceMockRecorder) Upload(ctx, image interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Upload", reflect.TypeOf((*MockclientImageService)(nil).Upload), ctx, image)
}
