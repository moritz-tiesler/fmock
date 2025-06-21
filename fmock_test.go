package main

import (
	"testing"
)

func ExecCallBack(cb func(int, int) int, arg int) {
	cb(arg, arg)
}

func TestXxx(t *testing.T) {
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
