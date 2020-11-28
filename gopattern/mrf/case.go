package mrf

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

/**
map/reduce/filter
*/
// map reduce 1
func MapUpCase(arr []string, fn func(s string) string) []string {
	var newArray = []string{}
	for _, item := range arr {
		newArray = append(newArray, fn(item))
	}

	return newArray
}

func MapLen(arr []string, fn func(s string) int) []int {
	var newArray = []int{}
	for _, item := range arr {
		newArray = append(newArray, fn(item))
	}

	return newArray
}

func Reduce(arr []string, fn func(s string) int) int {
	sum := 0
	for _, it := range arr {
		sum += fn(it)
	}

	return sum
}

func Filter(arr []int, fn func(n int) bool) []int {
	var newArray = []int{}
	for _, it := range arr {
		if fn(it) {
			newArray = append(newArray, it)
		}
	}

	return newArray
}

type Employee struct {
	Name string
	Age int
	Vacation int
	Salary int
}

var employeeList = []Employee{
	{"hufan", 22, 10, 22222},
	{"bob", 32, 23, 100203},
	{"alice", 44, 33, 322222},
}

func EmployeeCountIf(list []Employee, fn func(e *Employee) bool) int {
	count := 0
	for i, _ := range list {
		if fn(&list[i]) {
			count += 1
		}
	}

	return count
}

/**
generic map
*/
func transform(slice, function interface{}, inPalce bool) interface{} {
	sliceInType := reflect.ValueOf(slice)
	if sliceInType.Kind() != reflect.Slice {
		panic("not slice")
	}

	fn := reflect.ValueOf(function)
	elemType := sliceInType.Type().Elem()
	if !vertifyFuncSignature(fn, elemType, nil) {
		panic("transform: function must be of type func(" + sliceInType.Type().Elem().String() + ") outputElementType")
	}

	sliceOutType := sliceInType
	if !inPalce {
		sliceOutType = reflect.MakeSlice(reflect.SliceOf(fn.Type().Out(0)), sliceInType.Len(), sliceInType.Len())
	}
	for i := 0; i < sliceInType.Len(); i++ {
		sliceOutType.Index(i).Set(fn.Call([]reflect.Value{sliceInType.Index(i)})[0])
	}

	return sliceOutType.Interface()
}

func vertifyFuncSignature(fn reflect.Value, types ...reflect.Type) bool {
	if fn.Kind() != reflect.Func {
		return false
	}

	if (fn.Type().NumIn() != len(types)-1) || (fn.Type().NumOut() != 1) {
		return false
	}

	for i := 0; i < len(types)-1; i++ {
		if fn.Type().In(i) != types[i] {
			return false
		}
	}

	outType := types[len(types)-1]
	if outType != nil && fn.Type().Out(0) != outType {
		return false
	}

	return true
}

func Transform(slice, fn interface{}) interface{} {
	return transform(slice, fn ,false)
}

func TransformInPlace(slice, fn interface{}) interface{} {
	return transform(slice, fn, true)
}

/*
	generic reduce
*/
func Reduce1(slice, pairFunc, zero interface{}) interface{} {
	sliceInType := reflect.ValueOf(slice)
	if sliceInType.Kind() != reflect.Slice {
		panic("reduce: wrong type, not slice")
	}

	len := sliceInType.Len()
	if len == 0 {
		return zero
	} else if len == 1 {
		return sliceInType.Index(0)
	}

	elemType := sliceInType.Type().Elem()
	fn := reflect.ValueOf(pairFunc)
	if !vertifyFuncSignature(fn, elemType, elemType, elemType) {
		t := elemType.String()
		panic("reduce: function must be of type func(" + t + ", " + t + ") " + t)
	}

	var ins [2]reflect.Value
	ins[0] = sliceInType.Index(0)
	ins[1] = sliceInType.Index(1)
	out := fn.Call(ins[:])[0]

	for i := 2; i < len; i++ {
		ins[0] = out
		ins[1] = sliceInType.Index(i)
		out = fn.Call(ins[:])[0]
	}

	return out.Interface()
}

func mul(a, b int) int {
	return a * b
}

/**
generic filter
*/
var boolType = reflect.ValueOf(true).Type()

func filter1(slice, function interface{}, inPlace bool) (interface{}, int) {
	sliceInType := reflect.ValueOf(slice)
	if sliceInType.Kind() != reflect.Slice {
		panic("filter: wrong type, not a slice")
	}

	fn := reflect.ValueOf(function)
	elemType := sliceInType.Type().Elem()
	if !vertifyFuncSignature(fn, elemType, boolType) {
		panic("filter: function must be of type func(" + elemType.String() + ") bool")
	}

	var which []int
	for i := 0; i < sliceInType.Len(); i++ {
		if fn.Call([]reflect.Value{sliceInType.Index(i)})[0].Bool() {
			which = append(which, i)
		}
	}

	out := sliceInType

	if !inPlace {
		out = reflect.MakeSlice(sliceInType.Type(), len(which), len(which))
	}
	for i := range which {
		out.Index(i).Set(sliceInType.Index(which[i]))
	}

	return out.Interface(), len(which)
}

func Filter1(slice, fn interface{}) interface{}  {
	result, _ := filter1(slice, fn, false)
	return result
}

func FilterInPlace1(slicePtr, fn interface{}) {
	in := reflect.ValueOf(slicePtr)
	if in.Kind() != reflect.Ptr {
		panic("FilterInPlace: wrong type, " + "not a pointer to slice")
	}
	_, n := filter1(in.Elem().Interface(), fn, false)
	in.Elem().SetLen(n)
}

func isEven(a int) bool {
	return a%2 == 0
}

func isOddString(s string) bool {
	i, _ := strconv.ParseInt(s, 10, 32)
	return i%2 == 1
}


func main() {
	var list = []string{"heLlO", "Chen", "Hu"}

	// map
	x := MapUpCase(list, func(s string) string {
		return strings.ToUpper(s)
	})
	fmt.Printf("%v\n", x)

	y := MapLen(list, func(s string) int {
		return len(s)
	})
	fmt.Printf("%v\n", y)

	// reduce
	xlen := Reduce(list, func(s string) int {
		return len(s)
	})
	fmt.Printf("%v\n", xlen)

	// filter
	var inset = []int{1, 2, 3, 4, 5, 6}
	out := Filter(inset, func(n int) bool {
		return n % 2 == 1
	})
	fmt.Printf("%v\n", out)

	// generic map
	old := EmployeeCountIf(employeeList, func(e *Employee) bool {
		return e.Age > 30
	})
	fmt.Printf("old people(>30): %v\n", old)

	// generic reduce
	a := make([]int, 10)
	for i := range a {
		a[i] = i + 1
	}
	genericReduceOut := Reduce1(a, mul, 1).(int)
	fmt.Printf("genericReduceOut: %v\n", genericReduceOut)

	// generic filter
	a1 := []int{1, 2, 3, 4}
	result := Filter1(a1, isEven)
	fmt.Printf("%v\n", result)

	s1 := []string{"1", "2", "3", "4"}
	result = Filter1(s1, isOddString)
	fmt.Printf("%v\n", result)
}