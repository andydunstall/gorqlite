// Code generated by MockGen. DO NOT EDIT.
// Source: http_api.go

// Package mock_gorqlite is a generated GoMock package.
package mock_gorqlite

import (
	http "net/http"
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
)

// MockroundTripper is a mock of roundTripper interface.
type MockroundTripper struct {
	ctrl     *gomock.Controller
	recorder *MockroundTripperMockRecorder
}

// MockroundTripperMockRecorder is the mock recorder for MockroundTripper.
type MockroundTripperMockRecorder struct {
	mock *MockroundTripper
}

// NewMockroundTripper creates a new mock instance.
func NewMockroundTripper(ctrl *gomock.Controller) *MockroundTripper {
	mock := &MockroundTripper{ctrl: ctrl}
	mock.recorder = &MockroundTripperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockroundTripper) EXPECT() *MockroundTripperMockRecorder {
	return m.recorder
}

// RoundTrip mocks base method.
func (m *MockroundTripper) RoundTrip(arg0 *http.Request) (*http.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RoundTrip", arg0)
	ret0, _ := ret[0].(*http.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RoundTrip indicates an expected call of RoundTrip.
func (mr *MockroundTripperMockRecorder) RoundTrip(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RoundTrip", reflect.TypeOf((*MockroundTripper)(nil).RoundTrip), arg0)
}

// Mockclock is a mock of clock interface.
type Mockclock struct {
	ctrl     *gomock.Controller
	recorder *MockclockMockRecorder
}

// MockclockMockRecorder is the mock recorder for Mockclock.
type MockclockMockRecorder struct {
	mock *Mockclock
}

// NewMockclock creates a new mock instance.
func NewMockclock(ctrl *gomock.Controller) *Mockclock {
	mock := &Mockclock{ctrl: ctrl}
	mock.recorder = &MockclockMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mockclock) EXPECT() *MockclockMockRecorder {
	return m.recorder
}

// Sleep mocks base method.
func (m *Mockclock) Sleep(d time.Duration) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Sleep", d)
}

// Sleep indicates an expected call of Sleep.
func (mr *MockclockMockRecorder) Sleep(d interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Sleep", reflect.TypeOf((*Mockclock)(nil).Sleep), d)
}
