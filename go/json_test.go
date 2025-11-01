package na01_test

import (
	"na01"
	"testing"
)

type Base struct {
	P int
	Q int
}

type Ext struct {
	Base
	R int
}

func TestStringify(t *testing.T) {
	pi := 3.14
	five := 5
	tests := []struct {
		name    string // description of this test case
		data    any
		want    string
		wantErr bool
	}{{name: "Struct",
		data: struct{ A string }{A: "Hello, world!"},
		want: `{"A":"Hello, world!"}`},
		{name: "Array",
			data: []int{3, 14, 159, 2, 6, 535},
			want: "[3,14,159,2,6,535]"},
		{name: "Number",
			data: 3.14,
			want: "3.14"},
		{name: "Number PTR",
			data: &pi,
			want: "3.14"},
		{name: "NULL",
			data: nil,
			want: "null"},
		{name: "EmptyPTR",
			data: struct {
				M *int
				N *int `json:",omitempty"`
			}{M: &five},
			want: `{"M":5}`},
		{name: "Omit0",
			data: struct {
				M int
				N int `json:",omitzero"`
			}{M: 6 * 7},
			want: `{"M":42}`},
		{name: "Embedding",
			data: Ext{Base: Base{P: 7, Q: 8}, R: 9},
			want: `{"P":7,"Q":8,"R":9}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := na01.Stringify(tt.data)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Stringify() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Stringify() succeeded unexpectedly")
			}
			if got != tt.want {
				t.Errorf("Stringify() = %v, want %v", got, tt.want)
			}
		})
	}
}
