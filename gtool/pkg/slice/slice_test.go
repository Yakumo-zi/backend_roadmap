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

func TestDistinct(t *testing.T) {
	type args[T comparable] struct {
		arr []T
	}
	type testCase[T comparable] struct {
		name string
		args args[T]
		want []T
	}
	tests := []testCase[int]{
		{
			name: "distinct",
			args: args[int]{
				arr: []int{1, 2, 3, 1, 2, 4, 2, 3, 6, 1, 7},
			},
			want: []int{1, 2, 3, 4, 6, 7},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Distinct(tt.args.arr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Distinct() = %v, want %v", got, tt.want)
			}
		})
	}
}
