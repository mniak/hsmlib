// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/mniak/hsmlib/multi (interfaces: IDManager)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockIDManager is a mock of IDManager interface.
type MockIDManager struct {
	ctrl     *gomock.Controller
	recorder *MockIDManagerMockRecorder
}

// MockIDManagerMockRecorder is the mock recorder for MockIDManager.
type MockIDManagerMockRecorder struct {
	mock *MockIDManager
}

// NewMockIDManager creates a new mock instance.
func NewMockIDManager(ctrl *gomock.Controller) *MockIDManager {
	mock := &MockIDManager{ctrl: ctrl}
	mock.recorder = &MockIDManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIDManager) EXPECT() *MockIDManagerMockRecorder {
	return m.recorder
}

// CloseAllChannels mocks base method.
func (m *MockIDManager) CloseAllChannels() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "CloseAllChannels")
}

// CloseAllChannels indicates an expected call of CloseAllChannels.
func (mr *MockIDManagerMockRecorder) CloseAllChannels() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloseAllChannels", reflect.TypeOf((*MockIDManager)(nil).CloseAllChannels))
}

// FindChannel mocks base method.
func (m *MockIDManager) FindChannel(arg0 []byte) (chan<- []byte, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindChannel", arg0)
	ret0, _ := ret[0].(chan<- []byte)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// FindChannel indicates an expected call of FindChannel.
func (mr *MockIDManagerMockRecorder) FindChannel(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindChannel", reflect.TypeOf((*MockIDManager)(nil).FindChannel), arg0)
}

// IDLength mocks base method.
func (m *MockIDManager) IDLength() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IDLength")
	ret0, _ := ret[0].(int)
	return ret0
}

// IDLength indicates an expected call of IDLength.
func (mr *MockIDManagerMockRecorder) IDLength() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IDLength", reflect.TypeOf((*MockIDManager)(nil).IDLength))
}

// NewID mocks base method.
func (m *MockIDManager) NewID() ([]byte, <-chan []byte) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewID")
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(<-chan []byte)
	return ret0, ret1
}

// NewID indicates an expected call of NewID.
func (mr *MockIDManagerMockRecorder) NewID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewID", reflect.TypeOf((*MockIDManager)(nil).NewID))
}
