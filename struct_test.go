package gocsv

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type person struct {
	Name      string
	Age       uint
	Height    int
	IsTeacher bool `csv:"Is Teacher"`
}

func Test_ReadT(t *testing.T) {
	type args struct {
		file            string
		suppressError   bool
		caseInsensitive bool
	}
	tests := []struct {
		name  string
		file  string
		args  args
		want  []person
		error bool
	}{
		{
			name: "unmarshal to struct",
			args: args{
				file: "fixtures/1.csv",
			},
			want: []person{
				{
					Name:      "Jojo",
					Age:       uint(22),
					Height:    188,
					IsTeacher: false,
				},
			},
			error: false,
		},
		{
			name: "with empty value for some header[s]",
			args: args{
				file: "fixtures/2.csv",
			},
			want: []person{
				{
					Name:      "Jojo",
					Age:       0,
					Height:    188,
					IsTeacher: false,
				},
			},
			error: false,
		},
		{
			name: "with empty value for all headers",
			args: args{
				file:          "fixtures/5.csv",
				suppressError: false,
			},
			want:  []person{},
			error: false,
		},
		{
			name: "with empty value and empty headers",
			args: args{
				file:          "fixtures/6.csv",
				suppressError: false,
			},
			want:  []person{},
			error: false,
		},
		{
			name: "with invalid value and suppress error",
			args: args{
				file:          "fixtures/3.csv",
				suppressError: true,
			},
			want: []person{
				{
					Name:   "Jojo",
					Age:    0,
					Height: 0,
				},
			},
			error: false,
		},
		{
			name: "with invalid value and return error",
			args: args{
				file:          "fixtures/3.csv",
				suppressError: false,
			},
			want:  []person{},
			error: true,
		},
		{
			name: "with case insensitive header",
			args: args{
				file: "fixtures/4.csv",
			},
			want: []person{
				{
					Name:      "Jojo",
					Age:       0,
					Height:    188,
					IsTeacher: true,
				},
			},
			error: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := os.Open(tt.args.file)
			assert.NoError(t, err)
			defer f.Close()

			var persons []person
			if tt.args.suppressError {
				persons, err = Read[person](f, WithSuppressError(true))
			} else {
				persons, err = Read[person](f)
			}

			if tt.error {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, persons)
			}
		})
	}
}
