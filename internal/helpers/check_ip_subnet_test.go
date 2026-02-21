package helpers

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckIPSubnet(t *testing.T) {
	type args struct {
		ip            string
		trustedSubnet string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "valid ip",
			args: args{
				ip:            "192.168.1.1",
				trustedSubnet: "192.168.1.0/24",
			},
			want:    true,
			wantErr: assert.NoError,
		},
		{
			name: "invalid ip",
			args: args{
				ip:            "192.16811",
				trustedSubnet: "192.168.1.0/24",
			},
			want:    false,
			wantErr: assert.Error,
		},
		{
			name: "invalid trusted subnet",
			args: args{
				ip:            "192.168.1.1",
				trustedSubnet: "192",
			},
			want:    false,
			wantErr: assert.Error,
		},
		{
			name: "ip not in subnet",
			args: args{
				ip:            "92.68.1.1",
				trustedSubnet: "192.168.1.0/24",
			},
			want:    false,
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CheckIPSubnet(tt.args.ip, tt.args.trustedSubnet)
			if !tt.wantErr(t, err, fmt.Sprintf("CheckIPSubnet(%v, %v)", tt.args.ip, tt.args.trustedSubnet)) {
				return
			}
			assert.Equalf(t, tt.want, got, "CheckIPSubnet(%v, %v)", tt.args.ip, tt.args.trustedSubnet)
		})
	}
}
