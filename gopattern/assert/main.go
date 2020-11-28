package main

import (
	"fmt"
	"reflect"
)


/**
	type assert
 */
type Container []interface{}

func (c *Container) Put(elem interface{}) {
	*c = append(*c, elem)
}

func (c *Container) Get() interface{} {
	elem := (*c)[0]
	(*c)[0] = nil
	*c = (*c)[1:]

	return elem
}

/**
	reflection
 */
type Cabinet struct {
	s reflect.Value
}

func NewCabinet(t reflect.Type) *Cabinet {
	return &Cabinet{
		s: reflect.MakeSlice(reflect.SliceOf(t), 0, 10),
	}
}

func (c *Cabinet) Put(val interface{}) {
	if reflect.ValueOf(val).Type() != c.s.Type().Elem() {
		panic(fmt.Sprintf("Put: cannot put a %T into a slice of %s", val, c.s.Type().Elem()))
	}
	c.s = reflect.Append(c.s, reflect.ValueOf(val))
}

func (c *Cabinet) Get(retref interface{}) {
	retref = c.s.Index(0)
	c.s = c.s.Slice(1, c.s.Len())
}

func main() {
	// type assert
	iniContainer := &Container{}
	iniContainer.Put(3)
	iniContainer.Put(22)

	elem, ok := iniContainer.Get().(int)
	if !ok {
		fmt.Println("Unable to read an int from intContainer")
	}
	fmt.Printf("assertExample: %d (%T)\n", elem, elem)

	// reflection
	f := 3.33333
	c := NewCabinet(reflect.TypeOf(f))
	c.Put(f)
}
















