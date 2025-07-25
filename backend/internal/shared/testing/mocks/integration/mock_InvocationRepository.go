// Code generated by mockery v2.53.4. DO NOT EDIT.

package integration_mocks

import (
	context "context"

	domain "github.com/context-space/context-space/backend/internal/integration/domain"
	mock "github.com/stretchr/testify/mock"
)

// MockInvocationRepository is an autogenerated mock type for the InvocationRepository type
type MockInvocationRepository struct {
	mock.Mock
}

type MockInvocationRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *MockInvocationRepository) EXPECT() *MockInvocationRepository_Expecter {
	return &MockInvocationRepository_Expecter{mock: &_m.Mock}
}

// CountByOperationIdentifier provides a mock function with given fields: ctx, providerIdentifier, operationIdentifier
func (_m *MockInvocationRepository) CountByOperationIdentifier(ctx context.Context, providerIdentifier string, operationIdentifier string) (int64, error) {
	ret := _m.Called(ctx, providerIdentifier, operationIdentifier)

	if len(ret) == 0 {
		panic("no return value specified for CountByOperationIdentifier")
	}

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (int64, error)); ok {
		return rf(ctx, providerIdentifier, operationIdentifier)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) int64); ok {
		r0 = rf(ctx, providerIdentifier, operationIdentifier)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, providerIdentifier, operationIdentifier)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockInvocationRepository_CountByOperationIdentifier_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CountByOperationIdentifier'
type MockInvocationRepository_CountByOperationIdentifier_Call struct {
	*mock.Call
}

// CountByOperationIdentifier is a helper method to define mock.On call
//   - ctx context.Context
//   - providerIdentifier string
//   - operationIdentifier string
func (_e *MockInvocationRepository_Expecter) CountByOperationIdentifier(ctx interface{}, providerIdentifier interface{}, operationIdentifier interface{}) *MockInvocationRepository_CountByOperationIdentifier_Call {
	return &MockInvocationRepository_CountByOperationIdentifier_Call{Call: _e.mock.On("CountByOperationIdentifier", ctx, providerIdentifier, operationIdentifier)}
}

func (_c *MockInvocationRepository_CountByOperationIdentifier_Call) Run(run func(ctx context.Context, providerIdentifier string, operationIdentifier string)) *MockInvocationRepository_CountByOperationIdentifier_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *MockInvocationRepository_CountByOperationIdentifier_Call) Return(_a0 int64, _a1 error) *MockInvocationRepository_CountByOperationIdentifier_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockInvocationRepository_CountByOperationIdentifier_Call) RunAndReturn(run func(context.Context, string, string) (int64, error)) *MockInvocationRepository_CountByOperationIdentifier_Call {
	_c.Call.Return(run)
	return _c
}

// CountByProviderIdentifier provides a mock function with given fields: ctx, providerIdentifier
func (_m *MockInvocationRepository) CountByProviderIdentifier(ctx context.Context, providerIdentifier string) (int64, error) {
	ret := _m.Called(ctx, providerIdentifier)

	if len(ret) == 0 {
		panic("no return value specified for CountByProviderIdentifier")
	}

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (int64, error)); ok {
		return rf(ctx, providerIdentifier)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) int64); ok {
		r0 = rf(ctx, providerIdentifier)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, providerIdentifier)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockInvocationRepository_CountByProviderIdentifier_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CountByProviderIdentifier'
type MockInvocationRepository_CountByProviderIdentifier_Call struct {
	*mock.Call
}

// CountByProviderIdentifier is a helper method to define mock.On call
//   - ctx context.Context
//   - providerIdentifier string
func (_e *MockInvocationRepository_Expecter) CountByProviderIdentifier(ctx interface{}, providerIdentifier interface{}) *MockInvocationRepository_CountByProviderIdentifier_Call {
	return &MockInvocationRepository_CountByProviderIdentifier_Call{Call: _e.mock.On("CountByProviderIdentifier", ctx, providerIdentifier)}
}

