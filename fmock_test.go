package main

import "testing"

func ExecCallBack(cb func(int, int) int, arg int) {
	cb(arg, arg)
}

func TestXxx(t *testing.T) {

}
