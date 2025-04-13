// Code generated by MockGen. DO NOT EDIT.
// Source: authorization.go
//
// Generated by this command:
//
//	mockgen -source=authorization.go -package=usecase -destination=../../../../test/mock/usecase/authorization.go
//

// Package usecase is a generated GoMock package.
package usecase

import (
	context "context"
	reflect "reflect"

	dto "github.com/atsumarukun/holos-storage-api/internal/app/api/usecase/dto"
	gomock "go.uber.org/mock/gomock"
)

// MockAuthorizationUsecase is a mock of AuthorizationUsecase interface.
type MockAuthorizationUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockAuthorizationUsecaseMockRecorder
	isgomock struct{}
}

// MockAuthorizationUsecaseMockRecorder is the mock recorder for MockAuthorizationUsecase.
type MockAuthorizationUsecaseMockRecorder struct {
	mock *MockAuthorizationUsecase
}

// NewMockAuthorizationUsecase creates a new mock instance.
func NewMockAuthorizationUsecase(ctrl *gomock.Controller) *MockAuthorizationUsecase {
	mock := &MockAuthorizationUsecase{ctrl: ctrl}
	mock.recorder = &MockAuthorizationUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthorizationUsecase) EXPECT() *MockAuthorizationUsecaseMockRecorder {
	return m.recorder
}

// Authorize mocks base method.
func (m *MockAuthorizationUsecase) Authorize(arg0 context.Context, arg1 string) (*dto.AccountDTO, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Authorize", arg0, arg1)
	ret0, _ := ret[0].(*dto.AccountDTO)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Authorize indicates an expected call of Authorize.
func (mr *MockAuthorizationUsecaseMockRecorder) Authorize(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Authorize", reflect.TypeOf((*MockAuthorizationUsecase)(nil).Authorize), arg0, arg1)
}
