// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// StoreNoop is an autogenerated mock type for the StoreNoop type
type StoreNoop struct {
	mock.Mock
}

// NewStoreNoop creates a new instance of StoreNoop. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewStoreNoop(t interface {
	mock.TestingT
	Cleanup(func())
}) *StoreNoop {
	mock := &StoreNoop{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
