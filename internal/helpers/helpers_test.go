package helpers

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRandStringRunes(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "request length 1",
			args: args{
				n: 1,
			},
			want: 1,
		},
		{
			name: "request length 10",
			args: args{
				n: 10,
			},
			want: 10,
		},
		{
			name: "request length -1",
			args: args{
				n: -1,
			},
			want: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			str := RandStringRunes(tt.args.n)
			assert.Equal(t, tt.want, len(str))
		})
	}
}
