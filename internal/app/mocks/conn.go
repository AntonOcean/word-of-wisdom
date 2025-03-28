// Code generated by mockery v2.53.0. DO NOT EDIT.

package mocks

import (
	net "net"

	mock "github.com/stretchr/testify/mock"

	time "time"
)

// Conn is an autogenerated mock type for the Conn type
type Conn struct {
	mock.Mock
}

type Conn_Expecter struct {
	mock *mock.Mock
}

func (_m *Conn) EXPECT() *Conn_Expecter {
	return &Conn_Expecter{mock: &_m.Mock}
}

// Close provides a mock function with no fields
func (_m *Conn) Close() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Close")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Conn_Close_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Close'
type Conn_Close_Call struct {
	*mock.Call
}

// Close is a helper method to define mock.On call
func (_e *Conn_Expecter) Close() *Conn_Close_Call {
	return &Conn_Close_Call{Call: _e.mock.On("Close")}
}

func (_c *Conn_Close_Call) Run(run func()) *Conn_Close_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Conn_Close_Call) Return(_a0 error) *Conn_Close_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Conn_Close_Call) RunAndReturn(run func() error) *Conn_Close_Call {
	_c.Call.Return(run)
	return _c
}

// LocalAddr provides a mock function with no fields
func (_m *Conn) LocalAddr() net.Addr {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for LocalAddr")
	}

	var r0 net.Addr
	if rf, ok := ret.Get(0).(func() net.Addr); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(net.Addr)
		}
	}

	return r0
}

// Conn_LocalAddr_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'LocalAddr'
type Conn_LocalAddr_Call struct {
	*mock.Call
}

// LocalAddr is a helper method to define mock.On call
func (_e *Conn_Expecter) LocalAddr() *Conn_LocalAddr_Call {
	return &Conn_LocalAddr_Call{Call: _e.mock.On("LocalAddr")}
}

func (_c *Conn_LocalAddr_Call) Run(run func()) *Conn_LocalAddr_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Conn_LocalAddr_Call) Return(_a0 net.Addr) *Conn_LocalAddr_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Conn_LocalAddr_Call) RunAndReturn(run func() net.Addr) *Conn_LocalAddr_Call {
	_c.Call.Return(run)
	return _c
}

// Read provides a mock function with given fields: b
func (_m *Conn) Read(b []byte) (int, error) {
	ret := _m.Called(b)

	if len(ret) == 0 {
		panic("no return value specified for Read")
	}

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func([]byte) (int, error)); ok {
		return rf(b)
	}
	if rf, ok := ret.Get(0).(func([]byte) int); ok {
		r0 = rf(b)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func([]byte) error); ok {
		r1 = rf(b)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Conn_Read_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Read'
type Conn_Read_Call struct {
	*mock.Call
}

// Read is a helper method to define mock.On call
//   - b []byte
func (_e *Conn_Expecter) Read(b interface{}) *Conn_Read_Call {
	return &Conn_Read_Call{Call: _e.mock.On("Read", b)}
}

func (_c *Conn_Read_Call) Run(run func(b []byte)) *Conn_Read_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].([]byte))
	})
	return _c
}

func (_c *Conn_Read_Call) Return(n int, err error) *Conn_Read_Call {
	_c.Call.Return(n, err)
	return _c
}

func (_c *Conn_Read_Call) RunAndReturn(run func([]byte) (int, error)) *Conn_Read_Call {
	_c.Call.Return(run)
	return _c
}

// RemoteAddr provides a mock function with no fields
func (_m *Conn) RemoteAddr() net.Addr {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for RemoteAddr")
	}

	var r0 net.Addr
	if rf, ok := ret.Get(0).(func() net.Addr); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(net.Addr)
		}
	}

	return r0
}

// Conn_RemoteAddr_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RemoteAddr'
type Conn_RemoteAddr_Call struct {
	*mock.Call
}

// RemoteAddr is a helper method to define mock.On call
func (_e *Conn_Expecter) RemoteAddr() *Conn_RemoteAddr_Call {
	return &Conn_RemoteAddr_Call{Call: _e.mock.On("RemoteAddr")}
}

func (_c *Conn_RemoteAddr_Call) Run(run func()) *Conn_RemoteAddr_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Conn_RemoteAddr_Call) Return(_a0 net.Addr) *Conn_RemoteAddr_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Conn_RemoteAddr_Call) RunAndReturn(run func() net.Addr) *Conn_RemoteAddr_Call {
	_c.Call.Return(run)
	return _c
}

