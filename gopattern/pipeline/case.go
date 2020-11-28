package main

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"runtime"
)

/*
	basic pipeline
*/
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


/*
	generic pipeline
*/

// errType is the tyep of error interface
var errType = reflect.TypeOf((*error)(nil)).Elem()

// Pipeline is the func type for the pipeline result
type Pipeline func(...interface{}) (interface{}, error)

func empty(...interface{}) (interface{}, error) {
	return nil, nil
}

func getFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

// generic pipeline
func Pipe(fns ...interface{}) Pipeline {
	if len(fns) == 0 {
		return empty
	}

	return func(args ...interface{}) (interface{}, error) {
		var inputs []reflect.Value
		for _, arg := range args {
			inputs = append(inputs, reflect.ValueOf(arg))
		}

		for _, fn := range fns {
			outputs := reflect.ValueOf(fn).Call(inputs)
			inputs = inputs[:0] // clean inputs
			fnType := reflect.TypeOf(fn)

			for oIdx, output := range outputs {
				if fnType.Out(oIdx).Implements(errType) {
					if output.IsNil() {
						continue
					}
					err := fmt.Errorf("%s() failed: %w", getFunctionName(fn),
						output.Interface().(error))
					return nil, err
				}
				inputs = append(inputs, output)
			}
		}
		return inputs[0].Interface(), nil
	}
}

/*
	channel pipeline
*/
func echo(nums []int) <-chan int {
	out := make(chan int)

	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()

	return out
}

func sum(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		var sum = 0
		for n := range in {
			sum += n
		}
		out <- sum
		close(out)
	}()

	return out
}

func sq(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * n
		}
		close(out)
	}()

	return out
}


func main() {
	//http.HandleFunc("/v1/hello", Handler(hello, WithAuthCookie, WithHeader))
	//err := http.ListenAndServe(":8080", nil)
	//if err != nil {
	//	log.Fatal("LisenAndServer: ", err)
	//}

	var nums = []int{1, 2, 3, 4}
	for n := range sum(sq(echo(nums))) {
		fmt.Println(n)
	}
}
