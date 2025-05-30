// Code generated by MockGen. DO NOT EDIT.
// Source: entry.go
//
// Generated by this command:
//
//	mockgen -source=entry.go -package=repository -destination=../../../../../test/mock/domain/repository/entry.go
//

// Package repository is a generated GoMock package.
package repository

import (
	context "context"
	reflect "reflect"

	entity "github.com/atsumarukun/holos-storage-api/internal/app/api/domain/entity"
	uuid "github.com/google/uuid"
	gomock "go.uber.org/mock/gomock"
)

// MockEntryRepository is a mock of EntryRepository interface.
type MockEntryRepository struct {
	ctrl     *gomock.Controller
	recorder *MockEntryRepositoryMockRecorder
	isgomock struct{}
}

// MockEntryRepositoryMockRecorder is the mock recorder for MockEntryRepository.
type MockEntryRepositoryMockRecorder struct {
	mock *MockEntryRepository
}

// NewMockEntryRepository creates a new mock instance.
func NewMockEntryRepository(ctrl *gomock.Controller) *MockEntryRepository {
	mock := &MockEntryRepository{ctrl: ctrl}
	mock.recorder = &MockEntryRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEntryRepository) EXPECT() *MockEntryRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockEntryRepository) Create(arg0 context.Context, arg1 *entity.Entry) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockEntryRepositoryMockRecorder) Create(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockEntryRepository)(nil).Create), arg0, arg1)
}

// Delete mocks base method.
func (m *MockEntryRepository) Delete(arg0 context.Context, arg1 *entity.Entry) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockEntryRepositoryMockRecorder) Delete(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockEntryRepository)(nil).Delete), arg0, arg1)
}

// FindByVolumeIDAndAccountID mocks base method.
func (m *MockEntryRepository) FindByVolumeIDAndAccountID(arg0 context.Context, arg1, arg2 uuid.UUID, arg3 *string, arg4 *uint64) ([]*entity.Entry, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByVolumeIDAndAccountID", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].([]*entity.Entry)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByVolumeIDAndAccountID indicates an expected call of FindByVolumeIDAndAccountID.
func (mr *MockEntryRepositoryMockRecorder) FindByVolumeIDAndAccountID(arg0, arg1, arg2, arg3, arg4 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByVolumeIDAndAccountID", reflect.TypeOf((*MockEntryRepository)(nil).FindByVolumeIDAndAccountID), arg0, arg1, arg2, arg3, arg4)
}

// FindOneByKeyAndVolumeID mocks base method.
func (m *MockEntryRepository) FindOneByKeyAndVolumeID(arg0 context.Context, arg1 string, arg2 uuid.UUID) (*entity.Entry, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindOneByKeyAndVolumeID", arg0, arg1, arg2)
	ret0, _ := ret[0].(*entity.Entry)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindOneByKeyAndVolumeID indicates an expected call of FindOneByKeyAndVolumeID.
func (mr *MockEntryRepositoryMockRecorder) FindOneByKeyAndVolumeID(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindOneByKeyAndVolumeID", reflect.TypeOf((*MockEntryRepository)(nil).FindOneByKeyAndVolumeID), arg0, arg1, arg2)
}

// FindOneByKeyAndVolumeIDAndAccountID mocks base method.
func (m *MockEntryRepository) FindOneByKeyAndVolumeIDAndAccountID(arg0 context.Context, arg1 string, arg2, arg3 uuid.UUID) (*entity.Entry, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindOneByKeyAndVolumeIDAndAccountID", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(*entity.Entry)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindOneByKeyAndVolumeIDAndAccountID indicates an expected call of FindOneByKeyAndVolumeIDAndAccountID.
func (mr *MockEntryRepositoryMockRecorder) FindOneByKeyAndVolumeIDAndAccountID(arg0, arg1, arg2, arg3 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindOneByKeyAndVolumeIDAndAccountID", reflect.TypeOf((*MockEntryRepository)(nil).FindOneByKeyAndVolumeIDAndAccountID), arg0, arg1, arg2, arg3)
}

// Update mocks base method.
func (m *MockEntryRepository) Update(arg0 context.Context, arg1 *entity.Entry) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockEntryRepositoryMockRecorder) Update(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockEntryRepository)(nil).Update), arg0, arg1)
}
