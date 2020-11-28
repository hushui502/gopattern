package main

import (
	"errors"
	"fmt"
	"testing"
)

func TestPipeError(t *testing.T) {
	pipe := Pipe(
		func(x int) (int, error) {
			if x == 0 {
				return 0, errors.New("x should not be zero")
			}
			return x, nil
		},
		func(x int) float32 { return 100.0 / float32(x) },
		func(x float32) string { return fmt.Sprintf("%f\n", x) },)

	result, err := pipe(3)
	expect := "33.333332"
	if err != nil {
		t.Fatal(err)
	}
	if result.(string) != expect {
		t.Fatalf("pipeline failed! expect %v got %v", expect, result)
	}
}
