// Code generated by mockery v2.43.2. DO NOT EDIT.

package user

import mock "github.com/stretchr/testify/mock"

// MockStorer is an autogenerated mock type for the Storer type
type MockStorer struct {
	mock.Mock
}

type MockStorer_Expecter struct {
	mock *mock.Mock
}

func (_m *MockStorer) EXPECT() *MockStorer_Expecter {
	return &MockStorer_Expecter{mock: &_m.Mock}
}

// CreateUser provides a mock function with given fields: username, password
func (_m *MockStorer) CreateUser(username string, password string) error {
	ret := _m.Called(username, password)

	if len(ret) == 0 {
		panic("no return value specified for CreateUser")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(username, password)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockStorer_CreateUser_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateUser'
type MockStorer_CreateUser_Call struct {
	*mock.Call
}

// CreateUser is a helper method to define mock.On call
//   - username string
//   - password string
func (_e *MockStorer_Expecter) CreateUser(username interface{}, password interface{}) *MockStorer_CreateUser_Call {
	return &MockStorer_CreateUser_Call{Call: _e.mock.On("CreateUser", username, password)}
}

func (_c *MockStorer_CreateUser_Call) Run(run func(username string, password string)) *MockStorer_CreateUser_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string))
	})
	return _c
}

func (_c *MockStorer_CreateUser_Call) Return(_a0 error) *MockStorer_CreateUser_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockStorer_CreateUser_Call) RunAndReturn(run func(string, string) error) *MockStorer_CreateUser_Call {
	_c.Call.Return(run)
	return _c
}

// FindByCredential provides a mock function with given fields: username, password
func (_m *MockStorer) FindByCredential(username string, password string) (int32, error) {
	ret := _m.Called(username, password)

	if len(ret) == 0 {
		panic("no return value specified for FindByCredential")
	}

	var r0 int32
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string) (int32, error)); ok {
		return rf(username, password)
	}
	if rf, ok := ret.Get(0).(func(string, string) int32); ok {
		r0 = rf(username, password)
	} else {
		r0 = ret.Get(0).(int32)
	}

	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(username, password)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockStorer_FindByCredential_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindByCredential'
type MockStorer_FindByCredential_Call struct {
	*mock.Call
}

// FindByCredential is a helper method to define mock.On call
//   - username string
//   - password string
func (_e *MockStorer_Expecter) FindByCredential(username interface{}, password interface{}) *MockStorer_FindByCredential_Call {
	return &MockStorer_FindByCredential_Call{Call: _e.mock.On("FindByCredential", username, password)}
}

func (_c *MockStorer_FindByCredential_Call) Run(run func(username string, password string)) *MockStorer_FindByCredential_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string))
	})
	return _c
}

func (_c *MockStorer_FindByCredential_Call) Return(_a0 int32, _a1 error) *MockStorer_FindByCredential_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockStorer_FindByCredential_Call) RunAndReturn(run func(string, string) (int32, error)) *MockStorer_FindByCredential_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockStorer creates a new instance of MockStorer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockStorer(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockStorer {
	mock := &MockStorer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
