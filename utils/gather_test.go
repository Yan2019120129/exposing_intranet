package utils

import (
	"reflect"
	"testing"
)

func TestDiff(t *testing.T) {
	tests := []struct {
		name string
		a    []int
		b    []int
		want []int
	}{
		{
			name: "basic diff",
			a:    []int{1, 2, 3, 4, 5},
			b:    []int{2, 4},
			want: []int{1, 3, 5},
		},
		{
			name: "a empty",
			a:    []int{},
			b:    []int{1, 2},
			want: []int{},
		},
		{
			name: "b empty",
			a:    []int{1, 2},
			b:    []int{},
			want: []int{1, 2},
		},
		{
			name: "all in b",
			a:    []int{1, 2},
			b:    []int{1, 2, 3},
			want: []int{},
		},
		{
			name: "with duplicates in a",
			a:    []int{1, 1, 2, 3},
			b:    []int{2},
			want: []int{1, 1, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Diff(tt.a, tt.b)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Diff() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntersect(t *testing.T) {
	tests := []struct {
		name string
		a    []int
		b    []int
		want []int
	}{
		{
			name: "basic intersect",
			a:    []int{1, 2, 3, 4},
			b:    []int{2, 4, 6},
			want: []int{2, 4},
		},
		{
			name: "no intersection",
			a:    []int{1, 3, 5},
			b:    []int{2, 4, 6},
			want: []int{},
		},
		{
			name: "a empty",
			a:    []int{},
			b:    []int{1, 2},
			want: []int{},
		},
		{
			name: "b empty",
			a:    []int{1, 2},
			b:    []int{},
			want: []int{},
		},
		{
			name: "with duplicates in a",
			a:    []int{1, 2, 2, 3},
			b:    []int{2},
			want: []int{2, 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Intersect(tt.a, tt.b)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Intersect() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnion(t *testing.T) {
	tests := []struct {
		name string
		a    []int
		b    []int
		want []int
	}{
		{
			name: "basic union",
			a:    []int{1, 2, 3},
			b:    []int{3, 4, 5},
			want: []int{1, 2, 3, 4, 5},
		},
		{
			name: "a empty",
			a:    []int{},
			b:    []int{1, 2},
			want: []int{1, 2},
		},
		{
			name: "b empty",
			a:    []int{1, 2},
			b:    []int{},
			want: []int{1, 2},
		},
		{
			name: "both empty",
			a:    []int{},
			b:    []int{},
			want: []int{},
		},
		{
			name: "with duplicates",
			a:    []int{1, 1, 2},
			b:    []int{2, 2, 3},
			want: []int{1, 2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Union(tt.a, tt.b)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Union() = %v, want %v", got, tt.want)
			}
		})
	}
}
