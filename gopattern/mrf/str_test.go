package mrf

import (
	"bytes"
	"fmt"
	"math/rand"
	"reflect"
	"strconv"
	"testing"
)

/**
	perfer strconv over fmt
 */
// 143ns/op
func BenchmarkFmt(t *testing.B) {
	for i := 0; i < t.N; i++ {
		_ = fmt.Sprint(rand.Int())
	}
}

// 72.4ns/op
func BenchmarkStrconv(t *testing.B) {
	for i := 0; i < t.N; i++ {
		_ = strconv.Itoa(rand.Int())
	}
}

/**
	avoid string-to-byte conversion
 */
// 6.23ns/op
func BenchmarkStr2Byte(t *testing.B) {
	str := "hello"
	for i := 0; i < t.N; i++ {
		_ = []byte(str)
	}
}

// 4.25ns/op
func BenchmarkByte2Str(t *testing.B) {
	data := []byte("hello")
	for i := 0; i < t.N; i++ {
		_ = string(data)
	}
}


/**
	specify slice capacity
 */
// 4410ns/op
func BenchmarkNotSpecifyingCap(t *testing.B) {
	for i := 0; i < t.N; i++ {
		data := make([]int, 0)
		for k := 0; k < 1000; k++ {
			data = append(data, k)
		}
	}
}

// 698ns/op
func BenchmarkSpecifyingCap(t *testing.B) {
	for i := 0; i < t.N; i++ {
		data := make([]int, 0, 1000)
		for k := 0; k < 1000; k++ {
			data = append(data, k)
		}
	}
}


/*
	use stringbuilder or stringbuffer
*/
// 63542ns/op
func BenchmarkNoBufferorBuilder(t *testing.B) {
	var strLen = 3
	var str string
	for i := 0; i < t.N; i++ {
		for n := 0; n < strLen; n++ {
			str += "x"
		}
	}
}

// 8.12ns/op
// 21.5ns/op
func BenchmarkBufferorBuilder(t *testing.B) {
	var strLen = 3
	//var builder strings.Builder
	//for i := 0; i < t.N; i++ {
	//	for n := 0; n < strLen; n++ {
	//		builder.WriteString("x")
	//	}
	//}

	var buffer bytes.Buffer
	for i := 0; i < t.N; i++ {
		for n := 0; n < strLen; n++ {
			buffer.WriteString("x")
		}
	}
}

func TestTransform(t *testing.T) {
	list := []string{"1", "2", "3", "4"}
	expect := []string{"111", "222", "333", "444"}
	result := Transform(list, func(a string) string {
		return a + a + a
	})

	if !reflect.DeepEqual(expect, result) {
		t.Fatalf("Transform failed: expect %v got %v",
			expect, result)
	}
}

func TestMapEmployee(t *testing.T) {
	var list = []Employee{
		{"chen", 22, 11, 1111},
		{"hu", 24, 22, 2222},
		{"li", 34, 22, 3333},
	}
	var expect = []Employee{
		{"chen", 23, 12, 1111},
		{"hu", 25, 23, 2222},
		{"li", 35, 23, 3333},
	}

	result := TransformInPlace(list, func(e Employee) Employee {
		e.Age += 1
		e.Vacation += 1
		return e
	})

	if !reflect.DeepEqual(expect, result) {
		t.Fatalf("Transform failed: expect %v got %v",
			expect, result)
	}
}


















