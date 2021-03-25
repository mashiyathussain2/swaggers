// Code generated by MockGen. DO NOT EDIT.
// Source: go-app/app (interfaces: Content)

// Package mock is a generated GoMock package.
package mock

import (
	gomock "github.com/golang/mock/gomock"
	schema "go-app/schema"
	reflect "reflect"
)

// MockContent is a mock of Content interface
type MockContent struct {
	ctrl     *gomock.Controller
	recorder *MockContentMockRecorder
}

// MockContentMockRecorder is the mock recorder for MockContent
type MockContentMockRecorder struct {
	mock *MockContent
}

// NewMockContent creates a new mock instance
func NewMockContent(ctrl *gomock.Controller) *MockContent {
	mock := &MockContent{ctrl: ctrl}
	mock.recorder = &MockContentMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockContent) EXPECT() *MockContentMockRecorder {
	return m.recorder
}

// GenerateVideoUploadToken mocks base method
func (m *MockContent) GenerateVideoUploadToken(arg0 *schema.GenerateVideoUploadTokenOpts) (*schema.GenerateVideoUploadTokenResp, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateVideoUploadToken", arg0)
	ret0, _ := ret[0].(*schema.GenerateVideoUploadTokenResp)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateVideoUploadToken indicates an expected call of GenerateVideoUploadToken
func (mr *MockContentMockRecorder) GenerateVideoUploadToken(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateVideoUploadToken", reflect.TypeOf((*MockContent)(nil).GenerateVideoUploadToken), arg0)
}

// generateS3UploadToken mocks base method
func (m *MockContent) generateS3UploadToken(arg0, arg1 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "generateS3UploadToken", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// generateS3UploadToken indicates an expected call of generateS3UploadToken
func (mr *MockContentMockRecorder) generateS3UploadToken(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "generateS3UploadToken", reflect.TypeOf((*MockContent)(nil).generateS3UploadToken), arg0, arg1)
}
