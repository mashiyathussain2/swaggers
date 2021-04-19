// Code generated by MockGen. DO NOT EDIT.
// Source: go-app/app (interfaces: KeeperCatalog)

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	schema "go-app/schema"
	primitive "go.mongodb.org/mongo-driver/bson/primitive"
	reflect "reflect"
)

// MockKeeperCatalog is a mock of KeeperCatalog interface
type MockKeeperCatalog struct {
	ctrl     *gomock.Controller
	recorder *MockKeeperCatalogMockRecorder
}

// MockKeeperCatalogMockRecorder is the mock recorder for MockKeeperCatalog
type MockKeeperCatalogMockRecorder struct {
	mock *MockKeeperCatalog
}

// NewMockKeeperCatalog creates a new mock instance
func NewMockKeeperCatalog(ctrl *gomock.Controller) *MockKeeperCatalog {
	mock := &MockKeeperCatalog{ctrl: ctrl}
	mock.recorder = &MockKeeperCatalogMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockKeeperCatalog) EXPECT() *MockKeeperCatalogMockRecorder {
	return m.recorder
}

// AddCatalogContent mocks base method
func (m *MockKeeperCatalog) AddCatalogContent(arg0 *schema.AddCatalogContentOpts) (*schema.PayloadVideo, []error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddCatalogContent", arg0)
	ret0, _ := ret[0].(*schema.PayloadVideo)
	ret1, _ := ret[1].([]error)
	return ret0, ret1
}

// AddCatalogContent indicates an expected call of AddCatalogContent
func (mr *MockKeeperCatalogMockRecorder) AddCatalogContent(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddCatalogContent", reflect.TypeOf((*MockKeeperCatalog)(nil).AddCatalogContent), arg0)
}

// AddCatalogContentImage mocks base method
func (m *MockKeeperCatalog) AddCatalogContentImage(arg0 *schema.AddCatalogContentImageOpts) []error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddCatalogContentImage", arg0)
	ret0, _ := ret[0].([]error)
	return ret0
}

// AddCatalogContentImage indicates an expected call of AddCatalogContentImage
func (mr *MockKeeperCatalogMockRecorder) AddCatalogContentImage(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddCatalogContentImage", reflect.TypeOf((*MockKeeperCatalog)(nil).AddCatalogContentImage), arg0)
}