func (_c *MockInvocationRepository_CountByProviderIdentifier_Call) Run(run func(ctx context.Context, providerIdentifier string)) *MockInvocationRepository_CountByProviderIdentifier_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockInvocationRepository_CountByProviderIdentifier_Call) Return(_a0 int64, _a1 error) *MockInvocationRepository_CountByProviderIdentifier_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockInvocationRepository_CountByProviderIdentifier_Call) RunAndReturn(run func(context.Context, string) (int64, error)) *MockInvocationRepository_CountByProviderIdentifier_Call {
	_c.Call.Return(run)
	return _c
}

// CountByUserID provides a mock function with given fields: ctx, userID
func (_m *MockInvocationRepository) CountByUserID(ctx context.Context, userID string) (int64, error) {
	ret := _m.Called(ctx, userID)

	if len(ret) == 0 {
		panic("no return value specified for CountByUserID")
	}

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (int64, error)); ok {
		return rf(ctx, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) int64); ok {
		r0 = rf(ctx, userID)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockInvocationRepository_CountByUserID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CountByUserID'
type MockInvocationRepository_CountByUserID_Call struct {
	*mock.Call
}

// CountByUserID is a helper method to define mock.On call
//   - ctx context.Context
//   - userID string
func (_e *MockInvocationRepository_Expecter) CountByUserID(ctx interface{}, userID interface{}) *MockInvocationRepository_CountByUserID_Call {
	return &MockInvocationRepository_CountByUserID_Call{Call: _e.mock.On("CountByUserID", ctx, userID)}
}

func (_c *MockInvocationRepository_CountByUserID_Call) Run(run func(ctx context.Context, userID string)) *MockInvocationRepository_CountByUserID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockInvocationRepository_CountByUserID_Call) Return(_a0 int64, _a1 error) *MockInvocationRepository_CountByUserID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockInvocationRepository_CountByUserID_Call) RunAndReturn(run func(context.Context, string) (int64, error)) *MockInvocationRepository_CountByUserID_Call {
	_c.Call.Return(run)
	return _c
}

// Create provides a mock function with given fields: ctx, invocation
func (_m *MockInvocationRepository) Create(ctx context.Context, invocation *domain.Invocation) error {
	ret := _m.Called(ctx, invocation)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.Invocation) error); ok {
		r0 = rf(ctx, invocation)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockInvocationRepository_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type MockInvocationRepository_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - ctx context.Context
//   - invocation *domain.Invocation
func (_e *MockInvocationRepository_Expecter) Create(ctx interface{}, invocation interface{}) *MockInvocationRepository_Create_Call {
	return &MockInvocationRepository_Create_Call{Call: _e.mock.On("Create", ctx, invocation)}
}

func (_c *MockInvocationRepository_Create_Call) Run(run func(ctx context.Context, invocation *domain.Invocation)) *MockInvocationRepository_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*domain.Invocation))
	})
	return _c
}

func (_c *MockInvocationRepository_Create_Call) Return(_a0 error) *MockInvocationRepository_Create_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockInvocationRepository_Create_Call) RunAndReturn(run func(context.Context, *domain.Invocation) error) *MockInvocationRepository_Create_Call {
	_c.Call.Return(run)
	return _c
}

// GetByID provides a mock function with given fields: ctx, id
func (_m *MockInvocationRepository) GetByID(ctx context.Context, id string) (*domain.Invocation, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for GetByID")
	}

	var r0 *domain.Invocation
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*domain.Invocation, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *domain.Invocation); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Invocation)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockInvocationRepository_GetByID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetByID'
type MockInvocationRepository_GetByID_Call struct {
	*mock.Call
}

// GetByID is a helper method to define mock.On call
//   - ctx context.Context
//   - id string
func (_e *MockInvocationRepository_Expecter) GetByID(ctx interface{}, id interface{}) *MockInvocationRepository_GetByID_Call {
	return &MockInvocationRepository_GetByID_Call{Call: _e.mock.On("GetByID", ctx, id)}
}

