// Code generated by mockery v2.15.0. DO NOT EDIT.

package mocks

import (
	eval "github.com/DataDog/datadog-agent/pkg/compliance/eval"
	mock "github.com/stretchr/testify/mock"
)

// RegoConfiguration is an autogenerated mock type for the RegoConfiguration type
type RegoConfiguration struct {
	mock.Mock
}

// DumpInputPath provides a mock function with given fields:
func (_m *RegoConfiguration) DumpInputPath() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// ProvidedInput provides a mock function with given fields: ruleID
func (_m *RegoConfiguration) ProvidedInput(ruleID string) eval.RegoInputMap {
	ret := _m.Called(ruleID)

	var r0 eval.RegoInputMap
	if rf, ok := ret.Get(0).(func(string) eval.RegoInputMap); ok {
		r0 = rf(ruleID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(eval.RegoInputMap)
		}
	}

	return r0
}

// ShouldSkipRegoEval provides a mock function with given fields:
func (_m *RegoConfiguration) ShouldSkipRegoEval() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

type mockConstructorTestingTNewRegoConfiguration interface {
	mock.TestingT
	Cleanup(func())
}

// NewRegoConfiguration creates a new instance of RegoConfiguration. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewRegoConfiguration(t mockConstructorTestingTNewRegoConfiguration) *RegoConfiguration {
	mock := &RegoConfiguration{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
