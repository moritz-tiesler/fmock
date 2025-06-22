package main

import (
	"fmt"
	"reflect"
	"slices"
	"testing"
)

func ExecCallBack(cb func(int, int) int, arg int) {
	cb(arg, arg)
}

type Person struct {
	name string
}

func (p Person) Greeting(name string) string {
	return fmt.Sprintf("Hello %s from %s", name, p.name)
}

func TestFuncVar(t *testing.T) {
	mul := func(a, b int) int {
		return a * b
	}

	mock := MakeMock(&mul)

	for n := range 3 {
		ExecCallBack(mul, n)
		args := mock.CallArgs(n)
		if !args.Equals(n, n) {
			t.Errorf("wrong args in mock, expected=%v, got=%v\n", []int{n, n}, args)
		}
	}

	type data struct {
		i int
		m map[string]int
	}

	aFunction := func(i int, s string, d data) (int, string) {
		return i, s
	}

	callMe := func(f func(int, string, data) (int, string), i int, s string, d data) (int, string) {
		return f(i, s, d)
	}

	mock = MakeMock(&aFunction)

	argVec := []any{4, "yo", data{3, map[string]int{"one": 1}}}
	callMe(aFunction, argVec[0].(int), argVec[1].(string), argVec[2].(data))

	args := mock.CallArgs(0)
	if !args.Equals(argVec...) {
		t.Errorf("wrong args in mock, expected=%v, got=%v\n",
			argVec, args,
		)
	}
	results := mock.CallResults(0)
	if !results.Equals(4, "yo") {
		t.Errorf("wrong results in mock, expected=%v, got=%v\n",
			"4", results,
		)
	}

}

func TestMockReceiverFunc(t *testing.T) {
	p := Person{name: "Dave"}
	mockGreeting := p.Greeting
	mock := MakeMock(&mockGreeting)
	mockGreeting("Bob")
	if !mock.Called() {
		t.Errorf("expected at least one call")
	}
	args := mock.CallArgs(0)
	if !args.Equals("Bob") {
		t.Errorf("wrong args in mock, expected=%v, got=%v\n",
			"Bob", args,
		)
	}
	returns := mock.CallResults(0)
	if !returns.Equals("Hello Bob from Dave") {
		t.Errorf("wrong results in mock, expected=%v, got=%v\n",
			"Hello Bob from Dave", returns,
		)
	}

}

func TestMockGlobalReceiverFunc(t *testing.T) {

	tests := []struct {
		p      Person
		args   string
		result string
	}{
		{
			Person{"Dave"},
			"Bob",
			"Hello Bob from Dave",
		},
		{
			Person{"Bob"},
			"Dave",
			"Hello Dave from Bob",
		},
	}

	f := Person.Greeting
	mock := MakeMock(&f)

	for i, tt := range tests {
		f(tt.p, tt.args)
		args := mock.CallArgs(i)
		if !args.Equals(tt.p, tt.args) {
			t.Errorf("wrong args in mock, expected=%v, got=%v\n",
				tt.args, args,
			)
		}
		results := mock.CallResults(i)
		if !results.Equals(tt.result) {
			t.Errorf("wrong args in mock, expected=%v, got=%v\n",
				tt.result, results,
			)
		}
	}
}

func F(i int) string {
	return fmt.Sprintf("F with %d", i)
}

func TestMockFunc(t *testing.T) {
	mockedF := F
	mock := MakeMock(&mockedF)
	mockedF(3)
	args := mock.CallArgs(0)
	if !args.Equals(3) {
		t.Errorf("wrong args in mock, expected=%v, got=%v\n",
			3, args,
		)
	}
	returns := mock.CallResults(0)
	if !returns.Equals("F with 3") {
		t.Errorf("wrong results in mock, expected=%v, got=%v\n",
			"F with 3", returns,
		)
	}
}

func TestMockAsParameter(t *testing.T) {
	s := []int{0, 1, 2, 3}

	delF := func(i int) bool {
		return i > 1
	}

	mock := MakeMock(&delF)

	slices.DeleteFunc(s, delF)

	if !reflect.DeepEqual(s, []int{0, 1, 0, 0}) {
		t.Errorf("mock created wrong result, expected=%v, got=%v", []int{0, 1, 0, 0}, s)
	}

	if !mock.Called() {
		t.Errorf("expected mock to have been called")
	}
	if len(mock.Calls()) != len(s) {
		t.Errorf("wrong number of mock calls, expected=%d, got=%d", len(s), len(mock.Calls()))
	}
	for i, c := range mock.Calls() {
		if !c.Args.Equals(i) {
			t.Errorf("wrong args in mock call, expected=%v, got=%v", []int{0}, c.Args)
		}
	}
}