// SetDeadline provides a mock function with given fields: t
func (_m *Conn) SetDeadline(t time.Time) error {
	ret := _m.Called(t)

	if len(ret) == 0 {
		panic("no return value specified for SetDeadline")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(time.Time) error); ok {
		r0 = rf(t)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Conn_SetDeadline_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetDeadline'
type Conn_SetDeadline_Call struct {
	*mock.Call
}

// SetDeadline is a helper method to define mock.On call
//   - t time.Time
func (_e *Conn_Expecter) SetDeadline(t interface{}) *Conn_SetDeadline_Call {
	return &Conn_SetDeadline_Call{Call: _e.mock.On("SetDeadline", t)}
}

func (_c *Conn_SetDeadline_Call) Run(run func(t time.Time)) *Conn_SetDeadline_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(time.Time))
	})
	return _c
}

func (_c *Conn_SetDeadline_Call) Return(_a0 error) *Conn_SetDeadline_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Conn_SetDeadline_Call) RunAndReturn(run func(time.Time) error) *Conn_SetDeadline_Call {
	_c.Call.Return(run)
	return _c
}

// SetReadDeadline provides a mock function with given fields: t
func (_m *Conn) SetReadDeadline(t time.Time) error {
	ret := _m.Called(t)

	if len(ret) == 0 {
		panic("no return value specified for SetReadDeadline")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(time.Time) error); ok {
		r0 = rf(t)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Conn_SetReadDeadline_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetReadDeadline'
type Conn_SetReadDeadline_Call struct {
	*mock.Call
}

// SetReadDeadline is a helper method to define mock.On call
//   - t time.Time
func (_e *Conn_Expecter) SetReadDeadline(t interface{}) *Conn_SetReadDeadline_Call {
	return &Conn_SetReadDeadline_Call{Call: _e.mock.On("SetReadDeadline", t)}
}

func (_c *Conn_SetReadDeadline_Call) Run(run func(t time.Time)) *Conn_SetReadDeadline_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(time.Time))
	})
	return _c
}

func (_c *Conn_SetReadDeadline_Call) Return(_a0 error) *Conn_SetReadDeadline_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Conn_SetReadDeadline_Call) RunAndReturn(run func(time.Time) error) *Conn_SetReadDeadline_Call {
	_c.Call.Return(run)
	return _c
}

// SetWriteDeadline provides a mock function with given fields: t
func (_m *Conn) SetWriteDeadline(t time.Time) error {
	ret := _m.Called(t)

	if len(ret) == 0 {
		panic("no return value specified for SetWriteDeadline")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(time.Time) error); ok {
		r0 = rf(t)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Conn_SetWriteDeadline_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetWriteDeadline'
type Conn_SetWriteDeadline_Call struct {
	*mock.Call
}

// SetWriteDeadline is a helper method to define mock.On call
//   - t time.Time
func (_e *Conn_Expecter) SetWriteDeadline(t interface{}) *Conn_SetWriteDeadline_Call {
	return &Conn_SetWriteDeadline_Call{Call: _e.mock.On("SetWriteDeadline", t)}
}

func (_c *Conn_SetWriteDeadline_Call) Run(run func(t time.Time)) *Conn_SetWriteDeadline_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(time.Time))
	})
	return _c
}

func (_c *Conn_SetWriteDeadline_Call) Return(_a0 error) *Conn_SetWriteDeadline_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Conn_SetWriteDeadline_Call) RunAndReturn(run func(time.Time) error) *Conn_SetWriteDeadline_Call {
	_c.Call.Return(run)
	return _c
}

// Write provides a mock function with given fields: b
func (_m *Conn) Write(b []byte) (int, error) {
	ret := _m.Called(b)

	if len(ret) == 0 {
		panic("no return value specified for Write")
	}

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func([]byte) (int, error)); ok {
		return rf(b)
	}
	if rf, ok := ret.Get(0).(func([]byte) int); ok {
		r0 = rf(b)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func([]byte) error); ok {
		r1 = rf(b)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Conn_Write_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Write'
type Conn_Write_Call struct {
	*mock.Call
}

// Write is a helper method to define mock.On call
//   - b []byte
func (_e *Conn_Expecter) Write(b interface{}) *Conn_Write_Call {
	return &Conn_Write_Call{Call: _e.mock.On("Write", b)}
}

func (_c *Conn_Write_Call) Run(run func(b []byte)) *Conn_Write_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].([]byte))
	})
	return _c
}

func (_c *Conn_Write_Call) Return(n int, err error) *Conn_Write_Call {
	_c.Call.Return(n, err)
	return _c
}

func (_c *Conn_Write_Call) RunAndReturn(run func([]byte) (int, error)) *Conn_Write_Call {
	_c.Call.Return(run)
	return _c
}

// NewConn creates a new instance of Conn. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewConn(t interface {
	mock.TestingT
	Cleanup(func())
}) *Conn {
	mock := &Conn{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
