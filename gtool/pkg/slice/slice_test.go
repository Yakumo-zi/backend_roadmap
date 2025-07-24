package slice

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"
)

func TestReduce(t *testing.T) {
	type args[T any, V any] struct {
		arr  []T
		f    func(V, T) V
		init V
	}
	type testCase[T any, V any] struct {
		name string
		args args[T, V]
		want V
	}
	tests := []testCase[int, string]{{
		name: "reduce",
		args: args[int, string]{
			arr: []int{1, 2, 3},
			f: func(s string, i int) string {
				if s == "" {
					return strconv.Itoa(i)
				}
				return fmt.Sprintf("%s-%d", s, i)
			},
			init: "",
		},
		want: "1-2-3",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Reduce(tt.args.arr, tt.args.f, tt.args.init); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Reduce() = %v, want %v", got, tt.want)
			}
		})
	}
}
