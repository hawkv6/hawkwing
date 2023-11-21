// Code generated by MockGen. DO NOT EDIT.
// Source: client.go
//
// Generated by this command:
//
//	mockgen -source=client.go -destination=client_mock.go -package=client
//
// Package client is a generated GoMock package.
package client

import (
	reflect "reflect"

	ebpf "github.com/cilium/ebpf"
	gomock "go.uber.org/mock/gomock"
)

// MockClientBpfReader is a mock of ClientBpfReader interface.
type MockClientBpfReader struct {
	ctrl     *gomock.Controller
	recorder *MockClientBpfReaderMockRecorder
}

// MockClientBpfReaderMockRecorder is the mock recorder for MockClientBpfReader.
type MockClientBpfReaderMockRecorder struct {
	mock *MockClientBpfReader
}

// NewMockClientBpfReader creates a new mock instance.
func NewMockClientBpfReader(ctrl *gomock.Controller) *MockClientBpfReader {
	mock := &MockClientBpfReader{ctrl: ctrl}
	mock.recorder = &MockClientBpfReaderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClientBpfReader) EXPECT() *MockClientBpfReaderMockRecorder {
	return m.recorder
}

// ReadClientBpfObjects mocks base method.
func (m *MockClientBpfReader) ReadClientBpfObjects() (*clientObjects, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadClientBpfObjects")
	ret0, _ := ret[0].(*clientObjects)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadClientBpfObjects indicates an expected call of ReadClientBpfObjects.
func (mr *MockClientBpfReaderMockRecorder) ReadClientBpfObjects() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadClientBpfObjects", reflect.TypeOf((*MockClientBpfReader)(nil).ReadClientBpfObjects))
}

// ReadClientBpfSpecs mocks base method.
func (m *MockClientBpfReader) ReadClientBpfSpecs() (*ebpf.CollectionSpec, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadClientBpfSpecs")
	ret0, _ := ret[0].(*ebpf.CollectionSpec)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadClientBpfSpecs indicates an expected call of ReadClientBpfSpecs.
func (mr *MockClientBpfReaderMockRecorder) ReadClientBpfSpecs() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadClientBpfSpecs", reflect.TypeOf((*MockClientBpfReader)(nil).ReadClientBpfSpecs))
}