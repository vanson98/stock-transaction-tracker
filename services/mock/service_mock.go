// Code generated by MockGen. DO NOT EDIT.
// Source: stt/domain (interfaces: IAccountService,IInvestmentService)

// Package mock_service is a generated GoMock package.
package mock_service

import (
	context "context"
	reflect "reflect"
	db "stt/database/postgres/sqlc"
	dtos "stt/services/dtos"

	gomock "github.com/golang/mock/gomock"
)

// MockIAccountService is a mock of IAccountService interface.
type MockIAccountService struct {
	ctrl     *gomock.Controller
	recorder *MockIAccountServiceMockRecorder
}

// MockIAccountServiceMockRecorder is the mock recorder for MockIAccountService.
type MockIAccountServiceMockRecorder struct {
	mock *MockIAccountService
}

// NewMockIAccountService creates a new mock instance.
func NewMockIAccountService(ctrl *gomock.Controller) *MockIAccountService {
	mock := &MockIAccountService{ctrl: ctrl}
	mock.recorder = &MockIAccountServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIAccountService) EXPECT() *MockIAccountServiceMockRecorder {
	return m.recorder
}

// CreateNew mocks base method.
func (m *MockIAccountService) CreateNew(arg0 context.Context, arg1 db.CreateAccountParams) (db.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateNew", arg0, arg1)
	ret0, _ := ret[0].(db.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateNew indicates an expected call of CreateNew.
func (mr *MockIAccountServiceMockRecorder) CreateNew(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateNew", reflect.TypeOf((*MockIAccountService)(nil).CreateNew), arg0, arg1)
}

// Delete mocks base method.
func (m *MockIAccountService) Delete(arg0 context.Context, arg1 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockIAccountServiceMockRecorder) Delete(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockIAccountService)(nil).Delete), arg0, arg1)
}

// GetAllPaging mocks base method.
func (m *MockIAccountService) GetAllPaging(arg0 context.Context, arg1 db.ListAccountsParams) ([]db.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllPaging", arg0, arg1)
	ret0, _ := ret[0].([]db.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllPaging indicates an expected call of GetAllPaging.
func (mr *MockIAccountServiceMockRecorder) GetAllPaging(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllPaging", reflect.TypeOf((*MockIAccountService)(nil).GetAllPaging), arg0, arg1)
}

// TransferMoney mocks base method.
func (m *MockIAccountService) TransferMoney(arg0 context.Context, arg1 dtos.TransferMoneyTxParam) (dtos.TransferMoneyTxResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TransferMoney", arg0, arg1)
	ret0, _ := ret[0].(dtos.TransferMoneyTxResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// TransferMoney indicates an expected call of TransferMoney.
func (mr *MockIAccountServiceMockRecorder) TransferMoney(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TransferMoney", reflect.TypeOf((*MockIAccountService)(nil).TransferMoney), arg0, arg1)
}

// UpdateBalance mocks base method.
func (m *MockIAccountService) UpdateBalance(arg0 context.Context, arg1 db.AddAccountBalanceParams) (db.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateBalance", arg0, arg1)
	ret0, _ := ret[0].(db.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateBalance indicates an expected call of UpdateBalance.
func (mr *MockIAccountServiceMockRecorder) UpdateBalance(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateBalance", reflect.TypeOf((*MockIAccountService)(nil).UpdateBalance), arg0, arg1)
}

// MockIInvestmentService is a mock of IInvestmentService interface.
type MockIInvestmentService struct {
	ctrl     *gomock.Controller
	recorder *MockIInvestmentServiceMockRecorder
}

// MockIInvestmentServiceMockRecorder is the mock recorder for MockIInvestmentService.
type MockIInvestmentServiceMockRecorder struct {
	mock *MockIInvestmentService
}

// NewMockIInvestmentService creates a new mock instance.
func NewMockIInvestmentService(ctrl *gomock.Controller) *MockIInvestmentService {
	mock := &MockIInvestmentService{ctrl: ctrl}
	mock.recorder = &MockIInvestmentServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIInvestmentService) EXPECT() *MockIInvestmentServiceMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockIInvestmentService) Create(arg0 context.Context, arg1 db.CreateInvestmentParams) (db.Investment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(db.Investment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockIInvestmentServiceMockRecorder) Create(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockIInvestmentService)(nil).Create), arg0, arg1)
}

// Delete mocks base method.
func (m *MockIInvestmentService) Delete(arg0 context.Context, arg1 int32) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Delete", arg0, arg1)
}

// Delete indicates an expected call of Delete.
func (mr *MockIInvestmentServiceMockRecorder) Delete(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockIInvestmentService)(nil).Delete), arg0, arg1)
}

// GetAll mocks base method.
func (m *MockIInvestmentService) GetAll(arg0 context.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "GetAll", arg0)
}

// GetAll indicates an expected call of GetAll.
func (mr *MockIInvestmentServiceMockRecorder) GetAll(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockIInvestmentService)(nil).GetAll), arg0)
}

// GetById mocks base method.
func (m *MockIInvestmentService) GetById(arg0 context.Context, arg1 int32) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "GetById", arg0, arg1)
}

// GetById indicates an expected call of GetById.
func (mr *MockIInvestmentServiceMockRecorder) GetById(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetById", reflect.TypeOf((*MockIInvestmentService)(nil).GetById), arg0, arg1)
}
