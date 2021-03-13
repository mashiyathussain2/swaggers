// Code generated by MockGen. DO NOT EDIT.
// Source: go-app/app (interfaces: User)

// Package mock is a generated GoMock package.
package mock

import (
	gomock "github.com/golang/mock/gomock"
	schema "go-app/schema"
	auth "go-app/server/auth"
	reflect "reflect"
)

// MockUser is a mock of User interface
type MockUser struct {
	ctrl     *gomock.Controller
	recorder *MockUserMockRecorder
}

// MockUserMockRecorder is the mock recorder for MockUser
type MockUserMockRecorder struct {
	mock *MockUser
}

// NewMockUser creates a new mock instance
func NewMockUser(ctrl *gomock.Controller) *MockUser {
	mock := &MockUser{ctrl: ctrl}
	mock.recorder = &MockUserMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockUser) EXPECT() *MockUserMockRecorder {
	return m.recorder
}

// CreateUser mocks base method
func (m *MockUser) CreateUser(arg0 *schema.CreateUserOpts) (*schema.CreateUserResp, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", arg0)
	ret0, _ := ret[0].(*schema.CreateUserResp)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser
func (mr *MockUserMockRecorder) CreateUser(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockUser)(nil).CreateUser), arg0)
}

// EmailLoginCustomerUser mocks base method
func (m *MockUser) EmailLoginCustomerUser(arg0 *schema.EmailLoginCustomerOpts) (auth.Claim, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EmailLoginCustomerUser", arg0)
	ret0, _ := ret[0].(auth.Claim)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EmailLoginCustomerUser indicates an expected call of EmailLoginCustomerUser
func (mr *MockUserMockRecorder) EmailLoginCustomerUser(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EmailLoginCustomerUser", reflect.TypeOf((*MockUser)(nil).EmailLoginCustomerUser), arg0)
}

// ForgotPassword mocks base method
func (m *MockUser) ForgotPassword(arg0 *schema.ForgotPasswordOpts) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ForgotPassword", arg0)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ForgotPassword indicates an expected call of ForgotPassword
func (mr *MockUserMockRecorder) ForgotPassword(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ForgotPassword", reflect.TypeOf((*MockUser)(nil).ForgotPassword), arg0)
}

// GenerateMobileLoginOTP mocks base method
func (m *MockUser) GenerateMobileLoginOTP(arg0 *schema.GenerateMobileLoginOTPOpts) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateMobileLoginOTP", arg0)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateMobileLoginOTP indicates an expected call of GenerateMobileLoginOTP
func (mr *MockUserMockRecorder) GenerateMobileLoginOTP(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateMobileLoginOTP", reflect.TypeOf((*MockUser)(nil).GenerateMobileLoginOTP), arg0)
}

// GetUserByEMail mocks base method
func (m *MockUser) GetUserByEMail(arg0 string) (*schema.GetUserResp, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByEMail", arg0)
	ret0, _ := ret[0].(*schema.GetUserResp)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByEMail indicates an expected call of GetUserByEMail
func (mr *MockUserMockRecorder) GetUserByEMail(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByEMail", reflect.TypeOf((*MockUser)(nil).GetUserByEMail), arg0)
}

// LoginWithSocial mocks base method
func (m *MockUser) LoginWithSocial(arg0 *schema.LoginWithSocial) (auth.Claim, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoginWithSocial", arg0)
	ret0, _ := ret[0].(auth.Claim)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoginWithSocial indicates an expected call of LoginWithSocial
func (mr *MockUserMockRecorder) LoginWithSocial(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoginWithSocial", reflect.TypeOf((*MockUser)(nil).LoginWithSocial), arg0)
}

// MobileLoginCustomerUser mocks base method
func (m *MockUser) MobileLoginCustomerUser(arg0 *schema.MobileLoginCustomerUserOpts) (auth.Claim, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MobileLoginCustomerUser", arg0)
	ret0, _ := ret[0].(auth.Claim)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// MobileLoginCustomerUser indicates an expected call of MobileLoginCustomerUser
func (mr *MockUserMockRecorder) MobileLoginCustomerUser(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MobileLoginCustomerUser", reflect.TypeOf((*MockUser)(nil).MobileLoginCustomerUser), arg0)
}

// ResendConfirmationEmail mocks base method
func (m *MockUser) ResendConfirmationEmail(arg0 *schema.ResendVerificationEmailOpts) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ResendConfirmationEmail", arg0)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ResendConfirmationEmail indicates an expected call of ResendConfirmationEmail
func (mr *MockUserMockRecorder) ResendConfirmationEmail(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResendConfirmationEmail", reflect.TypeOf((*MockUser)(nil).ResendConfirmationEmail), arg0)
}

// ResetPassword mocks base method
func (m *MockUser) ResetPassword(arg0 *schema.ResetPasswordOpts) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ResetPassword", arg0)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ResetPassword indicates an expected call of ResetPassword
func (mr *MockUserMockRecorder) ResetPassword(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResetPassword", reflect.TypeOf((*MockUser)(nil).ResetPassword), arg0)
}

// VerifyEmail mocks base method
func (m *MockUser) VerifyEmail(arg0 *schema.VerifyEmailOpts) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerifyEmail", arg0)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// VerifyEmail indicates an expected call of VerifyEmail
func (mr *MockUserMockRecorder) VerifyEmail(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifyEmail", reflect.TypeOf((*MockUser)(nil).VerifyEmail), arg0)
}
