// Code generated by MockGen. DO NOT EDIT.
// Source: handler.go

// Package mock_handlers is a generated GoMock package.
package mock_handlers

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	models "github.com/sebasttiano/Owl/internal/models"
)

// MockAuthenticator is a mock of Authenticator interface.
type MockAuthenticator struct {
	ctrl     *gomock.Controller
	recorder *MockAuthenticatorMockRecorder
}

// MockAuthenticatorMockRecorder is the mock recorder for MockAuthenticator.
type MockAuthenticatorMockRecorder struct {
	mock *MockAuthenticator
}

// NewMockAuthenticator creates a new mock instance.
func NewMockAuthenticator(ctrl *gomock.Controller) *MockAuthenticator {
	mock := &MockAuthenticator{ctrl: ctrl}
	mock.recorder = &MockAuthenticatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthenticator) EXPECT() *MockAuthenticatorMockRecorder {
	return m.recorder
}

// Login mocks base method.
func (m *MockAuthenticator) Login(ctx context.Context, name, password string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login", ctx, name, password)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Login indicates an expected call of Login.
func (mr *MockAuthenticatorMockRecorder) Login(ctx, name, password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockAuthenticator)(nil).Login), ctx, name, password)
}

// Register mocks base method.
func (m *MockAuthenticator) Register(ctx context.Context, name, password string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Register", ctx, name, password)
	ret0, _ := ret[0].(error)
	return ret0
}

// Register indicates an expected call of Register.
func (mr *MockAuthenticatorMockRecorder) Register(ctx, name, password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockAuthenticator)(nil).Register), ctx, name, password)
}

// MockBinaryServ is a mock of BinaryServ interface.
type MockBinaryServ struct {
	ctrl     *gomock.Controller
	recorder *MockBinaryServMockRecorder
}

// MockBinaryServMockRecorder is the mock recorder for MockBinaryServ.
type MockBinaryServMockRecorder struct {
	mock *MockBinaryServ
}

// NewMockBinaryServ creates a new mock instance.
func NewMockBinaryServ(ctrl *gomock.Controller) *MockBinaryServ {
	mock := &MockBinaryServ{ctrl: ctrl}
	mock.recorder = &MockBinaryServMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBinaryServ) EXPECT() *MockBinaryServMockRecorder {
	return m.recorder
}

// DeleteBinary mocks base method.
func (m *MockBinaryServ) DeleteBinary(ctx context.Context, id int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteBinary", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteBinary indicates an expected call of DeleteBinary.
func (mr *MockBinaryServMockRecorder) DeleteBinary(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteBinary", reflect.TypeOf((*MockBinaryServ)(nil).DeleteBinary), ctx, id)
}

// GetAllBinaries mocks base method.
func (m *MockBinaryServ) GetAllBinaries(ctx context.Context) ([]models.Resource, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllBinaries", ctx)
	ret0, _ := ret[0].([]models.Resource)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllBinaries indicates an expected call of GetAllBinaries.
func (mr *MockBinaryServMockRecorder) GetAllBinaries(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllBinaries", reflect.TypeOf((*MockBinaryServ)(nil).GetAllBinaries), ctx)
}

// GetBinary mocks base method.
func (m *MockBinaryServ) GetBinary(ctx context.Context, id int) (models.Resource, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBinary", ctx, id)
	ret0, _ := ret[0].(models.Resource)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBinary indicates an expected call of GetBinary.
func (mr *MockBinaryServMockRecorder) GetBinary(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBinary", reflect.TypeOf((*MockBinaryServ)(nil).GetBinary), ctx, id)
}

// SetBinary mocks base method.
func (m *MockBinaryServ) SetBinary(ctx context.Context, data models.Resource) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetBinary", ctx, data)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetBinary indicates an expected call of SetBinary.
func (mr *MockBinaryServMockRecorder) SetBinary(ctx, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetBinary", reflect.TypeOf((*MockBinaryServ)(nil).SetBinary), ctx, data)
}

// MockResourceServ is a mock of ResourceServ interface.
type MockResourceServ struct {
	ctrl     *gomock.Controller
	recorder *MockResourceServMockRecorder
}

// MockResourceServMockRecorder is the mock recorder for MockResourceServ.
type MockResourceServMockRecorder struct {
	mock *MockResourceServ
}

// NewMockResourceServ creates a new mock instance.
func NewMockResourceServ(ctrl *gomock.Controller) *MockResourceServ {
	mock := &MockResourceServ{ctrl: ctrl}
	mock.recorder = &MockResourceServMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockResourceServ) EXPECT() *MockResourceServMockRecorder {
	return m.recorder
}

// DeleteResource mocks base method.
func (m *MockResourceServ) DeleteResource(ctx context.Context, res *models.Resource) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteResource", ctx, res)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteResource indicates an expected call of DeleteResource.
func (mr *MockResourceServMockRecorder) DeleteResource(ctx, res interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteResource", reflect.TypeOf((*MockResourceServ)(nil).DeleteResource), ctx, res)
}

// GetAllResources mocks base method.
func (m *MockResourceServ) GetAllResources(ctx context.Context, uid int) ([]*models.Resource, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllResources", ctx, uid)
	ret0, _ := ret[0].([]*models.Resource)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllResources indicates an expected call of GetAllResources.
func (mr *MockResourceServMockRecorder) GetAllResources(ctx, uid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllResources", reflect.TypeOf((*MockResourceServ)(nil).GetAllResources), ctx, uid)
}

// GetResource mocks base method.
func (m *MockResourceServ) GetResource(ctx context.Context, res *models.Resource) (*models.Resource, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetResource", ctx, res)
	ret0, _ := ret[0].(*models.Resource)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetResource indicates an expected call of GetResource.
func (mr *MockResourceServMockRecorder) GetResource(ctx, res interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetResource", reflect.TypeOf((*MockResourceServ)(nil).GetResource), ctx, res)
}

// SetResource mocks base method.
func (m *MockResourceServ) SetResource(ctx context.Context, res models.Resource) (*models.Resource, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetResource", ctx, res)
	ret0, _ := ret[0].(*models.Resource)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SetResource indicates an expected call of SetResource.
func (mr *MockResourceServMockRecorder) SetResource(ctx, res interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetResource", reflect.TypeOf((*MockResourceServ)(nil).SetResource), ctx, res)
}