func (_c *MockInvocationRepository_GetByID_Call) Run(run func(ctx context.Context, id string)) *MockInvocationRepository_GetByID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockInvocationRepository_GetByID_Call) Return(_a0 *domain.Invocation, _a1 error) *MockInvocationRepository_GetByID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockInvocationRepository_GetByID_Call) RunAndReturn(run func(context.Context, string) (*domain.Invocation, error)) *MockInvocationRepository_GetByID_Call {
	_c.Call.Return(run)
	return _c
}

// ListByUserID provides a mock function with given fields: ctx, userID, limit, offset
func (_m *MockInvocationRepository) ListByUserID(ctx context.Context, userID string, limit int, offset int) ([]*domain.Invocation, error) {
	ret := _m.Called(ctx, userID, limit, offset)

	if len(ret) == 0 {
		panic("no return value specified for ListByUserID")
	}

	var r0 []*domain.Invocation
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, int, int) ([]*domain.Invocation, error)); ok {
		return rf(ctx, userID, limit, offset)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, int, int) []*domain.Invocation); ok {
		r0 = rf(ctx, userID, limit, offset)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*domain.Invocation)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, int, int) error); ok {
		r1 = rf(ctx, userID, limit, offset)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockInvocationRepository_ListByUserID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListByUserID'
type MockInvocationRepository_ListByUserID_Call struct {
	*mock.Call
}

// ListByUserID is a helper method to define mock.On call
//   - ctx context.Context
//   - userID string
//   - limit int
//   - offset int
func (_e *MockInvocationRepository_Expecter) ListByUserID(ctx interface{}, userID interface{}, limit interface{}, offset interface{}) *MockInvocationRepository_ListByUserID_Call {
	return &MockInvocationRepository_ListByUserID_Call{Call: _e.mock.On("ListByUserID", ctx, userID, limit, offset)}
}

func (_c *MockInvocationRepository_ListByUserID_Call) Run(run func(ctx context.Context, userID string, limit int, offset int)) *MockInvocationRepository_ListByUserID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(int), args[3].(int))
	})
	return _c
}

func (_c *MockInvocationRepository_ListByUserID_Call) Return(_a0 []*domain.Invocation, _a1 error) *MockInvocationRepository_ListByUserID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockInvocationRepository_ListByUserID_Call) RunAndReturn(run func(context.Context, string, int, int) ([]*domain.Invocation, error)) *MockInvocationRepository_ListByUserID_Call {
	_c.Call.Return(run)
	return _c
}

// Update provides a mock function with given fields: ctx, invocation
func (_m *MockInvocationRepository) Update(ctx context.Context, invocation *domain.Invocation) error {
	ret := _m.Called(ctx, invocation)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.Invocation) error); ok {
		r0 = rf(ctx, invocation)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockInvocationRepository_Update_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Update'
type MockInvocationRepository_Update_Call struct {
	*mock.Call
}

// Update is a helper method to define mock.On call
//   - ctx context.Context
//   - invocation *domain.Invocation
func (_e *MockInvocationRepository_Expecter) Update(ctx interface{}, invocation interface{}) *MockInvocationRepository_Update_Call {
	return &MockInvocationRepository_Update_Call{Call: _e.mock.On("Update", ctx, invocation)}
}

func (_c *MockInvocationRepository_Update_Call) Run(run func(ctx context.Context, invocation *domain.Invocation)) *MockInvocationRepository_Update_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*domain.Invocation))
	})
	return _c
}

func (_c *MockInvocationRepository_Update_Call) Return(_a0 error) *MockInvocationRepository_Update_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockInvocationRepository_Update_Call) RunAndReturn(run func(context.Context, *domain.Invocation) error) *MockInvocationRepository_Update_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockInvocationRepository creates a new instance of MockInvocationRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockInvocationRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockInvocationRepository {
	mock := &MockInvocationRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
