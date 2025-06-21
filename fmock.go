package main

import (
	"reflect"
)

type Mock struct {
	calls      int
	calledWith [][]reflect.Value
}

type argVec []reflect.Value

func (av argVec) Equals(other ...any) bool {
	if len(av) != len(other) {
		return false
	}
	if len(av) == 0 {
		return true
	}
	equals := true
	for i, v := range av {
		o := other[i]
		equals = equals && reflect.DeepEqual(v.Interface(), o)
		if !equals {
			break
		}
	}
	return equals
}

func (m Mock) CallArgs(i int) argVec {
	return argVec(m.calledWith[i])
}

func (m Mock) Calls() int {
	return m.calls
}

func (m Mock) Called() bool {
	return m.calls > 0
}

func MakeMock(fptr any) *Mock {
	m := &Mock{}

	// get the value of the underlying function ptr
	funcValue := reflect.ValueOf(fptr).Elem()
	// create fresh pointer to underlying function type
	copy := reflect.New(funcValue.Type()).Elem()
	// assign the value to the fresh pointer. This is a copy of the function
	copy.Set(funcValue)

	wrapper := func(in []reflect.Value) []reflect.Value {
		res := copy.Call(in)
		// Track call data
		m.calledWith = append(m.calledWith, in)
		m.calls++
		return res
	}

	wrapFunc := func(fptr any) {
		fn := reflect.ValueOf(fptr).Elem()
		// swap the out original function for the wrapper
		v := reflect.MakeFunc(fn.Type(), wrapper)
		fn.Set(v)
	}
	wrapFunc(fptr)

	return m
}
