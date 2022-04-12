// Code generated by MockGen. DO NOT EDIT.
// Source: go-app/app (interfaces: User)

// Package mock is a generated GoMock package.
package mock

import (
	gomock "github.com/golang/mock/gomock"
	model "go-app/model"
	schema "go-app/schema"
	auth "go-app/server/auth"
	primitive "go.mongodb.org/mongo-driver/bson/primitive"
	reflect "reflect"
)

// MockUser is a mock of User interface
type MockUser struct {
	ctrl     *gomock.Controller
	recorder *MockUserMockRecorder
}

// GetUserIDByInfluencerID implements app.User
func (*MockUser) GetUserIDByInfluencerID(opts *schema.GetUserInfoByIDOpts) (primitive.ObjectID, error) {
	panic("unimplemented")
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

// CheckEmail mocks base method
func (m *MockUser) CheckEmail(arg0 *schema.CheckEmailOpts) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckEmail", arg0)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckEmail indicates an expected call of CheckEmail
func (mr *MockUserMockRecorder) CheckEmail(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckEmail", reflect.TypeOf((*MockUser)(nil).CheckEmail), arg0)
}

// CheckPhoneNo mocks base method
func (m *MockUser) CheckPhoneNo(arg0 *schema.CheckPhoneNoOpts) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckPhoneNo", arg0)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckPhoneNo indicates an expected call of CheckPhoneNo
func (mr *MockUserMockRecorder) CheckPhoneNo(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckPhoneNo", reflect.TypeOf((*MockUser)(nil).CheckPhoneNo), arg0)
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

// GetUserByID mocks base method
func (m *MockUser) GetUserByID(arg0 primitive.ObjectID) (*model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByID", arg0)
	ret0, _ := ret[0].(*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByID indicates an expected call of GetUserByID
func (mr *MockUserMockRecorder) GetUserByID(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByID", reflect.TypeOf((*MockUser)(nil).GetUserByID), arg0)
}

// GetUserClaim mocks base method
func (m *MockUser) GetUserClaim(arg0 *model.User, arg1 *model.Customer) auth.Claim {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserClaim", arg0, arg1)
	ret0, _ := ret[0].(auth.Claim)
	return ret0
}

// GetUserClaim indicates an expected call of GetUserClaim
func (mr *MockUserMockRecorder) GetUserClaim(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserClaim", reflect.TypeOf((*MockUser)(nil).GetUserClaim), arg0, arg1)
}

// GetUserInfoByID mocks base method
func (m *MockUser) GetUserInfoByID(arg0 *schema.GetUserInfoByIDOpts) (primitive.M, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserInfoByID", arg0)
	ret0, _ := ret[0].(primitive.M)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserInfoByID indicates an expected call of GetUserInfoByID
func (mr *MockUserMockRecorder) GetUserInfoByID(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserInfoByID", reflect.TypeOf((*MockUser)(nil).GetUserInfoByID), arg0)
}

// LoginWithApple mocks base method
func (m *MockUser) LoginWithApple(arg0 *schema.LoginWithApple) (auth.Claim, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoginWithApple", arg0)
	ret0, _ := ret[0].(auth.Claim)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoginWithApple indicates an expected call of LoginWithApple
func (mr *MockUserMockRecorder) LoginWithApple(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoginWithApple", reflect.TypeOf((*MockUser)(nil).LoginWithApple), arg0)
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

// UpdateUserEmail mocks base method
func (m *MockUser) UpdateUserEmail(arg0 *schema.UpdateUserEmailOpts) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserEmail", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUserEmail indicates an expected call of UpdateUserEmail
func (mr *MockUserMockRecorder) UpdateUserEmail(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserEmail", reflect.TypeOf((*MockUser)(nil).UpdateUserEmail), arg0)
}

// UpdateUserPhoneNo mocks base method
func (m *MockUser) UpdateUserPhoneNo(arg0 *schema.UpdateUserPhoneNoOpts) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserPhoneNo", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUserPhoneNo indicates an expected call of UpdateUserPhoneNo
func (mr *MockUserMockRecorder) UpdateUserPhoneNo(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserPhoneNo", reflect.TypeOf((*MockUser)(nil).UpdateUserPhoneNo), arg0)
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

// VerifyPhoneNo mocks base method
func (m *MockUser) VerifyPhoneNo(arg0 *schema.VerifyPhoneNoOpts) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerifyPhoneNo", arg0)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// VerifyPhoneNo indicates an expected call of VerifyPhoneNo
func (mr *MockUserMockRecorder) VerifyPhoneNo(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifyPhoneNo", reflect.TypeOf((*MockUser)(nil).VerifyPhoneNo), arg0)
}
