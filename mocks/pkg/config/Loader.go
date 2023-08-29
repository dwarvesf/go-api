// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	viper "github.com/spf13/viper"
	mock "github.com/stretchr/testify/mock"
)

// Loader is an autogenerated mock type for the Loader type
type Loader struct {
	mock.Mock
}

type Loader_Expecter struct {
	mock *mock.Mock
}

func (_m *Loader) EXPECT() *Loader_Expecter {
	return &Loader_Expecter{mock: &_m.Mock}
}

// Load provides a mock function with given fields: _a0
func (_m *Loader) Load(_a0 viper.Viper) (*viper.Viper, error) {
	ret := _m.Called(_a0)

	var r0 *viper.Viper
	var r1 error
	if rf, ok := ret.Get(0).(func(viper.Viper) (*viper.Viper, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(viper.Viper) *viper.Viper); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*viper.Viper)
		}
	}

	if rf, ok := ret.Get(1).(func(viper.Viper) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Loader_Load_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Load'
type Loader_Load_Call struct {
	*mock.Call
}

// Load is a helper method to define mock.On call
//   - _a0 viper.Viper
func (_e *Loader_Expecter) Load(_a0 interface{}) *Loader_Load_Call {
	return &Loader_Load_Call{Call: _e.mock.On("Load", _a0)}
}

func (_c *Loader_Load_Call) Run(run func(_a0 viper.Viper)) *Loader_Load_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(viper.Viper))
	})
	return _c
}

func (_c *Loader_Load_Call) Return(_a0 *viper.Viper, _a1 error) *Loader_Load_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Loader_Load_Call) RunAndReturn(run func(viper.Viper) (*viper.Viper, error)) *Loader_Load_Call {
	_c.Call.Return(run)
	return _c
}

// NewLoader creates a new instance of Loader. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewLoader(t interface {
	mock.TestingT
	Cleanup(func())
}) *Loader {
	mock := &Loader{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}