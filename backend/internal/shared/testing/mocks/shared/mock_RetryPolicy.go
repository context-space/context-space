// Code generated by mockery v2.53.4. DO NOT EDIT.

package shared_mocks

import (
	mock "github.com/stretchr/testify/mock"

	time "time"
)

// MockRetryPolicy is an autogenerated mock type for the RetryPolicy type
type MockRetryPolicy struct {
	mock.Mock
}

type MockRetryPolicy_Expecter struct {
	mock *mock.Mock
}

func (_m *MockRetryPolicy) EXPECT() *MockRetryPolicy_Expecter {
	return &MockRetryPolicy_Expecter{mock: &_m.Mock}
}

// MaxAttempts provides a mock function with no fields
func (_m *MockRetryPolicy) MaxAttempts() int {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for MaxAttempts")
	}

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// MockRetryPolicy_MaxAttempts_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'MaxAttempts'
type MockRetryPolicy_MaxAttempts_Call struct {
	*mock.Call
}

// MaxAttempts is a helper method to define mock.On call
func (_e *MockRetryPolicy_Expecter) MaxAttempts() *MockRetryPolicy_MaxAttempts_Call {
	return &MockRetryPolicy_MaxAttempts_Call{Call: _e.mock.On("MaxAttempts")}
}

func (_c *MockRetryPolicy_MaxAttempts_Call) Run(run func()) *MockRetryPolicy_MaxAttempts_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockRetryPolicy_MaxAttempts_Call) Return(_a0 int) *MockRetryPolicy_MaxAttempts_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockRetryPolicy_MaxAttempts_Call) RunAndReturn(run func() int) *MockRetryPolicy_MaxAttempts_Call {
	_c.Call.Return(run)
	return _c
}

// NextBackoff provides a mock function with given fields: attempt
func (_m *MockRetryPolicy) NextBackoff(attempt int) time.Duration {
	ret := _m.Called(attempt)

	if len(ret) == 0 {
		panic("no return value specified for NextBackoff")
	}

	var r0 time.Duration
	if rf, ok := ret.Get(0).(func(int) time.Duration); ok {
		r0 = rf(attempt)
	} else {
		r0 = ret.Get(0).(time.Duration)
	}

	return r0
}

// MockRetryPolicy_NextBackoff_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'NextBackoff'
type MockRetryPolicy_NextBackoff_Call struct {
	*mock.Call
}

// NextBackoff is a helper method to define mock.On call
//   - attempt int
func (_e *MockRetryPolicy_Expecter) NextBackoff(attempt interface{}) *MockRetryPolicy_NextBackoff_Call {
	return &MockRetryPolicy_NextBackoff_Call{Call: _e.mock.On("NextBackoff", attempt)}
}

func (_c *MockRetryPolicy_NextBackoff_Call) Run(run func(attempt int)) *MockRetryPolicy_NextBackoff_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(int))
	})
	return _c
}

func (_c *MockRetryPolicy_NextBackoff_Call) Return(_a0 time.Duration) *MockRetryPolicy_NextBackoff_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockRetryPolicy_NextBackoff_Call) RunAndReturn(run func(int) time.Duration) *MockRetryPolicy_NextBackoff_Call {
	_c.Call.Return(run)
	return _c
}

// ShouldRetry provides a mock function with given fields: err, attempt
func (_m *MockRetryPolicy) ShouldRetry(err error, attempt int) bool {
	ret := _m.Called(err, attempt)

	if len(ret) == 0 {
		panic("no return value specified for ShouldRetry")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func(error, int) bool); ok {
		r0 = rf(err, attempt)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// MockRetryPolicy_ShouldRetry_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ShouldRetry'
type MockRetryPolicy_ShouldRetry_Call struct {
	*mock.Call
}

// ShouldRetry is a helper method to define mock.On call
//   - err error
//   - attempt int
func (_e *MockRetryPolicy_Expecter) ShouldRetry(err interface{}, attempt interface{}) *MockRetryPolicy_ShouldRetry_Call {
	return &MockRetryPolicy_ShouldRetry_Call{Call: _e.mock.On("ShouldRetry", err, attempt)}
}

func (_c *MockRetryPolicy_ShouldRetry_Call) Run(run func(err error, attempt int)) *MockRetryPolicy_ShouldRetry_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(error), args[1].(int))
	})
	return _c
}

func (_c *MockRetryPolicy_ShouldRetry_Call) Return(_a0 bool) *MockRetryPolicy_ShouldRetry_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockRetryPolicy_ShouldRetry_Call) RunAndReturn(run func(error, int) bool) *MockRetryPolicy_ShouldRetry_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockRetryPolicy creates a new instance of MockRetryPolicy. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockRetryPolicy(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockRetryPolicy {
	mock := &MockRetryPolicy{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
