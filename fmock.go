package fmock

import (
	"fmt"
	"reflect"
	"strings"
)

type Mock struct {
	numCalls   int
	calledWith [][]reflect.Value
	calls      []functionCall
}

type valueSlice []reflect.Value

func (vs valueSlice) String() string {
	var sb strings.Builder
	sb.Grow(len(vs))
	sb.WriteString("[")
	for i, v := range vs {
		sb.WriteString(fmt.Sprintf("%+v", v))
		if i != len(vs)-1 {
			sb.WriteString(" ")
		}
	}
	sb.WriteString("]")
	return strings.TrimSpace(sb.String())
}

type functionCall struct {
	Args    valueSlice
	Results valueSlice
}

func (vs valueSlice) Equals(other ...any) bool {
	if len(vs) != len(other) {
		return false
	}
	if len(vs) == 0 {
		return true
	}
	equals := true
	for i, v := range vs {
		o := other[i]
		equals = equals && reflect.DeepEqual(v.Interface(), o)
		if !equals {
			break
		}
	}
	return equals
}

func (m Mock) CallArgs(i int) valueSlice {
	return valueSlice(m.calls[i].Args)
}

func (m Mock) CallResults(i int) valueSlice {
	return valueSlice(m.calls[i].Results)
}

func (m Mock) Calls() []functionCall {
	return m.calls
}

func (m Mock) Called() bool {
	return len(m.calls) > 0
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
		m.calls = append(m.calls, functionCall{in, res})
		m.numCalls++
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
