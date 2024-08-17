// Code generated by mockery v2.43.2. DO NOT EDIT.

package category

import (
	store "github.com/opplieam/bb-admin-api/internal/store"
	utils "github.com/opplieam/bb-admin-api/internal/utils"
	mock "github.com/stretchr/testify/mock"
)

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

// GetAllCategory provides a mock function with given fields:
func (_m *MockStorer) GetAllCategory() ([]store.AllCategoryResult, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetAllCategory")
	}

	var r0 []store.AllCategoryResult
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]store.AllCategoryResult, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []store.AllCategoryResult); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]store.AllCategoryResult)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockStorer_GetAllCategory_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAllCategory'
type MockStorer_GetAllCategory_Call struct {
	*mock.Call
}

// GetAllCategory is a helper method to define mock.On call
func (_e *MockStorer_Expecter) GetAllCategory() *MockStorer_GetAllCategory_Call {
	return &MockStorer_GetAllCategory_Call{Call: _e.mock.On("GetAllCategory")}
}

func (_c *MockStorer_GetAllCategory_Call) Run(run func()) *MockStorer_GetAllCategory_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockStorer_GetAllCategory_Call) Return(_a0 []store.AllCategoryResult, _a1 error) *MockStorer_GetAllCategory_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockStorer_GetAllCategory_Call) RunAndReturn(run func() ([]store.AllCategoryResult, error)) *MockStorer_GetAllCategory_Call {
	_c.Call.Return(run)
	return _c
}

// GetUnmatchedCategory provides a mock function with given fields: filter
func (_m *MockStorer) GetUnmatchedCategory(filter utils.Filter) ([]store.UnmatchedCategoryResult, utils.MetaData, error) {
	ret := _m.Called(filter)

	if len(ret) == 0 {
		panic("no return value specified for GetUnmatchedCategory")
	}

	var r0 []store.UnmatchedCategoryResult
	var r1 utils.MetaData
	var r2 error
	if rf, ok := ret.Get(0).(func(utils.Filter) ([]store.UnmatchedCategoryResult, utils.MetaData, error)); ok {
		return rf(filter)
	}
	if rf, ok := ret.Get(0).(func(utils.Filter) []store.UnmatchedCategoryResult); ok {
		r0 = rf(filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]store.UnmatchedCategoryResult)
		}
	}

	if rf, ok := ret.Get(1).(func(utils.Filter) utils.MetaData); ok {
		r1 = rf(filter)
	} else {
		r1 = ret.Get(1).(utils.MetaData)
	}

	if rf, ok := ret.Get(2).(func(utils.Filter) error); ok {
		r2 = rf(filter)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// MockStorer_GetUnmatchedCategory_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetUnmatchedCategory'
type MockStorer_GetUnmatchedCategory_Call struct {
	*mock.Call
}

// GetUnmatchedCategory is a helper method to define mock.On call
//   - filter utils.Filter
func (_e *MockStorer_Expecter) GetUnmatchedCategory(filter interface{}) *MockStorer_GetUnmatchedCategory_Call {
	return &MockStorer_GetUnmatchedCategory_Call{Call: _e.mock.On("GetUnmatchedCategory", filter)}
}

func (_c *MockStorer_GetUnmatchedCategory_Call) Run(run func(filter utils.Filter)) *MockStorer_GetUnmatchedCategory_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(utils.Filter))
	})
	return _c
}

func (_c *MockStorer_GetUnmatchedCategory_Call) Return(_a0 []store.UnmatchedCategoryResult, _a1 utils.MetaData, _a2 error) *MockStorer_GetUnmatchedCategory_Call {
	_c.Call.Return(_a0, _a1, _a2)
	return _c
}

func (_c *MockStorer_GetUnmatchedCategory_Call) RunAndReturn(run func(utils.Filter) ([]store.UnmatchedCategoryResult, utils.MetaData, error)) *MockStorer_GetUnmatchedCategory_Call {
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
