package filter

import "reflect"

func Apply(slice, function interface{}) interface{} {
	return apply(slice, function, false)
}

func ApplyInPlace(slice, function interface{}) {
	apply(slice, function, true)
}

func Choose(slice, function interface{}) interface{} {
	out, _ := chooseOrDrop(slice, function, false, true)
	return out
}

func Drop(slice, function interface{}) interface{} {
	out, _ := chooseOrDrop(slice, function, false, false)
	return out
}

func ChooseInPlace(pointerToSlice, function interface{}) {
	chooseOrDropInPlace(pointerToSlice, function, true)
}

func DropInPlace(pointerToSlice, function interface{}) {
	chooseOrDropInPlace(pointerToSlice, function, false)
}

func apply(slice, function interface{}, inPlace bool) interface{} {
	if strSlice, ok := slice.([]string); ok {
		if strFn, ok := function.(func(string) string); ok {
			r := strSlice
			if !inPlace {
				r = make([]string, len(strSlice))
			}
			for i, s := range strSlice {
				r[i] = strFn(s)
			}
			return r
		}
	}
	in := reflect.ValueOf(slice)
	if in.Kind() != reflect.Slice {
		panic("apply: not slice")
	}
	fn := reflect.ValueOf(function)
	elemType := in.Type().Elem()
	if !goodFunc(fn, elemType, nil) {
		panic("apply: function must be of type func(" + in.Type().Elem().String() + ") outputElemType")
	}
	out := in
	if !inPlace {
		out = reflect.MakeSlice(reflect.SliceOf(fn.Type().Out(0)), in.Len(), in.Len())
	}
	var ins [1]reflect.Value
	for i := 0; i < in.Len(); i++ {
		ins[0] = in.Index(i)
		out.Index(i).Set(fn.Call(ins[:])[0])
	}

	return out.Interface()
}

var boolType = reflect.ValueOf(true).Type()

func chooseOrDrop(slice, function interface{}, inPlace, truth bool) (interface{}, int) {
	if strSlice, ok := slice.([]string); ok {
		if strFn, ok := function.(func(string) bool); ok {
			var r []string
			if inPlace {
				r = strSlice[:0]
			}
			for _, v := range strSlice {
				if strFn(v) == truth {
					r = append(r, v)
				}
			}
			return r, len(r)
		}
	}

	in := reflect.ValueOf(slice)
	if in.Kind() != reflect.Slice {
		panic("choose/drop: not slice")
	}
	fn := reflect.ValueOf(function)
	elemType := in.Type().Elem();
	if !goodFunc(fn, elemType, boolType) {
		panic("choose/drop: function must be of type func(" + elemType.String() + ") bool")
	}
	var which []int
	var ins [1]reflect.Value
	for i := 0; i < in.Len(); i++ {
		ins[0] = in.Index(i)
		if fn.Call(ins[:])[0].Bool() == truth {
			which = append(which, i)
		}
	}

	out := in
	if !inPlace {
		out = reflect.MakeSlice(in.Type(), len(which), len(which))
	}

	for i, v := range which {
		out.Index(i).Set(in.Index(v))
	}

	return out.Interface(), len(which)
}

func chooseOrDropInPlace(slice, function interface{}, truth bool) {
	inp := reflect.ValueOf(slice)
	if inp.Kind() != reflect.Ptr {
		panic("choose/drop: not pointer to slice")
	}
	_, n := chooseOrDrop(inp.Elem().Interface(), function, true, truth)
	inp.Elem().SetLen(n)
}

func goodFunc(fn reflect.Value, types ...reflect.Type) bool {
	if fn.Kind() != reflect.Func {
		return false
	}

	if fn.Type().NumIn() != len(types)-1 || fn.Type().NumOut() != 1 {
		return false
	}

	for i := 0; i < len(types)-1; i++ {
		if fn.Type().In(i) != types[i] {
			return false
		}
	}

	outType := types[len(types)-1]
	if outType != nil && outType != fn.Type().Out(0) {
		return false
	}

	return true
}





















