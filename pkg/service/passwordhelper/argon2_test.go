package passwordhelper

import (
	"testing"
)

func Test_arg2_Compare(t *testing.T) {
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
				hashedPassword: "lTeA/erQXwNjcF/oA/NhTY24DfdS9/7gI7LsxNFUBm4",
				salt:           "vFN5CmKh+7K5zpSjYgK3Cg",
			},
			want: true,
		},
		"fail": {
			args: args{
				password:       "invalid",
				hashedPassword: "lTeA/erQXwNjcF/oA/NhTY24DfdS9/7gI7LsxNFUBm4",
				salt:           "vFN5CmKh+7K5zpSjYgK3Cg",
			},
			want: false,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			h := newArgon2Default()
			if got := h.Compare(tt.args.password, tt.args.hashedPassword, tt.args.salt); got != tt.want {
				t.Errorf("arg2.Compare() = %v, want %v", got, tt.want)
			}
		})
	}
}