// AddVariant mocks base method
func (m *MockKeeperCatalog) AddVariant(arg0 *schema.AddVariantOpts) (*schema.CreateVariantResp, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddVariant", arg0)
	ret0, _ := ret[0].(*schema.CreateVariantResp)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddVariant indicates an expected call of AddVariant
func (mr *MockKeeperCatalogMockRecorder) AddVariant(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddVariant", reflect.TypeOf((*MockKeeperCatalog)(nil).AddVariant), arg0)
}

// CheckCatalogIDsExists mocks base method
func (m *MockKeeperCatalog) CheckCatalogIDsExists(arg0 context.Context, arg1 []primitive.ObjectID) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckCatalogIDsExists", arg0, arg1)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckCatalogIDsExists indicates an expected call of CheckCatalogIDsExists
func (mr *MockKeeperCatalogMockRecorder) CheckCatalogIDsExists(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckCatalogIDsExists", reflect.TypeOf((*MockKeeperCatalog)(nil).CheckCatalogIDsExists), arg0, arg1)
}

// CreateCatalog mocks base method
func (m *MockKeeperCatalog) CreateCatalog(arg0 *schema.CreateCatalogOpts) (*schema.CreateCatalogResp, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCatalog", arg0)
	ret0, _ := ret[0].(*schema.CreateCatalogResp)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateCatalog indicates an expected call of CreateCatalog
func (mr *MockKeeperCatalogMockRecorder) CreateCatalog(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCatalog", reflect.TypeOf((*MockKeeperCatalog)(nil).CreateCatalog), arg0)
}

// DeleteVariant mocks base method
func (m *MockKeeperCatalog) DeleteVariant(arg0 *schema.DeleteVariantOpts) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteVariant", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteVariant indicates an expected call of DeleteVariant
func (mr *MockKeeperCatalogMockRecorder) DeleteVariant(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteVariant", reflect.TypeOf((*MockKeeperCatalog)(nil).DeleteVariant), arg0)
}

// EditCatalog mocks base method
func (m *MockKeeperCatalog) EditCatalog(arg0 *schema.EditCatalogOpts) (*schema.EditCatalogResp, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EditCatalog", arg0)
	ret0, _ := ret[0].(*schema.EditCatalogResp)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EditCatalog indicates an expected call of EditCatalog
func (mr *MockKeeperCatalogMockRecorder) EditCatalog(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EditCatalog", reflect.TypeOf((*MockKeeperCatalog)(nil).EditCatalog), arg0)
}

// EditVariantSKU mocks base method
func (m *MockKeeperCatalog) EditVariantSKU(arg0 *schema.EditVariantSKU) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EditVariantSKU", arg0)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EditVariantSKU indicates an expected call of EditVariantSKU
func (mr *MockKeeperCatalogMockRecorder) EditVariantSKU(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EditVariantSKU", reflect.TypeOf((*MockKeeperCatalog)(nil).EditVariantSKU), arg0)
}

// GetAllCatalogInfo mocks base method
func (m *MockKeeperCatalog) GetAllCatalogInfo(arg0 primitive.ObjectID) (*schema.GetAllCatalogInfoResp, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllCatalogInfo", arg0)
	ret0, _ := ret[0].(*schema.GetAllCatalogInfoResp)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllCatalogInfo indicates an expected call of GetAllCatalogInfo
func (mr *MockKeeperCatalogMockRecorder) GetAllCatalogInfo(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllCatalogInfo", reflect.TypeOf((*MockKeeperCatalog)(nil).GetAllCatalogInfo), arg0)
}

// GetBasicCatalogInfo mocks base method
func (m *MockKeeperCatalog) GetBasicCatalogInfo(arg0 *schema.GetBasicCatalogFilter) ([]schema.GetBasicCatalogResp, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBasicCatalogInfo", arg0)
	ret0, _ := ret[0].([]schema.GetBasicCatalogResp)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBasicCatalogInfo indicates an expected call of GetBasicCatalogInfo
func (mr *MockKeeperCatalogMockRecorder) GetBasicCatalogInfo(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBasicCatalogInfo", reflect.TypeOf((*MockKeeperCatalog)(nil).GetBasicCatalogInfo), arg0)
}

// GetCatalogByIDs mocks base method
func (m *MockKeeperCatalog) GetCatalogByIDs(arg0 context.Context, arg1 []primitive.ObjectID) ([]schema.GetCatalogResp, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCatalogByIDs", arg0, arg1)
	ret0, _ := ret[0].([]schema.GetCatalogResp)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCatalogByIDs indicates an expected call of GetCatalogByIDs
func (mr *MockKeeperCatalogMockRecorder) GetCatalogByIDs(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCatalogByIDs", reflect.TypeOf((*MockKeeperCatalog)(nil).GetCatalogByIDs), arg0, arg1)
}

// GetCatalogBySlug mocks base method
func (m *MockKeeperCatalog) GetCatalogBySlug(arg0 string) (*schema.GetCatalogResp, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCatalogBySlug", arg0)
	ret0, _ := ret[0].(*schema.GetCatalogResp)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCatalogBySlug indicates an expected call of GetCatalogBySlug
func (mr *MockKeeperCatalogMockRecorder) GetCatalogBySlug(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCatalogBySlug", reflect.TypeOf((*MockKeeperCatalog)(nil).GetCatalogBySlug), arg0)
}

// GetCatalogContent mocks base method
func (m *MockKeeperCatalog) GetCatalogContent(arg0 primitive.ObjectID) ([]schema.CatalogContentInfoResp, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCatalogContent", arg0)
	ret0, _ := ret[0].([]schema.CatalogContentInfoResp)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCatalogContent indicates an expected call of GetCatalogContent
func (mr *MockKeeperCatalogMockRecorder) GetCatalogContent(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCatalogContent", reflect.TypeOf((*MockKeeperCatalog)(nil).GetCatalogContent), arg0)
}

// GetCatalogFilter mocks base method
func (m *MockKeeperCatalog) GetCatalogFilter() (*schema.GetCatalogFilterResp, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCatalogFilter")
	ret0, _ := ret[0].(*schema.GetCatalogFilterResp)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCatalogFilter indicates an expected call of GetCatalogFilter
func (mr *MockKeeperCatalogMockRecorder) GetCatalogFilter() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCatalogFilter", reflect.TypeOf((*MockKeeperCatalog)(nil).GetCatalogFilter))
}

// GetCatalogVariant mocks base method
func (m *MockKeeperCatalog) GetCatalogVariant(arg0, arg1 primitive.ObjectID) (*schema.GetCatalogVariantResp, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCatalogVariant", arg0, arg1)
	ret0, _ := ret[0].(*schema.GetCatalogVariantResp)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCatalogVariant indicates an expected call of GetCatalogVariant
func (mr *MockKeeperCatalogMockRecorder) GetCatalogVariant(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCatalogVariant", reflect.TypeOf((*MockKeeperCatalog)(nil).GetCatalogVariant), arg0, arg1)
}

// GetCatalogsByFilter mocks base method
func (m *MockKeeperCatalog) GetCatalogsByFilter(arg0 *schema.GetCatalogsByFilterOpts) ([]schema.GetCatalogResp, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCatalogsByFilter", arg0)
	ret0, _ := ret[0].([]schema.GetCatalogResp)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCatalogsByFilter indicates an expected call of GetCatalogsByFilter
func (mr *MockKeeperCatalogMockRecorder) GetCatalogsByFilter(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCatalogsByFilter", reflect.TypeOf((*MockKeeperCatalog)(nil).GetCatalogsByFilter), arg0)
}

// GetCollectionCatalogInfo mocks base method
func (m *MockKeeperCatalog) GetCollectionCatalogInfo(arg0 []primitive.ObjectID) ([]schema.GetAllCatalogInfoResp, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCollectionCatalogInfo", arg0)
	ret0, _ := ret[0].([]schema.GetAllCatalogInfoResp)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCollectionCatalogInfo indicates an expected call of GetCollectionCatalogInfo
func (mr *MockKeeperCatalogMockRecorder) GetCollectionCatalogInfo(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCollectionCatalogInfo", reflect.TypeOf((*MockKeeperCatalog)(nil).GetCollectionCatalogInfo), arg0)
}

// GetKeeperCatalogContent mocks base method
func (m *MockKeeperCatalog) GetKeeperCatalogContent(arg0 primitive.ObjectID) ([]schema.CatalogContentInfoResp, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetKeeperCatalogContent", arg0)
	ret0, _ := ret[0].([]schema.CatalogContentInfoResp)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetKeeperCatalogContent indicates an expected call of GetKeeperCatalogContent
func (mr *MockKeeperCatalogMockRecorder) GetKeeperCatalogContent(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetKeeperCatalogContent", reflect.TypeOf((*MockKeeperCatalog)(nil).GetKeeperCatalogContent), arg0)
}

// GetPebbleCatalogInfo mocks base method
func (m *MockKeeperCatalog) GetPebbleCatalogInfo(arg0 []primitive.ObjectID) ([]schema.GetAllCatalogInfoResp, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPebbleCatalogInfo", arg0)
	ret0, _ := ret[0].([]schema.GetAllCatalogInfoResp)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPebbleCatalogInfo indicates an expected call of GetPebbleCatalogInfo
func (mr *MockKeeperCatalogMockRecorder) GetPebbleCatalogInfo(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPebbleCatalogInfo", reflect.TypeOf((*MockKeeperCatalog)(nil).GetPebbleCatalogInfo), arg0)
}

// KeeperSearchCatalog mocks base method
func (m *MockKeeperCatalog) KeeperSearchCatalog(arg0 *schema.KeeperSearchCatalogOpts) ([]schema.KeeperSearchCatalogResp, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "KeeperSearchCatalog", arg0)
	ret0, _ := ret[0].([]schema.KeeperSearchCatalogResp)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// KeeperSearchCatalog indicates an expected call of KeeperSearchCatalog
func (mr *MockKeeperCatalogMockRecorder) KeeperSearchCatalog(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "KeeperSearchCatalog", reflect.TypeOf((*MockKeeperCatalog)(nil).KeeperSearchCatalog), arg0)
}

// RemoveContent mocks base method
func (m *MockKeeperCatalog) RemoveContent(arg0 *schema.RemoveContentOpts) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveContent", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveContent indicates an expected call of RemoveContent
func (mr *MockKeeperCatalogMockRecorder) RemoveContent(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveContent", reflect.TypeOf((*MockKeeperCatalog)(nil).RemoveContent), arg0)
}

// SyncCatalog mocks base method
func (m *MockKeeperCatalog) SyncCatalog(arg0 primitive.ObjectID) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SyncCatalog", arg0)
}

// SyncCatalog indicates an expected call of SyncCatalog
func (mr *MockKeeperCatalogMockRecorder) SyncCatalog(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SyncCatalog", reflect.TypeOf((*MockKeeperCatalog)(nil).SyncCatalog), arg0)
}

// SyncCatalogContent mocks base method
func (m *MockKeeperCatalog) SyncCatalogContent(arg0 primitive.ObjectID) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SyncCatalogContent", arg0)
}

// SyncCatalogContent indicates an expected call of SyncCatalogContent
func (mr *MockKeeperCatalogMockRecorder) SyncCatalogContent(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SyncCatalogContent", reflect.TypeOf((*MockKeeperCatalog)(nil).SyncCatalogContent), arg0)
}

// SyncCatalogs mocks base method
func (m *MockKeeperCatalog) SyncCatalogs(arg0 []primitive.ObjectID) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SyncCatalogs", arg0)
}

// SyncCatalogs indicates an expected call of SyncCatalogs
func (mr *MockKeeperCatalogMockRecorder) SyncCatalogs(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SyncCatalogs", reflect.TypeOf((*MockKeeperCatalog)(nil).SyncCatalogs), arg0)
}

// UpdateCatalogStatus mocks base method
func (m *MockKeeperCatalog) UpdateCatalogStatus(arg0 *schema.UpdateCatalogStatusOpts) ([]schema.UpdateCatalogStatusResp, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateCatalogStatus", arg0)
	ret0, _ := ret[0].([]schema.UpdateCatalogStatusResp)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateCatalogStatus indicates an expected call of UpdateCatalogStatus
func (mr *MockKeeperCatalogMockRecorder) UpdateCatalogStatus(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateCatalogStatus", reflect.TypeOf((*MockKeeperCatalog)(nil).UpdateCatalogStatus), arg0)
}
