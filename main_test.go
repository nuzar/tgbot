package main

import (
	"testing"
)

func TestReverse(t *testing.T) {
	type args struct {
		in string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "zero",
			args: args{""},
			want: "",
		},
		{
			name: "one",
			args: args{"a"},
			want: "a",
		},
		{
			name: "two",
			args: args{"ab"},
			want: "ba",
		},
		{
			name: "three",
			args: args{"abc"},
			want: "cba",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := reverse(tt.args.in); got != tt.want {
				t.Errorf("Reverse() = %v, want %v", got, tt.want)
			}
		})
	}
}
