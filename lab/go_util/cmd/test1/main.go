package main

import (
	"fmt"
	"reflect"
)

func Compose(funcs ...interface{}) interface{} {
	if len(funcs) == 0 {
		panic("no functions provided")
	}

	// Validate all functions
	var funcVals []reflect.Value
	for _, f := range funcs {
		v := reflect.ValueOf(f)
		t := v.Type()
		if t.Kind() != reflect.Func {
			panic("Compose: all arguments must be functions")
		}
		if t.NumIn() != 1 || t.NumOut() != 1 {
			panic("Compose: each function must have exactly 1 input and 1 output")
		}
		funcVals = append(funcVals, v)
	}

	// The final function type: same as the first function
	finalType := funcVals[0].Type()

	// Build the composed function
	composed := reflect.MakeFunc(finalType, func(args []reflect.Value) (results []reflect.Value) {
		val := args[0]
		// apply in reverse order: last to first
		for i := len(funcVals) - 1; i >= 0; i-- {
			val = funcVals[i].Call([]reflect.Value{val})[0]
		}
		return []reflect.Value{val}
	})

	return composed.Interface()
}

func main() {
	// Example: int -> int functions
	double := func(x int) int { return x * 2 }
	inc := func(x int) int { return x + 1 }
	square := func(x int) int { return x * x }

	f := Compose(double, inc, square).(func(int) int)

	fmt.Println(f(2)) // output = double(inc(square(2))) = (2² + 1)*2 = 10
}
