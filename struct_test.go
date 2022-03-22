package gocsv

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type person struct {
	Name      string `csv:"Name"`
	Age       uint   `csv:"Age"`
	Height    int    `csv:"Height"`
	IsTeacher bool   `csv:"Is Teacher"`
}

func TestMapToStruct(t *testing.T) {
	type args struct {
		src             map[string]string
		suppressError   bool
		caseInsensitive bool
	}
	tests := []struct {
		name string
		args args
		want *person
	}{
		{
			name: "unmarshal to struct",
			args: args{
				src: map[string]string{
					"Name":       "Jojo",
					"Age":        "22",
					"Height":     "188",
					"Is Teacher": "false",
				},
			},
			want: &person{
				Name:      "Jojo",
				Age:       22,
				Height:    188,
				IsTeacher: false,
			},
		},
		{
			name: "with empty value",
			args: args{
				src: map[string]string{
					"Name":       "Jojo",
					"Age":        "",
					"Height":     "188",
					"Is Teacher": "false",
				},
			},
			want: &person{
				Name:      "Jojo",
				Age:       0,
				Height:    188,
				IsTeacher: false,
			},
		},
		{
			name: "with interface out",
			args: args{
				src: map[string]string{
					"Name":       "Jojo",
					"Age":        "",
					"Height":     "188",
					"Is Teacher": "false",
				},
			},
			want: &person{
				Name:      "Jojo",
				Age:       0,
				Height:    188,
				IsTeacher: false,
			},
		},
		{
			name: "with invalid value error",
			args: args{
				src: map[string]string{
					"Name":       "Jojo",
					"Age":        "",
					"Height":     "L",
					"Is Teacher": "false",
				},
				suppressError: true,
			},
			want: &person{
				Name:   "Jojo",
				Age:    0,
				Height: 0,
			},
		},
		{
			name: "with caseInsensitive true",
			args: args{
				src: map[string]string{
					"Name":       "Jojo",
					"age":        "",
					"height":     "188",
					"is teacher": "true",
				},
				caseInsensitive: true,
			},
			want: &person{
				Name:      "Jojo",
				Age:       0,
				Height:    188,
				IsTeacher: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			var out *person
			if tt.args.suppressError {
				out, err = MapToStruct[person](tt.args.src, WithSuppressError(true))
			} else if tt.args.caseInsensitive {
				out, err = MapToStruct[person](tt.args.src, WithCaseInsensitive(true))
			} else {
				out, err = MapToStruct[person](tt.args.src)
			}

			assert.NoError(t, err)
			assert.Equal(t, out, tt.want)
		})
	}
}
