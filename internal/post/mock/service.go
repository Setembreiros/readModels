// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package mock_post is a generated GoMock package.
package mock_post

import (
	post "readmodels/internal/post"
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

// AddNewPostMetadata mocks base method.
func (m *MockRepository) AddNewPostMetadata(data *post.PostMetadata) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddNewPostMetadata", data)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddNewPostMetadata indicates an expected call of AddNewPostMetadata.
func (mr *MockRepositoryMockRecorder) AddNewPostMetadata(data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddNewPostMetadata", reflect.TypeOf((*MockRepository)(nil).AddNewPostMetadata), data)
}

// GetPostMetadatasByUser mocks base method.
func (m *MockRepository) GetPostMetadatasByUser(username string) ([]*post.PostMetadata, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPostMetadatasByUser", username)
	ret0, _ := ret[0].([]*post.PostMetadata)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPostMetadatasByUser indicates an expected call of GetPostMetadatasByUser.
func (mr *MockRepositoryMockRecorder) GetPostMetadatasByUser(username interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPostMetadatasByUser", reflect.TypeOf((*MockRepository)(nil).GetPostMetadatasByUser), username)
}