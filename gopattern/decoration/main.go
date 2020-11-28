package main

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"runtime"
	"time"
)

// EXAMPLE 1
func decorator(f func(s string)) func(s string) {
	return func(s string) {
		fmt.Println("started")
		f(s)
		fmt.Println("done")
	}
}

func Hello(s string) {
	fmt.Println(s)
}

// EXAMPLE 2
func Sum1(start, end int64) int64 {
	var sum int64
	sum = 0
	if start > end {
		start, end = end, start
	}
	for i := start; i < end; i++ {
		sum += i
	}

	return sum
}

type SumFunc func(int64, int64) int64

func getFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func timedSumFunc(f SumFunc) SumFunc {
	return func(start int64, end int64) int64 {
		defer func(t time.Time) {
			fmt.Printf("---- Time Elapsed (%s): %v ---\n", getFunctionName(f), time.Since(t))
		}(time.Now())

		return f(start, end)
	}
}

// eg 3 httpserver middle
func WithAuthCookie(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("=======with auth cookie=========")
		cookie := &http.Cookie{Name:"Auth", Value:"Pass", Path:"/"}
		http.SetCookie(w, cookie)

		h(w, r)
	}
}

func WithHeader(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("=======with auth cookie=========")
		//cookie := &http.Cookie{Name:"Auth", Value:"Pass", Path:"/"}
		//http.SetCookie(w, cookie)
		w.Header().Set("", "")
		h(w, r)
	}
}

type HttpHandlerDecorator func(http.HandlerFunc) http.HandlerFunc

func Handler(h http.HandlerFunc, decors ...HttpHandlerDecorator) http.HandlerFunc {
	for i := range decors {
		d := decors[len(decors)-1-i]
		h = d(h)
	}
	return h
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello world" + r.URL.Path)
}

// ex 4
func Decorator(decoPtr, fn interface{}) (err error) {
	var decoratedFunc, targetFunc reflect.Value
	if decoPtr == nil ||
		reflect.TypeOf(decoPtr).Kind() != reflect.Ptr ||
			reflect.ValueOf(decoPtr).Elem().Kind() != reflect.Func {
		err = fmt.Errorf("Need a function pointer!")
		return
	}

	decoratedFunc = reflect.ValueOf(decoPtr).Elem()
	targetFunc = reflect.ValueOf(fn)
	if targetFunc.Kind() != reflect.Func {
		err = fmt.Errorf("Need a function!")
		return
	}

	v := reflect.MakeFunc(targetFunc.Type(),
		func(in []reflect.Value) (out []reflect.Value) {
			fmt.Println("before")
			if targetFunc.Type().IsVariadic() {
				out = targetFunc.CallSlice(in)
			} else {
				out = targetFunc.Call(in)
			}
			fmt.Println("after")
			return
		})

	decoratedFunc.Set(v)

	return
}




func main() {
	// ex 1
	decorator(Hello)("HELLO WORLD")

	// ex 2
	sum := timedSumFunc(Sum1)
	fmt.Printf("%d \n", sum(1, 1000000))

	// ex 3
	http.HandleFunc("/v1/hello", Handler(hello, WithAuthCookie, WithHeader))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("LisenAndServer: ", err)
	}
}


