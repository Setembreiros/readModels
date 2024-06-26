// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package mock_userprofile is a generated GoMock package.
package mock_userprofile

import (
	userprofile "readmodels/internal/userprofile"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// AddNewUserProfile mocks base method.
func (m *MockRepository) AddNewUserProfile(data *userprofile.UserProfile) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddNewUserProfile", data)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddNewUserProfile indicates an expected call of AddNewUserProfile.
func (mr *MockRepositoryMockRecorder) AddNewUserProfile(data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddNewUserProfile", reflect.TypeOf((*MockRepository)(nil).AddNewUserProfile), data)
}

// GetUserProfile mocks base method.
func (m *MockRepository) GetUserProfile(username string) (*userprofile.UserProfile, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserProfile", username)
	ret0, _ := ret[0].(*userprofile.UserProfile)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserProfile indicates an expected call of GetUserProfile.
func (mr *MockRepositoryMockRecorder) GetUserProfile(username interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserProfile", reflect.TypeOf((*MockRepository)(nil).GetUserProfile), username)
}

// UpdateUserProfile mocks base method.
func (m *MockRepository) UpdateUserProfile(data *userprofile.UserProfile) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserProfile", data)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUserProfile indicates an expected call of UpdateUserProfile.
func (mr *MockRepositoryMockRecorder) UpdateUserProfile(data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserProfile", reflect.TypeOf((*MockRepository)(nil).UpdateUserProfile), data)
}
