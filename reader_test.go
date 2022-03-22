package gocsv

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"strings"
	"testing"
)

var in = `first_name,last_name,username
"Rob","Pike",rob
Ken,Thompson,ken
"Robert","Griesemer","gri"
`
var f = strings.NewReader(in)
var want = [][]string{
	[]string{"first_name", "last_name", "username"},
	[]string{"Rob", "Pike", "rob"},
	[]string{"Ken", "Thompson", "ken"},
	[]string{"Robert", "Griesemer", "gri"},
}

func TestReadAll(t *testing.T) {
	type args struct {
		f io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    [][]string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "simple ReadAll",
			args:    args{f: f},
			want:    want,
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadAll(tt.args.f)
			if !tt.wantErr(t, err, fmt.Sprintf("ReadAll(%v)", tt.args.f)) {
				return
			}
			assert.Equalf(t, tt.want, got, "ReadAll(%v)", tt.args.f)
		})
	}
}
