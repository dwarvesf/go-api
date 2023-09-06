package passwordhelper

import (
	"testing"
)

func Test_scryptImpl_Compare(t *testing.T) {

	type args struct {
		password       string
		hashedPassword string
		salt           string
	}
	tests := map[string]struct {
		args args
		want bool
	}{
		"success": {
			args: args{
				password:       "password",
				hashedPassword: "8VHtKBiCms/r0gV790RTxP8JzPlFFguRCG1goGqXQJg",
				salt:           "TizK61lmaHY",
			},
			want: true,
		},
		"failure": {
			args: args{
				password:       "invalid",
				hashedPassword: "8VHtKBiCms/r0gV790RTxP8JzPlFFguRCG1goGqXQJg",
				salt:           "TizK61lmaHY",
			},
			want: false,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			h := newScryptDefault()
			if got := h.Compare(tt.args.password, tt.args.hashedPassword, tt.args.salt); got != tt.want {
				t.Errorf("case %v: scryptImpl.Compare() = %v, want %v", name, got, tt.want)
			}
		})
	}
}
