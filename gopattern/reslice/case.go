package main

import (
	"bytes"
	"fmt"
	"reflect"
)

/*
	slice reallocated
*/
func main() {
	path := []byte("AAAA/BBBBB")
	sepIndex := bytes.IndexByte(path, '/')

	dir1 := path[:sepIndex:sepIndex]
	dir2 := path[sepIndex+1:]

	fmt.Println("dir1 => ", string(dir1))
	fmt.Println("dir2 => ", string(dir2))


	dir1 = append(dir1, "suffix"...)
	fmt.Println("dir1 => ", string(dir1))
	fmt.Println("dir2 => ", string(dir2))
}


/*
	deep comparison
*/
type data struct {
	num int
	check [10]func() bool
	doit func() bool
	m map[string]string
	bytes []byte
}

func main() {
	v1 := data{}
	v2 := data{}
	fmt.Println("deep equal v1 == v2 ? ", reflect.DeepEqual(v1, v2))

	m1 := map[string]string{"one":"a", "two":"b"}
	m2 := map[string]string{"two":"b", "one":"a"}
	fmt.Println("v1 == v2 ? ", reflect.DeepEqual(m1, m2))

	s1 := []int{1, 2, 3}
	s2 := []int{1, 2, 3}
	fmt.Println("s1 == s2 ? ", reflect.DeepEqual(s1, s2))
}


/*
	Interface pattern
*/
type Country struct {
	WithName
}

type City struct {
	WithName
}

type WithName struct {
	Name string
}

type Printable interface {
	PrintStr()
}

func (w WithName) PrintStr() {
	fmt.Println(w.Name)
}

//func main() {
//	city := City{WithName{"BEIJING"}}
//	country := Country{WithName{"CHINA"}}
//	city.PrintStr()
//	country.PrintStr()
//}



/*
	var AInterface = (*AImpl)(nil)
*/
type Shape interface {
	Sides() int
	Area() int
}

type Square struct {
	len int
}

func (s *Square) Area() int {
	return s.len * s.len
}

func (s *Square) Sides() int {
	return 4
}

//func main() {
//	var _ Shape = (*Square)(nil)
//	s := Square{len:4}
//	fmt.Printf("%d\n", s.Sides())
//}