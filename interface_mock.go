// Code generated by MockGen. DO NOT EDIT.
// Source: interface.go
//
// Generated by this command:
//
//	mockgen -source interface.go -destination interface_mock.go -package ctxlog
//

// Package ctxlog is a generated GoMock package.
package ctxlog

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockILogger is a mock of ILogger interface.
type MockILogger struct {
	ctrl     *gomock.Controller
	recorder *MockILoggerMockRecorder
}

// MockILoggerMockRecorder is the mock recorder for MockILogger.
type MockILoggerMockRecorder struct {
	mock *MockILogger
}

// NewMockILogger creates a new mock instance.
func NewMockILogger(ctrl *gomock.Controller) *MockILogger {
	mock := &MockILogger{ctrl: ctrl}
	mock.recorder = &MockILoggerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockILogger) EXPECT() *MockILoggerMockRecorder {
	return m.recorder
}

// Debug mocks base method.
func (m *MockILogger) Debug(ctx context.Context, msg string, args ...any) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, msg}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Debug", varargs...)
}

// Debug indicates an expected call of Debug.
func (mr *MockILoggerMockRecorder) Debug(ctx, msg any, args ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, msg}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Debug", reflect.TypeOf((*MockILogger)(nil).Debug), varargs...)
}

// Error mocks base method.
func (m *MockILogger) Error(ctx context.Context, msg string, args ...any) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, msg}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Error", varargs...)
}

// Error indicates an expected call of Error.
func (mr *MockILoggerMockRecorder) Error(ctx, msg any, args ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, msg}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Error", reflect.TypeOf((*MockILogger)(nil).Error), varargs...)
}

// Info mocks base method.
func (m *MockILogger) Info(ctx context.Context, msg string, args ...any) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, msg}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Info", varargs...)
}

// Info indicates an expected call of Info.
func (mr *MockILoggerMockRecorder) Info(ctx, msg any, args ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, msg}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Info", reflect.TypeOf((*MockILogger)(nil).Info), varargs...)
}

// Warn mocks base method.
func (m *MockILogger) Warn(ctx context.Context, msg string, args ...any) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, msg}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Warn", varargs...)
}

// Warn indicates an expected call of Warn.
func (mr *MockILoggerMockRecorder) Warn(ctx, msg any, args ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, msg}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Warn", reflect.TypeOf((*MockILogger)(nil).Warn), varargs...)
}
