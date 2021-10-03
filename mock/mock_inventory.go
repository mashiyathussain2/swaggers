// Code generated by MockGen. DO NOT EDIT.
// Source: go-app/app (interfaces: Inventory)

// Package mock is a generated GoMock package.
package mock

import (
	gomock "github.com/golang/mock/gomock"
	schema "go-app/schema"
	primitive "go.mongodb.org/mongo-driver/bson/primitive"
	reflect "reflect"
)

// MockInventory is a mock of Inventory interface
type MockInventory struct {
	ctrl     *gomock.Controller
	recorder *MockInventoryMockRecorder
}

// MockInventoryMockRecorder is the mock recorder for MockInventory
type MockInventoryMockRecorder struct {
	mock *MockInventory
}

// NewMockInventory creates a new mock instance
func NewMockInventory(ctrl *gomock.Controller) *MockInventory {
	mock := &MockInventory{ctrl: ctrl}
	mock.recorder = &MockInventoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockInventory) EXPECT() *MockInventoryMockRecorder {
	return m.recorder
}

// CheckInventoryExists mocks base method
func (m *MockInventory) CheckInventoryExists(arg0, arg1 primitive.ObjectID, arg2 int) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckInventoryExists", arg0, arg1, arg2)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckInventoryExists indicates an expected call of CheckInventoryExists
func (mr *MockInventoryMockRecorder) CheckInventoryExists(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckInventoryExists", reflect.TypeOf((*MockInventory)(nil).CheckInventoryExists), arg0, arg1, arg2)
}

// CreateInventory mocks base method
func (m *MockInventory) CreateInventory(arg0 *schema.CreateInventoryOpts) (primitive.ObjectID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateInventory", arg0)
	ret0, _ := ret[0].(primitive.ObjectID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateInventory indicates an expected call of CreateInventory
func (mr *MockInventoryMockRecorder) CreateInventory(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateInventory", reflect.TypeOf((*MockInventory)(nil).CreateInventory), arg0)
}

// SetOutOfStock mocks base method
func (m *MockInventory) SetOutOfStock(arg0 primitive.ObjectID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetOutOfStock", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetOutOfStock indicates an expected call of SetOutOfStock
func (mr *MockInventoryMockRecorder) SetOutOfStock(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetOutOfStock", reflect.TypeOf((*MockInventory)(nil).SetOutOfStock), arg0)
}

// UpdateInventory mocks base method
func (m *MockInventory) UpdateInventory(arg0 *schema.UpdateInventoryOpts) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateInventory", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateInventory indicates an expected call of UpdateInventory
func (mr *MockInventoryMockRecorder) UpdateInventory(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateInventory", reflect.TypeOf((*MockInventory)(nil).UpdateInventory), arg0)
}

// UpdateInventoryInternal mocks base method
func (m *MockInventory) UpdateInventoryInternal(arg0 []schema.UpdateInventoryCVOpts) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateInventoryInternal", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateInventoryInternal indicates an expected call of UpdateInventoryInternal
func (mr *MockInventoryMockRecorder) UpdateInventoryInternal(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateInventoryInternal", reflect.TypeOf((*MockInventory)(nil).UpdateInventoryInternal), arg0)
}

// UpdateInventorybySKUs mocks base method
func (m *MockInventory) UpdateInventorybySKUs(arg0 []schema.UpdateInventoryBySKUOpt) (*schema.UpdateInventoryBySKUResp, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateInventorybySKUs", arg0)
	ret0, _ := ret[0].(*schema.UpdateInventoryBySKUResp)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateInventorybySKUs indicates an expected call of UpdateInventorybySKUs
func (mr *MockInventoryMockRecorder) UpdateInventorybySKUs(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateInventorybySKUs", reflect.TypeOf((*MockInventory)(nil).UpdateInventorybySKUs), arg0)
}
