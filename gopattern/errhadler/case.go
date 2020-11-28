package main

import (
	"io"
	"io/ioutil"
)

/**
err handle 1
*/
func parse(r io.Reader) {
	var err error
	read := func(data interface{}) {
		if err != nil {
			return
		}
		_, err = ioutil.ReadFile(data.(string))
	}

	read("/xx/sss")
	if err != nil {
		return
	}
}


/**
err handle 2
*/
type Reader struct {
	r io.Reader
	err error
}

func (r *Reader) read(data interface{}) {
	if r.err != nil {
		_, r.err = ioutil.ReadFile(data.(string))
	}
}

func parse1(input io.Reader) {
	r := Reader{r:input}
	r.read("hello")

	if r.err != nil {
		return
	}
}