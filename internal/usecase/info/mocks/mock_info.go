// Code generated by MockGen. DO NOT EDIT.
// Source: info/info.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	domain "github.com/Azreil-OFD/Avito-test/internal/infrastructure/domain"
	gomock "github.com/golang/mock/gomock"
)

// MockInfoRepository is a mock of InfoRepository interface.
type MockInfoRepository struct {
	ctrl     *gomock.Controller
	recorder *MockInfoRepositoryMockRecorder
}

// MockInfoRepositoryMockRecorder is the mock recorder for MockInfoRepository.
type MockInfoRepositoryMockRecorder struct {
	mock *MockInfoRepository
}

// NewMockInfoRepository creates a new mock instance.
func NewMockInfoRepository(ctrl *gomock.Controller) *MockInfoRepository {
	mock := &MockInfoRepository{ctrl: ctrl}
	mock.recorder = &MockInfoRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockInfoRepository) EXPECT() *MockInfoRepositoryMockRecorder {
	return m.recorder
}

// GetMerchByUserId mocks base method.
func (m *MockInfoRepository) GetMerchByUserId(ctx context.Context, userId int) ([]domain.MerchStore, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMerchByUserId", ctx, userId)
	ret0, _ := ret[0].([]domain.MerchStore)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMerchByUserId indicates an expected call of GetMerchByUserId.
func (mr *MockInfoRepositoryMockRecorder) GetMerchByUserId(ctx, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMerchByUserId", reflect.TypeOf((*MockInfoRepository)(nil).GetMerchByUserId), ctx, userId)
}

// GetTransactionsInfoById mocks base method.
func (m *MockInfoRepository) GetTransactionsInfoById(ctx context.Context, id int) ([]domain.TransactionFormat, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTransactionsInfoById", ctx, id)
	ret0, _ := ret[0].([]domain.TransactionFormat)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTransactionsInfoById indicates an expected call of GetTransactionsInfoById.
func (mr *MockInfoRepositoryMockRecorder) GetTransactionsInfoById(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTransactionsInfoById", reflect.TypeOf((*MockInfoRepository)(nil).GetTransactionsInfoById), ctx, id)
}

// GetUserById mocks base method.
func (m *MockInfoRepository) GetUserById(ctx context.Context, id int) (domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserById", ctx, id)
	ret0, _ := ret[0].(domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserById indicates an expected call of GetUserById.
func (mr *MockInfoRepositoryMockRecorder) GetUserById(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserById", reflect.TypeOf((*MockInfoRepository)(nil).GetUserById), ctx, id)
}
